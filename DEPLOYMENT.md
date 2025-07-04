# Deployment Guide

This guide will help you deploy your webhook test server to various free hosting platforms.

## Option 1: Railway (Recommended)

Railway offers a generous free tier and is perfect for Go applications.

### Steps:

1. **Sign up for Railway**
   - Go to [railway.app](https://railway.app)
   - Sign up with your GitHub account

2. **Deploy from GitHub**
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose your `webhook-test-server` repository
   - Railway will automatically detect the Dockerfile and deploy

3. **Access your application**
   - Railway will provide you with a URL like `https://your-app-name.railway.app`
   - Your webhook endpoint will be: `https://your-app-name.railway.app/webhook`
   - Web UI: `https://your-app-name.railway.app`

### Railway Configuration
The `railway.json` file is already configured for optimal deployment.

## Option 2: Render

Render is another excellent free option.

### Steps:

1. **Sign up for Render**
   - Go to [render.com](https://render.com)
   - Sign up with your GitHub account

2. **Create a new Web Service**
   - Click "New +" → "Web Service"
   - Connect your GitHub repository
   - Choose the `webhook-test-server` repository

3. **Configure the service**
   - **Name**: `webhook-test-server` (or your preferred name)
   - **Environment**: `Docker`
   - **Region**: Choose closest to you
   - **Branch**: `main`
   - **Build Command**: Leave empty (uses Dockerfile)
   - **Start Command**: Leave empty (uses Dockerfile)

4. **Deploy**
   - Click "Create Web Service"
   - Render will build and deploy your application

## Option 3: Fly.io

Fly.io offers a generous free tier with global deployment.

### Steps:

1. **Install Fly CLI**
   ```bash
   # macOS
   brew install flyctl
   
   # Or download from https://fly.io/docs/hands-on/install-flyctl/
   ```

2. **Sign up and login**
   ```bash
   fly auth signup
   # or
   fly auth login
   ```

3. **Deploy**
   ```bash
   fly launch
   ```
   - Follow the prompts
   - Choose a unique app name
   - Select a region
   - Deploy

## Option 4: Heroku

Heroku has a limited free tier but is still a good option.

### Steps:

1. **Install Heroku CLI**
   ```bash
   # macOS
   brew install heroku/brew/heroku
   ```

2. **Login and create app**
   ```bash
   heroku login
   heroku create your-app-name
   ```

3. **Deploy**
   ```bash
   git push heroku main
   ```

## Testing Your Deployed Application

Once deployed, test your webhook server:

### 1. Health Check
```bash
curl https://your-app-url/health
```

### 2. Web UI
Visit `https://your-app-url` in your browser

### 3. Webhook Test
```bash
curl -X POST https://your-app-url/webhook \
  -F "json_data={\"event\":\"test\",\"data\":{\"message\":\"Hello from deployed server!\"}}" \
  -F "description=Deployment test"
```

### 4. File Upload Test
```bash
curl -X POST https://your-app-url/webhook \
  -F "file_pdf=@test-files/sample.pdf" \
  -F "description=File upload test"
```

## Environment Variables

The application supports these environment variables:

- `PORT`: Server port (default: 8080)
- Railway/Render/Fly.io will automatically set this

## Monitoring

- **Health Check**: `https://your-app-url/health`
- **Web UI**: `https://your-app-url`
- **API**: `https://your-app-url/api/requests`

## Troubleshooting

### Common Issues:

1. **Build fails**: Check that all files are committed to GitHub
2. **Port issues**: The app automatically uses the `PORT` environment variable
3. **File uploads**: Ensure your hosting platform supports file uploads
4. **WebSocket**: Some platforms may require additional configuration for WebSocket support

### Logs:
- Railway: View logs in the Railway dashboard
- Render: View logs in the Render dashboard
- Fly.io: `fly logs`
- Heroku: `heroku logs --tail`

## Cost Considerations

All platforms mentioned offer free tiers:
- **Railway**: $5/month after free tier (very generous)
- **Render**: $7/month after free tier
- **Fly.io**: $1.94/month after free tier
- **Heroku**: $7/month after free tier

The free tiers should be sufficient for testing and development purposes. 