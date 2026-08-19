# 🎬 Go Video Streamer

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/sudowhiterose/go-video-streaming?display_name=tag)](https://github.com/sudowhiterose/go-video-streaming/releases)

A minimalist service for video streaming with built-in rewind support (Range Requests). The backend is powered by Go, metadata is managed in a Docker-hosted PostgreSQL database, and video files are streamed directly from local storage.

---

## 🚀 Features
* **Range Requests:** Full HTTP 206 Partial Content support. You can fast-forward or rewind to any point instantly without downloading the entire video file first.
* **Dockerized Database:** Spin up a clean PostgreSQL instance via Docker Compose without cluttering your host machine.
* **Native Go Implementation:** The server leverages Go's built-in `net/http` package paired with the high-performance `pgx` driver.

---

## ⏱️ Quick Start Guide

