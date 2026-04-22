package service

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/easzlab/ksk8s/internal/model"
	"gopkg.in/yaml.v3"
)

type ConfigGenerator struct{}

func NewConfigGenerator() *ConfigGenerator {
	return &ConfigGenerator{}
}

// ParamMeta holds display order, type and comment hint for a config.yml parameter.
type ParamMeta struct {
	Key     string `json:"key"`
	Type    string `json:"type"` // "scalar" or "list"
	Comment string `json:"comment"`
}

// HostsVarMeta holds display order and comment hint for an [all:vars] variable.
type HostsVarMeta struct {
	Key     string `json:"key"`
	Comment string `json:"comment"`
}

// HostsSectionMeta holds display order and comments for a hosts section.
type HostsSectionMeta struct {
	Name     string   `json:"name"`
	Comments []string `json:"comments"`
}

// GenerateHostsMeta renders hosts INI with template-derived section/var order and comments.
func (g *ConfigGenerator) GenerateHostsMeta(cluster *model.Cluster, nodes []model.ClusterNode, vars *model.HostsGlobalVars, sectionMeta []HostsSectionMeta, varMeta []HostsVarMeta) (string, error) {
	// Group nodes
	groups := make(map[string][]model.ClusterNode)
	for _, n := range nodes {
		groups[n.GroupName] = append(groups[n.GroupName], n)
	}

	var buf bytes.Buffer

	// Use template section order if available
	if len(sectionMeta) > 0 {
		for _, sec := range sectionMeta {
			for _, c := range sec.Comments {
				buf.WriteString(fmt.Sprintf("# %s\n", c))
			}
			buf.WriteString(fmt.Sprintf("[%s]\n", sec.Name))

			if sec.Name == "all:vars" {
				for _, vm := range varMeta {
					if vm.Comment != "" {
						for _, line := range strings.Split(vm.Comment, "\n") {
							buf.WriteString(fmt.Sprintf("# %s\n", line))
						}
					}
					line := formatHostsVarLine(vm.Key, vars)
					if line != "" {
						buf.WriteString(line + "\n")
					}
				}
				buf.WriteString("\n")
				continue
			}

			if members, ok := groups[sec.Name]; ok && len(members) > 0 {
				for _, n := range members {
					line := n.IPAddress
					if n.K8sNodename != "" {
						line += fmt.Sprintf(" k8s_nodename=%s", n.K8sNodename)
					}
					if n.NewInstall {
						line += " new_install=yes"
					}
					if n.LBRole != nil && *n.LBRole != "" {
						line += fmt.Sprintf(" lb_role=%s", *n.LBRole)
					}
					if n.ExApiserverVIP != nil && *n.ExApiserverVIP != "" {
						line += fmt.Sprintf(" ex_apiserver_vip=%s", *n.ExApiserverVIP)
					}
					if n.ExApiserverPort != nil && *n.ExApiserverPort != "" {
						line += fmt.Sprintf(" ex_apiserver_port=%s", *n.ExApiserverPort)
					}
					buf.WriteString(line + "\n")
				}
			}
			buf.WriteString("\n")
		}
		return buf.String(), nil
	}

	// Fallback: hard-coded order (no comments)
	return g.GenerateHosts(cluster, nodes, vars)
}

// GenerateHosts renders the hosts INI file from structured data (no comments).
func (g *ConfigGenerator) GenerateHosts(cluster *model.Cluster, nodes []model.ClusterNode, vars *model.HostsGlobalVars) (string, error) {
	var buf bytes.Buffer

	// Group nodes by group_name
	groups := make(map[string][]model.ClusterNode)
	for _, n := range nodes {
		groups[n.GroupName] = append(groups[n.GroupName], n)
	}

	// Order matters for Ansible
	// Always emit section headers (even empty) so optional groups preserve format and order
	groupOrder := []string{"etcd", "kube_master", "kube_node", "harbor", "ex_lb", "chrony"}
	for _, group := range groupOrder {
		buf.WriteString(fmt.Sprintf("[%s]\n", group))
		if members, ok := groups[group]; ok && len(members) > 0 {
			for _, n := range members {
				line := n.IPAddress
				if n.K8sNodename != "" {
					line += fmt.Sprintf(" k8s_nodename=%s", n.K8sNodename)
				}
				if n.NewInstall {
					line += " new_install=yes"
				}
				if n.LBRole != nil && *n.LBRole != "" {
					line += fmt.Sprintf(" lb_role=%s", *n.LBRole)
				}
				if n.ExApiserverVIP != nil && *n.ExApiserverVIP != "" {
					line += fmt.Sprintf(" ex_apiserver_vip=%s", *n.ExApiserverVIP)
				}
				if n.ExApiserverPort != nil && *n.ExApiserverPort != "" {
					line += fmt.Sprintf(" ex_apiserver_port=%s", *n.ExApiserverPort)
				}
				buf.WriteString(line + "\n")
			}
		}
		buf.WriteString("\n")
	}

	// [all:vars] section
	buf.WriteString("[all:vars]\n")
	buf.WriteString(fmt.Sprintf("SECURE_PORT=\"%s\"\n", vars.SecurePort))
	buf.WriteString(fmt.Sprintf("CONTAINER_RUNTIME=\"%s\"\n", vars.ContainerRuntime))
	buf.WriteString(fmt.Sprintf("CLUSTER_NETWORK=\"%s\"\n", vars.ClusterNetwork))
	buf.WriteString(fmt.Sprintf("PROXY_MODE=\"%s\"\n", vars.ProxyMode))
	buf.WriteString(fmt.Sprintf("SERVICE_CIDR=\"%s\"\n", vars.ServiceCIDR))
	buf.WriteString(fmt.Sprintf("CLUSTER_CIDR=\"%s\"\n", vars.ClusterCIDR))
	buf.WriteString(fmt.Sprintf("NODE_PORT_RANGE=\"%s\"\n", vars.NodePortRange))
	buf.WriteString(fmt.Sprintf("CLUSTER_DNS_DOMAIN=\"%s\"\n", vars.ClusterDNSDomain))
	buf.WriteString(fmt.Sprintf("bin_dir=\"%s\"\n", vars.BinDir))
	buf.WriteString(fmt.Sprintf("base_dir=\"%s\"\n", vars.BaseDir))
	buf.WriteString(fmt.Sprintf("cluster_dir=\"%s\"\n", vars.ClusterDir))
	buf.WriteString(fmt.Sprintf("ca_dir=\"%s\"\n", vars.CaDir))
	if vars.K8sNodename != "" {
		buf.WriteString(fmt.Sprintf("k8s_nodename=\"%s\"\n", vars.K8sNodename))
	}
	// Use defaults for ansible vars to avoid empty strings that break ansible parsing
	ansiblePython := vars.AnsiblePythonInterpreter
	if ansiblePython == "" {
		ansiblePython = "/usr/bin/python3"
	}
	buf.WriteString(fmt.Sprintf("ansible_python_interpreter=\"%s\"\n", ansiblePython))

	ansibleUser := vars.AnsibleUser
	if ansibleUser == "" {
		ansibleUser = "root"
	}
	buf.WriteString(fmt.Sprintf("ansible_user=\"%s\"\n", ansibleUser))

	ansibleBecome := vars.AnsibleBecome
	if ansibleBecome == "" {
		ansibleBecome = "no"
	}
	buf.WriteString(fmt.Sprintf("ansible_become=\"%s\"\n", ansibleBecome))

	return buf.String(), nil
}

func formatHostsVarLine(key string, vars *model.HostsGlobalVars) string {
	switch key {
	case "SECURE_PORT":
		return fmt.Sprintf("SECURE_PORT=\"%s\"", vars.SecurePort)
	case "CONTAINER_RUNTIME":
		return fmt.Sprintf("CONTAINER_RUNTIME=\"%s\"", vars.ContainerRuntime)
	case "CLUSTER_NETWORK":
		return fmt.Sprintf("CLUSTER_NETWORK=\"%s\"", vars.ClusterNetwork)
	case "PROXY_MODE":
		return fmt.Sprintf("PROXY_MODE=\"%s\"", vars.ProxyMode)
	case "SERVICE_CIDR":
		return fmt.Sprintf("SERVICE_CIDR=\"%s\"", vars.ServiceCIDR)
	case "CLUSTER_CIDR":
		return fmt.Sprintf("CLUSTER_CIDR=\"%s\"", vars.ClusterCIDR)
	case "NODE_PORT_RANGE":
		return fmt.Sprintf("NODE_PORT_RANGE=\"%s\"", vars.NodePortRange)
	case "CLUSTER_DNS_DOMAIN":
		return fmt.Sprintf("CLUSTER_DNS_DOMAIN=\"%s\"", vars.ClusterDNSDomain)
	case "bin_dir":
		return fmt.Sprintf("bin_dir=\"%s\"", vars.BinDir)
	case "base_dir":
		return fmt.Sprintf("base_dir=\"%s\"", vars.BaseDir)
	case "cluster_dir":
		return fmt.Sprintf("cluster_dir=\"%s\"", vars.ClusterDir)
	case "ca_dir":
		return fmt.Sprintf("ca_dir=\"%s\"", vars.CaDir)
	case "k8s_nodename":
		if vars.K8sNodename != "" {
			return fmt.Sprintf("k8s_nodename=\"%s\"", vars.K8sNodename)
		}
		return ""
	case "ansible_python_interpreter":
		v := vars.AnsiblePythonInterpreter
		if v == "" {
			v = "/usr/bin/python3"
		}
		return fmt.Sprintf("ansible_python_interpreter=\"%s\"", v)
	case "ansible_user":
		v := vars.AnsibleUser
		if v == "" {
			v = "root"
		}
		return fmt.Sprintf("ansible_user=\"%s\"", v)
	case "ansible_become":
		v := vars.AnsibleBecome
		if v == "" {
			v = "no"
		}
		return fmt.Sprintf("ansible_become=\"%s\"", v)
	}
	return ""
}

// GenerateConfigYAML renders config.yml from structured params (fallback when no base text)
func (g *ConfigGenerator) GenerateConfigYAML(cluster *model.Cluster, params []model.ClusterParam, paramLists []model.ClusterParamList) (string, error) {
	config := make(map[string]interface{})

	// Group params by param_group
	for _, p := range params {
		config[p.ParamKey] = p.ParamValue
	}

	// Handle list params (e.g. MASTER_CERT_HOSTS, INSECURE_REG)
	listMap := make(map[string][]string)
	for _, pl := range paramLists {
		listMap[pl.ParamKey] = append(listMap[pl.ParamKey], pl.ItemValue)
	}
	for key, values := range listMap {
		config[key] = values
	}

	// Ensure cluster_name is set
	config["CLUSTER_NAME"] = cluster.Name

	out, err := yaml.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("marshal config: %w", err)
	}
	return string(out), nil
}

// UpdateConfigYAML merges param changes into an existing YAML text, preserving comments, order and format.
func (g *ConfigGenerator) UpdateConfigYAML(original string, params map[string]string, lists map[string][]string, clusterName string) (string, error) {
	if original == "" {
		return "", nil
	}

	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(original), &doc); err != nil {
		return "", fmt.Errorf("unmarshal original config: %w", err)
	}

	if len(doc.Content) == 0 {
		return original, nil
	}
	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return original, nil
	}

	allKeys := make(map[string]bool)
	for k := range params {
		allKeys[k] = true
	}
	for k := range lists {
		allKeys[k] = true
	}
	// CLUSTER_NAME is always managed
	allKeys["CLUSTER_NAME"] = true

	// Update existing keys
	existingKeys := make(map[string]bool)
	for i := 0; i < len(root.Content); i += 2 {
		keyNode := root.Content[i]
		valNode := root.Content[i+1]
		key := keyNode.Value

		if key == "CLUSTER_NAME" {
			valNode.Value = clusterName
			existingKeys[key] = true
			continue
		}

		if newVal, ok := params[key]; ok {
			valNode.Value = newVal
			// Convert sequence back to scalar if needed
			if valNode.Kind == yaml.SequenceNode {
				valNode.Kind = yaml.ScalarNode
				valNode.Tag = "!!str"
				valNode.Content = nil
			}
			existingKeys[key] = true
		} else if newList, ok := lists[key]; ok {
			valNode.Kind = yaml.SequenceNode
			valNode.Tag = "!!seq"
			valNode.Content = nil
			for _, item := range newList {
				valNode.Content = append(valNode.Content, &yaml.Node{
					Kind:  yaml.ScalarNode,
					Tag:   "!!str",
					Value: item,
				})
			}
			existingKeys[key] = true
		}
	}

	// Add new keys at the end (sorted for stability)
	var newScalarKeys []string
	for key := range params {
		if !existingKeys[key] {
			newScalarKeys = append(newScalarKeys, key)
		}
	}
	if !existingKeys["CLUSTER_NAME"] {
		newScalarKeys = append(newScalarKeys, "CLUSTER_NAME")
	}
	sort.Strings(newScalarKeys)
	for _, key := range newScalarKeys {
		val := params[key]
		if key == "CLUSTER_NAME" {
			val = clusterName
		}
		root.Content = append(root.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key},
			&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: val},
		)
	}

	var newListKeys []string
	for key := range lists {
		if !existingKeys[key] {
			newListKeys = append(newListKeys, key)
		}
	}
	sort.Strings(newListKeys)
	for _, key := range newListKeys {
		seq := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, item := range lists[key] {
			seq.Content = append(seq.Content, &yaml.Node{
				Kind: yaml.ScalarNode, Tag: "!!str", Value: item,
			})
		}
		root.Content = append(root.Content,
			&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key},
			seq,
		)
	}

	out, err := yaml.Marshal(&doc)
	if err != nil {
		return "", fmt.Errorf("marshal updated config: %w", err)
	}
	return string(out), nil
}

// ParseHostsText parses a hosts INI text into structured nodes and vars
func (g *ConfigGenerator) ParseHostsText(content string) ([]model.ClusterNode, *model.HostsGlobalVars, error) {
	vars := &model.HostsGlobalVars{}
	var nodes []model.ClusterNode
	var currentGroup string

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentGroup = line[1 : len(line)-1]
			continue
		}

		if currentGroup == "all:vars" {
			parseVarLine(line, vars)
			continue
		}

		if currentGroup != "" {
			node := parseNodeLine(line, currentGroup)
			if node.IPAddress != "" {
				nodes = append(nodes, node)
			}
		}
	}

	return nodes, vars, nil
}

func parseNodeLine(line, group string) model.ClusterNode {
	node := model.ClusterNode{GroupName: group, NewInstall: false}
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return node
	}
	node.IPAddress = fields[0]
	for _, f := range fields[1:] {
		parts := strings.SplitN(f, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := parts[0], parts[1]
		switch key {
		case "k8s_nodename":
			node.K8sNodename = val
		case "new_install":
			node.NewInstall = val == "yes"
		case "lb_role":
			node.LBRole = &val
		case "ex_apiserver_vip":
			node.ExApiserverVIP = &val
		case "ex_apiserver_port":
			node.ExApiserverPort = &val
		}
	}
	return node
}

func parseVarLine(line string, vars *model.HostsGlobalVars) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return
	}
	key := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])
	val = strings.Trim(val, `"`)

	switch key {
	case "SECURE_PORT":
		vars.SecurePort = val
	case "CONTAINER_RUNTIME":
		vars.ContainerRuntime = val
	case "CLUSTER_NETWORK":
		vars.ClusterNetwork = val
	case "PROXY_MODE":
		vars.ProxyMode = val
	case "SERVICE_CIDR":
		vars.ServiceCIDR = val
	case "CLUSTER_CIDR":
		vars.ClusterCIDR = val
	case "NODE_PORT_RANGE":
		vars.NodePortRange = val
	case "CLUSTER_DNS_DOMAIN":
		vars.ClusterDNSDomain = val
	case "bin_dir":
		vars.BinDir = val
	case "base_dir":
		vars.BaseDir = val
	case "cluster_dir":
		vars.ClusterDir = val
	case "ca_dir":
		vars.CaDir = val
	case "k8s_nodename":
		vars.K8sNodename = val
	case "ansible_python_interpreter":
		vars.AnsiblePythonInterpreter = val
	case "ansible_user":
		vars.AnsibleUser = val
	case "ansible_become":
		vars.AnsibleBecome = val
	}
}

// ParseConfigYAML parses config.yml text into structured params
func (g *ConfigGenerator) ParseConfigYAML(content string) ([]model.ClusterParam, []model.ClusterParamList, error) {
	var raw map[string]interface{}
	if err := yaml.Unmarshal([]byte(content), &raw); err != nil {
		return nil, nil, fmt.Errorf("unmarshal config: %w", err)
	}

	var params []model.ClusterParam
	var lists []model.ClusterParamList

	for key, val := range raw {
		switch v := val.(type) {
		case []interface{}:
			if len(v) == 0 {
				// Preserve empty list with a marker entry
				lists = append(lists, model.ClusterParamList{
					ParamKey:  key,
					ItemValue: "",
					SortOrder: 0,
				})
			} else {
				for i, item := range v {
					lists = append(lists, model.ClusterParamList{
						ParamKey:  key,
						ItemValue: fmt.Sprintf("%v", item),
						SortOrder: i,
					})
				}
			}
		default:
			params = append(params, model.ClusterParam{
				ParamKey:   key,
				ParamValue: fmt.Sprintf("%v", v),
			})
		}
	}

	// Sort for deterministic output
	sort.Slice(params, func(i, j int) bool { return params[i].ParamKey < params[j].ParamKey })
	sort.Slice(lists, func(i, j int) bool {
		if lists[i].ParamKey == lists[j].ParamKey {
			return lists[i].SortOrder < lists[j].SortOrder
		}
		return lists[i].ParamKey < lists[j].ParamKey
	})

	return params, lists, nil
}

// ParseConfigMeta extracts ordered parameter metadata (key, type, comment) from config.yml text.
func (g *ConfigGenerator) ParseConfigMeta(content string) ([]ParamMeta, error) {
	if content == "" {
		return nil, nil
	}
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(content), &doc); err != nil {
		return nil, fmt.Errorf("unmarshal config meta: %w", err)
	}
	if len(doc.Content) == 0 {
		return nil, nil
	}
	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return nil, nil
	}

	var meta []ParamMeta
	for i := 0; i < len(root.Content); i += 2 {
		keyNode := root.Content[i]
		valNode := root.Content[i+1]

		typ := "scalar"
		if valNode.Kind == yaml.SequenceNode {
			typ = "list"
		}

		comment := cleanYAMLComment(keyNode.HeadComment)

		meta = append(meta, ParamMeta{
			Key:     keyNode.Value,
			Type:    typ,
			Comment: comment,
		})
	}
	return meta, nil
}

// ParseHostsMeta extracts section order/comments and [all:vars] variable order/comments from hosts text.
func (g *ConfigGenerator) ParseHostsMeta(content string) ([]HostsSectionMeta, []HostsVarMeta, error) {
	if content == "" {
		return nil, nil, nil
	}

	lines := strings.Split(content, "\n")
	var sections []HostsSectionMeta
	var varMeta []HostsVarMeta
	var pendingComments []string
	currentSection := ""

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			pendingComments = nil
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			pendingComments = append(pendingComments, strings.TrimSpace(trimmed[1:]))
			continue
		}
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			sectionName := trimmed[1 : len(trimmed)-1]
			sections = append(sections, HostsSectionMeta{
				Name:     sectionName,
				Comments: append([]string(nil), pendingComments...),
			})
			currentSection = sectionName
			pendingComments = nil
			continue
		}
		if currentSection == "all:vars" {
			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				varMeta = append(varMeta, HostsVarMeta{
					Key:     key,
					Comment: strings.Join(pendingComments, "\n"),
				})
			}
			pendingComments = nil
		}
	}

	return sections, varMeta, nil
}

func cleanYAMLComment(raw string) string {
	lines := strings.Split(raw, "\n")
	var cleaned []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			cleaned = append(cleaned, strings.TrimSpace(line[1:]))
		} else if line != "" {
			cleaned = append(cleaned, line)
		}
	}
	return strings.Join(cleaned, "\n")
}
