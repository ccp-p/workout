package main

import (
	"log"
	"workout-tracker/handlers"
	"workout-tracker/presenter"
	"workout-tracker/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 设置数据目录和上传目录
	dataDir := "../data"
	uploadDir := "../uploads"

	// 初始化仓库和呈现器
	repo := repository.NewFileRepository(dataDir)
	presenter := presenter.NewWorkoutPresenter()
	handler := handlers.NewWorkoutHandler(repo, presenter, uploadDir)

	// 设置路由
	r := gin.Default()

	// 跨域设置
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// 静态文件服务
	r.Static("/uploads", uploadDir)
	r.Static("/static", "../frontend")

	// API 路由
	api := r.Group("/api")
	{
		// 动作相关
		api.GET("/exercises", handler.GetExercises)
		api.POST("/exercises", handler.CreateExercise)
		api.PUT("/exercises/:id", handler.UpdateExercise)
		api.DELETE("/exercises/:id", handler.DeleteExercise)

		// 训练计划相关
		api.GET("/workouts", handler.GetWorkouts)
		api.POST("/workouts", handler.CreateWorkout)
		api.GET("/workouts/:id", handler.GetWorkout)
		api.PUT("/workouts/:id", handler.UpdateWorkout)
		api.DELETE("/workouts/:id", handler.DeleteWorkout)

		// 训练记录相关
		api.GET("/sessions", handler.GetSessions)
		api.POST("/sessions", handler.CreateSession)
		api.GET("/sessions/:id", handler.GetSession)
		api.PUT("/sessions/:id", handler.UpdateSession)

		// 统计相关
		api.GET("/statistics", handler.GetStatistics)

		// 文件上传
		api.POST("/upload", handler.UploadFile)
	}

	// 根路径重定向到后台管理页面
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/static/admin.html")
	})

	// 移动端页面
	r.GET("/mobile", func(c *gin.Context) {
		c.Redirect(302, "/static/mobile.html")
	})

	log.Println("服务器启动在 http://localhost:8080")
	log.Println("后台管理: http://localhost:8080")
	log.Println("移动端训练: http://localhost:8080/mobile")
	
	r.Run(":8080")
}
