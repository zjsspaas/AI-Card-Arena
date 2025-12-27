# 🎯 游戏引擎连接 Go 后端工作流程

## 📊 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                         前端层 (Vue.js)                          │
│  - 游戏 UI 界面                                                   │
│  - WebSocket 实时通信                                             │
│  - HTTP API 调用                                                 │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 │ HTTP/WebSocket
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Go 后端服务层 (Gin)                           │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Router (路由层)                                          │  │
│  │  - /api/users/*     用户相关                              │  │
│  │  - /api/games/*     游戏相关                              │  │
│  │  - /api/games/:id/ws WebSocket                           │  │
│  └────────────────┬─────────────────────────────────────────┘  │
│                   │                                              │
│  ┌────────────────▼─────────────────────────────────────────┐  │
│  │  Handlers (处理器层)                                      │  │
│  │  - user.go          用户注册/登录                         │  │
│  │  - game.go          游戏 CRUD                             │  │
│  │  - game_websocket.go WebSocket 处理                       │  │
│  └────────────────┬─────────────────────────────────────────┘  │
│                   │                                              │
│  ┌────────────────▼─────────────────────────────────────────┐  │
│  │  Services (业务逻辑层)                                     │  │
│  │  - game_service.go       游戏会话管理                     │  │
│  │  - game_engine_client.go Python 引擎客户端                │  │
│  └────────────────┬─────────────────────────────────────────┘  │
│                   │                                              │
│  ┌────────────────▼─────────────────────────────────────────┐  │
│  │  Models & Config                                          │  │
│  │  - model/user.go     用户模型                             │  │
│  │  - config/database.go 数据库配置                          │  │
│  └───────────────────────────────────────────────────────────┘  │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 │ HTTP (内部调用)
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                Python 游戏引擎层 (FastAPI)                        │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  FastAPI Routes                                           │  │
│  │  - POST /internal/game/init      初始化游戏               │  │
│  │  - POST /internal/game/step      执行动作                 │  │
│  │  - GET  /internal/game/:id/state 获取状态                 │  │
│  │  - DELETE /internal/game/:id     删除游戏                 │  │
│  └────────────────┬─────────────────────────────────────────┘  │
│                   │                                              │
│  ┌────────────────▼─────────────────────────────────────────┐  │
│  │  Game Logic (游戏逻辑)                                     │  │
│  │  - room.py   房间管理 (GameRoom)                          │  │
│  │  - engine.py 游戏引擎 (GameEngine)                        │  │
│  └────────────────┬─────────────────────────────────────────┘  │
│                   │                                              │
│  ┌────────────────▼─────────────────────────────────────────┐  │
│  │  RLCard Framework                                         │  │
│  │  - 斗地主游戏规则                                          │  │
│  │  - 动作空间管理                                            │  │
│  │  - 状态管理                                                │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🔄 完整游戏流程

### 1️⃣ 用户注册和登录

```
前端                    Go 后端                   MySQL
 │                        │                        │
 ├─ POST /api/users/register ──────────────────────┤
 │  {username, password, role}                     │
 │                        │                        │
 │                        ├─ bcrypt 加密密码       │
 │                        ├─ 创建用户 ─────────────▶│
 │                        │                        │
 │◀─ {user_id, username} ─┤                        │
 │                        │                        │
 ├─ POST /api/users/login ──────────────────────────┤
 │  {username, password}                           │
 │                        │                        │
 │                        ├─ 验证密码               │
 │                        ├─ 生成 JWT Token        │
 │                        │                        │
 │◀─ {token} ─────────────┤                        │
```

### 2️⃣ 创建游戏

```
前端                    Go 后端                   Python 引擎
 │                        │                        │
 ├─ POST /api/games/create ─────────────────────────┤
 │  {player_ids: [1,2,3], seed: 42}                │
 │                        │                        │
 │                        ├─ 验证玩家存在           │
 │                        │                        │
 │                        ├─ POST /internal/game/init ──▶│
 │                        │  {seed: 42}            │
 │                        │                        │
 │                        │                        ├─ 创建 GameRoom
 │                        │                        ├─ 初始化 RLCard
 │                        │                        ├─ 发牌、确定地主
 │                        │                        │
 │                        │◀─ {game_id, state} ────┤
 │                        │                        │
 │                        ├─ 创建 GameSession      │
 │                        ├─ 保存到内存/Redis      │
 │                        │                        │
 │◀─ {game_id, session, initial_state} ───────────┤
```

### 3️⃣ WebSocket 连接（推荐方式）

```
前端                    Go 后端                   Python 引擎
 │                        │                        │
 ├─ WS /api/games/:id/ws?player_id=0 ──────────────┤
 │                        │                        │
 │                        ├─ 升级为 WebSocket      │
 │                        ├─ 加入游戏房间          │
 │                        │                        │
 │◀─ {type: "connected"} ─┤                        │
 │                        │                        │
 ├─ {type: "get_state"} ──────────────────────────▶│
 │                        │                        │
 │                        ├─ GET /internal/game/:id/state ──▶│
 │                        │                        │
 │                        │◀─ {state} ─────────────┤
 │                        │                        │
 │◀─ {type: "state", data: {...}} ─────────────────┤
```

### 4️⃣ 游戏循环（出牌）

```
前端 (玩家0)            Go 后端                   Python 引擎
 │                        │                        │
 ├─ WS {type: "action", data: {action: "33"}} ────▶│
 │                        │                        │
 │                        ├─ 验证轮次               │
 │                        │                        │
 │                        ├─ POST /internal/game/step ──▶│
 │                        │  {game_id, player_id: 0, action: "33"}
 │                        │                        │
 │                        │                        ├─ 验证动作合法性
 │                        │                        ├─ 执行动作
 │                        │                        ├─ 更新游戏状态
 │                        │                        │
 │                        │◀─ {status: "playing", next_player: 1, state: {...}} ─┤
 │                        │                        │
 │                        ├─ 更新 GameSession      │
 │                        ├─ 广播给所有玩家        │
 │                        │                        │
 │◀─ {type: "next_turn", data: {next_player: 1, state: {...}}} ─┤
 │                        │                        │
 
前端 (玩家1)            Go 后端                   Python 引擎
 │                        │                        │
 │◀─ {type: "next_turn", data: {next_player: 1, state: {...}}} ─┤
 │                        │                        │
 ├─ WS {type: "action", data: {action: "pass"}} ──▶│
 │                        │                        │
 │                        ├─ POST /internal/game/step ──▶│
 │                        │  {game_id, player_id: 1, action: "pass"}
 │                        │                        │
 │                        │◀─ {status: "playing", next_player: 2, state: {...}} ─┤
 │                        │                        │
 │                        ├─ 广播给所有玩家        │
 │                        │                        │
 
... 循环直到游戏结束 ...
```

### 5️⃣ 游戏结束

```
前端 (玩家2)            Go 后端                   Python 引擎
 │                        │                        │
 ├─ WS {type: "action", data: {action: "K"}} ─────▶│
 │                        │                        │
 │                        ├─ POST /internal/game/step ──▶│
 │                        │  {game_id, player_id: 2, action: "K"}
 │                        │                        │
 │                        │                        ├─ 检测游戏结束
 │                        │                        ├─ 计算得分
 │                        │                        ├─ 删除房间
 │                        │                        │
 │                        │◀─ {status: "finished", payoffs: [1, -2, 1], winner: 0} ─┤
 │                        │                        │
 │                        ├─ 更新 GameSession      │
 │                        ├─ 广播给所有玩家        │
 │                        │                        │
 │◀─ {type: "game_over", data: {payoffs: [...], winner: 0}} ─┤
 │                        │                        │
 
所有玩家                Go 后端
 │                        │
 │◀─ {type: "game_over", data: {payoffs: [...], winner: 0}} ─┤
```

---

## 📝 详细开发步骤

### ✅ 阶段一：Python 引擎完善（已完成）

**文件**: `arena-backend/engine/main.py`

- [x] 完善 FastAPI 端点
- [x] 添加错误处理
- [x] 添加健康检查
- [x] 创建启动脚本

**测试**:
```bash
cd arena-backend/engine
uvicorn main:app --port 8081 --reload
curl http://localhost:8081/health
```

---

### ✅ 阶段二：Go 服务层创建（已完成）

**文件**: 
- `arena-backend/services/game_engine_client.go` - Python 引擎客户端
- `arena-backend/services/game_service.go` - 游戏业务逻辑

**功能**:
- [x] HTTP 客户端封装
- [x] 游戏会话管理
- [x] 错误处理

---

### ✅ 阶段三：Go 处理器创建（已完成）

**文件**:
- `arena-backend/handlers/game.go` - HTTP 处理器
- `arena-backend/handlers/game_websocket.go` - WebSocket 处理器

**功能**:
- [x] 创建游戏
- [x] 执行动作
- [x] 获取状态
- [x] WebSocket 实时通信

---

### ✅ 阶段四：路由配置（已完成）

**文件**: `arena-backend/router/router.go`

- [x] 添加游戏路由
- [x] 添加 WebSocket 路由

---

### ⬜ 阶段五：前端开发（待完成）

**目录**: `arena-frontend/doudizhu/`

**需要实现**:

1. **游戏大厅页面**
   - 创建游戏
   - 加入游戏
   - 游戏列表

2. **游戏界面**
   - 手牌显示
   - 出牌区域
   - 玩家信息
   - 动作按钮

3. **WebSocket 集成**
   - 连接管理
   - 消息处理
   - 状态同步

4. **游戏逻辑**
   - 出牌验证
   - 动作选择
   - 结果展示

---

## 🧪 测试流程

### 1. 单元测试

#### Python 引擎测试
```bash
cd arena-backend/engine
python -m pytest tests/
```

#### Go 服务测试
```bash
cd arena-backend
go test ./services/...
go test ./handlers/...
```

### 2. 集成测试

```bash
# 启动所有服务
./start_all.sh

# 运行测试脚本
./test_game_flow.sh
```

### 3. 手动测试

使用 Postman 或 curl 测试各个端点：

1. 注册用户
2. 创建游戏
3. 获取状态
4. 执行动作
5. 查看结果

---

## 🔧 开发工具推荐

### 1. API 测试
- **Postman**: REST API 测试
- **wscat**: WebSocket 测试
  ```bash
  npm install -g wscat
  wscat -c "ws://localhost:8080/api/games/game_xxx/ws?player_id=0"
  ```

### 2. 数据库管理
- **MySQL Workbench**: 可视化管理
- **DBeaver**: 跨平台数据库工具

### 3. 日志查看
- **GoLand**: Go 开发 IDE
- **VS Code**: 轻量级编辑器
- **PyCharm**: Python 开发 IDE

---

## 📊 性能监控

### 1. 添加日志

```go
// main.go
import "github.com/sirupsen/logrus"

func main() {
    logrus.SetLevel(logrus.InfoLevel)
    logrus.Info("服务启动...")
}
```

### 2. 添加指标

```go
// 使用 Prometheus
import "github.com/prometheus/client_golang/prometheus"

var (
    gamesCreated = prometheus.NewCounter(...)
    activeGames  = prometheus.NewGauge(...)
)
```

---

## 🚀 部署建议

### 开发环境
- Python 引擎: `uvicorn --reload`
- Go 后端: `go run main.go`
- 前端: `npm run dev`

### 生产环境
- Python 引擎: `gunicorn + uvicorn workers`
- Go 后端: 编译二进制 `go build`
- 前端: `npm run build` + Nginx
- 数据库: MySQL 主从复制
- 缓存: Redis 集群
- 负载均衡: Nginx

---

## 📚 相关文档

- [部署指南](./DEPLOYMENT.md)
- [功能分析](./功能分析.md)
- [游戏日志规范](../arena-docs/Game-log/Gamelog.md)
- [RLCard 文档](./rlcard_demo.md)

---

## 🎯 下一步计划

### 短期目标（1-2周）
1. ✅ 完成 Go-Python 连接
2. ⬜ 实现基础前端界面
3. ⬜ 完成端到端测试
4. ⬜ 添加 AI Agent 对战

### 中期目标（1个月）
1. ⬜ 实现房间匹配系统
2. ⬜ 添加游戏回放功能
3. ⬜ 实现排行榜
4. ⬜ 优化性能

### 长期目标（3个月）
1. ⬜ 支持多种卡牌游戏
2. ⬜ AI 训练平台
3. ⬜ 数据分析系统
4. ⬜ 移动端支持

---

## 💡 技术要点

### 1. 为什么使用 Go + Python 混合架构？

- **Go**: 高性能 HTTP 服务、并发处理、数据库操作
- **Python**: 游戏引擎、AI 算法、RLCard 生态

### 2. 为什么使用 WebSocket？

- 实时双向通信
- 减少 HTTP 轮询开销
- 更好的用户体验

### 3. 为什么需要游戏会话管理？

- 跟踪游戏状态
- 验证玩家权限
- 支持断线重连
- 游戏数据持久化

---

## 🤝 团队协作建议

### 1. 代码规范
- Go: 使用 `gofmt` 和 `golint`
- Python: 使用 `black` 和 `flake8`
- 提交前运行测试

### 2. Git 工作流
```bash
# 功能分支
git checkout -b feature/game-websocket

# 提交
git commit -m "feat: 添加游戏 WebSocket 支持"

# 合并前测试
go test ./...
python -m pytest

# 合并到主分支
git checkout main
git merge feature/game-websocket
```

### 3. 文档更新
- 每个新功能都要更新文档
- API 变更要更新接口文档
- 重要决策记录在 ADR (Architecture Decision Records)

---

祝开发顺利！🎉
