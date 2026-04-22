package service

import (
	"strings"
	"testing"

	"github.com/easzlab/ksk8s/internal/model"
)

func TestGenerateHosts(t *testing.T) {
	g := NewConfigGenerator()
	nodes := []model.ClusterNode{
		{GroupName: "etcd", IPAddress: "192.168.1.1"},
		{GroupName: "etcd", IPAddress: "192.168.1.2"},
		{GroupName: "kube_master", IPAddress: "192.168.1.1", K8sNodename: "master-01"},
		{GroupName: "kube_node", IPAddress: "192.168.1.3", K8sNodename: "worker-01"},
	}
	vars := &model.HostsGlobalVars{
		SecurePort:               "6443",
		ContainerRuntime:         "containerd",
		ClusterNetwork:           "cilium",
		ProxyMode:                "iptables",
		ServiceCIDR:              "10.19.0.0/16",
		ClusterCIDR:              "172.19.0.0/16",
		NodePortRange:            "30000-32767",
		ClusterDNSDomain:         "cluster.local",
		BinDir:                   "/opt/kube/bin",
		BaseDir:                  "/etc/kubeasz",
		CaDir:                    "/etc/kubernetes/ssl",
		AnsiblePythonInterpreter: "/usr/bin/python3",
		AnsibleUser:              "root",
		AnsibleBecome:            "yes",
	}

	hosts, err := g.GenerateHosts(&model.Cluster{Name: "test"}, nodes, vars)
	if err != nil {
		t.Fatalf("GenerateHosts failed: %v", err)
	}

	if !strings.Contains(hosts, "[etcd]") {
		t.Error("missing [etcd] section")
	}
	if !strings.Contains(hosts, "[kube_master]") {
		t.Error("missing [kube_master] section")
	}
	if !strings.Contains(hosts, "192.168.1.1 k8s_nodename=master-01") {
		t.Error("missing master node with nodename")
	}
	if !strings.Contains(hosts, "CONTAINER_RUNTIME=\"containerd\"") {
		t.Error("missing container runtime")
	}
	if !strings.Contains(hosts, "CLUSTER_NETWORK=\"cilium\"") {
		t.Error("missing cluster network")
	}
}

func TestGenerateConfigYAML(t *testing.T) {
	g := NewConfigGenerator()
	params := []model.ClusterParam{
		{ParamKey: "K8S_VER", ParamValue: "1.34.3"},
		{ParamKey: "CLUSTER_NAME", ParamValue: "test"},
	}
	paramLists := []model.ClusterParamList{
		{ParamKey: "MASTER_CERT_HOSTS", ItemValue: "192.168.1.1"},
		{ParamKey: "MASTER_CERT_HOSTS", ItemValue: "192.168.1.2"},
	}

	yaml, err := g.GenerateConfigYAML(&model.Cluster{Name: "test"}, params, paramLists)
	if err != nil {
		t.Fatalf("GenerateConfigYAML failed: %v", err)
	}

	if !strings.Contains(yaml, "K8S_VER: 1.34.3") {
		t.Error("missing K8S_VER")
	}
	if !strings.Contains(yaml, "MASTER_CERT_HOSTS:") {
		t.Error("missing MASTER_CERT_HOSTS list")
	}
	if !strings.Contains(yaml, "- 192.168.1.1") {
		t.Error("missing list item 192.168.1.1")
	}
}

func TestParseConfigYAML(t *testing.T) {
	g := NewConfigGenerator()
	yaml := `
K8S_VER: "1.34.3"
MASTER_CERT_HOSTS:
  - "192.168.1.1"
  - "192.168.1.2"
`
	params, lists, err := g.ParseConfigYAML(yaml)
	if err != nil {
		t.Fatalf("ParseConfigYAML failed: %v", err)
	}

	found := false
	for _, p := range params {
		if p.ParamKey == "K8S_VER" && p.ParamValue == "1.34.3" {
			found = true
			break
		}
	}
	if !found {
		t.Error("missing parsed K8S_VER")
	}

	if len(lists) != 2 {
		t.Errorf("expected 2 list items, got %d", len(lists))
	}
}

func TestParseHostsText(t *testing.T) {
	g := NewConfigGenerator()
	content := `[etcd]
192.168.1.1
[kube_master]
192.168.1.2 k8s_nodename=master-01 new_install=yes
[all:vars]
SECURE_PORT="6443"
CONTAINER_RUNTIME="containerd"
CLUSTER_NETWORK="calico"
`
	nodes, vars, err := g.ParseHostsText(content)
	if err != nil {
		t.Fatalf("ParseHostsText failed: %v", err)
	}
	if len(nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(nodes))
	}
	if vars.SecurePort != "6443" {
		t.Errorf("expected secure port 6443, got %s", vars.SecurePort)
	}
	if vars.ContainerRuntime != "containerd" {
		t.Errorf("expected containerd, got %s", vars.ContainerRuntime)
	}

	foundMaster := false
	for _, n := range nodes {
		if n.IPAddress == "192.168.1.2" {
			foundMaster = true
			if n.K8sNodename != "master-01" {
				t.Errorf("expected nodename master-01, got %s", n.K8sNodename)
			}
			if !n.NewInstall {
				t.Error("expected new_install=yes")
			}
		}
	}
	if !foundMaster {
		t.Error("expected to find kube_master node")
	}
}

func TestHostsRoundTrip(t *testing.T) {
	g := NewConfigGenerator()
	originalNodes := []model.ClusterNode{
		{GroupName: "etcd", IPAddress: "192.168.1.1"},
		{GroupName: "kube_master", IPAddress: "192.168.1.2", K8sNodename: "master-01"},
		{GroupName: "kube_node", IPAddress: "192.168.1.3", K8sNodename: "worker-01"},
	}
	originalVars := &model.HostsGlobalVars{
		SecurePort:               "6443",
		ContainerRuntime:         "containerd",
		ClusterNetwork:           "calico",
		ProxyMode:                "ipvs",
		ServiceCIDR:              "10.68.0.0/16",
		ClusterCIDR:              "172.20.0.0/16",
		NodePortRange:            "30000-32767",
		ClusterDNSDomain:         "cluster.local",
		BinDir:                   "/opt/kube/bin",
		BaseDir:                  "/etc/kubeasz",
		CaDir:                    "/etc/kubernetes/ssl",
		AnsiblePythonInterpreter: "/usr/bin/python3",
		AnsibleUser:              "root",
		AnsibleBecome:            "no",
	}

	hosts, err := g.GenerateHosts(&model.Cluster{Name: "test"}, originalNodes, originalVars)
	if err != nil {
		t.Fatalf("GenerateHosts failed: %v", err)
	}

	parsedNodes, parsedVars, err := g.ParseHostsText(hosts)
	if err != nil {
		t.Fatalf("ParseHostsText failed: %v", err)
	}

	if len(parsedNodes) != len(originalNodes) {
		t.Errorf("expected %d nodes, got %d", len(originalNodes), len(parsedNodes))
	}

	if parsedVars.SecurePort != originalVars.SecurePort {
		t.Errorf("round-trip secure port mismatch: %s vs %s", parsedVars.SecurePort, originalVars.SecurePort)
	}
	if parsedVars.ClusterNetwork != originalVars.ClusterNetwork {
		t.Errorf("round-trip cluster network mismatch: %s vs %s", parsedVars.ClusterNetwork, originalVars.ClusterNetwork)
	}
}

func TestConfigYAMLRoundTrip(t *testing.T) {
	g := NewConfigGenerator()
	originalParams := []model.ClusterParam{
		{ParamKey: "K8S_VER", ParamValue: "1.34.3"},
		{ParamKey: "ETCD_VERSION", ParamValue: "v3.5.15"},
	}
	originalLists := []model.ClusterParamList{
		{ParamKey: "MASTER_CERT_HOSTS", ItemValue: "192.168.1.1"},
		{ParamKey: "MASTER_CERT_HOSTS", ItemValue: "192.168.1.2"},
	}

	yaml, err := g.GenerateConfigYAML(&model.Cluster{Name: "test"}, originalParams, originalLists)
	if err != nil {
		t.Fatalf("GenerateConfigYAML failed: %v", err)
	}

	parsedParams, parsedLists, err := g.ParseConfigYAML(yaml)
	if err != nil {
		t.Fatalf("ParseConfigYAML failed: %v", err)
	}

	foundK8s := false
	for _, p := range parsedParams {
		if p.ParamKey == "K8S_VER" && p.ParamValue == "1.34.3" {
			foundK8s = true
			break
		}
	}
	if !foundK8s {
		t.Error("round-trip lost K8S_VER")
	}

	if len(parsedLists) != 2 {
		t.Errorf("expected 2 list items after round-trip, got %d", len(parsedLists))
	}
}

func TestUpdateConfigYAML_PreservesTemplateKeys(t *testing.T) {
	g := NewConfigGenerator()
	original := `
INSTALL_SOURCE: "online"
OS_HARDEN: false
K8S_VER: "1.30.0"
INSECURE_REG:
  - "http://old.local"
MASTER_CERT_HOSTS:
  - "10.0.0.1"
ETCD_DATA_DIR: "/var/lib/etcd"
`
	params := map[string]string{
		"K8S_VER": "1.34.3",
	}
	lists := map[string][]string{
		"INSECURE_REG": {"http://new.local"},
	}

	result, err := g.UpdateConfigYAML(original, params, lists, "test-cluster")
	if err != nil {
		t.Fatalf("UpdateConfigYAML failed: %v", err)
	}

	// Updated values
	if !strings.Contains(result, "K8S_VER: \"1.34.3\"") {
		t.Error("K8S_VER not updated")
	}
	if !strings.Contains(result, "- http://new.local") {
		t.Error("INSECURE_REG not updated")
	}
	if !strings.Contains(result, "CLUSTER_NAME: test-cluster") {
		t.Error("CLUSTER_NAME not set")
	}

	// Preserved untouched keys
	if !strings.Contains(result, "INSTALL_SOURCE: \"online\"") {
		t.Error("INSTALL_SOURCE not preserved")
	}
	if !strings.Contains(result, "OS_HARDEN: false") {
		t.Error("OS_HARDEN not preserved")
	}
	if !strings.Contains(result, "ETCD_DATA_DIR: \"/var/lib/etcd\"") {
		t.Error("ETCD_DATA_DIR not preserved")
	}
	if !strings.Contains(result, "MASTER_CERT_HOSTS:") {
		t.Error("MASTER_CERT_HOSTS not preserved")
	}
}

func TestUpdateConfigYAML_EmptyOriginal(t *testing.T) {
	g := NewConfigGenerator()
	result, err := g.UpdateConfigYAML("", map[string]string{"K8S_VER": "1.34.3"}, nil, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty result for empty original, got: %s", result)
	}
}

func TestParseHostsText_Empty(t *testing.T) {
	g := NewConfigGenerator()
	nodes, vars, err := g.ParseHostsText("")
	if err != nil {
		t.Fatalf("ParseHostsText(empty) failed: %v", err)
	}
	if len(nodes) != 0 {
		t.Errorf("expected 0 nodes for empty input, got %d", len(nodes))
	}
	if vars == nil {
		t.Fatal("expected non-nil vars for empty input")
	}
}
