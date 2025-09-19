package presenter

import (
	"strconv"
	"time"
	"workout-tracker/models"
)

type WorkoutPresenter struct{}

func NewWorkoutPresenter() *WorkoutPresenter {
	return &WorkoutPresenter{}
}

// 格式化训练时间
func (p *WorkoutPresenter) FormatDuration(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	if hours > 0 {
		return strconv.Itoa(hours) + "h " + strconv.Itoa(minutes) + "m " + strconv.Itoa(secs) + "s"
	} else if minutes > 0 {
		return strconv.Itoa(minutes) + "m " + strconv.Itoa(secs) + "s"
	}
	return strconv.Itoa(secs) + "s"
}

// 格式化日期
func (p *WorkoutPresenter) FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// 格式化时间
func (p *WorkoutPresenter) FormatTime(t time.Time) string {
	return t.Format("15:04:05")
}

// 训练统计响应
type StatisticsResponse struct {
	TodayStats   DayStats              `json:"todayStats"`
	WeekStats    WeekStats             `json:"weekStats"`
	MonthStats   MonthStats            `json:"monthStats"`
	BodyPartData []BodyPartStatistics  `json:"bodyPartData"`
}

type DayStats struct {
	Date         string `json:"date"`
	TotalTime    string `json:"totalTime"`
	WorkoutCount int    `json:"workoutCount"`
	Exercises    int    `json:"exercises"`
}

type WeekStats struct {
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	TotalTime    string `json:"totalTime"`
	WorkoutCount int    `json:"workoutCount"`
	AvgPerDay    string `json:"avgPerDay"`
}

type MonthStats struct {
	Month        string `json:"month"`
	TotalTime    string `json:"totalTime"`
	WorkoutCount int    `json:"workoutCount"`
	AvgPerDay    string `json:"avgPerDay"`
}

type BodyPartStatistics struct {
	BodyPart string `json:"bodyPart"`
	Count    int    `json:"count"`
	Percent  int    `json:"percent"`
}

// 格式化统计数据
func (p *WorkoutPresenter) FormatStatistics(sessions []models.WorkoutSession) StatisticsResponse {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := today.AddDate(0, 0, -int(today.Weekday()))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 今日统计
	todayStats := p.calculateDayStats(sessions, today)
	
	// 本周统计
	weekStats := p.calculateWeekStats(sessions, weekStart, today.AddDate(0, 0, 1))
	
	// 本月统计
	monthStats := p.calculateMonthStats(sessions, monthStart, today.AddDate(0, 0, 1))
	
	// 身体部位统计
	bodyPartData := p.calculateBodyPartStats(sessions)

	return StatisticsResponse{
		TodayStats:   todayStats,
		WeekStats:    weekStats,
		MonthStats:   monthStats,
		BodyPartData: bodyPartData,
	}
}

func (p *WorkoutPresenter) calculateDayStats(sessions []models.WorkoutSession, date time.Time) DayStats {
	var totalTime, workoutCount, exerciseCount int
	
	for _, session := range sessions {
		if p.isSameDay(session.Date, date) && session.IsCompleted {
			totalTime += session.TotalTime
			workoutCount++
			exerciseCount += len(session.Exercises)
		}
	}

	return DayStats{
		Date:         p.FormatDate(date),
		TotalTime:    p.FormatDuration(totalTime),
		WorkoutCount: workoutCount,
		Exercises:    exerciseCount,
	}
}

func (p *WorkoutPresenter) calculateWeekStats(sessions []models.WorkoutSession, start, end time.Time) WeekStats {
	var totalTime, workoutCount int
	
	for _, session := range sessions {
		if session.Date.After(start) && session.Date.Before(end) && session.IsCompleted {
			totalTime += session.TotalTime
			workoutCount++
		}
	}

	avgPerDay := 0
	if workoutCount > 0 {
		avgPerDay = totalTime / 7
	}

	return WeekStats{
		StartDate:    p.FormatDate(start),
		EndDate:      p.FormatDate(end.AddDate(0, 0, -1)),
		TotalTime:    p.FormatDuration(totalTime),
		WorkoutCount: workoutCount,
		AvgPerDay:    p.FormatDuration(avgPerDay),
	}
}

func (p *WorkoutPresenter) calculateMonthStats(sessions []models.WorkoutSession, start, end time.Time) MonthStats {
	var totalTime, workoutCount int
	
	for _, session := range sessions {
		if session.Date.After(start) && session.Date.Before(end) && session.IsCompleted {
			totalTime += session.TotalTime
			workoutCount++
		}
	}

	daysInMonth := end.AddDate(0, 0, -1).Day()
	avgPerDay := 0
	if workoutCount > 0 {
		avgPerDay = totalTime / daysInMonth
	}

	return MonthStats{
		Month:        start.Format("2006-01"),
		TotalTime:    p.FormatDuration(totalTime),
		WorkoutCount: workoutCount,
		AvgPerDay:    p.FormatDuration(avgPerDay),
	}
}

func (p *WorkoutPresenter) calculateBodyPartStats(sessions []models.WorkoutSession) []BodyPartStatistics {
	bodyPartCount := make(map[string]int)
	totalCount := 0

	for _, session := range sessions {
		if session.IsCompleted {
			for _, exercise := range session.Exercises {
				if exercise.IsCompleted {
					// 这里需要根据exerciseId获取bodyPart，简化处理
					bodyPartCount["全身"]++
					totalCount++
				}
			}
		}
	}

	var result []BodyPartStatistics
	for bodyPart, count := range bodyPartCount {
		percent := 0
		if totalCount > 0 {
			percent = (count * 100) / totalCount
		}
		result = append(result, BodyPartStatistics{
			BodyPart: bodyPart,
			Count:    count,
			Percent:  percent,
		})
	}

	return result
}

func (p *WorkoutPresenter) isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
