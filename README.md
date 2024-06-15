
# UserService
<!-- TODO add github action badge -->
<!-- ![GitHub Actions](https://github.com/sean-cunniffe/user-service/actions/workflows/main.yml/badge.svg) -->

## Overview

UserService is a microservice designed to handle CRUD (Create, Read, Update, Delete) operations for user management. It also includes functionality to verify user passwords. The service is written in Go and communicates over gRPC, making it efficient and scalable.

## Features

- **CRUD Operations**: Create, read, update, and delete user records.
- **Password Verification**: Securely verify user passwords.
- **gRPC Interface**: Efficient and structured communication protocol.

## Architecture

- **Language**: Go
- **Database**: Redis
- **Communication Protocol**: gRPC
- **Protobuf Definition**: Located at `interface/grpc/user-service.proto`

## Development Environment

The development environment is configured using a dev container. The configuration file is located at `.devcontainer/devcontainer.json`. This ensures a consistent development setup and simplifies dependency management.

## CI/CD Pipeline

The project utilizes GitHub Actions for continuous integration and continuous deployment (CI/CD). The pipeline automates testing, building, and deploying the microservice.

## Getting Started

### Prerequisites

- Docker
- Go 1.22+
- Redis

### Setting Up Development Environment

1. **Clone the repository:**
   ```bash
   git clone https://github.com/sean-cunniffe/user-service.git
   cd user-service
   ```

2. **Start the development container:**
   Ensure Docker is running and use your preferred IDE to open the repository. The IDE should automatically detect the `.devcontainer` folder and prompt you to reopen the folder in the container. Follow the prompts to start the dev container.

3. **Install dependencies:**
   Inside the dev container, run:
   ```bash
   go mod tidy
   ```

### Running the Service

1. **Start Redis:**
   ```bash
   docker run -d -p 6379:6379 redis
   ```

2. **Run the service:**
   ```bash
   go run main.go
   ```

### Running Tests

To run tests, use the following command:
```bash
go test ./...
```

## gRPC Interface

The gRPC service definitions are specified in `interface/grpc/user-service.proto`. To regenerate the Go code for the gRPC interface after making changes to the `.proto` file, use the following command:
```bash
protoc --go_out=. --go-grpc_out=. interface/grpc/user-service.proto
```

## Deploying with Helm

### Prerequisites

- Kubernetes cluster
- Helm 3.x

### Installing the Helm Chart

1. **Add the Helm repository:**
   ```bash
   helm repo add userservice https://your-username.github.io/UserService
   helm repo update
   ```

2. **Install the chart:**
   ```bash
   helm install my-userservice userservice/userservice
   ```

### Uninstalling the Helm Chart

To uninstall/delete the `my-userservice` deployment:
```bash
helm uninstall my-userservice
```

### Customizing the Chart

To customize the deployment, you can pass custom values using the `--set` or `-f` flag with your custom `values.yaml` file.

```bash
helm install my-userservice userservice/userservice --values custom-values.yaml
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Contact

For questions or feedback, please open an issue on GitHub.

---

Happy coding!
