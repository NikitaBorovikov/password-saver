# REST API for secure storage and generation of new passwords
## Description
This API is designed to help users securely store, manage, and generate strong passwords for their accounts. Built with security as a top priority, this project uses ```AES-128-GCM``` encryption to ensure that your sensitive information remains confidential and protected from unauthorized access. 

## Features
- 🔒 Secure password storage with AES-128-GCM encryption
- 🔑 Password generation with customizable complexity
- 🛡️ Protection at every stage
- 📦 Clean architecture implementation
- 🔄 Session-based authentication

## Security Overview

### Authentication & Sessions
- **Secure Session Management**: Uses Gorilla Sessions with cryptographically random session IDs
- **Cookie Security**: HttpOnly and SameSite flags enabled
- **No Password Storage**: Only encrypted password data is stored

### Data Protection
- **AES-128-GCM Encryption**: For both passwords and service names
- **Separate Encryption Keys**: Different keys for different data types
- **Secure Key Management**: Via environment variables (.env)

### API Security
- **Rate Limiting**: Protection against brute force and DDoS attacks
- **Input Validation**: Protection against injection attacks
- **Zero Sensitive Data Exposure**: In responses or errors

## Technology Stack
  - API in accordance with <b>REST</b> principles.
  - The structure of the application in accordance with the principles of the <b>Clean Architecture</b>.
  - Storing data using <b>Postgres</b>. Generation of migration files.
  - HTTP server <a href = https://github.com/go-chi/chi>go-chi</a>.
  - Configuration using <a href = https://github.com/ilyakaznacheev/cleanenv>cleanenv</a>. Working with environment variables.
  - Implemented registration and authentication using <a href = https://github.com/gorilla/sessions>Gorilla Sessions</a>.
  - Writing SQL queries using <a href = https://github.com/jmoiron/sqlx>sqlx</a>.
  - Data encryption using <a href = https://github.com/alpertayfun/crypto-aes>crypto/AES</a>.
  - Working with Dockerfile and docker-compose.
  - API documentation using <a href = https://github.com/go-swagger/go-swagger>swagger</a>.

## Possible improvements
  - switching to HTTPS

## Getting Started
- Go(version 1.23 or higher)
- Docker
- Make (optional)

## Installation
### 1. Clone the repository
```
git clone https://github.com/NikitaBorovikov/password-saver.git
cd password-saver
```
### 2. Environment Setup
This project uses two environment files:
- ```.env``` - for <b>production</b>
- ```.env.dev``` - for <b>development</b>

Copy the template from the repository:
```
cp env.example .env       # For production
cp env.example .env.dev   # For development
```
Open the .env file and fill in the required values:
Open ```.env``` and ```.env.dev``` files and fill in the required values.

### 3. Build and run the application:
Build the docker image:
```
 docker build -t password-saver-app .  
```
Run the application:
```
docker-compose up
```
The server will be running on port ```8081```

### 4. Run db migration
If the application is being launched for the first time, migrations must be applied to the database:
```
make migrate
```

## Documentation 

### API
- To view the API documentation, you can use swagger (go to http://localhost:8081/swagger/index.html after the server is started).
- [Swagger files](dosc/)

### Database
- [DB Diagram](docs/db_diagram.png)