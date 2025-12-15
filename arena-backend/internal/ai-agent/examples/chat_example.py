#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Chat Agent 对话示例 - Python 版本
"""

import requests
import json

BASE_URL = "http://localhost:8080"
AGENT_NAME = "chat"

def chat_with_agent(message, session_id="default"):
    """与 Agent 对话"""
    url = f"{BASE_URL}/api/v1/agents/{AGENT_NAME}/process"
    payload = {
        "message": message,
        "params": {
            "session_id": session_id
        }
    }
    
    response = requests.post(url, json=payload)
    return response.json()

def main():
    print("=== Chat Agent 对话示例 ===\n")
    
    # 使用同一个 session_id 进行多轮对话
    session_id = "python_user_001"
    
    # 第一轮对话
    print("用户: 你好")
    result = chat_with_agent("你好", session_id)
    print(f"Agent: {result['data']['result']}\n")
    
    # 第二轮对话
    print("用户: 你叫什么名字？")
    result = chat_with_agent("你叫什么名字？", session_id)
    print(f"Agent: {result['data']['result']}\n")
    
    # 第三轮对话
    print("用户: 你能帮我做什么？")
    result = chat_with_agent("你能帮我做什么？", session_id)
    print(f"Agent: {result['data']['result']}\n")
    
    # 显示对话历史
    print("=== 对话历史 ===")
    history = result['data']['history']
    for msg in history:
        role = "用户" if msg['role'] == 'user' else "Agent"
        print(f"{role}: {msg['content']}")

if __name__ == "__main__":
    main()


