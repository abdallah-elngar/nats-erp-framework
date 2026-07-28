#!/bin/bash

# ============================================
# نشر NATS Framework
# ============================================

set -e

echo "🚀 Deploying NATS Framework..."

# المتغيرات
ENV=${1:-production}
TAG=${2:-latest}

# بناء صورة Docker
echo "🐳 Building Docker image..."
docker build -t nats-framework:$TAG .

# دفع الصورة (اختياري)
if [ "$PUSH" = "true" ]; then
    echo "📤 Pushing Docker image..."
    docker push nats-framework:$TAG
fi

# نشر باستخدام Docker Compose
if [ "$ENV" = "production" ]; then
    echo "🌐 Deploying to production..."
    docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
else
    echo "🧪 Deploying to $ENV environment..."
    docker-compose -f docker-compose.yml up -d
fi

echo "✅ Deployment complete!"
echo "🌐 Application is running at http://localhost:8080"