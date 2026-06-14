# Cloud POS System - Backend

- A Go backend service for ticket sales, place access management, history tracking, and entry validation.

- This project utilizes hexagonal architecture inspired by https://www.tomkdickinson.co.uk/hexagonal-architecture-with-go-and-google-wire-e4344dd24b94

- For frontend please visit https://github.com/Maimeetang/Ticketing-system-frontend

<br>

## Development Setup

### Prerequisites

- Docker & Docker Compose

<br>

1. **Clone the Repository**

   ```bash
   git clone <repository-url> project_name
   cd project_name
   ```

   <br>

2. **Environment Configuration**
   Duplicate the example environment file and fill in your credentials:

   ```bash
   cp env.example .env
   ```

   _Open `.env` and configure your environment

   <br>

3. **Create External Docker Network**
   Since the configuration uses an external network for `pos-dev-net`, create it before running the containers:

   ```bash
   docker network create pos-dev-net
   ```

   <br>

4. **Run Development Environment via Docker**
   Spin up the required infrastructure (Database, Redis, etc.) and application containers using the development Compose file:

   ```bash
   docker compose -f compose-dev.yml up -d --build
   ```

   <br>

5. **Verify Container Status**
   ```bash
   docker compose -f compose-dev.yml ps
   ```
   _The application will dynamically watch your local file changes and sync via volumes._

<br>

## API Documentation

- Response Format: JSend (https://github.com/omniti-labs/jsend)

- More documents will be added in the future.
