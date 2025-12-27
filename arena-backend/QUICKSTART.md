# ⚡ 快速开始指南

## 🎯 5分钟启动游戏服务

### 前置准备

确保已安装：
- Go 1.19+
- Python 3.8+
- MySQL 8.0+

### 步骤 1: 安装依赖

```bash
# Python 依赖
cd arena-backend
pip install -r requirements.txt
pip install fastapi uvicorn

# Go 依赖
go mod tidy
go get github.com/gorilla/websocket
```

### 步骤 2: 配置数据库

```bash
# 登录 MySQL
mysql -u root -p

# 执行以下 SQL
CREATE DATABASE ddz CHARACTER SET utf8mb4;
```

更新 `config/database.go` 中的数据库连接信息。

### 步骤 3: 创建日志目录

```bash
mkdir -p logs
```

### 步骤 4: 启动服务

```bash
# 给启动脚本添加执行权限
chmod +x start_all.sh
chmod +x engine/start.sh
chmod +x test_game_flow.sh

# 启动所有服务
./start_all.sh
```

### 步骤 5: 测试

打开新终端，运行测试脚本：

```bash
./test_game_flow.sh
```

---

## 🎮 快速测试游戏流程

### 1. 注册三个玩家

```bash
# 玩家 1
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"player1","password":"123456","role":"developer"}'

# 玩家 2
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"player2","password":"123456","role":"developer"}'

# 玩家 3
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"player3","password":"123456","role":"developer"}'
```

### 2. 创建游戏

```bash
curl -X POST http://localhost:8080/api/games/create \
  -H "Content-Type: application/json" \
  -d '{"player_ids":[1,2,3],"seed":42}' | jq '.'
```

记录返回的 `game_id`，例如：`game_abc123`

### 3. 获取游戏状态

```bash
curl "http://localhost:8080/api/games/game_abc123/state?player_id=0" | jq '.'
```

你会看到：
- 当前玩家手牌
- 可用动作列表
- 地主信息
- 其他玩家剩余牌数

### 4. 执行出牌动作

```bash
# 假设合法动作包含 "pass"
curl -X POST http://localhost:8080/api/games/game_abc123/action \
  -H "Content-Type: application/json" \
  -d '{"player_id":0,"action":"pass"}' | jq '.'
```

### 5. 继续游戏

重复步骤 3 和 4，直到游戏结束。

---

## 🌐 WebSocket 测试

### 使用 wscat

```bash
# 安装 wscat
npm install -g wscat

# 连接到游戏
wscat -c "ws://localhost:8080/api/games/game_abc123/ws?player_id=0"

# 连接后发送消息
> {"type":"get_state"}

# 出牌
> {"type":"action","data":{"action":"pass"}}
```

### 使用浏览器控制台

```javascript
const ws = new WebSocket('ws://localhost:8080/api/games/game_abc123/ws?player_id=0');

ws.onopen = () => {
  console.log('✅ 已连接');
  ws.send(JSON.stringify({type: 'get_state'}));
};

ws.onmessage = (event) => {
  console.log('📨 收到消息:', JSON.parse(event.data));
};

// 出牌
ws.send(JSON.stringify({
  type: 'action',
  data: {action: 'pass'}
}));
```

---

## 📊 API 端点速查

### 用户相关
```bash
# 注册
POST /api/users/register
Body: {"username":"xxx","password":"xxx","role":"developer"}

# 登录
POST /api/users/login
Body: {"username":"xxx","password":"xxx"}

# 获取个人信息（需要 JWT）
GET /api/users/profile
Header: Authorization: Bearer <token>
```

### 游戏相关
```bash
# 创建游戏
POST /api/games/create
Body: {"player_ids":[1,2,3],"seed":42}

# 获取游戏状态
GET /api/games/:game_id/state?player_id=0

# 执行动作
POST /api/games/:game_id/action
Body: {"player_id":0,"action":"pass"}

# 获取游戏会话
GET /api/games/:game_id

# 列出所有游戏
GET /api/games

# WebSocket 连接
WS /api/games/:game_id/ws?player_id=0
```

---

## 🐛 常见问题

### Python 引擎启动失败

**错误**: `ModuleNotFoundError: No module named 'rlcard'`

**解决**:
```bash
pip install rlcard numpy termcolor fastapi uvicorn
```

### Go 编译错误

**错误**: `cannot find package "github.com/gorilla/websocket"`

**解决**:
```bash
go get github.com/gorilla/websocket
go mod tidy
```

### 数据库连接失败

**错误**: `Error 1045: Access denied`

**解决**:
1. 检查 MySQL 是否运行: `systemctl status mysql`
2. 验证用户名密码
3. 更新 `config/database.go` 中的连接字符串

### 端口被占用

**错误**: `bind: address already in use`

**解决**:
```bash
# 查找占用端口的进程
lsof -i :8080
lsof -i :8081

# 杀死进程
kill -9 <PID>
```

---

## 📚 下一步

- 📖 阅读 [完整工作流程](./WORKFLOW.md)
- 🚀 查看 [部署指南](./DEPLOYMENT.md)
- 🎮 开发前端界面
- 🤖 添加 AI Agent

---

## 💡 提示

### 查看实时日志

```bash
# Python 引擎日志
tail -f logs/python_engine.log

# Go 后端日志
tail -f logs/go_backend.log
```

### 停止服务

按 `Ctrl+C` 或：

```bash
# 查找进程
ps aux | grep uvicorn
ps aux | grep "go run"

# 杀死进程
kill <PID>
```

### 重启服务

```bash
# 停止当前服务（Ctrl+C）
# 重新运行
./start_all.sh
```

---

## 🎉 成功标志

如果看到以下输出，说明服务启动成功：

```
==========================================
✨ 所有服务已成功启动！
==========================================

📍 服务地址:
   Python 引擎: http://localhost:8081
   Go 后端:    http://localhost:8080
```

现在可以开始开发了！🚀
