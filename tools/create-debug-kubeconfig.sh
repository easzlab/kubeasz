#!/usr/bin/env bash
set +x
set -eo pipefail

usage() {
  cat <<'EOF'
用法: create-debug-kubeconfig.sh [选项]

  --admin-kubeconfig FILE  管理员 kubeconfig，默认 ~/.kube/config
  --output FILE           输出路径，默认 ./debug-kubeconfig，拒绝覆盖已有文件
  --duration TIME         令牌有效期，默认 30d，支持整数天(d)、小时(h)、分钟(m)
  --context NAME          管理员 kubeconfig 中的 context，默认 current-context
  --namespace NAME        新 ServiceAccount 所在的已有命名空间，默认 default
  -h, --help              显示帮助

权限作用于整个集群：内置 view，以及节点、命名空间、事件、存储、Secrets 和 RBAC 的读取权限。
可用 kubectl top 查看已启用的节点和 Pod 指标。
保留 pods/exec 的 get/create 以便进入容器排查。补充角色中的其余规则只授予读取权限。
Pod exec 可执行容器内命令，Secrets 可包含敏感凭证。请妥善保管输出文件。
生成的 kubeconfig 跳过 API Server 的 TLS 证书校验。
需要 kubectl 和 Python 3。集群未按请求签发令牌有效期时，不生成 kubeconfig。
每次成功签发都会创建独立的 ServiceAccount 和 RBAC 对象。令牌到期后对象仍保留。
EOF
}

die() {
  printf '%s\n' "$*" >&2
  exit 1
}

admin_config=~/.kube/config
output=./debug-kubeconfig
duration=30d
context=''
namespace=default
while (($#)); do
  case "$1" in
    --admin-kubeconfig|--output|--duration|--context|--namespace)
      (($# >= 2)) && [[ -n $2 && $2 != --* ]] || die "$1 缺少参数值"
      case "$1" in
        --admin-kubeconfig) admin_config=$2 ;;
        --output) output=$2 ;;
        --duration) duration=$2 ;;
        --context) context=$2 ;;
        --namespace) namespace=$2 ;;
      esac
      shift 2
      ;;
    -h|--help) usage; exit 0 ;;
    *) die "未知参数: $1（使用 --help 查看用法）" ;;
  esac
done

[[ -f $admin_config && -r $admin_config ]] || die "管理员 kubeconfig 不存在或不可读: $admin_config"
[[ $output != */ && ! -e $output && ! -L $output ]] || die "输出路径已存在或不是文件路径: $output"
[[ $namespace =~ ^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$ ]] || die "无效的命名空间: $namespace"
[[ $duration =~ ^([1-9][0-9]{0,8})([dhm])$ ]] || die '有效期须为正整数和单位，例如 30d、24h、30m'
amount=${BASH_REMATCH[1]}
case "${BASH_REMATCH[2]}" in
  d) seconds=$((amount * 86400)) ;;
  h) seconds=$((amount * 3600)) ;;
  m) seconds=$((amount * 60)) ;;
esac
((seconds >= 600 && seconds <= 9223372036)) || die '有效期超出范围，至少为 10m'
command -v kubectl >/dev/null || die '缺少 kubectl'
command -v python3 >/dev/null || die '缺少 python3'

admin=(kubectl --kubeconfig "$admin_config")
if [[ -z $context ]]; then context=$("${admin[@]}" config current-context); fi
admin+=(--context "$context")
umask 077
staged=$(mktemp "$(dirname -- "$output")/.debug-kubeconfig.XXXXXXXXXXXX")
name="debug-kubeconfig-$(printf '%s' "${staged##*.}" | LC_ALL=C tr '[:upper:]' '[:lower:]')"
created_account=0
created_rbac=()
published=0

cleanup() {
  if ((!published)); then
    if ((created_account)); then
      "${admin[@]}" -n "$namespace" delete "serviceaccount/$name" --ignore-not-found --wait=false >&2 ||
        printf '清理 ServiceAccount 失败: %s/%s\n' "$namespace" "$name" >&2
    fi
    if ((${#created_rbac[@]})); then
      "${admin[@]}" delete "${created_rbac[@]}" --ignore-not-found --wait=false >&2 ||
        printf '清理 RBAC 对象失败: %s\n' "${created_rbac[*]}" >&2
    fi
  fi
  rm -f -- "$staged"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

source_config=$("${admin[@]}" config view --minify -o json)
"${admin[@]}" -n "$namespace" create serviceaccount "$name" >/dev/null
created_account=1
token_response=$("${admin[@]}" -n "$namespace" create token "$name" --duration "${seconds}s" -o json)

expires=$(python3 - "$staged" "$name" "$namespace" "$seconds" 3<<<"$source_config" 4<<<"$token_response" <<'PY'
import base64
import datetime
import json
import os
import sys

path, name, namespace, requested = sys.argv[1:]
try:
    with os.fdopen(3) as config_file, os.fdopen(4) as response_file:
        cluster = json.load(config_file)["clusters"][0]
        status = json.load(response_file)["status"]
    if not cluster["name"] or not cluster["cluster"]["server"]:
        raise ValueError
    token = status["token"]
    payload = token.split(".")[1]
    claims = json.loads(base64.urlsafe_b64decode(payload + "=" * (-len(payload) % 4)))
    if type(claims["exp"]) is not int or type(claims["iat"]) is not int:
        raise ValueError
    # 令牌的 iat 和 exp 使用同一个签发时钟，避免响应时间戳取整产生的秒差。
    lifetime = claims["exp"] - claims["iat"]
    expires = datetime.datetime.fromtimestamp(claims["exp"], datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
except (KeyError, IndexError, TypeError, ValueError, OverflowError, OSError):
    sys.exit("无法读取集群连接信息或验证 TokenRequest 令牌的有效期。")

if lifetime != int(requested):
    sys.exit(
        f"集群签发的有效期为 {lifetime:g} 秒，到期时间 {expires}。"
        f"请求值为 {requested} 秒，请按集群允许的有效期调整 --duration。"
    )

# 重新组装 kubeconfig，只从源文件取集群连接信息。
cluster_settings = cluster["cluster"]
cluster_settings.pop("certificate-authority-data", None)
cluster_settings.pop("certificate-authority", None)
cluster_settings["insecure-skip-tls-verify"] = True
config = {
    "apiVersion": "v1",
    "kind": "Config",
    "clusters": [cluster],
    "users": [{"name": name, "user": {"token": token}}],
    "contexts": [{"name": name, "context": {
        "cluster": cluster["name"], "user": name, "namespace": namespace,
    }}],
    "current-context": name,
}
with open(path, "w", encoding="utf-8") as output:
    json.dump(config, output, ensure_ascii=False, indent=2)
    output.write("\n")
print(expires)
PY
)
unset source_config token_response

"${admin[@]}" create -f - >/dev/null <<EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: $name
rules:
  - apiGroups: [""]
    resources: ["pods/exec"]
    verbs: ["get", "create"]
  - apiGroups: [""]
    resources: ["secrets", "nodes", "namespaces", "events", "persistentvolumes"]
    verbs: ["get", "list", "watch"]
  - apiGroups: ["events.k8s.io"]
    resources: ["events"]
    verbs: ["get", "list", "watch"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["storageclasses", "csidrivers", "csinodes", "volumeattachments"]
    verbs: ["get", "list", "watch"]
  - apiGroups: ["rbac.authorization.k8s.io"]
    resources: ["roles", "rolebindings", "clusterroles", "clusterrolebindings"]
    verbs: ["get", "list", "watch"]
  - apiGroups: ["metrics.k8s.io"]
    resources: ["nodes", "pods"]
    verbs: ["get", "list", "watch"]
EOF
created_rbac+=("clusterrole/$name")
"${admin[@]}" create clusterrolebinding "$name-view" --clusterrole=view \
  --serviceaccount="$namespace:$name" >/dev/null
created_rbac+=("clusterrolebinding/$name-view")
"${admin[@]}" create clusterrolebinding "$name-extra" --clusterrole="$name" \
  --serviceaccount="$namespace:$name" >/dev/null
created_rbac+=("clusterrolebinding/$name-extra")

# 同目录硬链接原子发布，拒绝覆盖并发生成的文件。
python3 -c 'import os, sys; os.link(sys.argv[1], sys.argv[2])' "$staged" "$output"
published=1
printf '已生成 %s，到期时间 %s\n' "$output" "$expires"
printf 'ServiceAccount %s/%s，ClusterRole %s，ClusterRoleBindings %s-view 和 %s-extra\n' \
  "$namespace" "$name" "$name" "$name" "$name"

