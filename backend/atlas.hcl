// Atlas configuration for CV Builder database migrations
// This file configures how Atlas manages database schema migrations

// Define the environments (local, staging, production)
env "local" {
  src = "file://migrations"
  dst = "postgres://cv_user:cv_password@localhost:5432/cv_db?sslmode=disable"
  
  // Migration directory
  migration {
    dir = "file://migrations"
  }
  
  // Allow destructive changes in development
  allow_destructive = true
}

env "docker" {
  src = "file://migrations"
  dst = "postgres://cv_user:cv_password@sql-db:5432/cv_db?sslmode=disable"
  
  migration {
    dir = "file://migrations"
  }
  
  allow_destructive = true
}
