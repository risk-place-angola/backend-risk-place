#!/bin/bash

set -e

# Get version from argument or default to latest
VERSION=${1:-latest}
REPLICAS=${2:-3}  # Default: 3 replicas

echo "🚀 Starting deployment"
echo "  Version: $VERSION"
echo "  Replicas: $REPLICAS"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Docker is running${NC}"

# Determine docker compose command
COMPOSE_CMD="docker compose"
if ! docker compose version &> /dev/null 2>&1; then
    COMPOSE_CMD="docker-compose"
fi

echo -e "${BLUE}ℹ️  Using: $COMPOSE_CMD${NC}"

# Check if kilua_network exists
if ! docker network ls | grep -q kilua_network; then
    echo -e "${YELLOW}⚠️  Network kilua_network not found. Creating...${NC}"
    docker network create --driver overlay --attachable kilua_network
    echo -e "${GREEN}✅ Network created${NC}"
else
    echo -e "${GREEN}✅ Network kilua_network exists${NC}"
fi

# Pull specific version
echo -e "${YELLOW}📥 Pulling image version: $VERSION...${NC}"
docker pull riskplaceangola/backend-core:$VERSION

# Tag as latest locally for docker-compose
echo -e "${YELLOW}🏷️  Tagging as latest locally...${NC}"
docker tag riskplaceangola/backend-core:$VERSION riskplaceangola/backend-core:latest

# Stop old containers
echo -e "${YELLOW}🛑 Stopping old containers...${NC}"
$COMPOSE_CMD -f docker-compose.prod.yml down 2>/dev/null || true

$COMPOSE_CMD -f docker-compose.prod.yml up -d --scale backend_core=$REPLICAS --remove-orphans

# Wait for containers to be healthy
echo -e "${YELLOW}⏳ Waiting for containers to be healthy...${NC}"
TIMEOUT=120
ELAPSED=0
HEALTHY_COUNT=0

while [ $ELAPSED -lt $TIMEOUT ]; do
    # Count healthy containers (matches backend-risk-place-backend_core-1, etc)
    HEALTHY_COUNT=$(docker ps --filter "name=backend_core" --filter "health=healthy" --format "{{.Names}}" 2>/dev/null | wc -l | tr -d ' ')
    TOTAL_COUNT=$(docker ps --filter "name=backend_core" --format "{{.Names}}" 2>/dev/null | wc -l | tr -d ' ')
    
    if [ "$HEALTHY_COUNT" -eq "$REPLICAS" ]; then
        echo -e "${GREEN}✅ All $REPLICAS instances are healthy!${NC}"
        break
    fi
    
    echo -e "${YELLOW}⏳ Healthy: $HEALTHY_COUNT/$REPLICAS running ($TOTAL_COUNT started) - ${ELAPSED}s/${TIMEOUT}s${NC}"
    
    # Show which containers are not healthy yet
    if [ $ELAPSED -gt 30 ] && [ "$HEALTHY_COUNT" -lt "$REPLICAS" ]; then
        docker ps --filter "name=backend_core" --format "table {{.Names}}\t{{.Status}}" 2>/dev/null || true
    fi
    
    sleep 5
    ELAPSED=$((ELAPSED + 5))
done

# Check if we timed out
if [ "$HEALTHY_COUNT" -lt "$REPLICAS" ]; then
    echo -e "${RED}❌ Deployment failed: Only $HEALTHY_COUNT/$REPLICAS instances became healthy${NC}"
    echo -e "${YELLOW}Container status:${NC}"
    docker ps -a --filter "name=backend_core" --format "table {{.Names}}\t{{.Status}}" || true
    echo -e "\n${YELLOW}Recent logs from unhealthy containers:${NC}"
    docker ps --filter "name=backend_core" --format "{{.Names}}" | while read container; do
        HEALTH=$(docker inspect --format='{{.State.Health.Status}}' "$container" 2>/dev/null || echo "no-health")
        if [ "$HEALTH" != "healthy" ]; then
            echo -e "\n${YELLOW}Logs from $container:${NC}"
            docker logs --tail 30 "$container" 2>&1 || true
        fi
    done
    exit 1
fi

# Show deployment info
echo -e "\n${GREEN}✅ Deployment successful!${NC}"
echo -e "${GREEN}📦 Deployed version: $VERSION${NC}"
echo -e "${GREEN}🔢 Instances: $REPLICAS${NC}"

echo -e "\n${YELLOW}📊 Container status:${NC}"
docker ps --filter "name=backend_core" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" 2>/dev/null || docker ps | grep backend_core

echo -e "\n${YELLOW}📈 Resource usage:${NC}"
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}" $(docker ps --filter "name=backend_core" --format "{{.Names}}" | tr '\n' ' ') 2>/dev/null || echo "Stats not available"

echo -e "\n${YELLOW}Recent logs from first instance:${NC}"
FIRST_CONTAINER=$(docker ps --filter "name=backend_core" --format "{{.Names}}" | head -1)
if [ -n "$FIRST_CONTAINER" ]; then
    docker logs --tail 20 "$FIRST_CONTAINER" 2>&1 || true
fi

# Update Nginx upstream configuration if Nginx is running
echo -e "\n${YELLOW}🔄 Checking for Nginx load balancer...${NC}"
if docker ps | grep -q "nginx"; then
    echo -e "${GREEN}✅ Nginx found - updating upstream configuration${NC}"
    
    # Update Nginx to load balance across all replicas
    NGINX_CONTAINER=$(docker ps --filter "name=nginx" --format "{{.Names}}" | head -1)
    
    # Test if backends are reachable from Nginx
    REACHABLE=0
    for i in $(seq 1 $REPLICAS); do
        if docker exec $NGINX_CONTAINER wget --spider --tries=1 --timeout=2 http://backend_core-$i:8090/health 2>&1 | grep -q "200 OK"; then
            REACHABLE=$((REACHABLE + 1))
        fi
    done
    
    if [ $REACHABLE -eq $REPLICAS ]; then
        echo -e "${GREEN}✅ All $REPLICAS backends are reachable from Nginx${NC}"
        echo -e "${BLUE}ℹ️  Manual step: Update Nginx upstream to include all replicas${NC}"
        echo -e "${BLUE}   See: backend-config/NGINX-LOAD-BALANCING.md${NC}"
    else
        echo -e "${YELLOW}⚠️  Only $REACHABLE/$REPLICAS backends reachable from Nginx${NC}"
        echo -e "${YELLOW}   Wait for all backends to be healthy, then update Nginx config${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Nginx not found in this environment${NC}"
    echo -e "${BLUE}ℹ️  If using external Nginx, update upstream configuration manually${NC}"
fi

# Show deployed image info
echo -e "\n${YELLOW}🏷️  Image info:${NC}"
docker images | grep -E "backend-core.*(latest|$VERSION)" | head -3 || true

# Cleanup old images (keep last 3 versions)
echo -e "\n${YELLOW}🧹 Cleaning up old images...${NC}"
docker images | grep backend-core | tail -n +4 | awk '{print $3}' | xargs docker rmi -f 2>/dev/null || true
echo -e "${GREEN}✅ Cleanup complete${NC}"

# Security: Remove .env file after container is running
echo -e "\n${YELLOW}🔒 Securing environment file...${NC}"
if [ -f ".env" ]; then
    # Shred the file (overwrite with random data before deleting)
    if command -v shred &> /dev/null; then
        shred -vfz -n 3 .env
        echo -e "${GREEN}✅ .env file securely deleted (shredded)${NC}"
    else
        # Fallback: overwrite and delete
        dd if=/dev/urandom of=.env bs=1k count=10 2>/dev/null
        rm -f .env
        echo -e "${GREEN}✅ .env file securely deleted${NC}"
    fi
else
    echo -e "${YELLOW}ℹ️  No .env file found (already cleaned)${NC}"
fi

echo -e "\n${GREEN}✅ Deployment complete!${NC}"
