#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=========================================="
echo "🚀 启动 AI Card Arena 服务"
echo -e "==========================================${NC}"
echo ""

# 检查 Python 引擎健康状态
check_python_engine() {
    for i in {1..10}; do
        if curl -s http://localhost:8081/health > /dev/null 2>&1; then
            echo -e "${GREEN}✅ Python 引擎已就绪${NC}"
            return 0
        fi
        echo -e "${YELLOW}⏳ 等待 Python 引擎启动... ($i/10)${NC}"
        sleep 1
    done
    echo -e "${RED}❌ Python 引擎启动超时${NC}"
    return 1
}

# 清理函数
cleanup() {
    echo ""
    echo -e "${YELLOW}🛑 正在停止所有服务...${NC}"
    if [ ! -z "$PYTHON_PID" ]; then
        kill $PYTHON_PID 2>/dev/null
        echo -e "${GREEN}✅ Python 引擎已停止${NC}"
    fi
    if [ ! -z "$GO_PID" ]; then
        kill $GO_PID 2>/dev/null
        echo -e "${GREEN}✅ Go 后端已停止${NC}"
    fi
    exit 0
}

# 捕获中断信号
trap cleanup INT TERM

# 1. 启动 Python 游戏引擎
echo -e "${YELLOW}📦 启动 Python 游戏引擎...${NC}"
cd engine
uvicorn main:app --host 0.0.0.0 --port 8081 > ../logs/python_engine.log 2>&1 &
PYTHON_PID=$!
cd ..

# 等待 Python 引擎启动
if ! check_python_engine; then
    echo -e "${RED}❌ Python 引擎启动失败，请查看日志: logs/python_engine.log${NC}"
    cleanup
fi

echo ""

# 2. 启动 Go 后端服务
echo -e "${YELLOW}🔧 启动 Go 后端服务...${NC}"
go run main.go > logs/go_backend.log 2>&1 &
GO_PID=$!

# 等待 Go 后端启动
sleep 3

# 检查 Go 后端
if curl -s http://localhost:8080/api/games > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Go 后端已就绪${NC}"
else
    echo -e "${RED}❌ Go 后端启动失败，请查看日志: logs/go_backend.log${NC}"
    cleanup
fi

echo ""
echo -e "${GREEN}=========================================="
echo "✨ 所有服务已成功启动！"
echo "=========================================="
echo ""
echo "📍 服务地址:"
echo "   Python 引擎: http://localhost:8081"
echo "   Go 后端:    http://localhost:8080"
echo ""
echo "📝 日志文件:"
echo "   Python: logs/python_engine.log"
echo "   Go:     logs/go_backend.log"
echo ""
echo "🧪 快速测试:"
echo "   curl http://localhost:8081/health"
echo "   curl http://localhost:8080/api/games"
echo ""
echo "📚 查看文档:"
echo "   cat WORKFLOW.md"
echo "   cat DEPLOYMENT.md"
echo ""
echo -e "按 ${RED}Ctrl+C${NC} 停止所有服务"
echo -e "==========================================${NC}"

# 保持脚本运行
wait
