# URL Shortener

A simple and efficient URL shortening service.

## Features

- Shorten long URLs to compact aliases
- Redirect users from short URLs to original URLs
- RESTful API for URL management

## Getting Started

### Prerequisites

- Go 1.25+

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/tuchango/my-url-shortener.git
   ```

2. Navigate to the directory:
   ```bash
   cd my-url-shortener
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Generate mocks:
   ```bash
   mockery
   ```

5. Edit `config/local.yaml` with your configuration

6. Run the application:
   ```bash
   go run cmd/my-url-shortener/main.go
   ```

<!-- ### Environment Variables

You can set the environment variables in the .env file. Here are some important variables:

    HTTP_SERVER_PASSWORD -->

## Usage

## Endpoints

- `GET /{alias}`: Redirect to original URL
- `POST /url`: Create alias for URL
- `DELETE /url`: Delete alias for URL

## Authentication

POST and DELETE endpoints require Basic Authentication:
- username: `myusername`
- password: `mypassword`

## Usage Examples

#### Redirect to original URL
Visit: `http://localhost:8081/abc123`

#### Shorten a URL
```bash
curl -X POST http://localhost:8081/url \
  -H "Content-Type: application/json" \
  -u myusername:mypassword \
  -d '{"url": "https://example.com/very/long/url"}'
```

#### Delete a URL
```bash
curl -X DELETE http://localhost:8081/url/abc123 \
  -u myusername:mypassword
```

## Testing

Run tests:
```bash
go test ./...
```

## License

MIT License