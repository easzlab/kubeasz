# readme

openebs 使用启用lvm-localpv-controller，当创建StorageClass 时启用thinProvision: "yes"，需要注意vg_k8s_thinpool 的容量问题：

```
lvs
  LV                                       VG     Attr       LSize   Pool            Origin Data%  Meta%  Move Log Cpy%Sync Convert
  pvc-2214c3d8-83de-44e6-988c-6293277e9b1e vg_k8s Vwi-aotz--  5.00g vg_k8s_thinpool        2.91
  pvc-d8ea9413-5ddd-42db-a3f8-4363f745909a vg_k8s Vwi-aotz-- 20.00g vg_k8s_thinpool        2.23
  vg_k8s_thinpool                          vg_k8s twi-aotzD- 10.00g                        100.00 6.99
```

默认vg_k8s_thinpool 只有10G 大小，当如上图 Data% = 100% 时，就无法创建新的pv，需要扩容：

lvextend -L +190G vg_k8s/vg_k8s_thinpool

