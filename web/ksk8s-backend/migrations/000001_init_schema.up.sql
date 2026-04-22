-- ksk8s initial schema
-- Generated 2026-04-18

-- Users
CREATE TABLE users (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    username        VARCHAR(64) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    email           VARCHAR(128),
    role            ENUM('admin','viewer') DEFAULT 'viewer',
    created_at      DATETIME DEFAULT NOW(),
    updated_at      DATETIME DEFAULT NOW()
);

-- Clusters
CREATE TABLE clusters (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    name            VARCHAR(64) NOT NULL UNIQUE,
    description     VARCHAR(512),
    status          ENUM('draft','configuring','installing','installed','error','destroyed') DEFAULT 'draft',
    template_id     BIGINT,
    hosts_content   LONGTEXT,
    config_content  LONGTEXT,
    current_version INT DEFAULT 1,
    created_by      BIGINT NOT NULL,
    created_at      DATETIME DEFAULT NOW(),
    updated_at      DATETIME DEFAULT NOW()
);

-- Cluster nodes (structured hosts)
CREATE TABLE cluster_nodes (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    cluster_id      BIGINT NOT NULL,
    group_name      ENUM('etcd','kube_master','kube_node','harbor','ex_lb','chrony') NOT NULL,
    ip_address      VARCHAR(64) NOT NULL,
    k8s_nodename    VARCHAR(64),
    new_install     BOOLEAN DEFAULT TRUE,
    lb_role         ENUM('master','backup'),
    ex_apiserver_vip VARCHAR(64),
    ex_apiserver_port VARCHAR(8),
    sort_order      INT DEFAULT 0,
    created_at      DATETIME DEFAULT NOW(),
    FOREIGN KEY (cluster_id) REFERENCES clusters(id) ON DELETE CASCADE,
    INDEX idx_cluster_group (cluster_id, group_name)
);

-- Cluster scalar params (config.yml)
CREATE TABLE cluster_params (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    cluster_id      BIGINT NOT NULL,
    param_group     VARCHAR(32) DEFAULT 'general',
    param_key       VARCHAR(64) NOT NULL,
    param_value     TEXT,
    created_at      DATETIME DEFAULT NOW(),
    updated_at      DATETIME DEFAULT NOW(),
    FOREIGN KEY (cluster_id) REFERENCES clusters(id) ON DELETE CASCADE,
    UNIQUE KEY uk_cluster_param (cluster_id, param_key),
    INDEX idx_cluster_group_key (cluster_id, param_group, param_key)
);

-- Cluster list params (e.g. MASTER_CERT_HOSTS, INSECURE_REG)
CREATE TABLE cluster_param_lists (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    cluster_id      BIGINT NOT NULL,
    param_key       VARCHAR(64) NOT NULL,
    item_value      VARCHAR(256) NOT NULL,
    sort_order      INT DEFAULT 0,
    created_at      DATETIME DEFAULT NOW(),
    FOREIGN KEY (cluster_id) REFERENCES clusters(id) ON DELETE CASCADE,
    INDEX idx_cluster_key (cluster_id, param_key)
);

-- Hosts global vars [all:vars]
CREATE TABLE hosts_global_vars (
    id                          BIGINT PRIMARY KEY AUTO_INCREMENT,
    cluster_id                  BIGINT NOT NULL UNIQUE,
    secure_port                 VARCHAR(8) DEFAULT '6443',
    container_runtime           ENUM('containerd','docker') DEFAULT 'containerd',
    cluster_network             ENUM('calico','flannel','kube-router','cilium','kube-ovn') DEFAULT 'calico',
    proxy_mode                  ENUM('ipvs','iptables') DEFAULT 'ipvs',
    service_cidr                VARCHAR(20) DEFAULT '10.68.0.0/16',
    cluster_cidr                VARCHAR(20) DEFAULT '172.20.0.0/16',
    node_port_range             VARCHAR(20) DEFAULT '30000-32767',
    cluster_dns_domain          VARCHAR(64) DEFAULT 'cluster.local',
    bin_dir                     VARCHAR(128) DEFAULT '/opt/kube/bin',
    base_dir                    VARCHAR(128) DEFAULT '/etc/kubeasz',
    ca_dir                      VARCHAR(128) DEFAULT '/etc/kubernetes/ssl',
    k8s_nodename                VARCHAR(64) DEFAULT '',
    ansible_python_interpreter  VARCHAR(128) DEFAULT '/usr/bin/python3',
    ansible_user                VARCHAR(32) DEFAULT 'root',
    ansible_become              VARCHAR(8) DEFAULT 'no',
    FOREIGN KEY (cluster_id) REFERENCES clusters(id) ON DELETE CASCADE
);

-- Cluster version history
CREATE TABLE cluster_versions (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    cluster_id      BIGINT NOT NULL,
    version_number  INT NOT NULL,
    change_summary  JSON,
    hosts_content   LONGTEXT,
    config_content  LONGTEXT,
    created_by      BIGINT NOT NULL,
    created_at      DATETIME DEFAULT NOW(),
    FOREIGN KEY (cluster_id) REFERENCES clusters(id) ON DELETE CASCADE,
    UNIQUE KEY uk_cluster_version (cluster_id, version_number)
);

-- Templates
CREATE TABLE templates (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    name            VARCHAR(64) NOT NULL UNIQUE,
    description     VARCHAR(512),
    is_default      BOOLEAN DEFAULT FALSE,
    hosts_content   LONGTEXT NOT NULL,
    config_content  LONGTEXT NOT NULL,
    created_by      BIGINT NOT NULL,
    created_at      DATETIME DEFAULT NOW(),
    updated_at      DATETIME DEFAULT NOW()
);

-- Tasks
CREATE TABLE tasks (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    cluster_id      BIGINT NOT NULL,
    task_type       ENUM('setup','start','stop','upgrade','backup','restore','destroy','add-etcd','add-master','add-node','del-etcd','del-master','del-node','ca-renew','kcfg-adm') NOT NULL,
    step_number     VARCHAR(8),
    target_node_ip  VARCHAR(64),
    status          ENUM('pending','queued','running','awaiting_approval','success','failed','aborted','rolling_back') DEFAULT 'pending',
    worker_pid      INT,
    log_path        VARCHAR(512),
    started_at      DATETIME,
    completed_at    DATETIME,
    approved_by     BIGINT,
    approved_at     DATETIME,
    exit_code       INT,
    error_message   TEXT,
    created_at      DATETIME DEFAULT NOW(),
    FOREIGN KEY (cluster_id) REFERENCES clusters(id) ON DELETE CASCADE,
    INDEX idx_cluster_status (cluster_id, status),
    INDEX idx_worker_pid (worker_pid)
);

-- Logs (per-task line-by-line)
CREATE TABLE logs (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_id         BIGINT NOT NULL,
    line_number     INT NOT NULL,
    content         TEXT NOT NULL,
    timestamp       DATETIME DEFAULT NOW(),
    stream          ENUM('stdout','stderr') DEFAULT 'stdout',
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    INDEX idx_task_line (task_id, line_number)
);

-- Running tasks guard (prevent duplicate execution)
CREATE TABLE running_tasks (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    cluster_id      BIGINT NOT NULL,
    step_number     VARCHAR(8) NOT NULL,
    task_id         BIGINT NOT NULL,
    created_at      DATETIME DEFAULT NOW(),
    FOREIGN KEY (cluster_id) REFERENCES clusters(id) ON DELETE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    UNIQUE KEY uk_cluster_step (cluster_id, step_number)
);

-- User cluster bindings
CREATE TABLE user_cluster_bindings (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id         BIGINT NOT NULL,
    cluster_id      BIGINT NOT NULL,
    role            ENUM('admin','viewer') DEFAULT 'viewer',
    created_at      DATETIME DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (cluster_id) REFERENCES clusters(id) ON DELETE CASCADE,
    UNIQUE KEY uk_user_cluster (user_id, cluster_id)
);

-- Audit log
CREATE TABLE audits (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id         BIGINT,
    cluster_id      BIGINT,
    action          VARCHAR(64) NOT NULL,
    resource_type   VARCHAR(32),
    resource_id     BIGINT,
    details         JSON,
    ip_address      VARCHAR(64),
    created_at      DATETIME DEFAULT NOW(),
    INDEX idx_cluster_time (cluster_id, created_at),
    INDEX idx_user_time (user_id, created_at)
);
