# Webhook Test Environment API Documentation

## Overview

The Webhook Test Environment provides a comprehensive testing platform for webhook endpoints with real-time monitoring, file upload support, and detailed request logging.

## Base URL

```
http://localhost:8080
```

## Endpoints

### 1. Web UI

**GET /**  
Serves the main web interface for monitoring webhook requests.

**Response:** HTML page with real-time webhook monitoring interface.

---

### 2. Webhook Endpoint

**POST /webhook**  
Main webhook endpoint that accepts multipart form data with files and JSON metadata.

**Content-Type:** `multipart/form-data`

**Form Fields:**
- `json_data` (optional): JSON string containing event data
- `description` (optional): Human-readable description of the request
- `file_pdf` (optional): PDF file upload
- `file_png` (optional): PNG file upload
- `file_csv` (optional): CSV file upload
- `file_excel` (optional): Excel file upload

**Example Request:**
```bash
curl -X POST http://localhost:8080/webhook \
  -F "json_data={\"event\":\"test\",\"data\":{\"message\":\"Hello World\"}}" \
  -F "file_pdf=@sample.pdf" \
  -F "description=Test with PDF and JSON"
```

**Response:**
```json
{
  "status": "success",
  "message": "Webhook received successfully",
  "time": "2025-07-04T09:44:46.523201+05:30"
}
```

---

### 3. ThoughtSpot Webhook Endpoint

**POST /webhook/thoughtspot**  
Specialized endpoint for ThoughtSpot webhook testing.

**Content-Type:** `application/json`

**Example Request:**
```bash
curl -X POST http://localhost:8080/webhook/thoughtspot \
  -H "Content-Type: application/json" \
  -d '{
    "notificationType": "scheduledReportWebhookNotification",
    "scheduledReportWebhookNotification": {
      "reportName": "Test Report",
      "attachments": [
        {
          "fileName": "report.pdf",
          "fileSize": 1024,
          "contentType": "application/pdf"
        }
      ]
    }
  }'
```

---

### 4. Health Check

**GET /health**  
Health check endpoint for monitoring server status.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-07-04T09:44:46.523201+05:30",
  "uptime": "1h23m45s"
}
```

---

### 5. API Requests

**GET /api/requests**  
Returns all received webhook requests.

**Response:**
```json
{
  "requests": [
    {
      "id": "req-1751602486523034000",
      "timestamp": "2025-07-04T09:44:46.523201+05:30",
      "method": "POST",
      "headers": {...},
      "body": {...},
      "url": "/webhook",
      "query": {},
      "files": [
        {
          "filename": "sample.pdf",
          "content_type": "application/pdf",
          "size": 457554,
          "downloadURL": "/download/sample.pdf"
        }
      ],
      "remoteAddr": "127.0.0.1:12345",
      "contentType": "multipart/form-data; boundary=..."
    }
  ],
  "count": 1
}
```

---

### 6. File Download

**GET /download/{filename}**  
Downloads uploaded files.

**Parameters:**
- `filename`: Name of the file to download

**Response:** File content with appropriate headers.

**Example:**
```bash
curl -O http://localhost:8080/download/sample.pdf
```

---

### 7. WebSocket

**WebSocket /ws**  
Real-time updates for the web UI.

**Protocol:** WebSocket

**Messages:** JSON objects containing webhook request data.

**Example Message:**
```json
{
  "id": "req-1751602486523034000",
  "timestamp": "2025-07-04T09:44:46.523201+05:30",
  "method": "POST",
  "headers": {...},
  "body": {...},
  "url": "/webhook",
  "query": {},
  "files": [...],
  "remoteAddr": "127.0.0.1:12345",
  "contentType": "multipart/form-data; boundary=..."
}
```

## Data Structures

### WebhookRequest

```json
{
  "id": "string",
  "timestamp": "datetime",
  "method": "string",
  "headers": {
    "string": ["string"]
  },
  "body": "object",
  "url": "string",
  "query": {
    "string": ["string"]
  },
  "files": [
    {
      "filename": "string",
      "content_type": "string",
      "size": "number",
      "downloadURL": "string"
    }
  ],
  "remoteAddr": "string",
  "contentType": "string"
}
```

### FileInfo

```json
{
  "filename": "string",
  "content_type": "string",
  "size": "number",
  "downloadURL": "string"
}
```

### WebhookResponse

```json
{
  "status": "string",
  "message": "string",
  "time": "datetime"
}
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "Error parsing multipart form",
  "status": "error"
}
```

### 404 Not Found
```json
{
  "error": "File not found",
  "status": "error"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error",
  "status": "error"
}
```

## File Upload Limits

- **Maximum file size:** 32 MB per file
- **Maximum memory usage:** 32 MB for multipart parsing
- **Supported file types:** PDF, PNG, CSV, Excel

## Rate Limiting

Currently, there are no rate limits implemented. The server can handle multiple concurrent requests.

## CORS

The server includes CORS headers for WebSocket connections:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization`

## Logging

All requests are logged to both console and `webhook-server.log` file with detailed information including:
- Request headers
- Form data
- File uploads
- Processing steps
- Errors

## Testing

Use the provided test script `test-files/test-multipart-cases.sh` for comprehensive testing of all scenarios. 