# NeruBot Test Automation Script
Write-Host "`n=== NeruBot Service Testing Suite ===`n" -ForegroundColor Cyan

# Check if .env exists
if (Test-Path .env) {
    Write-Host "OK .env file found" -ForegroundColor Green
} else {
    Write-Host "ERROR .env file missing - creating from template" -ForegroundColor Red
    Copy-Item .env.example .env
    Write-Host "WARNING Please edit .env and add your API keys" -ForegroundColor Yellow
    exit 1
}

# Check for Discord token
$envContent = Get-Content .env -Raw
if ($envContent -match 'DISCORD_TOKEN=(?!your_discord_bot_token_here)(.+)') {
    Write-Host "OK Discord token configured" -ForegroundColor Green
} else {
    Write-Host "ERROR Discord token not configured in .env" -ForegroundColor Red
    exit 1
}

Write-Host "`n=== Building Services ===`n" -ForegroundColor Cyan

# Build main bot
Write-Host "Building main bot..." -ForegroundColor Yellow
go build -o build/bot.exe ./cmd/nerubot 2>&1 | Out-Null
if ($LASTEXITCODE -eq 0) {
    $botSize = [math]::Round((Get-Item build/bot.exe).Length/1MB, 2)
    Write-Host "OK Main bot built ($botSize MB)" -ForegroundColor Green
} else {
    Write-Host "ERROR Main bot build failed" -ForegroundColor Red
    exit 1
}

# Build microservices
$services = @("gateway", "music", "confession", "roast", "chatbot", "news", "whale")
Write-Host "Building microservices..." -ForegroundColor Yellow
foreach ($service in $services) {
    Write-Host "  $service..." -NoNewline
    go build -o "build/$service/$service.exe" "./services/$service/cmd" 2>&1 | Out-Null
    if ($LASTEXITCODE -eq 0) {
        $size = [math]::Round((Get-Item "build/$service/$service.exe").Length/1MB, 2)
        Write-Host " OK ($size MB)" -ForegroundColor Green
    } else {
        Write-Host " FAIL" -ForegroundColor Red
    }
}

Write-Host "`n=== Choose Test Mode ===`n" -ForegroundColor Cyan
Write-Host "1. Run main bot (easiest)" -ForegroundColor White
Write-Host "2. Test health endpoints" -ForegroundColor White
Write-Host "3. Exit" -ForegroundColor White
$choice = Read-Host "`nEnter choice (1-3)"

if ($choice -eq "1") {
    Write-Host "`nStarting main bot... Press Ctrl+C to stop`n" -ForegroundColor Yellow
    .\build\bot.exe
}
elseif ($choice -eq "2") {
    Write-Host "`nStarting services for health check...`n" -ForegroundColor Yellow
    
    $jobs = @()
    $services = @{
        "gateway" = 8080
        "music" = 8081
        "confession" = 8082
        "roast" = 8083
        "chatbot" = 8084
        "news" = 8085
        "whale" = 8086
    }
    
    foreach ($svc in $services.Keys) {
        $job = Start-Job -ScriptBlock {
            param($path)
            Set-Location $using:PWD
            & $path
        } -ArgumentList ".\build\$svc\$svc.exe"
        $jobs += $job
    }
    
    Start-Sleep -Seconds 8
    Write-Host "Testing health endpoints:`n" -ForegroundColor Cyan
    
    foreach ($svc in $services.Keys) {
        $port = $services[$svc]
        Write-Host "$svc (port $port)..." -NoNewline
        try {
            $response = Invoke-WebRequest -Uri "http://localhost:$port/health" -TimeoutSec 3 -ErrorAction Stop
            Write-Host " OK" -ForegroundColor Green
        } catch {
            Write-Host " FAIL" -ForegroundColor Red
        }
    }
    
    Write-Host "`nStopping services..." -ForegroundColor Yellow
    $jobs | Stop-Job
    $jobs | Remove-Job
    Write-Host "Done!`n" -ForegroundColor Green
}
else {
    Write-Host "Exiting..." -ForegroundColor Yellow
}

Write-Host "`nFor detailed testing, see TESTING_GUIDE.md`n" -ForegroundColor Cyan
