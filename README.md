# E-Commerce Backend API

A production-style RESTful backend for an e-commerce platform built with Go as a personal learning project — covering authentication, product management, shopping cart, and order checkout with full inventory control.

> Built to learn real-world backend patterns: layered architecture, JWT auth, database transactions, and role-based access control.

---

## Tech Stack

| Technology | Purpose |
|------------|---------|
| **Go 1.21** | Main language |
| **Fiber v2** | HTTP web framework |
| **GORM** | ORM for database operations |
| **PostgreSQL 15** | Relational database |
| **JWT** | Stateless authentication |
| **bcrypt** | Password hashing |
| **Docker** | Local database setup |

---

## Architecture

Follows a 3-layer architecture to separate concerns clearly:

```
HTTP Request
     ↓
  Handler        ← parse request, return JSON response
     ↓
  Service        ← business logic and validation
     ↓
 Repository      ← database queries only (GORM)
     ↓
  PostgreSQL
```

Each layer has one responsibility and never crosses into another's job. Handlers never touch the database. Repositories never contain business logic.

---

## Project Structure

```
ecommerce/
├── cmd/
│   └── main.go              ← entry point, wires all layers and routes
├── config/
│   └── database.go          ← PostgreSQL connection
└── internal/
    ├── middleware/
    │   └── auth.go          ← JWT helpers, RequireAuth, RequireAdmin
    ├── model/               ← GORM structs (maps to DB tables)
    ├── repository/          ← database queries only
    ├── service/             ← business logic and validation
    └── handler/             ← Fiber HTTP handlers
```

---

## Features

### Authentication
- User registration with bcrypt password hashing
- Login returns a signed JWT token (24hr expiry)
- Token carries `user_id` and `role` — verified on every protected request
- Passwords are never stored in plain text and never returned in responses (`json:"-"`)

### Role-Based Access Control (RBAC)
- Two roles: `user` and `admin`
- `RequireAuth` middleware — blocks requests with missing or invalid tokens (401)
- `RequireAdmin` middleware — blocks non-admin users from admin routes (403)
- Middleware runs before the handler — handler only executes if both checks pass

### Product Management
- Public read: anyone can browse products without a token
- Admin write: create, update, and soft delete (protected by `RequireAdmin`)
- Soft delete — deleted products stay in the database so order history stays intact

### Shopping Cart
- Add item with upsert logic — adding an existing product increments quantity instead of creating a duplicate row
- Update quantity — setting to 0 removes the item automatically
- Cart items include full product details via GORM `Preload`
- `user_id` always comes from the verified JWT, never from the request body

### Order Checkout (GORM Transaction)
- Entire checkout runs inside one atomic transaction:
  1. Validate stock for every cart item
  2. Create the order row
  3. Create order_item rows (price snapshot at checkout time)
  4. Deduct stock using `gorm.Expr("stock - ?", qty)` — atomic SQL update, prevents overselling
  5. Clear the cart
- If any step fails → full rollback — no orphaned orders, no incorrect stock

---

## API Endpoints

### Auth — public
```
POST   /api/auth/register    Create account
POST   /api/auth/login       Login, returns JWT token
```

### Products — public read
```
GET    /api/products         List all products
GET    /api/products/:id     Get one product
```

### Products — admin only
```
POST   /api/admin/products         Create product
PUT    /api/admin/products/:id     Update product
DELETE /api/admin/products/:id     Soft delete product
```

### Cart — login required
```
GET    /api/cart                   View cart (with product details)
POST   /api/cart                   Add item (upsert)
PUT    /api/cart/:product_id       Update quantity
DELETE /api/cart/:product_id       Remove item
```

### Orders — login required
```
POST   /api/orders           Checkout (atomic transaction)
GET    /api/orders           Order history
```

---

## Local Setup

**Prerequisites:** Go 1.21+, Docker

### 1. Clone the repository
```bash
git clone https://github.com/suphachok09790/ecommerce-backend.git
cd ecommerce-backend
```

### 2. Start PostgreSQL with Docker
```bash
docker compose up -d
```

### 3. Configure environment
```bash
cp .env.example .env
# .env is already pre-filled to match docker-compose — no changes needed locally
```

### 4. Run
```bash
go mod tidy
go run cmd/main.go
```

Server starts at `http://localhost:3000`. Tables are created automatically on first run.

---

## Example Usage

**Register and login:**
```bash
# register
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com","password":"secret"}'

# login — copy the token from the response
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"secret"}'
```

**Full checkout flow:**
```bash
TOKEN="your_token_here"

# add item to cart
curl -X POST http://localhost:3000/api/cart \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product_id":1,"quantity":2}'

# checkout
curl -X POST http://localhost:3000/api/orders \
  -H "Authorization: Bearer $TOKEN"
```

**Make a user admin (run in PostgreSQL):**
```sql
UPDATE users SET role = 'admin' WHERE email = 'john@example.com';
```
Then login again — the new token will carry `role: "admin"`.

---

## What I Learned Building This

- How to structure a Go backend with Clean 3-layer architecture
- How JWT authentication works end to end — signing, verifying, and passing claims through middleware
- Why database transactions matter — and how `gorm.Expr` prevents race conditions on stock updates
- The difference between soft delete and hard delete, and when to use each
- How dependency injection keeps layers loosely coupled and independently testable
