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
- **Zero Sensitive Data Exposure**: In responses or errors
- **Input Validation**: Protection against injection attacks

## Technology Stack
  - API in accordance with <b>REST</b> principles.
  - The structure of the application in accordance with the principles of the <b>Clean Architecture</b>.
  - Storing data using <b>Postgres</b>. Generation of migration files.
  - HTTP server <a href = https://github.com/go-chi/chi>go-chi</a>.
  - Configuration using <a href = https://github.com/ilyakaznacheev/cleanenv>cleanenv</a>. Working with environment variables.
  - Implemented registration and authentication using <a href = https://github.com/gorilla/sessions>Gorilla Sessions</a>.
  - Writing SQL queries using <a href = https://github.com/jmoiron/sqlx>sqlx</a>.
  - data encryption using <a href = https://github.com/alpertayfun/crypto-aes>crypto/AES</a>.

## Possible improvements
  - switching to HTTPS

## Getting Started
- Go(version 1.23 or higher)

## Installation
### 1. Clone the repository
```
git clone https://github.com/NikitaBorovikov/password-saver.git
cd password-saver
```
### 2. Set up environment variables

