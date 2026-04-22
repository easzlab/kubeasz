package model

import "time"

type Cluster struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"size:64;not null;uniqueIndex" json:"name"`
	Description    string    `gorm:"size:512" json:"description"`
	Status         string    `gorm:"size:16;default:'draft'" json:"status"`
	TemplateID     *int64    `json:"template_id"`
	HostsContent   string    `gorm:"type:longtext" json:"hosts_content,omitempty"`
	ConfigContent  string    `gorm:"type:longtext" json:"config_content,omitempty"`
	CurrentVersion   int       `gorm:"default:1" json:"current_version"`
	InstallStepIndex int       `gorm:"default:0" json:"install_step_index"`
	CreatedBy        int64     `gorm:"not null" json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Nodes       []ClusterNode       `gorm:"foreignKey:ClusterID" json:"nodes,omitempty"`
	Params      []ClusterParam      `gorm:"foreignKey:ClusterID" json:"params,omitempty"`
	ParamLists  []ClusterParamList  `gorm:"foreignKey:ClusterID" json:"param_lists,omitempty"`
	GlobalVars  *HostsGlobalVars    `gorm:"foreignKey:ClusterID" json:"global_vars,omitempty"`
}

type ClusterNode struct {
	ID               int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ClusterID        int64     `gorm:"not null;index:idx_cluster_group" json:"cluster_id"`
	GroupName        string    `gorm:"size:16;not null;index:idx_cluster_group" json:"group_name"`
	IPAddress        string    `gorm:"size:64;not null" json:"ip_address"`
	K8sNodename      string    `gorm:"size:64" json:"k8s_nodename"`
	NewInstall       bool      `gorm:"default:false" json:"new_install"`
	LBRole           *string   `gorm:"size:8" json:"lb_role,omitempty"`
	ExApiserverVIP   *string   `gorm:"size:64" json:"ex_apiserver_vip,omitempty"`
	ExApiserverPort  *string   `gorm:"size:8" json:"ex_apiserver_port,omitempty"`
	SortOrder        int       `gorm:"default:0" json:"sort_order"`
	CreatedAt        time.Time `json:"created_at"`
}

type ClusterParam struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ClusterID  int64     `gorm:"not null;uniqueIndex:uk_cluster_param" json:"cluster_id"`
	ParamGroup string    `gorm:"size:32;default:'general'" json:"param_group"`
	ParamKey   string    `gorm:"size:64;not null;uniqueIndex:uk_cluster_param" json:"param_key"`
	ParamValue string    `gorm:"type:text" json:"param_value"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ClusterParamList struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ClusterID  int64     `gorm:"not null;index:idx_cluster_key" json:"cluster_id"`
	ParamKey   string    `gorm:"size:64;not null;index:idx_cluster_key" json:"param_key"`
	ItemValue  string    `gorm:"size:256;not null" json:"item_value"`
	SortOrder  int       `gorm:"default:0" json:"sort_order"`
	CreatedAt  time.Time `json:"created_at"`
}

type HostsGlobalVars struct {
	ID                       int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	ClusterID                int64   `gorm:"not null;uniqueIndex" json:"cluster_id"`
	SecurePort               string  `gorm:"size:8;default:'6443'" json:"secure_port"`
	ContainerRuntime         string  `gorm:"size:16;default:'containerd'" json:"container_runtime"`
	ClusterNetwork           string  `gorm:"size:16;default:'calico'" json:"cluster_network"`
	ProxyMode                string  `gorm:"size:8;default:'ipvs'" json:"proxy_mode"`
	ServiceCIDR              string  `gorm:"size:20;default:'10.68.0.0/16'" json:"service_cidr"`
	ClusterCIDR              string  `gorm:"size:20;default:'172.20.0.0/16'" json:"cluster_cidr"`
	NodePortRange            string  `gorm:"size:20;default:'30000-32767'" json:"node_port_range"`
	ClusterDNSDomain         string  `gorm:"size:64;default:'cluster.local'" json:"cluster_dns_domain"`
	BinDir                   string  `gorm:"size:128;default:'/opt/kube/bin'" json:"bin_dir"`
	BaseDir                  string  `gorm:"size:128;default:'/etc/kubeasz'" json:"base_dir"`
	CaDir                    string  `gorm:"size:128;default:'/etc/kubernetes/ssl'" json:"ca_dir"`
	ClusterDir               string  `gorm:"size:128;default:''" json:"cluster_dir"`
	K8sNodename              string  `gorm:"size:64;default:''" json:"k8s_nodename"`
	AnsiblePythonInterpreter string  `gorm:"size:128;default:'/usr/bin/python3'" json:"ansible_python_interpreter"`
	AnsibleUser              string  `gorm:"size:32;default:'root'" json:"ansible_user"`
	AnsibleBecome            string  `gorm:"size:8;default:'no'" json:"ansible_become"`
}

type ClusterVersion struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ClusterID      int64     `gorm:"not null;uniqueIndex:uk_cluster_version" json:"cluster_id"`
	VersionNumber  int       `gorm:"not null;uniqueIndex:uk_cluster_version" json:"version_number"`
	ChangeSummary  string    `gorm:"type:json" json:"change_summary,omitempty"`
	HostsContent   string    `gorm:"type:longtext" json:"hosts_content,omitempty"`
	ConfigContent  string    `gorm:"type:longtext" json:"config_content,omitempty"`
	CreatedBy      int64     `gorm:"not null" json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
}
