package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/easzlab/ksk8s/internal/model"
	"github.com/easzlab/ksk8s/internal/repository"
	"github.com/gin-gonic/gin"
)

func setupClusterHandlerTest(t *testing.T) (*gin.Engine, *ClusterHandler) {
	repository.SetupTestDB(t)
	gin.SetMode(gin.TestMode)
	handler := NewClusterHandler()
	router := gin.New()
	return router, handler
}

func TestClusterHandler_Create(t *testing.T) {
	_, h := setupClusterHandlerTest(t)

	body, _ := json.Marshal(CreateClusterRequest{
		Name:        "test-cluster",
		Description: "a test cluster",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/clusters", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", int64(1))

	h.Create(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.Cluster
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Name != "test-cluster" {
		t.Errorf("expected name test-cluster, got %s", resp.Name)
	}
	if resp.Status != "draft" {
		t.Errorf("expected status draft, got %s", resp.Status)
	}
	if resp.CreatedBy != 1 {
		t.Errorf("expected created_by 1, got %d", resp.CreatedBy)
	}
}

func TestClusterHandler_Create_InvalidJSON(t *testing.T) {
	_, h := setupClusterHandlerTest(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/clusters", bytes.NewReader([]byte("not json")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", int64(1))

	h.Create(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestClusterHandler_Create_WithTemplate(t *testing.T) {
	_, h := setupClusterHandlerTest(t)

	// Create a template first
	tmpl := &model.Template{
		Name:          "default",
		HostsContent:  "[kube_master]\n192.168.1.1",
		ConfigContent: "k8s_ver: v1.30",
	}
	if err := h.templateRepo.Create(tmpl); err != nil {
		t.Fatalf("failed to create template: %v", err)
	}

	body, _ := json.Marshal(CreateClusterRequest{
		Name:       "from-template",
		TemplateID: &tmpl.ID,
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/clusters", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", int64(1))

	h.Create(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.Cluster
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.HostsContent != tmpl.HostsContent {
		t.Errorf("expected hosts content from template, got %s", resp.HostsContent)
	}
	if resp.ConfigContent != tmpl.ConfigContent {
		t.Errorf("expected config content from template, got %s", resp.ConfigContent)
	}
}

func TestClusterHandler_List(t *testing.T) {
	_, h := setupClusterHandlerTest(t)

	// Seed data
	h.clusterRepo.Create(&model.Cluster{Name: "c1", CreatedBy: 1, Status: "draft"})
	h.clusterRepo.Create(&model.Cluster{Name: "c2", CreatedBy: 1, Status: "draft"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/clusters", nil)

	h.List(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp []model.Cluster
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp) != 2 {
		t.Errorf("expected 2 clusters, got %d", len(resp))
	}
}

func TestClusterHandler_Get(t *testing.T) {
	_, h := setupClusterHandlerTest(t)

	cluster := &model.Cluster{Name: "get-me", CreatedBy: 1, Status: "draft"}
	h.clusterRepo.Create(cluster)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/clusters/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Get(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.Cluster
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Name != "get-me" {
		t.Errorf("expected name get-me, got %s", resp.Name)
	}
}

func TestClusterHandler_Get_NotFound(t *testing.T) {
	_, h := setupClusterHandlerTest(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/clusters/999", nil)
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	h.Get(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestClusterHandler_Update(t *testing.T) {
	_, h := setupClusterHandlerTest(t)

	cluster := &model.Cluster{Name: "old-name", CreatedBy: 1, Status: "draft"}
	h.clusterRepo.Create(cluster)

	body, _ := json.Marshal(map[string]string{
		"name":   "new-name",
		"status": "active",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/api/clusters/1", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Update(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp model.Cluster
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Name != "new-name" {
		t.Errorf("expected name new-name, got %s", resp.Name)
	}
	if resp.Status != "active" {
		t.Errorf("expected status active, got %s", resp.Status)
	}
}

func TestClusterHandler_Update_NotFound(t *testing.T) {
	_, h := setupClusterHandlerTest(t)

	body, _ := json.Marshal(map[string]string{"name": "x"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/api/clusters/999", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	h.Update(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestClusterHandler_Delete(t *testing.T) {
	_, h := setupClusterHandlerTest(t)

	cluster := &model.Cluster{Name: "delete-me", CreatedBy: 1, Status: "draft"}
	h.clusterRepo.Create(cluster)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/api/clusters/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Delete(c)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}

	_, err := h.clusterRepo.GetByID(cluster.ID)
	if err == nil {
		t.Error("cluster should be deleted")
	}
}

func TestClusterHandler_GenerateConfig(t *testing.T) {
	_, h := setupClusterHandlerTest(t)

	cluster := &model.Cluster{Name: "gen-test", CreatedBy: 1, Status: "draft"}
	h.clusterRepo.Create(cluster)

	// Add a node
	h.nodeRepo.Create(&model.ClusterNode{
		ClusterID: cluster.ID,
		GroupName: "kube_master",
		IPAddress: "192.168.1.1",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/clusters/1/config", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GenerateConfig(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["hosts"] == nil || resp["hosts"] == "" {
		t.Error("expected hosts to be generated")
	}
	if resp["config"] == nil || resp["config"] == "" {
		t.Error("expected config to be generated")
	}
}
