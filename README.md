# Webhook Test Environment

A comprehensive webhook testing environment with real-time monitoring, file upload support, and comprehensive test cases.

## Features

- **Real-time WebSocket monitoring** of webhook requests
- **File upload support** for PDF, PNG, CSV, Excel files
- **Multipart form data** handling with JSON metadata
- **Beautiful web UI** for request monitoring
- **Comprehensive test cases** for all combinations
- **File download functionality** for uploaded files
- **Detailed logging** for debugging
- **Delete All functionality** to clear all requests from server memory
- **Test Webhook Endpoint** button to create dummy requests for testing

## Live Demo

🚀 **Deployed Version Available**: [https://webhook-test-server-263n.onrender.com/ui](https://webhook-test-server-263n.onrender.com/ui)

You can test the webhook environment directly in your browser without setting up anything locally. The deployed version includes all the latest features including the new Delete All and Test Webhook Endpoint buttons.

## Project Structure

```
webhook-test-env/
├── cmd/
│   └── webhook-test-server.go    # Main server application
├── static/
│   └── webhook-ui.html           # Web UI for monitoring
├── test-files/
│   ├── test-multipart-cases.sh   # Test script
│   ├── sample.pdf                # Sample PDF file
│   ├── sample.png                # Sample PNG file
│   ├── sample.csv                # Sample CSV file
│   └── sample.xlsx               # Sample Excel file
├── docs/
│   └── API.md                    # API documentation
├── go.mod                        # Go module file
├── Makefile                      # Build and run commands
└── README.md                     # This file
```

## Quick Start

### Prerequisites

- Go 1.21 or later
- curl (for testing)

### Installation

1. **Clone or download the project**
2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

### Running the Server

```bash
# Run the webhook test server
go run cmd/webhook-test-server.go
```

The server will start on `http://localhost:8080`

**Note**: The web UI is now available directly at the root URL `http://localhost:8080`

### Access Points

#### Local Development
- **Web UI**: http://localhost:8080
- **Webhook Endpoint**: http://localhost:8080/webhook
- **ThoughtSpot Endpoint**: http://localhost:8080/webhook/thoughtspot
- **Health Check**: http://localhost:8080/health

#### Deployed Version
- **Web UI**: https://webhook-test-server-263n.onrender.com
- **Webhook Endpoint**: https://webhook-test-server-263n.onrender.com/webhook
- **ThoughtSpot Endpoint**: https://webhook-test-server-263n.onrender.com/webhook/thoughtspot
- **Health Check**: https://webhook-test-server-263n.onrender.com/health

## Testing

### Run All Test Cases

```bash
cd test-files
./test-multipart-cases.sh
```

### Test Cases Included

1. **JSON only** - Basic JSON data test
2. **PDF only** - Single PDF file upload
3. **PNG only** - Single PNG file upload
4. **PDF + JSON** - PDF file with JSON metadata
5. **PNG + JSON** - PNG file with JSON metadata
6. **PDF + PNG + JSON** - Multiple files with JSON
7. **CSV + JSON** - CSV file with JSON metadata
8. **Excel + JSON** - Excel file with JSON metadata

### Manual Testing

You can also test manually using curl:

```bash
# Test with JSON only
curl -X POST http://localhost:8080/webhook \
  -F "json_data={\"event\":\"test\",\"data\":{\"message\":\"Hello World\"}}" \
  -F "description=Test request"

# Test with file upload
curl -X POST http://localhost:8080/webhook \
  -F "file_pdf=@test-files/sample.pdf" \
  -F "description=File upload test"
```

## Web UI Features

The web UI provides:

- **Real-time request monitoring** via WebSocket
- **Request details** including headers, body, and files
- **File download links** for uploaded files
- **Statistics** (total requests, GET/POST counts, multipart requests)
- **Request history** with expandable details
- **Beautiful, responsive design**

## API Endpoints

### POST /webhook
Main webhook endpoint that accepts:
- JSON data in form fields
- File uploads (PDF, PNG, CSV, Excel)
- Multipart form data

### GET /
Serves the web UI for monitoring requests

### GET /api/requests
Returns JSON with all received requests

### GET /download/{filename}
Downloads uploaded files

### WebSocket /ws
Real-time updates for the web UI

## Configuration

The server runs on port 8080 by default. To change the port, modify the `port` variable in `cmd/webhook-test-server.go`.

## Logging

The server logs to both console and `webhook-server.log` file for debugging purposes.

## File Support

Supported file types:
- **PDF** (.pdf) - application/pdf
- **PNG** (.png) - image/png
- **CSV** (.csv) - text/csv
- **Excel** (.xlsx, .xls) - application/vnd.openxmlformats-officedocument.spreadsheetml.sheet

## Development

### Building

```bash
go build -o webhook-server cmd/webhook-test-server.go
```

### Running Tests

```bash
# Run the test suite
cd test-files
./test-multipart-cases.sh
```

## Troubleshooting

### Common Issues

1. **Port already in use**: Change the port in the server code
2. **File uploads not working**: Check file permissions and ensure files exist
3. **WebSocket connection issues**: Check browser console for errors

### Debugging

- Check `webhook-server.log` for detailed server logs
- Use browser developer tools to inspect WebSocket connections
- Verify file paths in test scripts

## License

This project is provided as-is for testing and development purposes.

## Contributing

Feel free to submit issues and enhancement requests!
