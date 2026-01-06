# Brief Bulletin API

A RESTful news/blog API built with Go and Gin. It provides user authentication, article management, category management, and comment functionality backed by PostgreSQL. Configuration is file-driven and supports optional TLS and SMTP for email flows.

## Features

- **User Authentication**: Create users, login, password reset, user verification with OTP
- **Category Management**: Create, update, delete, and list news categories
- **Comment System**: Create comments, approve/disapprove comments, get comments by article
- **JWT-based Authorization**: Request authorization with configurable bypass list
- **PostgreSQL Integration**: Using pgx pool for database connections
- **CORS Enabled**: Development-friendly defaults
- **Optional TLS Support**: Secure connections when needed
- **Optional SMTP Configuration**: Email functionality for OTP and notifications
- **sqlc-based Query Generation**: Type-safe database query code generation

## Tech Stack

- **Go** 1.23+
- **Gin** - Web framework
- **Logrus** - Structured logging
- **Viper** - Configuration management
- **pgx/v5** (pgxpool) - PostgreSQL driver
- **sqlc** - Type-safe query codegen
- **JWT** - Authentication tokens
- **Swagger** (optional) - API documentation

## Getting Started

### Prerequisites

- Go (recommended 1.23+)
- Make (optional, for convenience commands)
- sqlc (for regenerating query code)
- PostgreSQL instance accessible from your machine

### Configuration

Create or update a JSON config file. A sample is provided at `config/connection/dev-config.json`. **Do not commit secrets.**

Example configuration structure:

```json
{
  "dbhost": "localhost",
  "dbPort": 5432,
  "dbname": "postgres",
  "dbuid": "postgres",
  "dbpassword": "your-db-password",
  "timeout": 100,
  "connRetryCount": 1,
  "connRetryInterval": 5000,
  "jwtKey": "your-jwt-signing-key",
  "bypassAuth": [
    "/api/auth/create",
    "/api/auth/login",
    "/api/auth/resetpwd",
    "/api/send-otp",
    "/api/verify-otp",
    "/api/verify-user",
    "/api/category",
    "/api/comment",
    "/api/all-comments"
  ],
  "adminEmailId": "admin@user.com",
  "adminPassword": "admin4test",
  "adminEmpCode": "0000",
  "rptServiceLink": "http://localhost:9985/convert",
  "rptFilePath": "./",
  "rptAuthKey": "your-auth-key",
  "isTLS": false,
  "tlsKeyPath": "",
  "tlsCertPath": "",
  "senderEmail": "user@gmail.com",
  "password": "app-password-or-token",
  "smtp_host": "smtp.gmail.com",
  "smtp_port": 587,
  "url": {
    "uiurl": "http://localhost:3000"
  }
}
```

## Database Schema

The service uses PostgreSQL and requires the following schemas and tables:

### Create Schemas

```sql
CREATE SCHEMA IF NOT EXISTS news;
CREATE SCHEMA IF NOT EXISTS news;
```

### Users Table (news schema)

```sql
CREATE TABLE news.users (
    user_id serial4 NOT NULL,
    user_name text NOT NULL,
    email text NOT NULL,
    phone text NOT NULL,
    pass text NOT NULL,
    pss_valid bool DEFAULT true NOT NULL,
    otp text NOT NULL,
    user_valid bool DEFAULT false NOT NULL,
    otp_exp timestamp NOT NULL,
    "role" text NOT NULL
);
```

**Table Structure:**
- `user_id`: Auto-incrementing primary key (serial)
- `user_name`: User's username (text, required)
- `email`: User's email address (text, required)
- `phone`: User's phone number (text, required)
- `pass`: User's password hash (text, required)
- `pss_valid`: Password validity flag (boolean, default: true)
- `otp`: One-time password for verification (text, required)
- `user_valid`: User verification status (boolean, default: false)
- `otp_exp`: OTP expiration timestamp (timestamp, required)
- `role`: User role/permissions (text, required)

### Categories Table (news schema)

```sql
CREATE TABLE news.categories (
    id serial4 NOT NULL,
    "name" text NOT NULL,
    slug text NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    CONSTRAINT categories_unique UNIQUE (name),
    CONSTRAINT categories_unique_1 UNIQUE (slug)
);
```

**Table Structure:**
- `id`: Auto-incrementing primary key (serial)
- `name`: Category name (text, required, unique)
- `slug`: URL-friendly category identifier (text, required, unique)
- `created_at`: Timestamp of creation (timestamp, default: current timestamp)

### Articles Table (news schema)

```sql
CREATE TABLE news.articles (
    id serial4 NOT NULL,
    title text NOT NULL,
    summary text NOT NULL,
    "content" text NOT NULL,
    featured_image text NULL,
    category_id int4 NOT NULL,
    status text DEFAULT 'draft'::text NOT NULL,
    published_at timestamp NULL,
    views_count int8 DEFAULT 0 NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    source_url text NOT NULL,
    CONSTRAINT articles_pk PRIMARY KEY (id),
    CONSTRAINT articles_categories_fk FOREIGN KEY (category_id) REFERENCES news.categories(id)
);
```

**Table Structure:**
- `id`: Auto-incrementing primary key (serial)
- `title`: Article title (text, required)
- `summary`: Article summary (text, required)
- `content`: Full article content (text, required)
- `featured_image`: URL to featured image (text, nullable)
- `category_id`: Foreign key to categories table (int4, required)
- `status`: Publication status - 'draft' or 'published' (text, default: 'draft')
- `published_at`: Publication timestamp (timestamp, nullable)
- `views_count`: Number of views (int8, default: 0)
- `created_at`: Creation timestamp (timestamp, default: current timestamp)
- `updated_at`: Last update timestamp (timestamp, default: current timestamp)
- `source_url`: Source URL for the article (text, required)

### Comments Table (news schema)

```sql
CREATE TABLE news."comments" (
    id serial4 NOT NULL,
    article_id int4 NOT NULL,
    user_name text NOT NULL,
    user_email text NOT NULL,
    "content" text NOT NULL,
    is_approved bool DEFAULT false NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    CONSTRAINT comments_unique UNIQUE (id),
    CONSTRAINT comments_articles_fk FOREIGN KEY (article_id) REFERENCES news.articles(id) ON DELETE CASCADE
);
```

**Table Structure:**
- `id`: Auto-incrementing primary key (serial)
- `article_id`: Foreign key to articles table (int4, required)
- `user_name`: Commenter's name (text, required)
- `user_email`: Commenter's email (text, required)
- `content`: Comment content (text, required)
- `is_approved`: Approval status (boolean, default: false)
- `created_at`: Creation timestamp (timestamp, default: current timestamp)

## Build and Run

The `Makefile` provides convenient targets.

### Build for Windows

```bash
make winBuild
```

This compiles `build/exec/briefbuletin.exe`.

### Build and Run in Development (Windows)

```bash
make dev
```

This compiles `build/exec/briefbuletin.exe` and runs it with:

```bash
./build/exec/briefbuletin.exe -c ./config/connection/dev-config.json --port 7070
```

### Cross-compile (Windows and Linux binaries)

```bash
make server
```

This creates both Windows (`briefbuletin.exe`) and Linux (`briefbuletin`) binaries.

### Regenerate sqlc Query Code

```bash
make sqlc
```

### Generate Swagger Docs (Optional)

```bash
make swag
```

### Run Directly (without Make)

```bash
go mod download
go build -o build/exec/briefbuletin.exe main.go
./build/exec/briefbuletin.exe -c ./config/connection/dev-config.json --port 7070
```

**Command Line Flags:**
- `-c` - Path to config JSON (default: `./config.json`)
- `--port` - Server port (default: `7070`)
- `-v` - Verbose logs

## API

### Base URL

Default: `http://localhost:7070`

### Endpoints

#### Public Endpoints (No Authentication Required)

- `GET /` - Health check
- `POST /api/auth/create` - Create new user
- `POST /api/auth/login` - Login and get JWT token
- `POST /api/auth/resetpwd` - Reset password
- `POST /api/send-otp` - Send OTP for user verification
- `POST /api/verify-otp` - Verify OTP
- `POST /api/verify-user` - Verify user account
- `GET /api/category` - Get all categories (public)
- `POST /api/comment` - Create a comment (public)
- `GET /api/all-comments` - Get all comments for an article (public, requires `article_id` query parameter)

#### Protected Endpoints (Require JWT Token)

**User Management:**
- `PUT /api/auth/update` - Update user information
- `GET /api/auth/users` - Get all users

**Category Management:**
- `POST /api/category-service` - Create a new category
- `PUT /api/category-service` - Update a category
- `DELETE /api/category-service` - Delete a category

**Comment Management:**
- `GET /api/active-comment` - Approve a comment (requires `id` query parameter)
- `GET /api/disable-comment` - Disable a comment (requires `id` query parameter)
- `GET /api/approval-due-comments` - Get list of comments pending approval

**Notes:**
- Requests are intercepted by an auth middleware. Paths in `bypassAuth` are accessible without a token.
- Static API docs (if generated/copied) are served from `/apidoc`.
- For protected endpoints, include the JWT token in the `Authorization` header: `Authorization: Bearer <token>`

## Development

- CORS is open by default for development. Harden before production.
- sqlc config lives at `config/sqlc/db_query/db.yaml`. Update queries under `config/sqlc/db_query` and run `make sqlc`.
- Generated query code is in `internal/dbmodel/db_query/`.
- Maximum file upload size is 16 MB (configurable in `server.go`).

## Debugging (VS Code)

Use a launch configuration similar to:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Brief Bulletin API",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/main.go",
      "args": ["-c", "./config/connection/dev-config.json", "--port", "7070"],
      "cwd": "${workspaceFolder}"
    }
  ]
}
```

## Testing with Postman

### Prerequisites

1. Start the server (see [Build and run](#build-and-run) section)
2. Ensure the database is set up with the required tables
3. Have Postman installed

### Step 1: Create a User (No Token Required)

**Request:**
- **Method:** `POST`
- **URL:** `http://localhost:7070/api/auth/create`
- **Headers:**
  ```
  Content-Type: application/json
  ```
- **Body (raw JSON):**
  ```json
  {
    "email": "test@example.com",
    "password": "testpassword123",
    "phone": "1234567890",
    "userName": "testuser",
    "role": "USER"
  }
  ```

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "User created successfully",
  "isSuccess": true,
  "ts": "2024-01-15-10:30:45.123"
}
```

### Step 2: Login and Get JWT Token

**Request:**
- **Method:** `POST`
- **URL:** `http://localhost:7070/api/auth/login`
- **Headers:**
  ```
  Content-Type: application/json
  ```
- **Body (raw JSON):**
  ```json
  {
    "login": "test@example.com",
    "pwd": "testpassword123"
  }
  ```
  > **Note:** The `login` field can be username, email, or phone number.

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "Login successful",
  "isSuccess": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "payload": {
    "user_id": 1,
    "user_name": "testuser",
    "email": "test@example.com",
    "phone": "1234567890",
    "role": "USER"
  },
  "ts": "2024-01-15-10:31:00.456"
}
```

**Important:** Copy the `token` value from the response. You'll need it for authenticated requests.

### Step 3: Set Up Postman Environment Variable (Optional but Recommended)

1. In Postman, click on **Environments** (left sidebar)
2. Create a new environment or use the default
3. Add a variable:
   - **Variable:** `auth_token`
   - **Initial Value:** (leave empty)
   - **Current Value:** (leave empty)
4. Save the environment

### Step 4: Use Token for Authenticated Requests

#### Option A: Using Authorization Header (Recommended)

For each authenticated request, add this header:
```
Authorization: Bearer <your-token-here>
```

Replace `<your-token-here>` with the token you received from the login response.

#### Option B: Using Postman Environment Variable

1. After login, in the **Tests** tab of the login request, add:
   ```javascript
   if (pm.response.code === 200) {
       var jsonData = pm.response.json();
       pm.environment.set("auth_token", jsonData.token);
   }
   ```
2. In authenticated requests, set the Authorization header as:
   ```
   Authorization: Bearer {{auth_token}}
   ```

### Step 5: Test Protected Endpoints

#### Get All Users

**Request:**
- **Method:** `GET`
- **URL:** `http://localhost:7070/api/auth/users`
- **Headers:**
  ```
  Authorization: Bearer <your-token-here>
  Content-Type: application/json
  ```

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "Users retrieved successfully",
  "isSuccess": true,
  "payload": [
    {
      "code": 1,
      "name": "testuser",
      "email": "test@example.com"
    }
  ],
  "ts": "2024-01-15-10:32:00.789"
}
```

#### Create Category

**Request:**
- **Method:** `POST`
- **URL:** `http://localhost:7070/api/category-service`
- **Headers:**
  ```
  Authorization: Bearer <your-token-here>
  Content-Type: application/json
  ```
- **Body (raw JSON):**
  ```json
  {
    "name": "Technology",
    "slug": "technology"
  }
  ```

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "Category created successfully",
  "isSuccess": true,
  "ts": "2024-01-15-10:33:00.123"
}
```

#### Get All Categories (Public)

**Request:**
- **Method:** `GET`
- **URL:** `http://localhost:7070/api/category`
- **Headers:**
  ```
  Content-Type: application/json
  ```

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "Category List",
  "isSuccess": true,
  "payload": [
    {
      "id": 1,
      "name": "Technology",
      "slug": "technology",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "ts": "2024-01-15-10:34:00.456"
}
```

#### Create Comment (Public)

**Request:**
- **Method:** `POST`
- **URL:** `http://localhost:7070/api/comment`
- **Headers:**
  ```
  Content-Type: application/json
  ```
- **Body (raw JSON):**
  ```json
  {
    "article_id": 1,
    "user_name": "John Doe",
    "user_email": "john@example.com",
    "content": "Great article!"
  }
  ```

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "Category List",
  "isSuccess": true,
  "payload": {
    "id": 1,
    "article_id": 1,
    "user_name": "John Doe",
    "user_email": "john@example.com",
    "content": "Great article!",
    "is_approved": false,
    "created_at": "2024-01-15T10:35:00Z"
  },
  "ts": "2024-01-15-10:35:00.789"
}
```

#### Get Comments for Article (Public)

**Request:**
- **Method:** `GET`
- **URL:** `http://localhost:7070/api/all-comments?article_id=1`
- **Headers:**
  ```
  Content-Type: application/json
  ```

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "Comments for the article",
  "isSuccess": true,
  "payload": [
    {
      "id": 1,
      "article_id": 1,
      "user_name": "John Doe",
      "user_email": "john@example.com",
      "content": "Great article!",
      "is_approved": true,
      "created_at": "2024-01-15T10:35:00Z"
    }
  ],
  "ts": "2024-01-15-10:36:00.123"
}
```

#### Approve Comment (Protected)

**Request:**
- **Method:** `GET`
- **URL:** `http://localhost:7070/api/active-comment?id=1`
- **Headers:**
  ```
  Authorization: Bearer <your-token-here>
  Content-Type: application/json
  ```

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "Comment approved successfully",
  "isSuccess": true,
  "ts": "2024-01-15-10:37:00.456"
}
```

#### Get Approval Due Comments (Protected)

**Request:**
- **Method:** `GET`
- **URL:** `http://localhost:7070/api/approval-due-comments`
- **Headers:**
  ```
  Authorization: Bearer <your-token-here>
  Content-Type: application/json
  ```

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "Approval due comments retrieved",
  "isSuccess": true,
  "payload": [
    {
      "comment_id": 1,
      "title": "Article Title",
      "news": "Article content...",
      "user_name": "John Doe",
      "user_email": "john@example.com",
      "comment": "Great article!"
    }
  ],
  "ts": "2024-01-15-10:38:00.789"
}
```

### Step 6: Reset Password (No Token Required)

**Request:**
- **Method:** `POST`
- **URL:** `http://localhost:7070/api/auth/resetpwd`
- **Headers:**
  ```
  Content-Type: application/json
  ```
- **Body (raw JSON):**
  ```json
  {
    "email": "test@example.com",
    "newPwd": "newpassword123"
  }
  ```

**Expected Response (200 OK):**
```json
{
  "statusCode": 200,
  "serviceMessage": "Password reset successfully",
  "isSuccess": true,
  "ts": "2024-01-15-10:38:00.789"
}
```

### Troubleshooting

**401 Unauthorized / "Unauthorized" response:**
- Check that you've included the `Authorization: Bearer <token>` header
- Verify the token is still valid (tokens expire after 1 hour)
- Make sure there's no extra spaces in the token
- Try logging in again to get a fresh token

**400 Bad Request:**
- Verify the JSON body is valid
- Check that all required fields are present
- Ensure field names match exactly (case-sensitive)

**404 Not Found:**
- Verify the endpoint URL is correct
- Check that the resource ID exists (for GET/PUT/DELETE by ID)
- Ensure the server is running on the correct port

**500 Internal Server Error:**
- Check server logs for detailed error messages
- Verify database connection is working
- Ensure database tables exist and schema is correct

## Notes

- Keep secrets out of git. Use environment-specific config files.
- For TLS, set `isTLS: true` and provide `tlsKeyPath` and `tlsCertPath`.
- SMTP fields are optional, required only if email flows are enabled.
- The service supports graceful shutdown on SIGINT or SIGTERM signals.
- Maximum multipart memory is set to 16 MB for file uploads.
