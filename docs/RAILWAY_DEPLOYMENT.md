# Railway Deployment Guide for NeruBot Microservices

## Prerequisites

1. **Railway Account**
   - Sign up at [Railway.app](https://railway.app)
   - Install Railway CLI: `npm i -g @railway/cli` or `brew install railway`
   - Login: `railway login`

2. **GitHub Repository**
   - Push your code to GitHub
   - Connect GitHub account to Railway

3. **Required API Keys**
   - Discord Bot Token
   - OpenAI API Key (optional)
   - Anthropic API Key (optional)
   - Gemini API Key (optional)
   - Whale Alert API Key (optional)

---

## Deployment Steps

### Step 1: Create Railway Project

```bash
# Login to Railway
railway login

# Create a new project
railway init

# Link to your repository
railway link
```

Or use the Railway Dashboard:
1. Go to [Railway Dashboard](https://railway.app/dashboard)
2. Click "New Project"
3. Select "Deploy from GitHub repo"
4. Choose `nerufuyo/nerubot` repository

### Step 2: Setup PostgreSQL

1. In Railway Dashboard, click "New" → "Database" → "PostgreSQL"
2. Railway will automatically provision a PostgreSQL instance
3. Note the connection string (available in Variables tab)
4. Create multiple databases for each service:

```sql
-- Connect to PostgreSQL using Railway's psql or any client
CREATE DATABASE music_db;
CREATE DATABASE confession_db;
CREATE DATABASE roast_db;
CREATE DATABASE chat_db;
CREATE DATABASE news_db;
CREATE DATABASE whale_db;
```

### Step 3: Setup Redis

1. In Railway Dashboard, click "New" → "Database" → "Redis"
2. Railway will automatically provision a Redis instance
3. Note the connection string (available in Variables tab)

### Step 4: Deploy Services

Railway can auto-detect services from your repository structure. Here's how to deploy each service:

#### Option A: Automatic Deployment (Monorepo)

Railway will detect multiple services in the `services/` directory:

1. Go to Project Settings
2. Enable "Monorepo Support"
3. Railway will create a service for each directory with a `Dockerfile`

#### Option B: Manual Service Creation

For each service, create a new Railway service:

**API Gateway Service:**
```bash
railway service create api-gateway
railway up --service api-gateway
```

**Music Service:**
```bash
railway service create music-service
railway up --service music-service
```

**Confession Service:**
```bash
railway service create confession-service
railway up --service confession-service
```

**Roast Service:**
```bash
railway service create roast-service
railway up --service roast-service
```

**Chatbot Service:**
```bash
railway service create chatbot-service
railway up --service chatbot-service
```

**News Service:**
```bash
railway service create news-service
railway up --service news-service
```

**Whale Service:**
```bash
railway service create whale-service
railway up --service whale-service
```

### Step 5: Configure Environment Variables

#### Shared Variables (Set in Project Settings)

```env
# Discord Configuration
DISCORD_TOKEN=your_discord_bot_token
DISCORD_GUILD_ID=your_guild_id

# Environment
ENVIRONMENT=production
LOG_LEVEL=INFO

# AI Providers (Optional)
OPENAI_API_KEY=your_openai_key
ANTHROPIC_API_KEY=your_anthropic_key
GEMINI_API_KEY=your_gemini_key

# External Services (Optional)
WHALE_ALERT_API_KEY=your_whale_alert_key
```

#### Service-Specific Variables

**API Gateway:**
```env
PORT=8080
MUSIC_SERVICE_URL=${{music-service.RAILWAY_PRIVATE_DOMAIN}}:8081
CONFESSION_SERVICE_URL=${{confession-service.RAILWAY_PRIVATE_DOMAIN}}:8082
ROAST_SERVICE_URL=${{roast-service.RAILWAY_PRIVATE_DOMAIN}}:8083
CHATBOT_SERVICE_URL=${{chatbot-service.RAILWAY_PRIVATE_DOMAIN}}:8084
NEWS_SERVICE_URL=${{news-service.RAILWAY_PRIVATE_DOMAIN}}:8085
WHALE_SERVICE_URL=${{whale-service.RAILWAY_PRIVATE_DOMAIN}}:8086
REDIS_URL=${{Redis.REDIS_URL}}
```

**Music Service:**
```env
PORT=8081
DATABASE_URL=${{Postgres.DATABASE_URL}}/music_db
REDIS_URL=${{Redis.REDIS_URL}}
FFMPEG_PATH=/usr/bin/ffmpeg
YTDLP_PATH=/usr/local/bin/yt-dlp
MAX_QUEUE_SIZE=100
```

**Confession Service:**
```env
PORT=8082
DATABASE_URL=${{Postgres.DATABASE_URL}}/confession_db
STORAGE_PATH=/data/confessions
MAX_IMAGE_SIZE=10485760
```

**Roast Service:**
```env
PORT=8083
DATABASE_URL=${{Postgres.DATABASE_URL}}/roast_db
REDIS_URL=${{Redis.REDIS_URL}}
ROAST_COOLDOWN_SECONDS=300
```

**Chatbot Service:**
```env
PORT=8084
DATABASE_URL=${{Postgres.DATABASE_URL}}/chat_db
REDIS_URL=${{Redis.REDIS_URL}}
SESSION_TIMEOUT_MINUTES=30
MAX_CONTEXT_MESSAGES=10
```

**News Service:**
```env
PORT=8085
DATABASE_URL=${{Postgres.DATABASE_URL}}/news_db
REDIS_URL=${{Redis.REDIS_URL}}
FETCH_INTERVAL_MINUTES=15
MAX_ARTICLES_PER_FETCH=50
```

**Whale Service:**
```env
PORT=8086
DATABASE_URL=${{Postgres.DATABASE_URL}}/whale_db
CHECK_INTERVAL_SECONDS=60
DEFAULT_MIN_AMOUNT_USD=1000000
```

### Step 6: Configure Networking

Railway automatically provides:
- **Private Networking:** Services communicate via `servicename.railway.internal`
- **Public Domain:** Only API Gateway needs a public domain
- **HTTPS:** Automatically provisioned SSL certificates

**Enable Public Domain for API Gateway:**
1. Go to API Gateway service settings
2. Click "Settings" → "Networking"
3. Click "Generate Domain"
4. Copy the generated domain (e.g., `api-gateway-production-xxxx.up.railway.app`)

**Internal Service URLs:**
All services communicate internally without exposing public endpoints.

### Step 7: Setup Health Checks

Railway automatically uses the `healthcheckPath` defined in `railway.toml`:

```toml
[deploy]
healthcheckPath = "/health"
healthcheckTimeout = 100
```

Ensure each service implements the `/health` endpoint:

```go
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "healthy",
        "service": "service-name",
        "timestamp": time.Now().Format(time.RFC3339),
    })
}
```

### Step 8: Setup Monitoring

**Built-in Railway Monitoring:**
1. View service metrics in Railway Dashboard
2. Monitor CPU, Memory, Network usage
3. View logs in real-time
4. Set up deployment notifications

**Custom Monitoring (Optional):**
```bash
# Add Prometheus metrics endpoint
# Install Grafana for dashboards
# Setup alerting via webhook
```

### Step 9: Configure Volumes (Persistent Storage)

For services that need persistent file storage:

**Confession Service (for images):**
```bash
railway volume create confession-data
railway volume mount /data/confessions --service confession-service
```

**Roast Service (for data backup):**
```bash
railway volume create roast-data
railway volume mount /data/roasts --service roast-service
```

### Step 10: Setup CI/CD

Railway automatically deploys on git push. Configure deployment triggers:

**In Railway Dashboard:**
1. Go to Service Settings
2. Click "Deployments"
3. Configure:
   - Deploy on: `main` branch (or specify branches)
   - Auto-deploy: Enabled
   - Deploy on PR: Optional

**GitHub Actions (Alternative):**
Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy to Railway

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install Railway
        run: npm i -g @railway/cli
      
      - name: Deploy to Railway
        run: railway up --service ${{ matrix.service }}
        env:
          RAILWAY_TOKEN: ${{ secrets.RAILWAY_TOKEN }}
        strategy:
          matrix:
            service:
              - api-gateway
              - music-service
              - confession-service
              - roast-service
              - chatbot-service
              - news-service
              - whale-service
```

---

## Verification Checklist

After deployment, verify each component:

### ✅ Infrastructure
- [ ] PostgreSQL instance running and accessible
- [ ] Redis instance running and accessible
- [ ] All databases created
- [ ] Volumes mounted correctly

### ✅ Services
- [ ] API Gateway deployed and healthy
- [ ] Music Service deployed and healthy
- [ ] Confession Service deployed and healthy
- [ ] Roast Service deployed and healthy
- [ ] Chatbot Service deployed and healthy
- [ ] News Service deployed and healthy
- [ ] Whale Service deployed and healthy

### ✅ Configuration
- [ ] All environment variables set
- [ ] Discord token configured
- [ ] Service URLs configured correctly
- [ ] Database connections working
- [ ] Redis connections working

### ✅ Networking
- [ ] API Gateway has public domain
- [ ] Internal services accessible via private network
- [ ] Health checks passing
- [ ] No port conflicts

### ✅ Functionality
- [ ] Discord bot connects successfully
- [ ] Commands route to correct services
- [ ] Music playback working
- [ ] Confession submission working
- [ ] Roast generation working
- [ ] All features operational

---

## Monitoring and Maintenance

### View Logs

**Via Railway Dashboard:**
1. Select service
2. Click "Logs" tab
3. Filter by time, level, or search

**Via CLI:**
```bash
# View logs for specific service
railway logs --service api-gateway

# Follow logs in real-time
railway logs --service music-service --tail
```

### View Metrics

**Via Railway Dashboard:**
1. Select service
2. Click "Metrics" tab
3. View CPU, Memory, Network graphs

### Restart Service

```bash
# Restart specific service
railway restart --service api-gateway

# Restart all services
railway restart
```

### Update Environment Variables

```bash
# Set variable for specific service
railway variables set KEY=VALUE --service api-gateway

# Bulk update from .env file
railway variables set --service music-service < .env
```

### Database Backups

Railway automatically backs up PostgreSQL:
- Daily backups retained for 7 days
- Manual backups via Dashboard

**Manual Backup:**
```bash
# Export database
railway run pg_dump $DATABASE_URL > backup.sql

# Restore database
railway run psql $DATABASE_URL < backup.sql
```

---

## Troubleshooting

### Service Won't Start

1. Check logs: `railway logs --service service-name`
2. Verify environment variables
3. Check database connection
4. Verify Dockerfile builds locally

### High Memory Usage

1. Check metrics in Dashboard
2. Optimize Go code (reduce allocations)
3. Increase memory limit in Railway settings
4. Add memory limits to Dockerfile

### Slow Response Times

1. Check service latency in metrics
2. Verify network connectivity between services
3. Add Redis caching
4. Optimize database queries
5. Enable connection pooling

### Database Connection Issues

1. Verify DATABASE_URL is correct
2. Check database is running
3. Verify connection limits
4. Check network connectivity
5. Review connection pooling settings

### Deployment Failures

1. Check build logs
2. Verify Dockerfile syntax
3. Ensure all dependencies in go.mod
4. Check for port conflicts
5. Verify health check endpoint

---

## Cost Optimization

### Free Tier Limits
- $5/month credit
- 500 GB egress
- 512 MB memory per service

### Optimization Tips
1. **Use single PostgreSQL instance** with multiple databases
2. **Use single Redis instance** shared across services
3. **Optimize Docker images** - use multi-stage builds
4. **Enable sleep mode** for non-critical services during low traffic
5. **Monitor usage** regularly in Dashboard

### Estimated Costs

**Hobby Plan ($5/month):**
- Small Discord bot: Free tier sufficient
- Medium Discord bot: $5-10/month
- Large Discord bot: $20-30/month

**Services:**
- PostgreSQL: ~$5/month
- Redis: ~$5/month
- Each service: ~$5/month (512 MB RAM)

---

## Security Best Practices

1. **Never commit secrets** - Use Railway environment variables
2. **Use Railway's built-in secrets** - Encrypted at rest
3. **Enable private networking** - Services not exposed publicly
4. **Use HTTPS** - Automatically provided by Railway
5. **Rotate API keys** regularly
6. **Implement rate limiting** in API Gateway
7. **Add authentication** for admin endpoints
8. **Regular security updates** - Keep dependencies updated

---

## Rollback Procedure

If deployment fails or causes issues:

```bash
# Via CLI
railway rollback --service api-gateway

# Via Dashboard
# 1. Go to service
# 2. Click "Deployments"
# 3. Find previous successful deployment
# 4. Click "Redeploy"
```

---

## Support and Resources

- **Railway Documentation:** https://docs.railway.app
- **Railway Discord:** https://discord.gg/railway
- **Railway Status:** https://status.railway.app
- **NeruBot Issues:** https://github.com/nerufuyo/nerubot/issues

---

## Next Steps

After successful deployment:

1. ✅ Test all features thoroughly
2. ✅ Monitor logs and metrics for 24 hours
3. ✅ Setup alerting for critical issues
4. ✅ Document any production-specific configuration
5. ✅ Create runbook for common issues
6. ✅ Schedule regular maintenance windows
7. ✅ Plan for scaling as usage grows

---

**Last Updated:** December 5, 2025  
**Author:** @nerufuyo  
**Version:** 1.0.0
