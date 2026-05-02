# Macauth

## A lightweight SSO service designed for centralized authentication.


## 🚀 Quick Start

### Prerequisites
- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) installed on your machine.

### 1. Clone the repository
```bash
git clone https://github.com/dmi3midd/lw-macauth.git
cd lw-macauth
```

### 2. Initialize the environment
Run the setup script.
```bash
./setup.sh
```

### 3. Start the service
Start the API in the background using Docker Compose:
```bash
docker compose up -d --build
```
Your SSO service is now up and running at `http://localhost:2700` (default port).

---

## 📖 API Documentation

The service exposes the following endpoints. By default, the base URL is `http://localhost:2700`.

### 1. Generate Tokens
`POST /generate`
Generates a new access and refresh token pair for a user.

**Request Body:**
```json
{
  "user": {
    "userId": "string",
    "username": "string",
    "email": "string",
    "isAdmin": "boolean"
  },
  "serviceId": "string"
}
```

**Response:**
```json
{
  "tokens": {
    "accessToken": "string",
    "refreshToken": "string"
  },
  "tokenId": "string"
}
```

### 2. Validate Access Token
`GET /validate/access`
Validates an existing access token and returns the user information.

**Headers:**
- `Authorization: Bearer <access_token>`

**Response:**
```json
{
  "user": {
    "userId": "string",
    "username": "string",
    "email": "string",
    "isAdmin": "boolean"
  }
}
```

### 3. Validate Refresh Token
`POST /validate/refresh`
Validates a refresh token and returns the associated `userId` and `tokenId`.

**Request Body:**
```json
{
  "refreshToken": "string"
}
```

**Response:**
```json
{
  "userId": "string",
  "tokenId": "string"
}
```

### 4. Get Public Key
`GET /public`
Returns the RSA Public Key used to sign the JWT tokens. You can fetch this key to perform **local validation** of access tokens in your other microservices without making network requests to `lw-macauth`.


## Will be improved in the future with more features...