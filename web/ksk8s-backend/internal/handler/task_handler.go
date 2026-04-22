package handler

import (
	"net/http"
	"strconv"

	"github.com/easzlab/ksk8s/internal/model"
	"github.com/easzlab/ksk8s/internal/service"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// List returns all tasks for a cluster.
func (h *TaskHandler) List(c *gin.Context) {
	clusterID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	tasks, err := h.taskService.ListTasksByCluster(clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// Get returns a single task.
func (h *TaskHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("taskId"), 10, 64)
	task, err := h.taskService.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// RunStep starts a task for a specific step.
func (h *TaskHandler) RunStep(c *gin.Context) {
	clusterID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	step := c.Param("step")

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	taskType := "setup"
	if step == "start" || step == "stop" || step == "upgrade" || step == "backup" || step == "restore" || step == "destroy" || step == "kca-renew" {
		taskType = step
	}

	// cluster_admin cannot destroy or restore clusters
	if model.NormalizeRole(role.(string)) == model.RoleClusterAdmin && (taskType == "destroy" || taskType == "restore") {
		c.JSON(http.StatusForbidden, gin.H{"error": "destroy/restore requires platform_admin"})
		return
	}

	task, err := h.taskService.StartTask(clusterID, step, taskType, userID.(int64))
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, task)
}

// Abort aborts a running task.
func (h *TaskHandler) Abort(c *gin.Context) {
	taskID, _ := strconv.ParseInt(c.Param("taskId"), 10, 64)
	if err := h.taskService.AbortTask(taskID); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task aborted"})
}

// Approve approves a completed task to allow next step.
func (h *TaskHandler) Approve(c *gin.Context) {
	taskID, _ := strconv.ParseInt(c.Param("taskId"), 10, 64)
	userID, _ := c.Get("user_id")

	if err := h.taskService.ApproveTask(taskID, userID.(int64)); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task approved"})
}

// Retry retries a failed or aborted task.
func (h *TaskHandler) Retry(c *gin.Context) {
	clusterID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	step := c.Param("step")
	userID, _ := c.Get("user_id")

	taskType := "setup"
	if step == "start" || step == "stop" || step == "upgrade" || step == "backup" || step == "restore" || step == "destroy" || step == "kca-renew" {
		taskType = step
	}

	task, err := h.taskService.RetryTask(clusterID, step, taskType, userID.(int64))
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, task)
}

// Logs returns historical logs for a task from DB.
func (h *TaskHandler) Logs(c *gin.Context) {
	taskID, _ := strconv.ParseInt(c.Param("taskId"), 10, 64)
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)
	limitStr := c.DefaultQuery("limit", "1000")
	limit, _ := strconv.Atoi(limitStr)

	logRepo := h.taskService.NewLogRepo()
	logs, err := logRepo.ListByTask(taskID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch logs"})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// Status returns the current status of a task including worker info.
func (h *TaskHandler) Status(c *gin.Context) {
	taskID, _ := strconv.ParseInt(c.Param("taskId"), 10, 64)
	task, err := h.taskService.GetTask(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	worker := h.taskService.GetWorker(taskID)
	pid := 0
	if worker != nil {
		pid = worker.PID()
	}

	c.JSON(http.StatusOK, gin.H{
		"task":    task,
		"pid":     pid,
		"running": worker != nil,
	})
}
