package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"

	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

type WebhookRequest struct {
	ID          string              `json:"id"`
	Timestamp   time.Time           `json:"timestamp"`
	Method      string              `json:"method"`
	Headers     map[string][]string `json:"headers"`
	Body        interface{}         `json:"body,omitempty"`
	URL         string              `json:"url"`
	Query       map[string][]string `json:"query"`
	Files       []FileInfo          `json:"files,omitempty"`
	RemoteAddr  string              `json:"remoteAddr"`
	ContentType string              `json:"contentType"`
}

type FileInfo struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	Content     string `json:"content,omitempty"`
	DownloadURL string `json:"downloadURL,omitempty"`
}

type WebhookResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

// ThoughtSpot webhook response structure
type ThoughtSpotWebhookData struct {
	Users []struct {
		DisplayName string `json:"displayName"`
		Email       string `json:"email"`
		UserID      string `json:"userId"`
	} `json:"users"`
	SchemaVersion                      string `json:"schemaVersion"`
	SchemaType                         string `json:"schemaType"`
	NotificationType                   string `json:"notificationType"`
	ScheduledReportWebhookNotification struct {
		ReportID     string `json:"reportId"`
		ReportName   string `json:"reportName"`
		ScheduleInfo struct {
			ScheduleID     string `json:"scheduleId"`
			ScheduleString string `json:"scheduleString"`
			NextRunTime    string `json:"nextRunTime"`
			Timezone       string `json:"timezone"`
		} `json:"scheduleInfo"`
		DeliveryInfo struct {
			DeliveryID     string `json:"deliveryId"`
			DeliveryTime   string `json:"deliveryTime"`
			DeliveryStatus string `json:"deliveryStatus"`
			RecipientCount int    `json:"recipientCount"`
		} `json:"deliveryInfo"`
		ReportMetadata struct {
			PinboardID   string `json:"pinboardId"`
			PinboardName string `json:"pinboardName"`
			ReportURL    string `json:"reportUrl"`
		} `json:"reportMetadata"`
		Attachments []struct {
			AttachmentID string `json:"attachmentId"`
			FileName     string `json:"fileName"`
			FileSize     int    `json:"fileSize"`
			ContentType  string `json:"contentType"`
			Disposition  string `json:"disposition"`
			Checksum     string `json:"checksum"`
			PartNumber   int    `json:"partNumber"`
		} `json:"attachments"`
	} `json:"scheduledReportWebhookNotification"`
}

// Global state for storing requests and WebSocket connections
var (
	requests       []WebhookRequest
	requestsMux    sync.RWMutex
	clients        = make(map[*websocket.Conn]*sync.Mutex) // Store client with its own mutex
	clientsMux     sync.RWMutex
	fileStorage    = make(map[string][]byte) // Store file content by filename
	fileStorageMux sync.RWMutex
	upgrader       = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}

	// Embedded HTML UI
	webhookUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Webhook Test Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; }
        .header h1 { color: #333; margin-bottom: 10px; }
        .status { text-align: center; margin: 20px 0; padding: 15px; border-radius: 5px; }
        .status.connected { background: #d4edda; color: #155724; }
        .status.disconnected { background: #f8d7da; color: #721c24; }
        .requests { margin-top: 30px; }
        .request { border: 1px solid #ddd; margin: 10px 0; padding: 15px; border-radius: 5px; }
        .request h3 { margin: 0 0 10px 0; color: #333; }
        .request pre { background: #f8f9fa; padding: 10px; border-radius: 3px; overflow-x: auto; }
        .btn { padding: 10px 20px; margin: 5px; border: none; border-radius: 5px; cursor: pointer; }
        .btn-primary { background: #007bff; color: white; }
        .btn-success { background: #28a745; color: white; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Webhook Test Server</h1>
            <p>Real-time webhook request monitoring</p>
        </div>
        
        <div id="status" class="status disconnected">
            <strong>Status:</strong> <span id="statusText">Disconnected</span>
        </div>
        
        <div style="text-align: center; margin: 20px 0;">
            <button class="btn btn-primary" onclick="clearRequests()">Clear Requests</button>
            <button class="btn btn-success" onclick="testWebhook()">Test Webhook</button>
        </div>
        
        <div class="requests">
            <h2>📨 Recent Requests</h2>
            <div id="requestsList"></div>
        </div>
    </div>

    <script>
        let ws = null;
        let requests = [];
        
        function connect() {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = protocol + '//' + window.location.host + '/ws';
            
            ws = new WebSocket(wsUrl);
            
            ws.onopen = function() {
                document.getElementById('status').className = 'status connected';
                document.getElementById('statusText').textContent = 'Connected';
            };
            
            ws.onclose = function() {
                document.getElementById('status').className = 'status disconnected';
                document.getElementById('statusText').textContent = 'Disconnected';
                setTimeout(connect, 3000);
            };
            
            ws.onmessage = function(event) {
                const request = JSON.parse(event.data);
                requests.unshift(request);
                if (requests.length > 50) requests = requests.slice(0, 50);
                updateRequestsList();
            };
        }
        
        function updateRequestsList() {
            const container = document.getElementById('requestsList');
            container.innerHTML = requests.map((req, index) => 
                '<div class="request">' +
                '<h3>Request #' + (index + 1) + ' - ' + req.method + ' ' + req.url + '</h3>' +
                '<p><strong>Time:</strong> ' + new Date(req.timestamp).toLocaleString() + '</p>' +
                '<p><strong>Remote:</strong> ' + req.remoteAddr + '</p>' +
                '<pre>' + JSON.stringify(req, null, 2) + '</pre>' +
                '</div>'
            ).join('');
        }
        
        function clearRequests() {
            requests = [];
            updateRequestsList();
        }
        
        function testWebhook() {
            fetch('/webhook', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ test: 'data', message: 'Hello from UI!' })
            });
        }
        
        connect();
    </script>
</body>
</html>`
)

// Helper function to get keys from a map
func getKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Helper function to get keys from file headers map
func getFileKeys(m map[string][]*multipart.FileHeader) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func main() {
	// Set up logging to console only for containerized environments
	log.SetOutput(os.Stdout)

	log.Printf("=== Webhook Server Starting ===")
	log.Printf("Logging to console")

	// Create a new mux to handle routing properly
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/ui", handleUI)
	mux.HandleFunc("/webhook", handleWebhook)
	mux.HandleFunc("/webhook/thoughtspot", handleThoughtSpotWebhook)
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/requests", handleAPIRequests)
	mux.HandleFunc("/ws", handleWebSocket)
	mux.HandleFunc("/download/", handleFileDownload)
	mux.HandleFunc("/test", handleTest)
	mux.HandleFunc("/ping", handlePing)

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port

	// Also try to bind to all interfaces explicitly
	log.Printf("Attempting to bind to port %s on all interfaces", port)

	log.Printf("Starting webhook test server on port %s", port)
	log.Printf("Webhook endpoint: http://0.0.0.0%s/webhook", port)
	log.Printf("ThoughtSpot webhook endpoint: http://0.0.0.0%s/webhook/thoughtspot", port)
	log.Printf("Web UI: http://0.0.0.0%s", port)
	log.Printf("Health check: http://0.0.0.0%s/health", port)

	log.Printf("Server listening on %s", port)
	log.Printf("Ready to accept connections...")

	// Try to bind to all interfaces explicitly
	server := &http.Server{
		Addr:    "0.0.0.0" + port,
		Handler: mux,
	}

	log.Printf("Starting HTTP server on 0.0.0.0%s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func handleUI(w http.ResponseWriter, r *http.Request) {
	log.Printf("UI request received: %s", r.URL.Path)
	if r.URL.Path == "/ui" {
		log.Printf("Serving embedded webhook UI")
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("X-Go-App", "webhook-server")
		w.Write([]byte(webhookUIHTML))
		return
	}
	log.Printf("Path not found: %s", r.URL.Path)
	http.NotFound(w, r)
}

func handleFileDownload(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/download/")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	fileStorageMux.RLock()
	fileContent, exists := fileStorage[filename]
	fileStorageMux.RUnlock()

	if !exists {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set appropriate headers for download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", getContentType(filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileContent)))

	w.Write(fileContent)
}

func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".csv":
		return "text/csv"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".json":
		return "application/json"
	case ".txt":
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers for WebSocket
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// Register client
	clientsMux.Lock()
	clients[conn] = &sync.Mutex{}
	clientsMux.Unlock()

	// Send existing requests to new client
	requestsMux.RLock()
	existingRequests := make([]WebhookRequest, len(requests))
	copy(existingRequests, requests)
	requestsMux.RUnlock()

	clientMux := clients[conn]
	for _, request := range existingRequests {
		clientMux.Lock()
		if err := conn.WriteJSON(request); err != nil {
			log.Printf("Error sending existing request to client: %v", err)
			clientMux.Unlock()
			break
		}
		clientMux.Unlock()
	}

	// Keep connection alive and handle disconnection
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			clientsMux.Lock()
			delete(clients, conn)
			clientsMux.Unlock()
			break
		}
	}
}

func handleAPIRequests(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	requestsMux.RLock()
	defer requestsMux.RUnlock()

	response := map[string]interface{}{
		"requests": requests,
		"count":    len(requests),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func broadcastRequest(request WebhookRequest) {
	log.Printf("Broadcasting request - Body: %+v", request.Body)
	log.Printf("Broadcasting request - Files: %+v", request.Files)

	clientsMux.RLock()
	// Create a copy of clients to avoid holding the lock during writes
	clientsCopy := make(map[*websocket.Conn]*sync.Mutex)
	for client, clientMux := range clients {
		clientsCopy[client] = clientMux
	}
	clientsMux.RUnlock()

	for client, clientMux := range clientsCopy {
		clientMux.Lock()
		if err := client.WriteJSON(request); err != nil {
			log.Printf("Error broadcasting to client: %v", err)
			client.Close()
			clientMux.Unlock()

			// Remove client from the main map
			clientsMux.Lock()
			delete(clients, client)
			clientsMux.Unlock()
		} else {
			clientMux.Unlock()
		}
	}
}

func addRequest(request WebhookRequest) {
	requestsMux.Lock()
	defer requestsMux.Unlock()

	// Add request to the beginning of the slice
	requests = append([]WebhookRequest{request}, requests...)

	// Keep only last 100 requests
	if len(requests) > 100 {
		requests = requests[:100]
	}

	// Broadcast to WebSocket clients
	go broadcastRequest(request)
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())

	log.Printf("=== Webhook Request Received ===")
	log.Printf("Request ID: %s", requestID)
	log.Printf("Method: %s", r.Method)
	log.Printf("URL: %s", r.URL.String())
	log.Printf("Remote Address: %s", r.RemoteAddr)
	log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))

	// Create webhook request record
	webhookReq := WebhookRequest{
		ID:          requestID,
		Timestamp:   time.Now(),
		Method:      r.Method,
		Headers:     r.Header,
		URL:         r.URL.String(),
		Query:       r.URL.Query(),
		RemoteAddr:  r.RemoteAddr,
		ContentType: r.Header.Get("Content-Type"),
	}

	// Log headers
	log.Printf("Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("  %s: %s", name, value)
		}
	}

	// Handle multipart form data
	contentType := r.Header.Get("Content-Type")
	if contentType != "" && len(contentType) > 20 && strings.HasPrefix(contentType, "multipart/form-data") {
		log.Printf("Processing multipart/form-data request...")
		log.Printf("Content-Type header: %s", r.Header.Get("Content-Type"))

		// Parse multipart form
		err := r.ParseMultipartForm(32 << 20) // 32 MB max memory
		if err != nil {
			log.Printf("Error parsing multipart form: %v", err)
			http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
			return
		}

		log.Printf("MultipartForm.Value keys: %v", getKeys(r.MultipartForm.Value))
		log.Printf("MultipartForm.File keys: %v", getFileKeys(r.MultipartForm.File))

		// Process form values (JSON content)
		formData := make(map[string]interface{})
		log.Printf("Form Values:")
		for key, values := range r.MultipartForm.Value {
			if len(values) == 1 {
				formData[key] = values[0]
				log.Printf("  %s: %s", key, values[0])
			} else {
				formData[key] = values
				log.Printf("  %s: %v", key, values)
			}
		}

		// Add file field names to form data for UI matching
		for fieldName := range r.MultipartForm.File {
			formData[fieldName] = fmt.Sprintf("FILE_UPLOADED_%s", fieldName)
			log.Printf("  %s: FILE_UPLOADED_%s", fieldName, fieldName)
		}

		webhookReq.Body = formData

		// Process uploaded files - ensure only one file type per request
		var files []FileInfo
		log.Printf("Uploaded Files:")

		// Track file types to ensure only one type per request
		var firstContentType string
		var hasFiles bool

		for fieldName, fileHeaders := range r.MultipartForm.File {
			for _, fileHeader := range fileHeaders {
				log.Printf("  Field: %s", fieldName)
				log.Printf("    Filename: %s", fileHeader.Filename)
				log.Printf("    Content-Type: %s", fileHeader.Header.Get("Content-Type"))
				log.Printf("    Size: %d bytes", fileHeader.Size)

				contentType := fileHeader.Header.Get("Content-Type")

				// Check if this is the first file or same type as first file
				if !hasFiles {
					firstContentType = contentType
					hasFiles = true
				} else if contentType != firstContentType {
					log.Printf("    Skipping file - different type already exists (first: %s, current: %s)", firstContentType, contentType)
					continue
				}

				fileInfo := FileInfo{
					Filename:    fileHeader.Filename,
					ContentType: contentType,
					Size:        fileHeader.Size,
					DownloadURL: fmt.Sprintf("/download/%s", fileHeader.Filename),
				}

				// Read file content
				file, err := fileHeader.Open()
				if err != nil {
					log.Printf("    Error opening file: %v", err)
					continue
				}
				defer file.Close()

				// Read entire file content
				fileContent, err := io.ReadAll(file)
				if err != nil {
					log.Printf("    Error reading file: %v", err)
					continue
				}

				// Store file content for download
				fileStorageMux.Lock()
				fileStorage[fileHeader.Filename] = fileContent
				fileStorageMux.Unlock()

				log.Printf("    File processed successfully")
				files = append(files, fileInfo)
			}
		}
		webhookReq.Files = files
		log.Printf("Multipart processing completed - Files count: %d", len(files))
	} else {
		log.Printf("Not a multipart request, processing as regular body")
		log.Printf("Content-Type was: '%s'", contentType)
		// Handle regular request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "Error reading body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if len(body) > 0 {
			log.Printf("Body: %s", string(body))

			// Try to parse as JSON for pretty printing
			var jsonBody interface{}
			if err := json.Unmarshal(body, &jsonBody); err == nil {
				log.Printf("Parsed JSON Body: %v", jsonBody)
				webhookReq.Body = jsonBody
			} else {
				webhookReq.Body = string(body)
			}
		}
	}

	// Log query parameters
	if len(r.URL.Query()) > 0 {
		log.Printf("Query Parameters:")
		for key, values := range r.URL.Query() {
			for _, value := range values {
				log.Printf("  %s: %s", key, value)
			}
		}
	}

	log.Printf("=== End Webhook Request ===")
	log.Printf("Final webhookReq.Body: %+v", webhookReq.Body)
	log.Printf("Final webhookReq.Files: %+v", webhookReq.Files)

	// Add request to storage and broadcast
	addRequest(webhookReq)

	// Send response
	response := WebhookResponse{
		Status:  "success",
		Message: "Webhook received successfully",
		Time:    time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	w.Write(responseJSON)
}

func handleThoughtSpotWebhook(w http.ResponseWriter, r *http.Request) {
	requestID := fmt.Sprintf("thoughtspot-%d", time.Now().UnixNano())

	log.Printf("=== ThoughtSpot Webhook Request Received ===")
	log.Printf("Request ID: %s", requestID)
	log.Printf("Method: %s", r.Method)
	log.Printf("URL: %s", r.URL.String())
	log.Printf("Remote Address: %s", r.RemoteAddr)

	// Create webhook request record
	webhookReq := WebhookRequest{
		ID:          requestID,
		Timestamp:   time.Now(),
		Method:      r.Method,
		Headers:     r.Header,
		URL:         r.URL.String(),
		Query:       r.URL.Query(),
		RemoteAddr:  r.RemoteAddr,
		ContentType: r.Header.Get("Content-Type"),
	}

	// Log headers
	log.Printf("Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("  %s: %s", name, value)
		}
	}

	// Log request body if any
	body, err := io.ReadAll(r.Body)
	if err == nil && len(body) > 0 {
		log.Printf("Request Body: %s", string(body))

		// Try to parse as JSON
		var jsonBody interface{}
		if err := json.Unmarshal(body, &jsonBody); err == nil {
			webhookReq.Body = jsonBody
		} else {
			webhookReq.Body = string(body)
		}
	}
	defer r.Body.Close()

	log.Printf("=== End ThoughtSpot Webhook Request ===")

	// Add request to storage and broadcast
	addRequest(webhookReq)

	// Create ThoughtSpot-style response
	now := time.Now()
	nextMonth := now.AddDate(0, 1, 0)

	webhookData := ThoughtSpotWebhookData{
		Users: []struct {
			DisplayName string `json:"displayName"`
			Email       string `json:"email"`
			UserID      string `json:"userId"`
		}{
			{
				DisplayName: "John Doe",
				Email:       "john.doe@thoughtspot.com",
				UserID:      "user-12345",
			},
		},
		SchemaVersion:    "v1",
		SchemaType:       "SCHEDULED_REPORT",
		NotificationType: "DELIVERY",
	}

	webhookData.ScheduledReportWebhookNotification.ReportID = "report-67890"
	webhookData.ScheduledReportWebhookNotification.ReportName = "Monthly Sales Dashboard"
	webhookData.ScheduledReportWebhookNotification.ScheduleInfo.ScheduleID = "schedule-abc123"
	webhookData.ScheduledReportWebhookNotification.ScheduleInfo.ScheduleString = "monthly on 1st day at 9:00 AM"
	webhookData.ScheduledReportWebhookNotification.ScheduleInfo.NextRunTime = nextMonth.Format(time.RFC3339)
	webhookData.ScheduledReportWebhookNotification.ScheduleInfo.Timezone = "America/New_York"
	webhookData.ScheduledReportWebhookNotification.DeliveryInfo.DeliveryID = "delivery-xyz789"
	webhookData.ScheduledReportWebhookNotification.DeliveryInfo.DeliveryTime = now.Format(time.RFC3339)
	webhookData.ScheduledReportWebhookNotification.DeliveryInfo.DeliveryStatus = "SUCCESS"
	webhookData.ScheduledReportWebhookNotification.DeliveryInfo.RecipientCount = 5
	webhookData.ScheduledReportWebhookNotification.ReportMetadata.PinboardID = "22a8f618-0b4f-4401-92db-ba029ee13486"
	webhookData.ScheduledReportWebhookNotification.ReportMetadata.PinboardName = "Sales Performance Dashboard"
	webhookData.ScheduledReportWebhookNotification.ReportMetadata.ReportURL = "http://thoughtspot.company.com/?utm_source=scheduled_report&utm_medium=webhook/#/pinboard/22a8f618-0b4f-4401-92db-ba029ee13486"

	webhookData.ScheduledReportWebhookNotification.Attachments = []struct {
		AttachmentID string `json:"attachmentId"`
		FileName     string `json:"fileName"`
		FileSize     int    `json:"fileSize"`
		ContentType  string `json:"contentType"`
		Disposition  string `json:"disposition"`
		Checksum     string `json:"checksum"`
		PartNumber   int    `json:"partNumber"`
	}{
		{
			AttachmentID: "att-001",
			FileName:     fmt.Sprintf("Monthly_Sales_Dashboard_%s.pdf", now.Format("2006-01")),
			FileSize:     2048576,
			ContentType:  "application/pdf",
			Disposition:  "attachment",
			Checksum:     "sha256:abc123def456ghi789jkl012mno345pqr678stu901vwx234yz",
			PartNumber:   2,
		},
	}

	// Create multipart/mixed response
	boundary := "----WebKitFormBoundary7MA4YWxkTrZu0gW"
	w.Header().Set("Content-Type", fmt.Sprintf("multipart/mixed; boundary=\"%s\"", boundary))
	w.Header().Set("Server", "COMS-Webhook/1.0")
	w.Header().Set("X-Webhook-Signature", "sha256=abc123def456ghi789jkl012mno345pqr678stu901vwx234yz")
	w.WriteHeader(http.StatusOK)

	// Write multipart response
	writer := multipart.NewWriter(w)
	writer.SetBoundary(boundary)

	// Part 1: JSON metadata
	jsonPart, err := writer.CreatePart(map[string][]string{
		"Content-Type":        {"application/json"},
		"Content-Disposition": {"inline"},
	})
	if err != nil {
		log.Printf("Error creating JSON part: %v", err)
		return
	}

	jsonData := map[string]interface{}{
		"data": webhookData,
	}

	jsonBytes, _ := json.MarshalIndent(jsonData, "", "  ")
	jsonPart.Write(jsonBytes)

	// Part 2: PDF attachment
	pdfPart, err := writer.CreatePart(map[string][]string{
		"Content-Type":              {"application/pdf"},
		"Content-Disposition":       {fmt.Sprintf("attachment; filename=\"%s\"", webhookData.ScheduledReportWebhookNotification.Attachments[0].FileName)},
		"Content-Transfer-Encoding": {"base64"},
		"Content-Length":            {"2048576"},
	})
	if err != nil {
		log.Printf("Error creating PDF part: %v", err)
		return
	}

	// Create a sample PDF content (base64 encoded)
	samplePDF := `%PDF-1.4
1 0 obj
<<
/Type /Catalog
/Pages 2 0 R
>>
endobj

2 0 obj
<<
/Type /Pages
/Kids [3 0 R]
/Count 1
>>
endobj

3 0 obj
<<
/Type /Page
/MediaBox [0 0 595.44 841.92]
/Resources <<
/ExtGState <<
/GS5 7 0 R
/GS8 8 0 R
>>
/Font <<
/F1 9 0 R
/F2 13 0 R
/F3 17 0 R
/F4 20 0 R
/F5 24 0 R
>>
/XObject <<
/Image17 27 0 R
>>
/ProcSet [/PDF /Text /ImageB /ImageC /ImageI]
>>
/Contents 28 0 R
/Group <<
/Type /Group
/S /Transparency
/CS /DeviceRGB
>>
/Tabs /S
/StructParents 0
/Parent 2 0 R
>>
endobj

4 0 obj
<<
/Footnote /Note
/Endnote /Note
/Textbox /Sect
/Header /Sect
/Footer /Sect
/InlineShape /Sect
/Annotation /Sect
/Artifact /Sect
/Workbook /Document
/Worksheet /Part
/Macrosheet /Part
/Chartsheet /Part
/Dialogsheet /Part
/Slide /Part
/Chart /Sect
/Diagram /Figure
>>
endobj

5 0 obj
<<
/Type /StructTreeRoot
/RoleMap 4 0 R
/K [29 0 R]
/ParentTree 106 0 R
/ParentTreeNextKey 2
>>
endobj

6 0 obj
<<
/Type /Pages
/Kids [5 0 R 6 0 R]
/Count 2
>>
endobj

xref
0 7
0000000000 65535 f
0000000015 00000 n
0000000071 00000 n
0000000127 00000 n
0000000183 00000 n
0000000239 00000 n
0000000295 00000 n
trailer
<<
/Size 7
/Root 1 0 R
>>
startxref
351
%%EOF`

	// Encode PDF content as base64
	encodedPDF := base64.StdEncoding.EncodeToString([]byte(samplePDF))
	pdfPart.Write([]byte(encodedPDF))

	writer.Close()
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	log.Printf("Health check request received")
	response := map[string]interface{}{
		"service":   "webhook-test-server",
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   "This is the actual Go application running!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Test endpoint request received")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("GO APPLICATION IS RUNNING! If you see this, your Go app is working correctly."))
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	log.Printf("Ping endpoint request received from %s", r.RemoteAddr)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Go-App", "webhook-server")
	w.Write([]byte("PONG - Your Go app is responding! This proves Railway is routing to your container."))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Printf("Root endpoint request received: %s", r.URL.Path)
	if r.URL.Path == "/" {
		// Return a simple text response to test if Railway routes to our app
		log.Printf("Serving root response")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("X-Go-App", "webhook-server")
		w.Header().Set("X-Request-Path", r.URL.Path)
		w.Write([]byte("GO APPLICATION ROOT - If you see this, Railway is routing to your Go app!"))
		return
	}
	http.NotFound(w, r)
}
