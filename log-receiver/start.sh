#!/bin/bash

echo "Starting Log Receiver Stack..."

docker-compose up -d

echo "Stack started. Waiting for services to be ready..."

sleep 5

echo "Services started:"
echo "- Redis: localhost:6379"
echo "- Go Consumer (gRPC): localhost:50051"
echo "- Node API: http://localhost:3000"
echo "- C++ Log Receiver: UDP localhost:514, TCP localhost:515"

echo ""
echo "To send test logs:"
echo "  UDP: echo '2024-01-01T00:00:00Z ERROR test message' | nc -u localhost 514"
echo "  TCP: echo '2024-01-01T00:00:00Z INFO test message' | nc localhost 515"

echo ""
echo "API endpoints:"
echo "  GET  /api/metrics - Get aggregated metrics"
echo "  GET  /api/alerts - Get alert rules"
echo "  POST /api/alerts - Create alert rule"
echo "  GET  /api/logs - Get recent logs"
echo "  GET  /api/logs/search?q=query - Search logs"