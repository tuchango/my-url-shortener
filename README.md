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

## API Endpoints

### URL Management

#### Create Short URL
```http
POST /url
Content-Type: application/json

{
  "url": "https://example.com/very/long/url",
  "alias": "optional-custom-alias"  // Optional
}
```

**Response:**
```json
{
  "alias": "abc123",
  "url": "https://example.com/very/long/url",
  "short_url": "http://localhost:8080/abc123"
}
```

#### Get URL Info
```http
GET /url/{alias}
```

**Response:**
```json
{
  "alias": "abc123",
  "url": "https://example.com/very/long/url",
  "short_url": "http://localhost:8080/abc123",
  "created_at": "2024-01-01T10:00:00Z"
}
```

#### Redirect to Original URL
```http
GET /{alias}
```

**Response:** 302 Redirect to original URL

#### Delete URL
```http
DELETE /url/{alias}
```

**Response:** 204 No Content

### Health Check

#### Application Status
```http
GET /ping
```

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T10:00:00Z"
}
```

## Usage Examples

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