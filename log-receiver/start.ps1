Write-Host "Starting Log Receiver Stack..."

docker-compose up -d

Write-Host "Stack started. Waiting for services to be ready..."

Start-Sleep -Seconds 5

Write-Host "Services started:"
Write-Host "- Redis: localhost:6379"
Write-Host "- Go Consumer (gRPC): localhost:50051"
Write-Host "- Node API: http://localhost:3000"
Write-Host "- C++ Log Receiver: UDP localhost:514, TCP localhost:515"

Write-Host ""
Write-Host "To send test logs:"
Write-Host "  UDP: Use PowerShell or nc to send to localhost:514"
Write-Host "  TCP: Use PowerShell or nc to send to localhost:515"

Write-Host ""
Write-Host "API endpoints:"
Write-Host "  GET  /api/metrics - Get aggregated metrics"
Write-Host "  GET  /api/alerts - Get alert rules"
Write-Host "  POST /api/alerts - Create alert rule"
Write-Host "  GET  /api/logs - Get recent logs"
Write-Host "  GET  /api/logs/search?q=query - Search logs"