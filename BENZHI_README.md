基于 Go 实现的 TrainWash 列车自动洗车机控制平台项目，一款后端服务，完成列车定位落盘、刷组下放门控、水洗风干顺序控制、浓度告警与持久化恢复。

## 构建

```bash
go build -mod=vendor -o bin/trainwash.exe ./cmd/trainwash
```

## 运行

```bash
bin/trainwash.exe -addr :8080 -data ./data
```

服务启动后访问 http://localhost:8080/ 查看控制台总览页，接口文档见下方路由列表。

## 主要路由

- GET /api/state 全量运行状态
- POST /api/plan/wash 启动洗车
- POST /api/plan/complete 完成洗车
- POST /api/plan/stop 紧急停机
- POST /api/brush/lower 侧刷下放
- POST /api/chem/spray 洗涤剂喷洒
- POST /api/dry/start 风干机启动
- POST /api/pos/calibrate 位置零点标定
- GET /api/audit 操作记录

## 前端页面

- / 总览页（实时状态灯与告警列表）
- /wash 操作页
- /brushes 刷组页
- /water 水系统页
- /alarms 告警与药剂页

## 容器构建

```bash
sh build_benzhi_docker.sh
```
