-- ksk8s schema rollback

DROP TABLE IF EXISTS audits;
DROP TABLE IF EXISTS user_cluster_bindings;
DROP TABLE IF EXISTS running_tasks;
DROP TABLE IF EXISTS logs;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS templates;
DROP TABLE IF EXISTS cluster_versions;
DROP TABLE IF EXISTS hosts_global_vars;
DROP TABLE IF EXISTS cluster_param_lists;
DROP TABLE IF EXISTS cluster_params;
DROP TABLE IF EXISTS cluster_nodes;
DROP TABLE IF EXISTS clusters;
DROP TABLE IF EXISTS users;
