package handler

import (
	"net/http"
	"strconv"

	"github.com/easzlab/ksk8s/internal/service"
	"github.com/gin-gonic/gin"
)

type NodeOpsHandler struct {
	nodeOpsService *service.NodeOpsService
}

func NewNodeOpsHandler(taskService *service.TaskService) *NodeOpsHandler {
	return &NodeOpsHandler{
		nodeOpsService: service.NewNodeOpsService(taskService),
	}
}

func (h *NodeOpsHandler) AddNode(c *gin.Context) {
	clusterID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	userID, _ := c.Get("user_id")

	var req struct {
		GroupName   string `json:"group_name" binding:"required"`
		IPAddress   string `json:"ip_address" binding:"required"`
		K8sNodename string `json:"k8s_nodename"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.nodeOpsService.AddNode(clusterID, req.GroupName, req.IPAddress, req.K8sNodename, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if task != nil {
		c.JSON(http.StatusAccepted, gin.H{"message": "node added, ezctl task started", "task": task})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "node added to config"})
	}
}

func (h *NodeOpsHandler) RemoveNode(c *gin.Context) {
	clusterID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	userID, _ := c.Get("user_id")

	var req struct {
		GroupName string `json:"group_name" binding:"required"`
		IPAddress string `json:"ip_address" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.nodeOpsService.RemoveNode(clusterID, req.GroupName, req.IPAddress, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if task != nil {
		c.JSON(http.StatusAccepted, gin.H{"message": "node removal started", "task": task})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "node removed from config"})
	}
}
