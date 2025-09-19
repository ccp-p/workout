package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"workout-tracker/models"
	"workout-tracker/presenter"
	"workout-tracker/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WorkoutHandler struct {
	repo      *repository.FileRepository
	presenter *presenter.WorkoutPresenter
	uploadDir string
}

func NewWorkoutHandler(repo *repository.FileRepository, presenter *presenter.WorkoutPresenter, uploadDir string) *WorkoutHandler {
	return &WorkoutHandler{
		repo:      repo,
		presenter: presenter,
		uploadDir: uploadDir,
	}
}

// Exercise handlers
func (h *WorkoutHandler) GetExercises(c *gin.Context) {
	exercises, err := h.repo.GetAllExercises()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exercises)
}

func (h *WorkoutHandler) CreateExercise(c *gin.Context) {
	var exercise models.Exercise
	if err := c.ShouldBindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exercise.ID = uuid.New().String()
	exercise.CreatedAt = time.Now()

	if err := h.repo.SaveExercise(exercise); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, exercise)
}

func (h *WorkoutHandler) UpdateExercise(c *gin.Context) {
	id := c.Param("id")
	var exercise models.Exercise
	if err := c.ShouldBindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exercise.ID = id
	if err := h.repo.SaveExercise(exercise); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exercise)
}

func (h *WorkoutHandler) DeleteExercise(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.DeleteExercise(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Exercise deleted successfully"})
}

// Workout handlers
func (h *WorkoutHandler) GetWorkouts(c *gin.Context) {
	workouts, err := h.repo.GetAllWorkouts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workouts)
}

func (h *WorkoutHandler) CreateWorkout(c *gin.Context) {
	var workout models.Workout
	if err := c.ShouldBindJSON(&workout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workout.ID = uuid.New().String()
	workout.CreatedAt = time.Now()

	if err := h.repo.SaveWorkout(workout); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, workout)
}

func (h *WorkoutHandler) GetWorkout(c *gin.Context) {
	id := c.Param("id")
	workout, err := h.repo.GetWorkoutByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, workout)
}

func (h *WorkoutHandler) UpdateWorkout(c *gin.Context) {
	id := c.Param("id")
	var workout models.Workout
	if err := c.ShouldBindJSON(&workout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workout.ID = id
	if err := h.repo.SaveWorkout(workout); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workout)
}

func (h *WorkoutHandler) DeleteWorkout(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.DeleteWorkout(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Workout deleted successfully"})
}

// Session handlers
func (h *WorkoutHandler) CreateSession(c *gin.Context) {
	var session models.WorkoutSession
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session.ID = uuid.New().String()
	session.Date = time.Now()
	session.StartTime = time.Now()

	if err := h.repo.SaveSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, session)
}

func (h *WorkoutHandler) GetSession(c *gin.Context) {
	id := c.Param("id")
	session, err := h.repo.GetSessionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, session)
}

func (h *WorkoutHandler) UpdateSession(c *gin.Context) {
	id := c.Param("id")
	var session models.WorkoutSession
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session.ID = id
	if session.IsCompleted && session.EndTime.IsZero() {
		session.EndTime = time.Now()
		session.TotalTime = int(session.EndTime.Sub(session.StartTime).Seconds())
	}

	if err := h.repo.SaveSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *WorkoutHandler) GetSessions(c *gin.Context) {
	// 获取查询参数
	startDate := c.Query("start")
	endDate := c.Query("end")
	
	sessions, err := h.repo.GetAllSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 如果有日期过滤参数
	if startDate != "" && endDate != "" {
		start, err1 := time.Parse("2006-01-02", startDate)
		end, err2 := time.Parse("2006-01-02", endDate)
		if err1 == nil && err2 == nil {
			filteredSessions, err := h.repo.GetSessionsByDateRange(start, end.AddDate(0, 0, 1))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			sessions = filteredSessions
		}
	}

	c.JSON(http.StatusOK, sessions)
}

// Statistics handler
func (h *WorkoutHandler) GetStatistics(c *gin.Context) {
	sessions, err := h.repo.GetAllSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stats := h.presenter.FormatStatistics(sessions)
	c.JSON(http.StatusOK, stats)
}

// Upload handler
func (h *WorkoutHandler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// 生成唯一文件名
	filename := uuid.New().String() + filepath.Ext(header.Filename)
	filePath := filepath.Join(h.uploadDir, filename)

	// 确保上传目录存在
	os.MkdirAll(h.uploadDir, 0755)

	// 保存文件
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 返回文件URL
	fileURL := "/uploads/" + filename
	c.JSON(http.StatusOK, gin.H{"url": fileURL})
}
