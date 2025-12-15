#!/bin/bash

# Chat Agent 对话示例脚本

BASE_URL="http://localhost:8080"
AGENT_NAME="chat"

echo "=== Chat Agent 对话示例 ==="
echo ""

# 1. 列出所有 Agent
echo "1. 列出所有可用的 Agent:"
curl -s "$BASE_URL/api/v1/agents" | jq '.'
echo ""
echo ""

# 2. 开始第一轮对话
echo "2. 第一轮对话 - 打招呼:"
curl -s -X POST "$BASE_URL/api/v1/agents/$AGENT_NAME/process" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "你好",
    "params": {
      "session_id": "user123"
    }
  }' | jq '.'
echo ""
echo ""

# 3. 继续对话（使用相同的 session_id）
echo "3. 第二轮对话 - 询问名字:"
curl -s -X POST "$BASE_URL/api/v1/agents/$AGENT_NAME/process" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "你叫什么名字？",
    "params": {
      "session_id": "user123"
    }
  }' | jq '.'
echo ""
echo ""

# 4. 第三轮对话
echo "4. 第三轮对话 - 询问帮助:"
curl -s -X POST "$BASE_URL/api/v1/agents/$AGENT_NAME/process" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "你能帮我做什么？",
    "params": {
      "session_id": "user123"
    }
  }' | jq '.'
echo ""
echo ""

# 5. 查看对话历史（最后一次请求会返回完整历史）
echo "5. 对话历史已保存在最后一次响应的 data.history 字段中"
echo ""


