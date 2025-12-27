# 🚀 AI Card Arena 部署指南

## 📋 前置要求

### 1. 环境依赖
- **Go**: 1.19+
- **Python**: 3.8+
- **MySQL**: 8.0+
- **Node.js**: 16+ (前端)

### 2. 安装 Python 依赖
```bash
cd arena-backend/engine
pip install -r ../requirements.txt
pip install fastapi uvicorn
```

### 3. 安装 Go 依赖
```bash
cd arena-backend
go mod tidy
go get github.com/gorilla/websocket
```

---

## 🔧 配置数据库

### 1. 创建数据库和用户
```sql
CREATE DATABASE ddz CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'ddz_user'@'localhost' IDENTIFIED BY 'ddz_password';
GRANT ALL PRIVILEGES ON ddz.* TO 'ddz_user'@'localhost';
FLUSH PRIVILEGES;
```

### 2. 更新配置
编辑 `arena-backend/config/database.go`，确保连接字符串正确：
```go
dsn := "root:root123456@tcp(127.0.0.1:3306)/ddz?charset=utf8mb4&parseTime=True&loc=Local"
```

---

## 🎮 启动服务

### 方式一：分别启动（推荐开发环境）

#### 1. 启动 Python 游戏引擎
```bash
cd arena-backend/engine
chmod +x start.sh
./start.sh

# 或手动启动
uvicorn main:app --host 0.0.0.0 --port 8081 --reload
```

验证启动：
```bash
curl http://localhost:8081/health
# 应返回: {"status":"healthy","active_games":0,"service":"AI Card Arena - Python Engine"}
```

#### 2. 启动 Go 后端服务
```bash
cd arena-backend
go run main.go
```

验证启动：
```bash
curl http://localhost:8080/api/games
# 应返回: {"games":[],"total":0}
```

---

### 方式二：使用启动脚本（推荐生产环境）

创建 `start_all.sh`:
```bash
#!/bin/bash

echo "🚀 启动 AI Card Arena 服务..."

# 启动 Python 引擎
cd arena-backend/engine
uvicorn main:app --host 0.0.0.0 --port 8081 &
PYTHON_PID=$!
echo "✅ Python 引擎已启动 (PID: $PYTHON_PID)"

# 等待 Python 引擎启动
sleep 3

# 启动 Go 后端
cd ..
go run main.go &
GO_PID=$!
echo "✅ Go 后端已启动 (PID: $GO_PID)"

echo ""
echo "=========================================="
echo "🎮 所有服务已启动"
echo "=========================================="
echo "Python 引擎: http://localhost:8081"
echo "Go 后端:    http://localhost:8080"
echo ""
echo "按 Ctrl+C 停止所有服务"
echo "=========================================="

# 等待中断信号
trap "kill $PYTHON_PID $GO_PID; exit" INT
wait
```

---

## 🧪 测试游戏流程

### 1. 使用测试脚本
```bash
cd arena-backend
chmod +x test_game_flow.sh
./test_game_flow.sh
```

### 2. 手动测试

#### 步骤 1: 注册用户
```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "player1",
    "password": "123456",
    "role": "developer"
  }'
```

重复三次，创建 player1, player2, player3

#### 步骤 2: 创建游戏
```bash
curl -X POST http://localhost:8080/api/games/create \
  -H "Content-Type: application/json" \
  -d '{
    "player_ids": [1, 2, 3],
    "seed": 42
  }'
```

记录返回的 `game_id`

#### 步骤 3: 获取游戏状态
```bash
curl "http://localhost:8080/api/games/{game_id}/state?player_id=0"
```

#### 步骤 4: 执行动作
```bash
curl -X POST http://localhost:8080/api/games/{game_id}/action \
  -H "Content-Type: application/json" \
  -d '{
    "player_id": 0,
    "action": "pass"
  }'
```

---

## 🌐 前端集成

### WebSocket 连接示例 (JavaScript)

```javascript
// 连接到游戏 WebSocket
const gameId = 'game_xxx';
const playerId = 0;
const ws = new WebSocket(`ws://localhost:8080/api/games/${gameId}/ws?player_id=${playerId}`);

// 连接成功
ws.onopen = () => {
  console.log('✅ 已连接到游戏');
  
  // 获取初始状态
  ws.send(JSON.stringify({
    type: 'get_state'
  }));
};

// 接收消息
ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  
  switch(msg.type) {
    case 'connected':
      console.log('欢迎消息:', msg.data);
      break;
      
    case 'state':
      console.log('游戏状态:', msg.data);
      updateGameUI(msg.data);
      break;
      
    case 'next_turn':
      console.log('下一回合:', msg.data);
      if (msg.data.next_player === playerId) {
        enablePlayerActions();
      }
      break;
      
    case 'game_over':
      console.log('游戏结束:', msg.data);
      showGameResult(msg.data);
      break;
      
    case 'error':
      console.error('错误:', msg.data.message);
      break;
  }
};

// 出牌
function playCard(action) {
  ws.send(JSON.stringify({
    type: 'action',
    data: {
      action: action  // 例如: "33", "pass", "AAAKK"
    }
  }));
}
```

### HTTP API 示例 (Vue.js)

```javascript
// api/game.js
import axios from 'axios';

const API_BASE = 'http://localhost:8080/api';

export const gameAPI = {
  // 创建游戏
  createGame(playerIds, seed = null) {
    return axios.post(`${API_BASE}/games/create`, {
      player_ids: playerIds,
      seed: seed
    });
  },
  
  // 获取游戏状态
  getGameState(gameId, playerId) {
    return axios.get(`${API_BASE}/games/${gameId}/state`, {
      params: { player_id: playerId }
    });
  },
  
  // 执行动作
  executeAction(gameId, playerId, action) {
    return axios.post(`${API_BASE}/games/${gameId}/action`, {
      player_id: playerId,
      action: action
    });
  },
  
  // 获取游戏会话
  getGameSession(gameId) {
    return axios.get(`${API_BASE}/games/${gameId}`);
  },
  
  // 列出所有游戏
  listGames() {
    return axios.get(`${API_BASE}/games`);
  }
};
```

---

## 📊 API 端点总览

### 用户相关
```
POST   /api/users/register      # 注册
POST   /api/users/login         # 登录
GET    /api/users/profile       # 获取用户信息 (需JWT)
```

### 游戏相关
```
POST   /api/games/create        # 创建游戏
POST   /api/games/:id/action    # 执行动作
GET    /api/games/:id/state     # 获取状态
GET    /api/games/:id           # 获取会话
DELETE /api/games/:id           # 删除游戏
GET    /api/games               # 列出所有游戏
GET    /api/games/:id/ws        # WebSocket 连接
```

---

## 🐛 常见问题

### 1. Python 引擎连接失败
**错误**: `请求 Python 引擎失败: dial tcp 127.0.0.1:8081: connect: connection refused`

**解决**:
- 确保 Python 引擎已启动: `curl http://localhost:8081/health`
- 检查端口是否被占用: `lsof -i :8081`

### 2. 数据库连接失败
**错误**: `数据库连接失败: Error 1045: Access denied`

**解决**:
- 检查 MySQL 是否运行: `systemctl status mysql`
- 验证用户名密码: `mysql -u root -p`
- 更新 `config/database.go` 中的连接字符串

### 3. CORS 错误
**错误**: 前端请求被 CORS 策略阻止

**解决**:
在 `main.go` 中添加 CORS 中间件:
```go
import "github.com/gin-contrib/cors"

func main() {
    r := gin.Default()
    
    // 配置 CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
    }))
    
    // ... 其他配置
}
```

---

## 📈 性能优化建议

### 1. 使用 Redis 存储游戏会话
当前使用内存存储，生产环境建议使用 Redis：
```go
// services/game_service.go
import "github.com/go-redis/redis/v8"

type GameService struct {
    engineClient *GameEngineClient
    redisClient  *redis.Client
}
```

### 2. 添加连接池
为 Python 引擎客户端添加连接池：
```go
httpClient: &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
    },
}
```

### 3. 添加日志系统
使用 `logrus` 或 `zap` 替代 `log`：
```bash
go get github.com/sirupsen/logrus
```

---

## 🔒 安全建议

1. **JWT 密钥**: 将 `utils/jwt.go` 中的密钥移到环境变量
2. **HTTPS**: 生产环境使用 HTTPS
3. **输入验证**: 严格验证所有用户输入
4. **速率限制**: 添加 API 速率限制中间件
5. **WebSocket 认证**: 为 WebSocket 连接添加 JWT 验证

---

## 📝 下一步开发建议

1. ✅ 完成基础游戏流程
2. ⬜ 添加 AI Agent 对战
3. ⬜ 实现游戏回放功能
4. ⬜ 添加排行榜系统
5. ⬜ 实现房间匹配系统
6. ⬜ 添加游戏统计和分析
7. ⬜ 开发前端 UI

---

## 📞 技术支持

如有问题，请查看：
- 项目文档: `arena-docs/`
- 功能分析: `arena-backend/功能分析.md`
- 游戏日志规范: `arena-docs/Game-log/Gamelog.md`
