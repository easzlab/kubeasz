# LVM 操作

以下是使用 parted 对 /dev/sdb 进行分区并配置 LVM 的完整操作流程，包含详细解释和注意事项。

## 1. 磁盘分区 (使用 parted)

``` bash
sudo parted /dev/sdb
# 在交互界面执行：
(parted) mklabel gpt           # 创建 GPT 分区表（兼容大容量磁盘）
(parted) mkpart lvm 0% 100%    # 创建占用整个磁盘的 LVM 分区
(parted) set 1 lvm on          # 设置分区类型为 LVM
(parted) print                 # 验证分区信息
(parted) quit                  # 退出
# 验证分区结果：应看到 sdb1 分区
lsblk /dev/sdb
```

## 2. 创建物理卷 (PV)

``` bash
sudo pvcreate /dev/sdb1         # 将分区初始化为物理卷
sudo pvs                        # 查看已创建的物理卷
# 输出示例：
  PV         VG Fmt  Attr PSize   PFree
  /dev/sdb1     lvm2 ---  100.00g 100.00g
```

## 3. 创建卷组 (VG)

``` bash
sudo vgcreate vg_data /dev/sdb1  # 创建名为 vg_data 的卷组
sudo vgs                         # 查看卷组信息
# 参数说明：
● vg_data：自定义卷组名称
● -s 4M：可指定 PE 大小（默认 4MB）
```

## 4. 创建逻辑卷 (LV)

```
sudo lvcreate -n lv_app -L 50G vg_data  # 创建 50G 的逻辑卷
sudo lvs                               # 查看逻辑卷
可选参数：
● -l 100%FREE：使用全部剩余空间
● -i 3 -I 4：创建条带化卷（需多磁盘）
```

## 5. 创建文件系统

```
sudo mkfs.xfs /dev/vg_data/lv_app  # 创建 XFS 文件系统
# 或使用 ext4：
# sudo mkfs.ext4 /dev/vg_data/lv_app
格式选择建议：
● XFS：适合大文件/高并发
● ext4：兼容性好
```

## 6. 挂载文件系统

```
sudo mkdir /data            	     		# 创建挂载点
sudo mount /dev/vg_data/lv_app /data  		# 临时挂载
df -hT /data                          		# 验证挂载
# 持久化挂载：
echo '/dev/mapper/vg_data-lv_app /data xfs defaults 0 0' | sudo tee -a /etc/fstab
sudo mount -a                    # 测试 fstab 配置
```

## 完整操作流程图
graph TD A[磁盘/dev/sdb] --> B[parted创建GPT分区] B --> C[pvcreate创建物理卷] C --> D[vgcreate创建卷组] D --> E[lvcreate创建逻辑卷] E --> F[mkfs创建文件系统] F --> G[mount挂载使用]
mermaidgraph TD
    A[磁盘/dev/sdb] --> B[parted创建GPT分区]
    B --> C[pvcreate创建物理卷]
    C --> D[vgcreate创建卷组]
    D --> E[lvcreate创建逻辑卷]
    E --> F[mkfs创建文件系统]
    F --> G[mount挂载使用]

关键命令速查表

| 操作 | 命令 |
|:---|:---|
|查看块设备|lsblk|
|验证分区表|parted /dev/sdb print|
|扩展逻辑卷|lvextend -L +10G /dev/vg_data/lv_app|
|扩展文件系统 (XFS)|xfs_growfs /data|
|扩展文件系统 (ext4)|resize2fs /dev/vg_data/lv_app|
|删除卷组|vgremove vg_data|

注意事项
● 数据备份：操作前确认磁盘无重要数据
● 容量对齐：生产环境建议保持 1MB 对齐（parted 使用 % 单位自动对齐）
● 在线扩展：XFS 支持在线扩容，但不支持缩小
● RAID 整合：可在 LVM 层整合多个 PV 实现软 RAID
● 快照功能：使用 lvcreate -s 创建快照卷实现备份
通过以上步骤，您已成功将原始磁盘配置为可弹性管理的存储空间。后续可通过 LVM 的动态调整特性，实现无需卸载的存储扩容。
