# 健身训练Web应用

一个简单易用的健身训练管理系统，包含后台管理和移动端训练功能。

## 功能特性

### 后台管理功能
- ✅ 动作管理：添加、编辑、删除训练动作
- ✅ 图片上传：支持JPG、PNG、GIF格式的动作演示图
- ✅ 训练计划：创建和管理训练计划，设置组数、次数、重量、休息时间
- ✅ 训练记录：查看历史训练记录和进度
- ✅ 数据统计：查看日、周、月训练统计数据

### 移动端训练功能  
- ✅ 训练执行：按照计划进行训练，实时记录
- ✅ 计时功能：总训练时间计时和组间休息计时
- ✅ 进度追踪：可视化训练进度条
- ✅ 震动提醒：休息结束震动提醒（支持的设备）
- ✅ 响应式设计：适配移动设备屏幕

## 技术架构

### 后端
- **语言**: Go 1.21+
- **框架**: Gin Web Framework  
- **存储**: JSON文件存储（简单快速）
- **架构**: MVP模式
  - Model: 数据模型定义
  - Presenter: 数据格式化和业务逻辑
  - View: HTTP API接口

### 前端
- **框架**: Vue 3 (选项式API)
- **HTTP客户端**: Axios
- **样式**: 原生CSS (简约设计)
- **部署**: 纯HTML文件，无需构建

## 项目结构

```
workout/
├── backend/                 # Go后端服务
│   ├── main.go             # 主程序入口
│   ├── go.mod              # Go模块依赖
│   ├── models/             # 数据模型
│   ├── repository/         # 数据访问层
│   ├── presenter/          # 业务逻辑层
│   └── handlers/           # HTTP处理器
├── frontend/               # 前端页面
│   ├── admin.html         # 后台管理页面
│   └── mobile.html        # 移动端训练页面
├── data/                   # JSON数据存储
├── uploads/                # 上传文件存储
└── start.bat              # Windows启动脚本
```

## 快速开始

### 环境要求
- Go 1.21 或更高版本
- 现代浏览器（Chrome、Firefox、Safari、Edge）

### 安装步骤

1. **安装Go环境**
   - 访问 https://golang.org/dl/ 下载并安装Go

2. **启动应用**
   - 双击 `start.bat` 启动服务器
   - 或在命令行中运行：
     ```bash
     cd backend
     go mod tidy
     go run main.go
     ```

3. **访问应用**
   - 后台管理：http://localhost:8080
   - 移动端训练：http://localhost:8080/mobile

## 使用指南

### 1. 动作管理
1. 在后台管理页面点击"动作管理"标签
2. 点击"添加新动作"按钮
3. 填写动作信息：
   - 动作名称（如：俯卧撑）
   - 身体部位（如：胸部）
   - 动作描述（如：标准俯卧撑动作要领）
   - 上传动作图片或GIF演示
4. 保存动作

### 2. 创建训练计划
1. 点击"训练计划"标签
2. 点击"创建新计划"按钮
3. 设置计划信息：
   - 计划名称（如：胸部训练）
   - 目标部位（如：胸部）
   - 计划描述
4. 添加训练动作：
   - 选择动作
   - 设置组数（默认4组）
   - 设置每组次数（默认12次）
   - 设置重量（可选）
   - 设置组间休息时间（默认60秒）
5. 保存计划

### 3. 开始训练
1. 在训练计划中点击"开始训练"按钮
2. 系统会打开移动端训练页面
3. 点击"开始训练"开始总计时
4. 按顺序完成每个动作的每组训练：
   - 点击"开始"按钮开始一组
   - 完成动作后点击"完成"按钮
   - 自动进入组间休息倒计时
   - 休息结束后继续下一组
5. 完成所有动作后查看训练总结

### 4. 查看统计数据
1. 在后台管理页面点击"数据统计"标签
2. 查看今日、本周、本月的训练数据
3. 在"训练记录"标签中查看详细的历史记录

## API接口

### 动作管理
- `GET /api/exercises` - 获取所有动作
- `POST /api/exercises` - 创建新动作
- `PUT /api/exercises/:id` - 更新动作
- `DELETE /api/exercises/:id` - 删除动作

### 训练计划
- `GET /api/workouts` - 获取所有训练计划
- `POST /api/workouts` - 创建新训练计划
- `GET /api/workouts/:id` - 获取特定训练计划
- `PUT /api/workouts/:id` - 更新训练计划
- `DELETE /api/workouts/:id` - 删除训练计划

### 训练记录
- `GET /api/sessions` - 获取训练记录
- `POST /api/sessions` - 创建新训练记录
- `GET /api/sessions/:id` - 获取特定训练记录
- `PUT /api/sessions/:id` - 更新训练记录

### 数据统计
- `GET /api/statistics` - 获取统计数据

### 文件上传
- `POST /api/upload` - 上传文件

## 数据模型

### 动作(Exercise)
```json
{
  "id": "uuid",
  "name": "动作名称",
  "description": "动作描述", 
  "imageUrl": "图片URL",
  "bodyPart": "身体部位",
  "createdAt": "创建时间"
}
```

### 训练计划(Workout)
```json
{
  "id": "uuid",
  "name": "计划名称",
  "description": "计划描述",
  "bodyPart": "目标部位",
  "exercises": [
    {
      "exerciseId": "动作ID",
      "sets": 4,
      "reps": 12,
      "weight": 20,
      "restTime": 60
    }
  ],
  "createdAt": "创建时间"
}
```

### 训练记录(WorkoutSession)
```json
{
  "id": "uuid",
  "workoutId": "训练计划ID",
  "date": "训练日期",
  "startTime": "开始时间",
  "endTime": "结束时间", 
  "totalTime": 1800,
  "exercises": [
    {
      "exerciseId": "动作ID",
      "completedSets": 4,
      "completedReps": [12,12,10,8],
      "actualRestTimes": [60,65,70,0],
      "isCompleted": true
    }
  ],
  "isCompleted": true
}
```

## 扩展计划

- [ ] 用户系统和权限管理
- [ ] 数据库支持（MySQL/PostgreSQL）
- [ ] 训练计划模板分享
- [ ] 社交功能和排行榜
- [ ] 移动端PWA支持
- [ ] 数据导出功能
- [ ] 更多图表和分析功能

## 常见问题

### Q: 如何修改服务器端口？
A: 在 `backend/main.go` 文件中修改 `r.Run(":8080")` 中的端口号。

### Q: 训练数据存储在哪里？
A: 数据以JSON格式存储在 `data/` 目录中，包括 `exercises.json`、`workouts.json`、`sessions.json`。

### Q: 如何备份数据？
A: 复制整个 `data/` 目录即可备份所有训练数据。

### Q: 移动端页面可以离线使用吗？
A: 当前版本需要网络连接，后续版本将支持PWA离线功能。

## 开发者信息

本项目采用MVP架构设计，便于快速开发和后期扩展。如有问题或建议，欢迎反馈。

---

**愉快训练！💪**
