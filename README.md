# Id Trials

Welcome to the Id Trials repository! This repository is a comprehensive setup for a project utilizing Go, React, and Docker. It's designed to provide a streamlined development environment with all the necessary configurations in place.

## Technologies

- **Go 1.21+**: A popular programming language designed for simplicity and efficiency.
- **React 19**: A platform for building mobile and desktop web applications using JavaScript.
- **Docker**: A set of platform-as-a-service products that use OS-level virtualization to deliver software in packages called containers.
- **Makefile**: used to run basic site dev functionality
- **Nginx**: used to host React and the go server

## Getting Started

Before you begin, make sure you have Docker installed on your machine. If not, you can download and install it from [Docker's official website](https://www.docker.com/).

### Initial Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/theCompanyDream/id-trials.git
   cd id-trials
   ```

2. Check for the `.env` file:
   - The Makefile will automatically check for the existence of a `.env` file. If it's not found, it will create one by copying the `.env.example` template.
   - Ensure that you modify the `.env` file with your specific configurations.

### Commands

The Makefile includes several commands to manage the development environment:

- **dev**: Starts the development environment using Docker.

  ```bash
  make dev
  ```

- **restart**: Restarts all the services.

  ```bash
  make restart
  ```

- **down**: Shuts down the development environment and removes volumes.

  ```bash
  make down
  ```

- **build**: Builds or rebuilds services.

  ```bash
  make build
  ```

- **stop**: Stops running containers without removing them.

  ```bash
  make stop
  ```

- **test**: Runs tests for the User application.

  ```bash
  make test
  ```

## Development Workflow

1. Start the development environment:

   ```bash
   make dev
   ```

   This will spin up the necessary Docker containers for the Go backend and the React frontend.

2. You can access the application at many different points depending on what your style there is a nginx server that runs at localhost port 80 where you can access at port 80. However the application do run indpendently on frontend on port http://localhost:4200, and api on http://localhost:3000 respectively. One more point is that the api is also accessable at http://localhost/api from the nginx server

3. To make any changes or additions to the project, you can modify the source code in the respective directories for the Go and React applications.

4. Run tests to ensure everything is functioning correctly:

   ```bash
   make test
   ```

5. Once you're done with development, you can stop the services:

   ```bash
   make stop
   ```

   or completely shut down the environment:

   ```bash
   make down
   ```

Thank you for using the User Test repository! If you encounter any issues or have suggestions for improvements, feel free to open an issue or submit a pull request.

## Future

if I come back to this I might create multiple user tables and compare

- uuid
- ulid
- cuid
- nanoid
- snowflakeId
- ksuid
