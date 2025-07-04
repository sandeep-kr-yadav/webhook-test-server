# Quick Start Guide

## 🚀 Get Started in 3 Steps

### 1. Start the Server
```bash
./start.sh
```
Or use make:
```bash
make run
```

### 2. Open the Web UI
Open your browser and go to: **http://localhost:8080**

### 3. Run Tests
In a new terminal:
```bash
cd test-files
./test-multipart-cases.sh
```

## 📁 Project Structure

```
webhook-test-env/
├── cmd/webhook-test-server.go    # Main server
├── static/webhook-ui.html        # Web UI
├── test-files/                   # Test scripts & sample files
├── docs/API.md                   # API documentation
├── start.sh                      # Quick start script
├── Makefile                      # Build commands
└── README.md                     # Full documentation
```

## 🧪 Test Cases

The test suite includes 8 comprehensive test cases:
1. JSON only
2. PDF only  
3. PNG only
4. PDF + JSON
5. PNG + JSON
6. PDF + PNG + JSON
7. CSV + JSON
8. Excel + JSON

## 🔧 Available Commands

```bash
make help      # Show all commands
make run       # Start server
make test      # Run test cases
make build     # Build binary
make clean     # Clean artifacts
make setup     # Complete setup
```

## 🌐 Access Points

- **Web UI**: http://localhost:8080
- **Webhook**: http://localhost:8080/webhook
- **Health**: http://localhost:8080/health
- **API**: http://localhost:8080/api/requests

## 📝 Features

✅ Real-time WebSocket monitoring  
✅ File upload support (PDF, PNG, CSV, Excel)  
✅ Multipart form data handling  
✅ Beautiful web UI  
✅ Comprehensive test cases  
✅ File download functionality  
✅ Detailed logging  
✅ API documentation  

## 🐛 Troubleshooting

- **Port in use**: Change port in `cmd/webhook-test-server.go`
- **File uploads fail**: Check file permissions
- **WebSocket issues**: Check browser console
- **Logs**: Check `webhook-server.log`

## 📚 Documentation

- [Full README](README.md)
- [API Documentation](docs/API.md)
- [Makefile Commands](Makefile)

---

**Ready to test webhooks! 🎉** 