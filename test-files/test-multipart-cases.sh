#!/bin/bash

# Change to the directory where this script is located
cd "$(dirname "$0")"

# Function to detect if local server is running
detect_server_url() {
    # Try to connect to localhost:8080/health
    if curl -s --connect-timeout 3 http://localhost:8080/health > /dev/null 2>&1; then
        echo "http://localhost:8080"
    else
        echo "https://webhook-test-server-263n.onrender.com"
    fi
}

# Get the appropriate server URL
SERVER_URL=$(detect_server_url)
echo "=== Testing Multipart Webhook Cases ==="
if [[ "$SERVER_URL" == "http://localhost:8080" ]]; then
    echo "✅ Using LOCAL server: $SERVER_URL"
else
    echo "🌐 Using DEPLOYED server: $SERVER_URL"
fi
echo ""

echo "1. Testing multipart with JSON only..."
curl -X POST $SERVER_URL/webhook \
  -F "json_data={\"event\":\"test_event\",\"data\":{\"message\":\"JSON only test\",\"timestamp\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"}}" \
  -F "description=Test with JSON data only"

echo -e "\n\n=== Test case 1 completed ===\nCheck the web UI at $SERVER_URL to see the result.\n"

echo "2. Testing multipart with PDF only..."
curl -X POST $SERVER_URL/webhook \
  -F "file_pdf=@sample.pdf" \
  -F "description=Test with PDF file only"

echo -e "\n\n=== Test case 2 completed ===\nCheck the web UI at $SERVER_URL to see the result.\n"

echo "3. Testing multipart with PNG only..."
curl -X POST $SERVER_URL/webhook \
  -F "file_png=@sample.png" \
  -F "description=Test with PNG file only"

echo -e "\n\n=== Test case 3 completed ===\nCheck the web UI at $SERVER_URL to see the result.\n"

echo "4. Testing multipart with PDF + JSON..."
curl -X POST $SERVER_URL/webhook \
  -F "file_pdf=@sample.pdf" \
  -F "json_data={\"event\":\"pdf_json_test\",\"data\":{\"message\":\"PDF + JSON test\",\"timestamp\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",\"file_type\":\"pdf\"}}" \
  -F "description=Test with PDF file and JSON data"

echo -e "\n\n=== Test case 4 completed ===\nCheck the web UI at $SERVER_URL to see the result.\n"

echo "5. Testing multipart with PNG + JSON..."
curl -X POST $SERVER_URL/webhook \
  -F "file_png=@sample.png" \
  -F "json_data={\"event\":\"png_json_test\",\"data\":{\"message\":\"PNG + JSON test\",\"timestamp\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",\"file_type\":\"png\"}}" \
  -F "description=Test with PNG file and JSON data"

echo -e "\n\n=== Test case 5 completed ===\nCheck the web UI at $SERVER_URL to see the result.\n"

echo "6. Testing multipart with PDF + PNG + JSON..."
curl -X POST $SERVER_URL/webhook \
  -F "file_pdf=@sample.pdf" \
  -F "file_png=@sample.png" \
  -F "json_data={\"event\":\"pdf_png_json_test\",\"data\":{\"message\":\"PDF + PNG + JSON test\",\"timestamp\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",\"file_types\":[\"pdf\",\"png\"]}}" \
  -F "description=Test with PDF file, PNG file, and JSON data"

echo -e "\n\n=== Test case 6 completed ===\nCheck the web UI at $SERVER_URL to see the result.\n"

echo "7. Testing multipart with CSV + JSON..."
curl -X POST $SERVER_URL/webhook \
  -F "file_csv=@sample.csv" \
  -F "json_data={\"event\":\"csv_json_test\",\"data\":{\"message\":\"CSV + JSON test\",\"timestamp\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",\"file_type\":\"csv\"}}" \
  -F "description=Test with CSV file and JSON data"

echo -e "\n\n=== Test case 7 completed ===\nCheck the web UI at $SERVER_URL to see the result.\n"

echo "8. Testing multipart with Excel + JSON..."
curl -X POST $SERVER_URL/webhook \
  -F "file_excel=@sample.xlsx" \
  -F "json_data={\"event\":\"excel_json_test\",\"data\":{\"message\":\"Excel + JSON test\",\"timestamp\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",\"file_type\":\"xlsx\"}}" \
  -F "description=Test with Excel file and JSON data"

echo -e "\n\n=== Test case 8 completed ===\nCheck the web UI at $SERVER_URL to see the result.\n"

echo "=== All test cases completed ==="
echo "Check the web UI at $SERVER_URL to see all results."
echo "Note: The server is configured to only process one file type per request,"
echo "so in test case 6 (PDF + PNG + JSON), only one file type will be processed." 