<div align="center">

<!-- Banner -->
<img width="100%" src="https://capsule-render.vercel.app/api?type=waving&color=0:1a1a2e,50:16213e,100:0f3460&height=200&section=header&text=👟%20SoleStore&fontSize=60&fontColor=e94560&animation=fadeIn&fontAlignY=38&desc=Modern%20Shoe%20E-Commerce%20API&descAlignY=55&descSize=18&descColor=a8b2d8"/>

<br/>

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://docker.com)
[![Redis](https://img.shields.io/badge/Redis-Cache-DC382D?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io)
[![MySQL](https://img.shields.io/badge/MySQL-Database-4479A1?style=for-the-badge&logo=mysql&logoColor=white)](https://mysql.com)
[![License](https://img.shields.io/badge/License-MIT-e94560?style=for-the-badge)](LICENSE)

<br/>

> **SoleStore** adalah backend REST API untuk platform e-commerce sepatu yang dibangun dengan **Go (Golang)** — performant, scalable, dan containerized.

<br/>

</div>

---

## 📖 Table of Contents

- [✨ Features](#-features)
- [🏗️ Architecture](#️-architecture)
- [🚀 Quick Start](#-quick-start)
- [⚙️ Environment Variables](#️-environment-variables)
- [📡 API Endpoints](#-api-endpoints)
- [🗄️ Database Schema](#️-database-schema)
- [🐳 Docker Setup](#-docker-setup)
- [📁 Project Structure](#-project-structure)
- [🧪 Running Tests](#-running-tests)
- [📬 Contact](#-contact)

---

## ✨ Features

| Feature | Description |
|---|---|
| 🔐 **Auth & JWT** | Register, login, refresh token, dan middleware autentikasi berbasis JWT |
| 👟 **Product Catalog** | CRUD produk sepatu dengan filter brand, ukuran, warna, dan harga |
| 🛒 **Shopping Cart** | Tambah, update, dan hapus item dari keranjang belanja |
| 💳 **Payment Gateway** | Integrasi payment dengan penanganan webhook dan status transaksi |
| 📦 **Order Tracking** | Lacak status pesanan dari pending hingga delivered |
| ⚡ **Redis Cache** | Caching product catalog dan session untuk performa optimal |

---

## 🏗️ Architecture

```
┌─────────────┐     HTTP/REST     ┌──────────────────────────────────────┐
│   Client    │ ◄───────────────► │             Go HTTP Server            │
└─────────────┘                   │                                      │
                                  │  ┌──────────┐    ┌─────────────────┐ │
                                  │  │  Router  │ ►  │   Middleware    │ │
                                  │  │  (Mux)   │    │  JWT • Logger   │ │
                                  │  └──────────┘    └─────────────────┘ │
                                  │         │                             │
                                  │  ┌──────▼──────────────────────────┐ │
                                  │  │           Handlers               │ │
                                  │  │  auth • product • cart • order  │ │
                                  │  └──────────────┬──────────────────┘ │
                                  │                 │                     │
                                  │  ┌──────────────▼──────────────────┐ │
                                  │  │           Services               │ │
                                  │  │        (Business Logic)          │ │
                                  │  └──────┬──────────────┬───────────┘ │
                                  └─────────┼──────────────┼─────────────┘
                                            │              │
                              ┌─────────────▼──┐     ┌─────▼──────────┐
                              │     MySQL       │     │     Redis      │
                              │   (Primary DB)  │     │    (Cache)     │
                              └─────────────────┘     └────────────────┘
```

---

## 🚀 Quick Start

### Prerequisites

Pastikan sudah terinstall:
- [Go 1.25+](https://golang.org/dl/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [Git](https://git-scm.com/)

### Installation

```bash
# 1. Clone repository
git clone https://github.com/username/solestore.git
cd solestore

# 2. Copy environment file
cp .env.example .env

# 3. Jalankan dengan Docker Compose (recommended)
docker compose up -d

# 4. Atau jalankan secara lokal
go mod tidy
go run cmd/main.go
```

Server akan berjalan di `http://localhost:8080` 🎉

---

## ⚙️ Environment Variables

Buat file `.env` di root project:

```env
# App
APP_PORT=8080
APP_ENV=development

# JWT
JWT_SECRET=your_super_secret_key_here
JWT_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=solestore

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Payment Gateway
PAYMENT_API_KEY=your_payment_api_key
PAYMENT_SECRET=your_payment_secret
PAYMENT_WEBHOOK_SECRET=your_webhook_secret
```

---

## 📡 API Endpoints

### 🔐 Auth
```
POST   /api/v1/auth/register       Register akun baru
POST   /api/v1/auth/login          Login & dapat JWT token
POST   /api/v1/auth/refresh        Refresh access token
POST   /api/v1/auth/logout         Logout & invalidate token
```

### 👟 Products
```
GET    /api/v1/products            List produk (support filter & pagination)
GET    /api/v1/products/:id        Detail produk
POST   /api/v1/products            Tambah produk (admin)
PUT    /api/v1/products/:id        Update produk (admin)
DELETE /api/v1/products/:id        Hapus produk (admin)
```

Query params yang didukung: `?brand=Nike&size=42&color=black&min_price=100000&max_price=500000&page=1&limit=10`

### 🛒 Cart
```
GET    /api/v1/cart                Lihat isi cart
POST   /api/v1/cart                Tambah item ke cart
PUT    /api/v1/cart/:item_id       Update quantity item
DELETE /api/v1/cart/:item_id       Hapus item dari cart
```

### 💳 Orders & Payment
```
POST   /api/v1/orders              Buat order dari cart
GET    /api/v1/orders              Riwayat order user
GET    /api/v1/orders/:id          Detail order
POST   /api/v1/payments/checkout   Buat payment link
POST   /api/v1/payments/webhook    Webhook handler dari payment gateway
```

---

## 🗄️ Database Schema

```sql
users          → id, name, email, password_hash, role, created_at
products       → id, name, brand, description, price, stock
product_sizes  → id, product_id, size, stock
carts          → id, user_id, created_at
cart_items     → id, cart_id, product_id, size, quantity
orders         → id, user_id, total_price, status, created_at
order_items    → id, order_id, product_id, size, quantity, price
payments       → id, order_id, amount, status, gateway_ref, created_at
```

---

## 🐳 Docker Setup

Project ini fully containerized. Semua service dijalankan via Docker Compose:

```yaml
# docker-compose.yml (overview)
services:
  app:      # Go API server  → port 8080
  mysql:    # MySQL 8.0      → port 3306
  redis:    # Redis 7        → port 6379
```

```bash
# Start semua service
docker compose up -d

# Lihat logs
docker compose logs -f app

# Stop semua service
docker compose down

# Rebuild image setelah perubahan kode
docker compose up -d --build
```

---

## 📁 Project Structure

```
solestore/
├── cmd/
│   └── main.go                  # Entry point
├── internal/
│   ├── handler/                 # HTTP handlers
│   │   ├── auth.go
│   │   ├── product.go
│   │   ├── cart.go
│   │   └── order.go
│   ├── service/                 # Business logic
│   ├── repository/              # Database queries
│   ├── middleware/              # JWT, logger, CORS
│   └── model/                  # Structs & domain models
├── pkg/
│   ├── database/                # MySQL & Redis connection
│   ├── jwt/                     # JWT helper
│   └── response/                # Standarisasi HTTP response
├── migrations/                  # SQL migration files
├── docker-compose.yml
├── Dockerfile
├── .env.example
└── README.md
```

---

## 🧪 Running Tests

```bash
# Unit tests
go test ./...

# Dengan coverage
go test ./... -cover

# Test package spesifik
go test ./internal/service/... -v
```

---

## 📬 Contact

<div align="center">

Dibuat dengan ☕ dan 💻 oleh **Bagus**

[![GitHub](https://img.shields.io/badge/GitHub-@username-181717?style=for-the-badge&logo=github)](https://github.com/username)
[![LinkedIn](https://img.shields.io/badge/LinkedIn-Connect-0A66C2?style=for-the-badge&logo=linkedin)](https://linkedin.com/in/username)

</div>

---

<div align="center">
<img width="100%" src="https://capsule-render.vercel.app/api?type=waving&color=0:1a1a2e,50:16213e,100:0f3460&height=100&section=footer"/>
</div>
