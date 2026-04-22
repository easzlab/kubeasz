package service

import (
	"fmt"
	"log"

	"github.com/easzlab/ksk8s/internal/model"
	"github.com/easzlab/ksk8s/internal/repository"
)

// NodeOpsService handles node add/remove operations via ezctl.
type NodeOpsService struct {
	taskService *TaskService
	clusterRepo *repository.ClusterRepository
}

// NewNodeOpsService creates a NodeOpsService.
func NewNodeOpsService(taskService *TaskService) *NodeOpsService {
	return &NodeOpsService{
		taskService: taskService,
		clusterRepo: repository.NewClusterRepository(),
	}
}

// AddNode adds a node to the cluster config and optionally runs ezctl add-*.
func (s *NodeOpsService) AddNode(clusterID int64, group string, ip string, nodename string, userID int64) (*model.Task, error) {
	cluster, err := s.clusterRepo.GetByID(clusterID)
	if err != nil {
		return nil, fmt.Errorf("cluster not found: %w", err)
	}

	// 1. Append node to config
	if err := s.appendNodeToConfig(cluster, group, ip, nodename); err != nil {
		return nil, fmt.Errorf("append node to config: %w", err)
	}

	// 2. If cluster is active, run ezctl add-* command
	if cluster.Status == "active" {
		action := map[string]string{
			"kube_node":   "add-node",
			"kube_master": "add-master",
			"etcd":        "add-etcd",
		}[group]
		if action == "" {
			return nil, fmt.Errorf("unsupported group for ezctl add: %s", group)
		}

		task, err := s.taskService.StartTask(clusterID, action, action, userID, ip)
		if err != nil {
			return nil, fmt.Errorf("start add task: %w", err)
		}
		return task, nil
	}

	return nil, nil
}

// RemoveNode removes a node from the cluster via ezctl del-* and then from config.
func (s *NodeOpsService) RemoveNode(clusterID int64, group string, ip string, userID int64) (*model.Task, error) {
	cluster, err := s.clusterRepo.GetByID(clusterID)
	if err != nil {
		return nil, fmt.Errorf("cluster not found: %w", err)
	}

	// Safety check: keep at least 1 etcd and 1 kube_master
	nodeRepo := repository.NewNodeRepository()
	nodes, _ := nodeRepo.GetByClusterID(clusterID)
	if group == "etcd" {
		etcdCount := 0
		for _, n := range nodes {
			if n.GroupName == "etcd" {
				etcdCount++
			}
		}
		if etcdCount <= 1 {
			return nil, fmt.Errorf("cannot remove the last etcd node")
		}
	}
	if group == "kube_master" {
		masterCount := 0
		for _, n := range nodes {
			if n.GroupName == "kube_master" {
				masterCount++
			}
		}
		if masterCount <= 1 {
			return nil, fmt.Errorf("cannot remove the last kube_master node")
		}
	}

	// If cluster is active, run ezctl del-* first, then remove from config asynchronously
	if cluster.Status == "active" {
		action := map[string]string{
			"kube_node":   "del-node",
			"kube_master": "del-master",
			"etcd":        "del-etcd",
		}[group]
		if action == "" {
			return nil, fmt.Errorf("unsupported group for ezctl del: %s", group)
		}

		task, err := s.taskService.StartTask(clusterID, action, action, userID, ip)
		if err != nil {
			return nil, fmt.Errorf("start del task: %w", err)
		}

		// Watch task completion and remove from config on success
		go s.watchAndRemoveNode(task.ID, clusterID, group, ip)
		return task, nil
	}

	// Cluster not active — just remove from config
	if err := s.removeNodeFromConfig(clusterID, group, ip); err != nil {
		return nil, fmt.Errorf("remove node from config: %w", err)
	}
	return nil, nil
}

func (s *NodeOpsService) appendNodeToConfig(cluster *model.Cluster, group string, ip string, nodename string) error {
	nodeRepo := repository.NewNodeRepository()
	maxOrder, _ := nodeRepo.GetMaxSortOrder(cluster.ID, group)

	node := &model.ClusterNode{
		ClusterID:   cluster.ID,
		GroupName:   group,
		IPAddress:   ip,
		K8sNodename: nodename,
		NewInstall:  false,
		SortOrder:   maxOrder + 1,
	}
	if err := nodeRepo.Create(node); err != nil {
		return err
	}

	// Regenerate hosts content
	return s.regenerateHosts(cluster.ID)
}

func (s *NodeOpsService) removeNodeFromConfig(clusterID int64, group string, ip string) error {
	nodeRepo := repository.NewNodeRepository()
	if err := nodeRepo.DeleteByClusterGroupIP(clusterID, group, ip); err != nil {
		return err
	}

	// Regenerate hosts content
	return s.regenerateHosts(clusterID)
}

func (s *NodeOpsService) regenerateHosts(clusterID int64) error {
	cluster, err := s.clusterRepo.GetByID(clusterID)
	if err != nil {
		return err
	}

	gen := NewConfigGenerator()
	hostsText, err := gen.GenerateHosts(cluster, cluster.Nodes, cluster.GlobalVars)
	if err != nil {
		return err
	}

	cluster.HostsContent = hostsText
	return s.clusterRepo.Update(cluster)
}

func (s *NodeOpsService) watchAndRemoveNode(taskID int64, clusterID int64, group string, ip string) {
	task, err := s.taskService.GetTask(taskID)
	if err != nil {
		log.Printf("[node-ops] failed to get task %d: %v", taskID, err)
		return
	}

	// Poll until task completes
	for task.Status == "running" {
		task, _ = s.taskService.GetTask(taskID)
	}

	if task.Status == "success" {
		if err := s.removeNodeFromConfig(clusterID, group, ip); err != nil {
			log.Printf("[node-ops] failed to remove node %s/%s from config: %v", group, ip, err)
		} else {
			log.Printf("[node-ops] removed node %s/%s from config after successful del task", group, ip)
		}
	} else {
		log.Printf("[node-ops] del task %d failed (status=%s), keeping node in config", taskID, task.Status)
	}
}
