#!/bin/bash

# Python 游戏引擎测试脚本
# 只测试 Python 部分，不涉及 Go 后端

BASE_URL="http://localhost:8081"

echo "=========================================="
echo "🎮 Python 游戏引擎测试"
echo "=========================================="
echo ""

# 检查服务是否运行
echo "📡 检查服务状态..."
if ! curl -s "$BASE_URL/health" > /dev/null 2>&1; then
    echo "❌ Python 引擎未运行！"
    echo "请先启动引擎："
    echo "  cd arena-backend"
    echo "  uvicorn engine.main:app --host 0.0.0.0 --port 8081 --reload"
    exit 1
fi

echo "✅ Python 引擎正在运行"
echo ""

# 测试 1: 健康检查
echo "=========================================="
echo "测试 1: 健康检查"
echo "=========================================="
curl -s "$BASE_URL/health" | python3 -m json.tool
echo ""
echo ""

# 测试 2: 初始化游戏
echo "=========================================="
echo "测试 2: 初始化游戏"
echo "=========================================="
RESPONSE=$(curl -s -X POST "$BASE_URL/internal/game/init" \
  -H "Content-Type: application/json" \
  -d '{"seed": 42}')

echo "$RESPONSE" | python3 -m json.tool
GAME_ID=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['game_id'])")
echo ""
echo "🎲 游戏ID: $GAME_ID"
echo ""

# 测试 3: 获取游戏状态
echo "=========================================="
echo "测试 3: 获取玩家 0 的游戏状态"
echo "=========================================="
curl -s "$BASE_URL/internal/game/$GAME_ID/state?player_id=0" | python3 -m json.tool
echo ""
echo ""

# 测试 4: 执行动作 (pass)
echo "=========================================="
echo "测试 4: 玩家 0 执行动作 (pass)"
echo "=========================================="
STEP_RESPONSE=$(curl -s -X POST "$BASE_URL/internal/game/step" \
  -H "Content-Type: application/json" \
  -d "{
    \"game_id\": \"$GAME_ID\",
    \"player_id\": 0,
    \"action\": \"pass\"
  }")

echo "$STEP_RESPONSE" | python3 -m json.tool
echo ""
echo ""

# 测试 5: 再次获取状态（应该轮到下一个玩家）
echo "=========================================="
echo "测试 5: 获取下一个玩家的状态"
echo "=========================================="
NEXT_PLAYER=$(echo "$STEP_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('next_player', 0))")
echo "下一个玩家: $NEXT_PLAYER"
curl -s "$BASE_URL/internal/game/$GAME_ID/state?player_id=$NEXT_PLAYER" | python3 -m json.tool
echo ""
echo ""

# 测试 6: 删除游戏
echo "=========================================="
echo "测试 6: 删除游戏"
echo "=========================================="
curl -s -X DELETE "$BASE_URL/internal/game/$GAME_ID" | python3 -m json.tool
echo ""
echo ""

# 测试 7: 验证游戏已删除
echo "=========================================="
echo "测试 7: 验证游戏已删除"
echo "=========================================="
echo "尝试获取已删除的游戏状态（应该返回 404）："
curl -s "$BASE_URL/internal/game/$GAME_ID/state?player_id=0" | python3 -m json.tool
echo ""
echo ""

echo "=========================================="
echo "✅ 所有测试完成！"
echo "=========================================="
