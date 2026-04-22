package handler

import (
	"net/http"
	"strconv"

	"github.com/easzlab/ksk8s/internal/model"
	"github.com/easzlab/ksk8s/internal/repository"
	"github.com/easzlab/ksk8s/internal/service"
	"github.com/gin-gonic/gin"
)

type ClusterHandler struct {
	clusterRepo  *repository.ClusterRepository
	nodeRepo     *repository.NodeRepository
	paramRepo    *repository.ParamRepository
	hostsVarRepo *repository.HostsVarRepository
	templateRepo *repository.TemplateRepository
	bindingRepo  *repository.BindingRepository
	configGen    *service.ConfigGenerator
}

func NewClusterHandler() *ClusterHandler {
	return &ClusterHandler{
		clusterRepo:  repository.NewClusterRepository(),
		nodeRepo:     repository.NewNodeRepository(),
		paramRepo:    repository.NewParamRepository(),
		hostsVarRepo: repository.NewHostsVarRepository(),
		templateRepo: repository.NewTemplateRepository(),
		bindingRepo:  repository.NewBindingRepository(),
		configGen:    service.NewConfigGenerator(),
	}
}

type CreateClusterRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	TemplateID  *int64 `json:"template_id"`
}

func (h *ClusterHandler) Create(c *gin.Context) {
	var req CreateClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	cluster := &model.Cluster{
		Name:        req.Name,
		Description: req.Description,
		TemplateID:  req.TemplateID,
		Status:      "draft",
		CreatedBy:   userID.(int64),
	}

	// If template specified, copy its content
	if req.TemplateID != nil {
		tmpl, err := h.templateRepo.GetByID(*req.TemplateID)
		if err == nil {
			cluster.HostsContent = tmpl.HostsContent
			cluster.ConfigContent = tmpl.ConfigContent
		}
	}

	if err := h.clusterRepo.Create(cluster); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create cluster"})
		return
	}

	// Create default global vars
	_, _ = h.hostsVarRepo.GetOrCreate(cluster.ID)

	c.JSON(http.StatusCreated, cluster)
}

func (h *ClusterHandler) List(c *gin.Context) {
	role, _ := c.Get("role")
	userID, _ := c.Get("user_id")
	uid := userID.(int64)

	var clusters []model.Cluster
	var err error

	switch role {
	case model.RolePlatformAdmin, model.RoleSecurityAuditor:
		clusters, err = h.clusterRepo.List()
	case model.RoleClusterAdmin:
		// Own clusters + bound clusters
		ownClusters, _ := h.clusterRepo.ListByUser(uid)
		bindings, _ := h.bindingRepo.ListByUser(uid)
		boundIDs := make([]int64, 0, len(bindings))
		for _, b := range bindings {
			boundIDs = append(boundIDs, b.ClusterID)
		}
		boundClusters, _ := h.clusterRepo.ListByIDs(boundIDs)
		// Merge
		seen := make(map[int64]bool)
		clusters = make([]model.Cluster, 0, len(ownClusters)+len(boundClusters))
		for _, cl := range ownClusters {
			if !seen[cl.ID] {
				seen[cl.ID] = true
				clusters = append(clusters, cl)
			}
		}
		for _, cl := range boundClusters {
			if !seen[cl.ID] {
				seen[cl.ID] = true
				clusters = append(clusters, cl)
			}
		}
	default:
		clusters, err = h.clusterRepo.ListByUser(uid)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list clusters"})
		return
	}
	c.JSON(http.StatusOK, clusters)
}

func (h *ClusterHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cluster, err := h.clusterRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster not found"})
		return
	}
	c.JSON(http.StatusOK, cluster)
}

func (h *ClusterHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cluster, err := h.clusterRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster not found"})
		return
	}

	var req struct {
		Name             string `json:"name"`
		Description      string `json:"description"`
		Status           string `json:"status"`
		InstallStepIndex *int   `json:"install_step_index"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		cluster.Name = req.Name
	}
	if req.Description != "" {
		cluster.Description = req.Description
	}
	if req.Status != "" {
		cluster.Status = req.Status
	}
	if req.InstallStepIndex != nil {
		cluster.InstallStepIndex = *req.InstallStepIndex
	}

	if err := h.clusterRepo.Update(cluster); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update cluster"})
		return
	}
	c.JSON(http.StatusOK, cluster)
}

func (h *ClusterHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.clusterRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete cluster"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *ClusterHandler) GetConfig(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cluster, err := h.clusterRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster not found"})
		return
	}

	nodes, _ := h.nodeRepo.GetByClusterID(id)
	vars, _ := h.hostsVarRepo.GetOrCreate(id)
	params, _ := h.paramRepo.GetByClusterID(id)
	paramLists, _ := h.paramRepo.GetListByClusterID(id)

	// If structured data is empty but text content exists, parse from text
	if len(nodes) == 0 && cluster.HostsContent != "" {
		parsedNodes, parsedVars, _ := h.configGen.ParseHostsText(cluster.HostsContent)
		nodes = parsedNodes
		if parsedVars != nil {
			vars = parsedVars
		}
	}

	// Group nodes by group_name
	nodeGroups := make(map[string][]gin.H)
	for _, n := range nodes {
		nodeGroups[n.GroupName] = append(nodeGroups[n.GroupName], gin.H{
			"id":                n.ID,
			"ip_address":        n.IPAddress,
			"k8s_nodename":      n.K8sNodename,
			"new_install":       n.NewInstall,
			"lb_role":           n.LBRole,
			"ex_apiserver_vip":  n.ExApiserverVIP,
			"ex_apiserver_port": n.ExApiserverPort,
		})
	}

	// Build params map: DB values take precedence; fill missing keys from template
	paramMap := make(map[string]string)
	for _, p := range params {
		paramMap[p.ParamKey] = p.ParamValue
	}
	if cluster.ConfigContent != "" {
		parsedParams, _, _ := h.configGen.ParseConfigYAML(cluster.ConfigContent)
		for _, p := range parsedParams {
			if _, exists := paramMap[p.ParamKey]; !exists {
				paramMap[p.ParamKey] = p.ParamValue
			}
		}
	}

	// Build param_lists map: DB values take precedence; fill missing keys from template
	listMap := make(map[string][]string)
	for _, pl := range paramLists {
		if pl.ItemValue == "" {
			if _, ok := listMap[pl.ParamKey]; !ok {
				listMap[pl.ParamKey] = []string{}
			}
			continue
		}
		listMap[pl.ParamKey] = append(listMap[pl.ParamKey], pl.ItemValue)
	}
	if cluster.ConfigContent != "" {
		_, parsedLists, _ := h.configGen.ParseConfigYAML(cluster.ConfigContent)
		for _, pl := range parsedLists {
			if _, exists := listMap[pl.ParamKey]; exists {
				continue
			}
			if pl.ItemValue == "" {
				listMap[pl.ParamKey] = []string{}
				continue
			}
			listMap[pl.ParamKey] = append(listMap[pl.ParamKey], pl.ItemValue)
		}
	}

	// Parse template metadata for display order and hints
	var paramMeta []service.ParamMeta
	var sectionMeta []service.HostsSectionMeta
	var varMeta []service.HostsVarMeta

	var hostsBase, configBase string
	if cluster.TemplateID != nil {
		tmpl, err := h.templateRepo.GetByID(*cluster.TemplateID)
		if err == nil {
			hostsBase = tmpl.HostsContent
			configBase = tmpl.ConfigContent
		}
	}
	if hostsBase == "" {
		hostsBase = cluster.HostsContent
	}
	if configBase == "" {
		configBase = cluster.ConfigContent
	}

	if configBase != "" {
		paramMeta, _ = h.configGen.ParseConfigMeta(configBase)
	}
	if hostsBase != "" {
		sectionMeta, varMeta, _ = h.configGen.ParseHostsMeta(hostsBase)
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes":          nodeGroups,
		"global_vars":    vars,
		"params":         paramMap,
		"param_lists":    listMap,
		"param_meta":     paramMeta,
		"global_var_meta": varMeta,
		"section_meta":   sectionMeta,
		"hosts_content":  cluster.HostsContent,
		"config_content": cluster.ConfigContent,
	})
}

type SaveConfigRequest struct {
	Nodes       map[string][]NodeRequest `json:"nodes"`
	GlobalVars  *model.HostsGlobalVars   `json:"global_vars"`
	Params      map[string]string        `json:"params"`
	ParamLists  map[string][]string      `json:"param_lists"`
}

type NodeRequest struct {
	IPAddress        string  `json:"ip_address"`
	K8sNodename      string  `json:"k8s_nodename"`
	NewInstall       bool    `json:"new_install"`
	LBRole           *string `json:"lb_role"`
	ExApiserverVIP   *string `json:"ex_apiserver_vip"`
	ExApiserverPort  *string `json:"ex_apiserver_port"`
}

func (h *ClusterHandler) SaveConfig(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cluster, err := h.clusterRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster not found"})
		return
	}

	var req SaveConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation: required hosts groups
	if len(req.Nodes["etcd"]) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "etcd nodes are required"})
		return
	}
	if len(req.Nodes["kube_master"]) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kube_master nodes are required"})
		return
	}

	// Validation: required config params
	if req.Params == nil {
		req.Params = make(map[string]string)
	}
	if req.ParamLists == nil {
		req.ParamLists = make(map[string][]string)
	}

	if len(req.ParamLists["INSECURE_REG"]) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INSECURE_REG is required"})
		return
	}
	if len(req.ParamLists["MASTER_CERT_HOSTS"]) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MASTER_CERT_HOSTS is required"})
		return
	}
	if req.Params["nfs_server"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nfs_server is required"})
		return
	}
	if req.Params["nfs_path"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nfs_path is required"})
		return
	}

	// Clear existing records
	_ = h.nodeRepo.DeleteByClusterID(id)
	_ = h.paramRepo.DeleteByClusterID(id)
	_ = h.paramRepo.DeleteListByClusterID(id)

	// Save nodes
	groupOrder := []string{"etcd", "kube_master", "kube_node", "harbor", "ex_lb", "chrony"}
	for _, group := range groupOrder {
		members, ok := req.Nodes[group]
		if !ok {
			continue
		}
		for i, n := range members {
			if n.IPAddress == "" {
				continue
			}
			node := &model.ClusterNode{
				ClusterID:       id,
				GroupName:       group,
				IPAddress:       n.IPAddress,
				K8sNodename:     n.K8sNodename,
				NewInstall:      n.NewInstall,
				LBRole:          n.LBRole,
				ExApiserverVIP:  n.ExApiserverVIP,
				ExApiserverPort: n.ExApiserverPort,
				SortOrder:       i,
			}
			if err := h.nodeRepo.Create(node); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save node"})
				return
			}
		}
	}

	// Save global vars
	if req.GlobalVars != nil {
		existingVars, _ := h.hostsVarRepo.GetOrCreate(id)
		existingVars.SecurePort = req.GlobalVars.SecurePort
		existingVars.ContainerRuntime = req.GlobalVars.ContainerRuntime
		existingVars.ClusterNetwork = req.GlobalVars.ClusterNetwork
		existingVars.ProxyMode = req.GlobalVars.ProxyMode
		existingVars.ServiceCIDR = req.GlobalVars.ServiceCIDR
		existingVars.ClusterCIDR = req.GlobalVars.ClusterCIDR
		existingVars.NodePortRange = req.GlobalVars.NodePortRange
		existingVars.ClusterDNSDomain = req.GlobalVars.ClusterDNSDomain
		existingVars.BinDir = req.GlobalVars.BinDir
		existingVars.BaseDir = req.GlobalVars.BaseDir
		existingVars.ClusterDir = req.GlobalVars.ClusterDir
		existingVars.CaDir = req.GlobalVars.CaDir
		existingVars.K8sNodename = req.GlobalVars.K8sNodename
		existingVars.AnsiblePythonInterpreter = req.GlobalVars.AnsiblePythonInterpreter
		existingVars.AnsibleUser = req.GlobalVars.AnsibleUser
		existingVars.AnsibleBecome = req.GlobalVars.AnsibleBecome
		if err := h.hostsVarRepo.Update(existingVars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save global vars"})
			return
		}
	}

	// Save scalar params
	for key, val := range req.Params {
		param := &model.ClusterParam{
			ClusterID:  id,
			ParamKey:   key,
			ParamValue: val,
		}
		if err := h.paramRepo.Create(param); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save param"})
			return
		}
	}

	// Save list params
	for key, items := range req.ParamLists {
		if len(items) == 0 {
			// Store empty list marker so the key is not lost
			item := &model.ClusterParamList{
				ClusterID: id,
				ParamKey:  key,
				ItemValue: "",
				SortOrder: 0,
			}
			if err := h.paramRepo.CreateListItem(item); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save param list"})
				return
			}
			continue
		}
		for i, val := range items {
			if val == "" {
				continue
			}
			item := &model.ClusterParamList{
				ClusterID: id,
				ParamKey:  key,
				ItemValue: val,
				SortOrder: i,
			}
			if err := h.paramRepo.CreateListItem(item); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save param list"})
				return
			}
		}
	}

	// Reload to regenerate text content
	nodes, _ := h.nodeRepo.GetByClusterID(id)
	vars, _ := h.hostsVarRepo.GetOrCreate(id)
	params, _ := h.paramRepo.GetByClusterID(id)
	paramLists, _ := h.paramRepo.GetListByClusterID(id)

	// Use template hosts as base to preserve comments and section order
	var hostsBase string
	if cluster.TemplateID != nil {
		tmpl, err := h.templateRepo.GetByID(*cluster.TemplateID)
		if err == nil && tmpl.HostsContent != "" {
			hostsBase = tmpl.HostsContent
		}
	}
	if hostsBase == "" {
		hostsBase = cluster.HostsContent
	}

	var sectionMeta []service.HostsSectionMeta
	var varMeta []service.HostsVarMeta
	if hostsBase != "" {
		sectionMeta, varMeta, _ = h.configGen.ParseHostsMeta(hostsBase)
	}

	hostsContent, err := h.configGen.GenerateHostsMeta(cluster, nodes, vars, sectionMeta, varMeta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate hosts"})
		return
	}

	// Build maps for config merge
	paramMap := make(map[string]string)
	for _, p := range params {
		paramMap[p.ParamKey] = p.ParamValue
	}
	listMap := make(map[string][]string)
	for _, pl := range paramLists {
		listMap[pl.ParamKey] = append(listMap[pl.ParamKey], pl.ItemValue)
	}

	// Use template config as merge base to preserve original format, order and all parameters
	var configBase string
	if cluster.TemplateID != nil {
		tmpl, err := h.templateRepo.GetByID(*cluster.TemplateID)
		if err == nil && tmpl.ConfigContent != "" {
			configBase = tmpl.ConfigContent
		}
	}
	if configBase == "" {
		configBase = cluster.ConfigContent
	}

	var configContent string
	if configBase != "" {
		configContent, err = h.configGen.UpdateConfigYAML(configBase, paramMap, listMap, cluster.Name)
	} else {
		configContent, err = h.configGen.GenerateConfigYAML(cluster, params, paramLists)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate config"})
		return
	}

	cluster.HostsContent = hostsContent
	cluster.ConfigContent = configContent
	if err := h.clusterRepo.Update(cluster); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update cluster"})
		return
	}

	// Write generated files to disk so ezctl can use them
	if err := service.WriteClusterConfigFiles(cluster); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "config saved but failed to write files: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "config saved"})
}

func (h *ClusterHandler) GenerateConfig(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cluster, err := h.clusterRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cluster not found"})
		return
	}

	nodes, _ := h.nodeRepo.GetByClusterID(id)
	vars, _ := h.hostsVarRepo.GetOrCreate(id)
	params, _ := h.paramRepo.GetByClusterID(id)
	paramLists, _ := h.paramRepo.GetListByClusterID(id)

	var hostsBase string
	if cluster.TemplateID != nil {
		tmpl, err := h.templateRepo.GetByID(*cluster.TemplateID)
		if err == nil && tmpl.HostsContent != "" {
			hostsBase = tmpl.HostsContent
		}
	}
	if hostsBase == "" {
		hostsBase = cluster.HostsContent
	}

	var sectionMeta []service.HostsSectionMeta
	var varMeta []service.HostsVarMeta
	if hostsBase != "" {
		sectionMeta, varMeta, _ = h.configGen.ParseHostsMeta(hostsBase)
	}

	hosts, err := h.configGen.GenerateHostsMeta(cluster, nodes, vars, sectionMeta, varMeta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate hosts"})
		return
	}

	paramMap := make(map[string]string)
	for _, p := range params {
		paramMap[p.ParamKey] = p.ParamValue
	}
	listMap := make(map[string][]string)
	for _, pl := range paramLists {
		listMap[pl.ParamKey] = append(listMap[pl.ParamKey], pl.ItemValue)
	}

	var configYAML string
	var configBase string
	if cluster.TemplateID != nil {
		tmpl, err := h.templateRepo.GetByID(*cluster.TemplateID)
		if err == nil && tmpl.ConfigContent != "" {
			configBase = tmpl.ConfigContent
		}
	}
	if configBase == "" {
		configBase = cluster.ConfigContent
	}
	if configBase != "" {
		configYAML, err = h.configGen.UpdateConfigYAML(configBase, paramMap, listMap, cluster.Name)
	} else {
		configYAML, err = h.configGen.GenerateConfigYAML(cluster, params, paramLists)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hosts":  hosts,
		"config": configYAML,
	})
}
