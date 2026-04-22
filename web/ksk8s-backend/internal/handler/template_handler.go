package handler

import (
	"net/http"
	"strconv"

	"github.com/easzlab/ksk8s/internal/model"
	"github.com/easzlab/ksk8s/internal/repository"
	"github.com/easzlab/ksk8s/internal/service"
	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	repo      *repository.TemplateRepository
	configGen *service.ConfigGenerator
}

func NewTemplateHandler() *TemplateHandler {
	return &TemplateHandler{
		repo:      repository.NewTemplateRepository(),
		configGen: service.NewConfigGenerator(),
	}
}

type CreateTemplateRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	HostsContent  string `json:"hosts_content" binding:"required"`
	ConfigContent string `json:"config_content" binding:"required"`
}

func (h *TemplateHandler) Create(c *gin.Context) {
	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	tmpl := &model.Template{
		Name:          req.Name,
		Description:   req.Description,
		HostsContent:  req.HostsContent,
		ConfigContent: req.ConfigContent,
		CreatedBy:     userID.(int64),
	}

	if err := h.repo.Create(tmpl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create template"})
		return
	}
	c.JSON(http.StatusCreated, tmpl)
}

func (h *TemplateHandler) List(c *gin.Context) {
	templates, err := h.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list templates"})
		return
	}
	c.JSON(http.StatusOK, templates)
}

func (h *TemplateHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	tmpl, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}
	c.JSON(http.StatusOK, tmpl)
}

func (h *TemplateHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	tmpl, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}

	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tmpl.Name = req.Name
	tmpl.Description = req.Description
	tmpl.HostsContent = req.HostsContent
	tmpl.ConfigContent = req.ConfigContent

	if err := h.repo.Update(tmpl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update template"})
		return
	}
	c.JSON(http.StatusOK, tmpl)
}

func (h *TemplateHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete template"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *TemplateHandler) SetDefault(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.repo.ClearDefault(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear default"})
		return
	}
	if err := h.repo.Update(&model.Template{ID: id, IsDefault: true}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set default"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "default template updated"})
}

// GetParsed returns a template with its hosts and config content parsed into structured data.
func (h *TemplateHandler) GetParsed(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	tmpl, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}

	nodes, vars, _ := h.configGen.ParseHostsText(tmpl.HostsContent)
	params, paramLists, _ := h.configGen.ParseConfigYAML(tmpl.ConfigContent)

	// Group nodes by group_name
	nodeGroups := make(map[string][]gin.H)
	for _, n := range nodes {
		nodeGroups[n.GroupName] = append(nodeGroups[n.GroupName], gin.H{
			"ip_address":        n.IPAddress,
			"k8s_nodename":      n.K8sNodename,
			"new_install":       n.NewInstall,
			"lb_role":           n.LBRole,
			"ex_apiserver_vip":  n.ExApiserverVIP,
			"ex_apiserver_port": n.ExApiserverPort,
		})
	}

	paramMap := make(map[string]string)
	for _, p := range params {
		paramMap[p.ParamKey] = p.ParamValue
	}

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

	c.JSON(http.StatusOK, gin.H{
		"id":             tmpl.ID,
		"name":           tmpl.Name,
		"description":    tmpl.Description,
		"is_default":     tmpl.IsDefault,
		"nodes":          nodeGroups,
		"global_vars":    vars,
		"params":         paramMap,
		"param_lists":    listMap,
		"hosts_content":  tmpl.HostsContent,
		"config_content": tmpl.ConfigContent,
	})
}
