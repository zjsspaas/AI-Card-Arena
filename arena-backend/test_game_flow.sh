#!/bin/bash
# 游戏流程测试脚本

BASE_URL="http://localhost:8080"

echo "=========================================="
echo "🎮 AI Card Arena 游戏流程测试"
echo "=========================================="
echo ""

# 1. 注册三个测试用户
echo "📝 步骤 1: 注册三个测试玩家..."
echo ""

USER1=$(curl -s -X POST "$BASE_URL/api/users/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "player1",
    "password": "123456",
    "role": "developer"
  }')
echo "玩家1: $USER1"

USER2=$(curl -s -X POST "$BASE_URL/api/users/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "player2",
    "password": "123456",
    "role": "developer"
  }')
echo "玩家2: $USER2"

USER3=$(curl -s -X POST "$BASE_URL/api/users/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "player3",
    "password": "123456",
    "role": "developer"
  }')
echo "玩家3: $USER3"
echo ""

# 2. 创建游戏
echo "🎲 步骤 2: 创建游戏..."
echo ""

GAME=$(curl -s -X POST "$BASE_URL/api/games/create" \
  -H "Content-Type: application/json" \
  -d '{
    "player_ids": [1, 2, 3],
    "seed": 42
  }')

echo "$GAME" | jq '.'
GAME_ID=$(echo "$GAME" | jq -r '.game_id')
echo ""
echo "游戏ID: $GAME_ID"
echo ""

# 3. 获取初始状态
echo "📊 步骤 3: 获取玩家0的游戏状态..."
echo ""

STATE=$(curl -s "$BASE_URL/api/games/$GAME_ID/state?player_id=0")
echo "$STATE" | jq '.'
echo ""

# 4. 执行第一步动作（假设玩家0先手）
echo "🃏 步骤 4: 玩家0出牌 (pass)..."
echo ""

ACTION1=$(curl -s -X POST "$BASE_URL/api/games/$GAME_ID/action" \
  -H "Content-Type: application/json" \
  -d '{
    "player_id": 0,
    "action": "pass"
  }')

echo "$ACTION1" | jq '.'
echo ""

# 5. 查看游戏会话
echo "📋 步骤 5: 查看游戏会话信息..."
echo ""

SESSION=$(curl -s "$BASE_URL/api/games/$GAME_ID")
echo "$SESSION" | jq '.'
echo ""

echo "=========================================="
echo "✅ 测试完成！"
echo "=========================================="
