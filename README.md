# 🚚 Goldship Logistic API (Golang + Gin)


<img width="580" height="360" alt="goldship-1ede5bfe090d45ae97ab135d875534dc" src="https://github.com/user-attachments/assets/232df093-b56d-41d1-b077-72903bfea7f3" />


This repository contains the backend API for the Goldship Logistic service.


## 📦 About Goldship Logistic API

<img width="580" height="360" alt="goldship-6d1d95f616e34c9abe072d8f54478cba" src="https://github.com/user-attachments/assets/3adde587-210b-442f-a9e0-0b34717e2257" />



Goldship Logistic is a logistic service that provides efficient and reliable delivery solutions for businesses and individuals. In this project, we aim to provide a robust and scalable API that supports all core functionalities of a modern logistic operation. This includes managing shipments, tracking deliveries, handling payments, and providing real-time updates to users.


## 🛠️ Tech Stack

*   **Language**: [Go](https://go.dev/)
*   **Web Framework**: [Gin](https://gin-gonic.com/)
*   **Database**: [PostgreSQL](https://www.postgresql.org/)
*   **Authentication**: JWT (JSON Web Tokens)
*   **External Services**:
    *   [Google Maps API](https://developers.google.com/maps): For calculating distances and shipment pricing.
    *   [Xendit](https://www.xendit.co/): For handling payment processing and invoices.
*   **Documentatiions**:
    *   [OpenAPI 3.0.8](https://www.openapis.org/): For API specification.
    *   [Scalar Go](https://github.com/bdpiprava/scalar-go): For Scalar interactive API documentation.
    *   [Gin Swagger](https://github.com/swaggo/gin-swagger): For Swagger UI interactive API documentation.



## 📁 Project's Directory Structure


```cli
golang-gin-goldship-api/
├── db/
│   └── migrations/
├── docs/
├── helpers/
│   ├── commons/
│   ├── googlemap/
│   ├── middlewares/
│   ├── models/
│   └── xendit-service/
└── modules/
    ├── auth/
    ├── branches/
    ├── shipments/
    ├── users/
    │   └── roles/
    └── webhooks/
```

<!-- Full tree (with files) -->
<!--
```cli
.golang-gin-goldship-api/
├── .env
├── .env.example
├── .gitignore
├── README.md
├── db/
│   ├── db-connection.go
│   └── migrations/
│       ├── 000001_create_initial_tables.down.sql
│       └── 000001_create_initial_tables.up.sql
├── docs/
│   └── openapi-specs.json
├── go.mod
├── go.sum
├── helpers/
│   ├── bcrypt.go
│   ├── commons/
│   │   └── model.go
│   ├── googlemap/
│   │   └── google.go
│   ├── jwt.go
│   ├── middlewares/
│   │   └── middlewares.go
│   ├── models/
│   │   ├── branches.go
│   │   ├── payments.go
│   │   ├── shipments.go
│   │   └── users.go
│   ├── pagination.go
│   ├── tracking-number.go
│   ├── validators.go
│   └── xendit-service/
│       └── xendit-service.go
├── main.go
└── modules/
    ├── auth/
    │   ├── controllers.go
    │   ├── dto.go
    │   └── routes.go
    ├── branches/
    │   ├── controllers.go
    │   ├── dto.go
    │   └── routers.go
    ├── shipments/
    │   ├── controllers.go
    │   ├── dto.go
    │   └── routes.go
    ├── users/
    │   ├── controllers.go
    │   ├── dto.go
    │   ├── roles/
    │   │   └── roles.go
    │   └── routes.go
    └── webhooks/
        ├── controllers.go
        ├── dto.go
        └── routes.go
```
-->

## 🚀 Getting Started

To get a local copy up and running, follow these simple steps.

### Prerequisites

*   Go (version 1.24 or later)
*   PostgreSQL
*   A `.env` file with the necessary environment variables.

### Installation & Running

1.  **Clone the repository**
    ```sh
    git clone https://github.com/masadamsahid/golang-gin-goldship-api.git
    cd golang-gin-goldship-api
    ```

2.  **Create an environment file**
    Create a `.env` file in the root directory and add the required environment variables for the database, JWT secret, Google Maps API key, and Xendit credentials.
    ```env
    # Server
    APP_PORT=8080

    # Database
    DB_HOST=
    DB_USER=
    DB_PASSWORD=
    DB_NAME=goldship_go_db
    DB_PORT=5432
    DB_SSL_MODE=disable

    DB_URL="postgres://{USER}:{PASSWORD}@{HOST}:{PORT}/{NAME}?sslmode={SSL_MODE}"

    # JWT
    JWT_SECRET_KEY=""

    # Xendit
    XENDIT_SECRET_API_KEY=""
    XENDIT_WEBHOOK_VERIFICATION_TOKEN=""

    # Google Maps
    GOOGLE_MAP_API_KEY=""
    ```

3.  **Install dependencies**
    ```sh
    go mod tidy
    ```

4.  **Run the application**
    ```sh
    go run main.go
    ```
    The server will start on the port specified in your `.env` file (e.g., `http://localhost:8080`).



## 🏃‍♂️ Guide for Running DB Migrations

Install [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) first. Then run:

```bash
migrate -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" -path database/migrations up
```


## 📖 API Documentation

The full, interactive API documentation is generated from the `docs/openapi-specs.json` file. You can view it by running the application and navigating to the appropriate endpoint, typically `/docs` (Scalar) or `/swagger/index.html` (Swagger).

## 🛣️ API Endpoints

Here is a summary of the available API endpoints.

**Auth**
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :---: |
| `POST` | `/api/auth/register` | Register as new user | No |
| `POST` | `/api/auth/login` | Login as a user | No |

**Users**
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :---: |
| `GET` | `/api/users/my-shipments` | Get all shipments for the authenticated user | Yes |
| `POST` | `/api/{username}/change-role` | Change user role (SUPERADMIN/ADMIN only) | Yes |

**Branches**
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :---: |
| `POST` | `/api/branches` | Create a new branch (ADMIN/SUPERADMIN only) | Yes |
| `GET` | `/api/branches` | Get all branches | No |
| `GET` | `/api/branches/{id}` | Get a branch by ID | No |
| `PUT` | `/api/branches/{id}` | Update a branch (ADMIN/SUPERADMIN only) | Yes |
| `DELETE` | `/api/branches/{id}` | Delete a branch (ADMIN/SUPERADMIN only) | Yes |

**Shipments**
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :---: |
| `POST` | `/api/shipments` | Create a new shipment | Yes |
| `GET` | `/api/shipments` | Get all shipments (Staff/Courier only) | Yes |
| `POST` | `/api/shipments/{id}/cancel` | Cancel a shipment (Sender only) | Yes |
| `POST` | `/api/shipments/{id}/pick-up` | Mark a shipment as picked up (Staff/Courier only) | Yes |
| `POST` | `/api/shipments/{id}/transit` | Mark a shipment as in transit (Staff/Courier only) | Yes |
| `POST` | `/api/shipments/{id}/deliver` | Mark a shipment as delivered (Staff/Courier only) | Yes |
| `GET` | `/api/shipments/track/{tracking_number}` | Get shipment history by tracking number | No |

**Webhooks**
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :---: |
| `POST` | `/api/webhooks/xendit` | Handles incoming payment status from Xendit | Header Token |

**Health Check**
| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :---: |
| `GET` | `/health-check` | Retrieve the health status of the service | No |


---

Thank you for your interest in the Goldship Logistic API! 👋

<img width="524" height="360" alt="goldship-9ef1643b412a4ce1afb93a9c54fea427" src="https://github.com/user-attachments/assets/f580869d-4f53-4a27-a0d9-02ce5de6cb40" />

<br>
<br>

---

🧙‍♂️✨ Wizardly created by [masadamsahid](https://github.com/masadamsahid)