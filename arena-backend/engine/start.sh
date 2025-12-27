#!/bin/bash
# Python 游戏引擎启动脚本

cd "$(dirname "$0")"

echo "🚀 启动 Python 游戏引擎..."
uvicorn main:app --host 0.0.0.0 --port 8081 --reload
