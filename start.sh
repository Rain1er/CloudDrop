#!/bin/bash

# 检查并终止已存在的进程
echo "Checking for existing processes..."
pkill -f clouddrop || true
lsof -ti:8989 | xargs kill -9 2>/dev/null || true

# 生成随机 JWT 密钥
echo "Generating random JWT secret..."
JWT_SECRET=$(openssl rand -base64 32)
if [ $? -ne 0 ]; then
    echo "Failed to generate JWT secret!"
    exit 1
fi
echo "JWT secret generated successfully"

# 启动后端服务
echo "Starting backend server..."
cd backend
export JWT_SECRET="$JWT_SECRET"
export REACT_APP_API_URL=http://localhost:8989/api/v1
go build -o clouddrop ./cmd/main.go
if [ $? -ne 0 ]; then
    echo "Backend build failed!"
    exit 1
fi

./clouddrop &
backend_pid=$!
echo "Backend started with PID: $backend_pid"

# 等待后端启动并测试连接
echo "Waiting for backend to start..."
for i in {1..30}; do
    if curl -s http://localhost:8989/health > /dev/null; then
        echo "Backend is ready!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "Backend failed to start!"
        kill $backend_pid
        exit 1
    fi
    sleep 1
done
