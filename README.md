# Go Video Streamer

A minimalist service for video streaming with built-in rewind support (Range Requests). The backend is powered by Go, metadata is managed in a Docker-hosted PostgreSQL database, and video files are streamed directly from local storage.

## Features
* **Range Requests:** Full HTTP 206 Partial Content support. You can fast-forward or rewind to any point instantly without downloading the entire video file first.
* **Dockerized Database:** Spin up a clean PostgreSQL instance via Docker Compose without cluttering your host machine.
* **Native Go Implementation:** The server leverages Go's built-in `net/http` package paired with the high-performance `pgx` driver.

---

## Project Structure

```text
├── main.go               # Main Go server application logic
├── index.html            # Frontend HTML5 video player UI
├── docker-compose.yml    # PostgreSQL container configuration
├── go.mod                # Go module dependencies
├── static/
│   └── css/
│       └── style.css     # Custom player styles
└── storage/
    └── videos/
        └── video1.mp4    # Local directory for video files (gitignored)
```

---

## Prerequisites
* **Go** (version 1.20 or higher)
* **Docker** & **Docker Compose** installed on your system

---

## Quick Start Guide

### 1. Start the Database
Spin up the PostgreSQL container in detached mode:
```bash
docker-compose up -d
```
*Note: The database is mapped to local port `5433` to prevent port collision with any existing local PostgreSQL instances.*

### 2. Initialize the Database Schema
Access the PostgreSQL command-line utility inside your active container:
```bash
docker exec -it video_player_db psql -U user -d videodb
```

Execute the following SQL commands to build the table structure and insert the metadata for your first video track:
```sql
CREATE TABLE videos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    filepath TEXT NOT NULL
);

INSERT INTO videos (title, filepath) VALUES ('First Video', './storage/videos/video1.mp4');
```
Type `\q` and press Enter to exit the PostgreSQL terminal interface.

### 3. Place Your Video File
Ensure you have a valid MP4 video placed exactly under this relative directory path:
`./storage/videos/video1.mp4`

### 4. Run the Go Server
Fetch the required database driver dependencies and spin up your backend service:
```bash
go get ://github.com
go run main.go
```

Once running, navigate to the web player in your browser: [http://localhost:8080](http://localhost:8080)

---
