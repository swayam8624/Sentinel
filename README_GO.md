# Sentinel Gateway - Go Implementation

This is the Go implementation of the Sentinel Gateway, a self-healing LLM firewall with cryptographic data protection.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Docker (for development and testing)
- PostgreSQL (for metadata storage)
- Redis (for caching)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/sentinel-platform/sentinel.git
cd sentinel
```

2. Install dependencies:

```bash
make deps
```

3. Build the application:

```bash
make build
```

### Running the Application

#### Development Mode

1. Start the required services:

```bash
docker-compose up -d
```

2. Run the application:

```bash
make run
```

#### Production Mode

1. Build the Docker image:

```bash
make docker-build
```

2. Run with Docker Compose:

```bash
docker-compose up -d
```

### Configuration

The application is configured through the `config.yaml` file. See [config.yaml](config.yaml) for all available options.

Environment variables can also be used to override configuration values:

- `DATABASE_URL` - Database connection string
- `REDIS_URL` - Redis connection string
- `PORT` - HTTP port to listen on

### API Endpoints

- `GET /health` - Health check endpoint
- `GET /version` - Version information
- `POST /v1/chat/completions` - OpenAI-compatible chat completions endpoint
- `GET /sentinel/admin/policies` - Get policies
- `POST /sentinel/admin/policies` - Create policy
- `GET /sentinel/admin/tenants` - Get tenants
- `POST /sentinel/admin/tenants` - Create tenant
- `GET /sentinel/admin/logs` - Get logs

### Development

#### Running Tests

```bash
make test
```

#### Running Tests with Coverage

```bash
make test-coverage
```

#### Linting

```bash
make lint
```

#### Generating Documentation

```bash
make docs
```

### Project Structure

```
.
├── main.go                 # Application entry point
├── config.yaml             # Configuration file
├── Dockerfile              # Docker image definition
├── docker-compose.yml      # Docker Compose configuration
├── Makefile                # Build and development tasks
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── docs/                   # Documentation
├── sentinel/               # Core components
│   ├── ciphermesh/         # Data detection and redaction
│   ├── sentinel/           # Security detection and response
│   ├── policy/             # Policy engine
│   ├── crypto/             # Cryptography and vault
│   ├── admin/              # Admin console and APIs
│   └── sdk/                # Language SDKs
└── adapters/               # LLM provider adapters
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests to ensure nothing is broken
5. Submit a pull request

### License

[License information to be added]
