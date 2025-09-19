package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	"workout-tracker/models"
)

type FileRepository struct {
	dataDir string
}

func NewFileRepository(dataDir string) *FileRepository {
	return &FileRepository{dataDir: dataDir}
}

// 确保数据目录存在
func (r *FileRepository) ensureDataDir() error {
	return os.MkdirAll(r.dataDir, 0755)
}

// 读取JSON文件
func (r *FileRepository) readJSONFile(filename string, data interface{}) error {
	filePath := filepath.Join(r.dataDir, filename)
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，返回空数据
		}
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if len(bytes) == 0 {
		return nil // 空文件
	}

	return json.Unmarshal(bytes, data)
}

// 写入JSON文件
func (r *FileRepository) writeJSONFile(filename string, data interface{}) error {
	if err := r.ensureDataDir(); err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	filePath := filepath.Join(r.dataDir, filename)
	return ioutil.WriteFile(filePath, bytes, 0644)
}

// Exercise 相关方法
func (r *FileRepository) GetAllExercises() ([]models.Exercise, error) {
	var exercises []models.Exercise
	err := r.readJSONFile("exercises.json", &exercises)
	return exercises, err
}

func (r *FileRepository) SaveExercise(exercise models.Exercise) error {
	exercises, err := r.GetAllExercises()
	if err != nil {
		return err
	}

	// 检查是否已存在，更新或添加
	found := false
	for i, ex := range exercises {
		if ex.ID == exercise.ID {
			exercises[i] = exercise
			found = true
			break
		}
	}
	if !found {
		exercises = append(exercises, exercise)
	}

	return r.writeJSONFile("exercises.json", exercises)
}

func (r *FileRepository) GetExerciseByID(id string) (*models.Exercise, error) {
	exercises, err := r.GetAllExercises()
	if err != nil {
		return nil, err
	}

	for _, exercise := range exercises {
		if exercise.ID == id {
			return &exercise, nil
		}
	}
	return nil, fmt.Errorf("exercise not found")
}

func (r *FileRepository) DeleteExercise(id string) error {
	exercises, err := r.GetAllExercises()
	if err != nil {
		return err
	}

	for i, exercise := range exercises {
		if exercise.ID == id {
			exercises = append(exercises[:i], exercises[i+1:]...)
			return r.writeJSONFile("exercises.json", exercises)
		}
	}
	return fmt.Errorf("exercise not found")
}

// Workout 相关方法
func (r *FileRepository) GetAllWorkouts() ([]models.Workout, error) {
	var workouts []models.Workout
	err := r.readJSONFile("workouts.json", &workouts)
	return workouts, err
}

func (r *FileRepository) SaveWorkout(workout models.Workout) error {
	workouts, err := r.GetAllWorkouts()
	if err != nil {
		return err
	}

	found := false
	for i, w := range workouts {
		if w.ID == workout.ID {
			workouts[i] = workout
			found = true
			break
		}
	}
	if !found {
		workouts = append(workouts, workout)
	}

	return r.writeJSONFile("workouts.json", workouts)
}

func (r *FileRepository) GetWorkoutByID(id string) (*models.Workout, error) {
	workouts, err := r.GetAllWorkouts()
	if err != nil {
		return nil, err
	}

	for _, workout := range workouts {
		if workout.ID == id {
			return &workout, nil
		}
	}
	return nil, fmt.Errorf("workout not found")
}

func (r *FileRepository) DeleteWorkout(id string) error {
	workouts, err := r.GetAllWorkouts()
	if err != nil {
		return err
	}

	for i, workout := range workouts {
		if workout.ID == id {
			workouts = append(workouts[:i], workouts[i+1:]...)
			return r.writeJSONFile("workouts.json", workouts)
		}
	}
	return fmt.Errorf("workout not found")
}

// WorkoutSession 相关方法
func (r *FileRepository) GetAllSessions() ([]models.WorkoutSession, error) {
	var sessions []models.WorkoutSession
	err := r.readJSONFile("sessions.json", &sessions)
	return sessions, err
}

func (r *FileRepository) SaveSession(session models.WorkoutSession) error {
	sessions, err := r.GetAllSessions()
	if err != nil {
		return err
	}

	found := false
	for i, s := range sessions {
		if s.ID == session.ID {
			sessions[i] = session
			found = true
			break
		}
	}
	if !found {
		sessions = append(sessions, session)
	}

	return r.writeJSONFile("sessions.json", sessions)
}

func (r *FileRepository) GetSessionsByDateRange(start, end time.Time) ([]models.WorkoutSession, error) {
	sessions, err := r.GetAllSessions()
	if err != nil {
		return nil, err
	}

	var result []models.WorkoutSession
	for _, session := range sessions {
		if session.Date.After(start) && session.Date.Before(end) {
			result = append(result, session)
		}
	}
	return result, nil
}

func (r *FileRepository) GetSessionByID(id string) (*models.WorkoutSession, error) {
	sessions, err := r.GetAllSessions()
	if err != nil {
		return nil, err
	}

	for _, session := range sessions {
		if session.ID == id {
			return &session, nil
		}
	}
	return nil, fmt.Errorf("session not found")
}
