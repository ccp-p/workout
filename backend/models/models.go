package models

import (
	"time"
)

// Exercise 动作模型
type Exercise struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"imageUrl"`
	BodyPart    string    `json:"bodyPart"`    // 身体部位：胸、背、腿、肩、臂等
	CreatedAt   time.Time `json:"createdAt"`
}

// ExerciseSet 组模型
type ExerciseSet struct {
	ExerciseID string `json:"exerciseId"`
	Sets       int    `json:"sets"`       // 组数
	Reps       int    `json:"reps"`       // 每组次数
	Weight     int    `json:"weight"`     // 重量
	RestTime   int    `json:"restTime"`   // 组间休息时间(秒)
}

// Workout 训练计划模型
type Workout struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	BodyPart    string        `json:"bodyPart"`
	Exercises   []ExerciseSet `json:"exercises"`
	CreatedAt   time.Time     `json:"createdAt"`
}

// WorkoutSession 训练记录模型
type WorkoutSession struct {
	ID           string              `json:"id"`
	WorkoutID    string              `json:"workoutId"`
	Date         time.Time           `json:"date"`
	StartTime    time.Time           `json:"startTime"`
	EndTime      time.Time           `json:"endTime"`
	TotalTime    int                 `json:"totalTime"`    // 总时间(秒)
	Exercises    []CompletedExercise `json:"exercises"`
	Notes        string              `json:"notes"`
	IsCompleted  bool                `json:"isCompleted"`
}

// CompletedExercise 完成的动作记录
type CompletedExercise struct {
	ExerciseID       string `json:"exerciseId"`
	CompletedSets    int    `json:"completedSets"`
	CompletedReps    []int  `json:"completedReps"`    // 每组实际完成次数
	ActualRestTimes  []int  `json:"actualRestTimes"`  // 每组实际休息时间
	IsCompleted      bool   `json:"isCompleted"`
}

// Statistics 统计数据模型
type Statistics struct {
	Date         time.Time          `json:"date"`
	TotalTime    int                `json:"totalTime"`
	WorkoutCount int                `json:"workoutCount"`
	BodyParts    map[string]int     `json:"bodyParts"`    // 各部位训练次数
	Exercises    map[string]int     `json:"exercises"`    // 各动作训练次数
}
