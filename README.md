# REST API for secure storage and generation of new passwords
## Description
This API is designed to help users securely store, manage, and generate strong passwords for their accounts. Built with security as a top priority, this project uses ```AES-128-GCM``` encryption to ensure that your sensitive information remains confidential and protected from unauthorized access. 
## Security
The security of the REST API is a top priority in this project. Below are the key measures implemented to ensure the confidentiality, integrity, and availability of user data:
### 1. Authentication
<b>Session-Based Authentication:</b> The API uses <a href = https://github.com/gorilla/sessions>Gorilla Sessions</a>. to manage user authentication securely. Each session ID generated using <a href = https://pkg.go.dev/crypto/rand>crypto/rand</a>.

<b>Secure Cookies:</b> Session cookies are configured with the HttpOnly and SameSite flags.

### 2. Data Encryption
<b>AES-128-GCM Encryption:</b> Sensitive data, such as passwords and associated service names, is encrypted using the AES-128-GCM algorithm. This provides both confidentiality and data integrity.

<b>Separate Encryption Keys:</b> Different keys are used to encrypt passwords and service names, adding an extra layer of security.

### 3. Secure Key Management
<b>Environment Variables:</b> Encryption keys and session values are stored in a protected .env file, which is excluded from version control (.gitignore).

<b>No Hardcoded Secrets:</b> The API avoids hardcoding sensitive information in the codebase.

### 4. API Response Security
<b>No Sensitive Data Exposure:</b> The API ensures that no confidential information (passwords, encryption keys) is included in responses.

<b>Custom Error Messages:</b> Error responses are carefully designed to avoid revealing internal details,  which could be exploited by attackers.

### 5. Input Validation
<b>Validation:</b> All user inputs are sanitized and validated to prevent injection attacks (SQL injection, XSS).

## Technologies and frameworks
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
