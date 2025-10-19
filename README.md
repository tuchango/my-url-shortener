# URL Shortener

A simple and efficient URL shortening service.

## Features

- Shorten long URLs to compact aliases
- Redirect users from short URLs to original URLs
- RESTful API for URL management

## Getting Started

### Prerequisites

- Go 1.25 or higher

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/tuchango/my-url-shortener.git
   cd my-url-shortener
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Generate mocks:
   ```bash
   mockery
   ```

4. Set up configuration:
   - Edit `config/local.yaml` with your configuration

5. Run the application:
   ```bash
   go run cmd/my-url-shortener/main.go
   ```

### Usage

#### Shorten a URL
```bash
curl -X POST http://localhost:8080/url \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/very/long/url"}'
```

<!-- #### Redirect to original URL -->
<!-- Visit: `http://localhost:8080/{alias}` -->

## API Endpoints

- `POST /url` - Create a short URL
<!-- - `GET /{alias}` - Redirect to original URL -->

## Testing

Run tests:
```bash
go test ./...
```

## License

MIT License