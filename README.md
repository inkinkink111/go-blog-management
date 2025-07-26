# Go Blog Management API

A high-performance RESTful API for blog management built with Go, Fiber, MongoDB, and Redis caching.

## 🚀 Features

- **User Authentication**: JWT-based registration and login
- **Blog Management**: Full CRUD operations for blog posts
- **Redis Caching**: Optimized performance
- **Pagination & Filtering**: Efficient data retrieval with tag-based filtering
- **API Documentation**: Complete Swagger/OpenAPI documentation

## 🛠️ Tech Stack

- **Framework**: [Fiber](https://github.com/gofiber/fiber) - Express-inspired web framework
- **Database**: [MongoDB](https://www.mongodb.com/) - Document-based NoSQL database
- **Caching**: [Redis](https://redis.io/) - In-memory data structure store
- **Authentication**: JWT (JSON Web Tokens)
- **Documentation**: [Swagger](https://swagger.io/) - API documentation and testing

## 🚦 Getting Started

### Prerequisites

- Go 1.19 or higher
- MongoDB (local or cloud instance)
- Redis server
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/inkinkink111/go-blog-management.git
   cd go-blog-management
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   # Create .env file
   cp .env.example .env
   ```
   
   Update `.env` with your configurations:
   ```env
   # Database
   MONGODB_URI=
   DB_NAME=
   
   # Redis
   REDIS_URL=localhost:6379
   REDIS_USERNAME=
   REDIS_PASSWORD=
   
   # JWT
   JWT_SECRET_KEY=your-secret-key
   
   # Server
   PORT=3000
   ```

4. **Generate Swagger documentation**
   ```bash
   swag init
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

The API will be available at `http://localhost:3000`
