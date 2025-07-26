# Go GORM Learning Project

A comprehensive learning project demonstrating GORM (Go Object-Relational Mapping) with PostgreSQL database operations, relationships, and transactions.

## Prerequisites

- Go 1.23.4+
- PostgreSQL
- Git

## Installation

1. Clone repository:
```bash
git clone <repository-url>
cd <this folder>
```

2. Install dependencies:
```bash
go mod download
```

3. Setup PostgreSQL database and update connection string in `gorm_test.go`

4. Optional: Create `.env` file for environment variables

## Database Models

- **User** - Main user entity with relationships
- **Todo** - User tasks with soft delete support
- **Wallet** - User wallet information
- **Address** - User addresses
- **Product** - Products for many-to-many relationships
- **UserLog** - Activity logging
- **UserLikeProduct** - Junction table for user-product relationships

## Testing

Run tests:
```bash
go test -v
```

Run specific test:
```bash
go test -v -run TestName
```

## Dependencies

- **gorm.io/gorm** - ORM framework
- **gorm.io/driver/postgres** - PostgreSQL driver
- **github.com/stretchr/testify** - Testing framework
- **github.com/joho/godotenv** - Environment variables

## Project Structure
```plaintext
Go_Gorm/
├── user.go # Model definitions
├── gorm_test.go # Test suite
├── go.mod # Module definition
```
