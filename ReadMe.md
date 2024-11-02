# E-commerce API

## Overview

This is a Go-based e-commerce API that provides a set of endpoints for managing products, orders, and users. The API is designed to be RESTful and follows standard HTTP methods for CRUD operations. The application uses JSON Web Tokens (JWT) for secure authentication and authorization.

**Base URL**: `http://localhost:4000`

## Directory Structure

The project is organized into the following directories:

* `config`: Contains configuration files for the application.
* `docs`: Contains generated Swagger documentation files for the API.
* `middleware`: Contains middleware functions for authentication and authorization.
* `models`: Contains data models for the application.
* `products`: Contains product-related business logic and endpoints.
* `orders`: Contains order-related business logic and endpoints.
* `auth`: Contains authentication logic for the application
* `routes`: Contains route definitions for the API.
* `services`: Contains business service logic for products, orders, and users.
* `utils`: Contains utility functions and helpers for the application.

## Endpoints

The API provides the following endpoints:

### Authentication
* **`POST /api/v1/auth/register`**: Registers a new user.
* **`POST /api/v1/auth/login`**: Logs in a user and returns a JWT token.

### Product Management
* **`GET /api/v1/products`**: Retrieves a list of all products (publicly accessible).
* **`GET /api/v1/products/{id}`**: Retrieves details of a specific product by ID.
* **`POST /api/v1/products`**: Creates a new product (Admin only).
* **`PUT /api/v1/products/{id}`**: Updates an existing product by ID (Admin only).
* **`DELETE /api/v1/products/{id}`**: Deletes a product by ID (Admin only).

### Order Management
* **`GET /api/v1/orders`**: Retrieves a list of orders for the authenticated user.
* **`POST /api/v1/orders`**: Creates a new order for the authenticated user.
* **`PUT /api/v1/orders/{id}/cancel`**: Cancels a pending order by ID for the authenticated user.
* **`PUT /api/v1/orders/{id}/status`**: Updates the status of an order (Admin only).

---

### Additional Notes
- **Authentication**: Most routes require a JWT token. To access protected endpoints, include the token in the `Authorization` header as `Bearer <token>`.
- **Admin Access**: Only admin users can access certain routes, such as creating, updating, or deleting products, managing users, and updating order statuses.
- **Swagger Documentation**: For detailed documentation and testing of each endpoint, refer to the Swagger UI available at `http://localhost:4000/swagger/index.html`.

## Authentication and Authorization

The API uses JSON Web Tokens (JWT) for authentication and authorization. The `middleware/auth.go` file contains the authentication middleware that checks for a valid JWT token in the `Authorization` header.

## Models

The application uses the following data models:

* **Product**: Represents a product with fields for `ID`, `Name`, `Description`, `Price`, and `Stock`.
* **Order**: Represents an order with fields for `ID`, `UserID`, `ProductID`, and `Quantity`.
* **User**: Represents a user with fields for `ID`, `Name`, and `Email`.

## Services

The application uses the following services:

* **ProductService**: Provides business logic for product-related operations.
* **OrderService**: Provides business logic for order-related operations.
* **UserService**: Provides business logic for user-related operations.

## Configuration

The application uses a configuration file (`config/config.go`) to manage environment variables and configuration settings. The main configuration values include:

- `DB_CONNECTION_STRING`: The PostgreSQL connection string.
- `PORT`: The port on which the server will run.
- `JWT_SECRET`: Secret key for JWT authentication.
- `SWAGGER_SERVER_URL`: URL for serving Swagger documentation.

Set these environment variables in a `.env` file in the root directory.

## Running the Application

To run the application, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/Hinata-Akiro/ecommerce-api
   cd ecommerce-api
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Create a .env file in the root directory and add your configuration details:
   ```bash
   DB_CONNECTION_STRING=postgres://<username>:<password>@localhost:5432/ecommerce-api
   PORT=:4000
   JWT_SECRET=<your-secret-key>
   SWAGGER_SERVER_URL=localhost:4000
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

5. Access the Swagger documentation at:
   ```bash
   http://localhost:4000/swagger/index.html
   ```