# Todo List RESTful CRUD API

A secure and efficient RESTful CRUD Todo List application built with Golang, Fiber, and PostgreSQL (via PGX driver), designed with security and usability best practices.

![cover](https://github.com/vadimfedulov035/vorterumilo/blob/main/cover.png)

## 🚀 Features

### Security Highlights
- **Minimal Container Footprint**  
    Dockerized using Alpine Linux with:
    - Non-root user (`nobody:nogroup`)
    - All Linux capabilities dropped (`cap_drop: ALL`)
    - Read-only filesystem (except `/tmp` and `/var/run`)
    - Privilege escalation prevention (`no_new_privileges: true`)

- **Database Hardening**  
    - Separate PostgreSQL superuser (`postgres_admin`) and application user (`task_manager`)
    - Schema isolation (`private` schema by default)
    - Least-privilege permissions for application user

- **Controlled Error Exposure**  
    Only specific safe errors (e.g., validation issues) are shown to users. All others return a generic "Something went wrong" message to prevent information leakage.

- **Strict Input Validation**  
    Enforced type constraints with clear error messages:
    ```json
    {
        "error": "Invalid status: must be string from (\"new\", \"in_progress\", \"done\"), got string 'pending'"
    }
    ```

### Usability Highlights
- **Zero-Config Setup**  
    Fully Dockerized with automatic:
    - Database initialization
    - Schema migrations (`0001_create_tasks_table.up.sql`)
    - Environment variable loading (via `.env`)

- **Self-Documenting Errors**  
    Validation failures return human-readable explanations:
    ```json
    {
        "error": "Invalid title: must be non-empty string, got string ''"
    }
    ```

- **Versioned API Endpoints**  
    Clear endpoint structure with `/api/v1` versioning.

## 🛠️ Getting Started

### Prerequisites
- Docker & Docker Compose
    ```bash
    sudo apt install docker-ce docker compose -y  # Install Docker (add official repo)
    ```

### 🐳 Setup with Docker
1. **Clone Repository**
    ```bash
    git clone https://github.com/vadimfedulov035/go-todo-api.git
    cd go-todo-api
    ```

2. **Create `.env` file**
    ```ini
    # Database superuser (initialization only)
    POSTGRES_SUPERUSER=postgres_admin
    POSTGRES_SUPERPASSWORD=super_secret

    # Application credentials
    DB_USER=task_manager
    DB_PASSWORD=app_secret
    DB_HOST=postgres
    DB_PORT=5432
    DB_NAME=tasks_db
    DB_SCHEMA=private
    DB_TABLENAME=private.tasks
    ```

3. **Manage Containers**
    ```bash
    docker compose up -d   # Start the containers
    sudo docker compose logs -f --tail 25  # Follow last 25 log lines
    sudo docker compose down  # Stop the containers
    ```

The API will be available at `http://localhost:8080`.

## 📡 API Endpoints

| Method | Endpoint                | Description                     |
|--------|-------------------------|---------------------------------|
| GET    | `/api/v1/tasks`         | Retrieve all tasks              |
| POST   | `/api/v1/tasks`         | Create new task                 |
| GET    | `/api/v1/tasks/:id`     | Get single task by ID           |
| PUT    | `/api/v1/tasks/:id`     | Update existing task by ID      |
| DELETE | `/api/v1/tasks/:id`     | Delete task by ID               |

**Task Schema**
```json
{
    "id": 1,
    "title": "Buy groceries",
    "description": "Milk, eggs, bread",
    "status": "new",
    "created_at": "2024-02-20T12:34:56Z",
    "updated_at": "2024-02-20T12:34:56Z"
}
```

## 🔒 Security Implementation Details

### Database Permissions
```sql
-- Application user permissions according to default .env variables
GRANT CONNECT ON DATABASE tasks_db TO task_manager;
GRANT USAGE, CREATE ON SCHEMA private TO task_manager;
ALTER DEFAULT PRIVILEGES IN SCHEMA private
GRANT SELECT, INSERT, UPDATE, DELETE ON private.tasks TO task_manager;
--- and so on...
```

### Container Security
```dockerfile
# Final stage: Minimal secure base
FROM scratch
# ...
ENV TZ=UTC
EXPOSE 8080
USER 65534:65534
```

## 🌐 Environment Variables

| Variable                 | Purpose                                  |
|--------------------------|------------------------------------------|
| `POSTGRES_SUPERUSER`     | PostgreSQL admin username (setup only)   |
| `POSTGRES_SUPERPASSWORD` | PostgreSQL admin password (setup only)   |
| `DB_USER`                | Application database user                |
| `DB_PASSWORD`            | Application database password            |
| `DB_SCHEMA`              | Isolated schema for data protection      |

## 🛑 Error Handling

**Exposed Errors**
- `TaskError`: Field validation failures (HTTP 422)
- `fiber.ErrUnprocessableEntity`: Invalid request format
- `pgx.ErrNoRows`: Missing resources (HTTP 404)

**All other errors return:**
```json
{
    "error": "Something went wrong"
}
```

---

[![License](https://img.shields.io/badge/license-GPLv3-blue.svg)](#)  
**License**: GNU General Public License v3.0  
