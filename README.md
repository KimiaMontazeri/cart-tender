# Cart Tender

This is a simple shopping cart microservice implemented in Go using the Echo web framework, GORM, and PostgreSQL as the database for the Internet Engineering course at AUT. 
The application provides endpoints to manage user accounts and their associated shopping carts.

## Endpoints
- `POST /register`: Registers a user.
- `POST /login`: Logs the user in.
- `GET /basket`: Retrieve a list of shopping baskets.
- `POST /basket`: Create a new shopping basket.
- `PATCH /basket/:id`: Update an existing shopping basket.
- `GET /basket/:id`: Retrieve details of a specific shopping basket.
- `DELETE /basket/:id`: Delete a shopping basket.

## Authentication
The application uses JWT for authentication.

## Database Migrations
The application performs automatic database migrations on startup to ensure the tables are created and updated.
