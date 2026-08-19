#!/bin/bash

if [ "$EUID" -ne 0 ]; then
  echo "❌ Error: Run this script with sudo: sudo $0"
  exit 1
fi

echo "🚀 Setting up Go Video Streamer environment..."

if ! command -v docker &> /dev/null; then
    echo "📦 Docker not found. Installing..."
    curl -fsSL https://docker.com | sh
fi

mkdir -p ./storage/videos

if [ ! "$(docker ps -q -f name=video_player_db)" ]; then
    if [ "$(docker ps -aq -f status=exited -f name=video_player_db)" ]; then
        docker start video_player_db
    else
        echo "💾 Starting PostgreSQL container..."
        docker run -d \
          --name video_player_db \
          -e POSTGRES_USER=user \
          -e POSTGRES_PASSWORD=password \
          -e POSTGRES_DB=videodb \
          -p 5433:5432 \
          -v pgdata:/var/lib/postgresql/data \
          --restart unless-stopped \
          postgres:15-alpine

        echo "⏳ Waiting for database initialization..."
        sleep 5

        docker exec -i video_player_db psql -U user -d videodb -c "
        CREATE TABLE IF NOT EXISTS videos (
            id SERIAL PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            filepath TEXT NOT NULL
        );
        INSERT INTO videos (title, filepath) VALUES ('First Video', './storage/videos/video1.mp4') ON CONFLICT DO NOTHING;
        "
    fi
fi

echo "🎬 Starting video server on port 8080..."
export DATABASE_URL="postgres://user:password@localhost:5433/videodb?sslmode=disable"

chmod +x ./video-server
./video-server
