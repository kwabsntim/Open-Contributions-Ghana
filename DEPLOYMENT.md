# Deploy Open Source Ghana to Render

## Prerequisites

1. Turso Database
   - Sign up at https://turso.tech
   - Create a new database
   - Get your database URL and auth token

2. GitHub Account for Render
   - Connect your repository to Render

## Setup Steps

### 1. Create Turso Database

```bash
# Install Turso CLI
curl -sSfL https://get.tur.so/install.sh | bash

# Create a new database
turso db create open-source-ghana

# Get database URL
turso db show open-source-ghana --url

# Create auth token
turso db tokens create open-source-ghana
```

### 2. Deploy to Render

1. Go to https://render.com and sign in with GitHub

2. Click "New +" and select "Web Service"

3. Connect your GitHub repository

4. Configure the service:
   - **Name**: open-source-ghana
   - **Environment**: Go
   - **Build Command**: `go build -o server cmd/server.go`
   - **Start Command**: `./server`
   - **Instance Type**: Free

5. Add Environment Variables:
   - Click "Advanced" → "Add Environment Variable"
   - Add the following:

   ```
   TURSO_DATABASE_URL=libsql://your-database.turso.io
   TURSO_AUTH_TOKEN=your-turso-auth-token
   PORT=8080
   GITHUB_TOKEN=your-github-token (optional)
   ```

6. Click "Create Web Service"

### 3. Update Frontend URL

Once deployed, Render will give you a URL like:
`https://open-source-ghana.onrender.com`

Update `web/app.js` to use your production URL:

```javascript
// Change from:
const response = await fetch('http://localhost:8080/api/projects');

// To:
const API_URL = window.location.hostname === 'localhost' 
  ? 'http://localhost:8080' 
  : 'https://open-source-ghana.onrender.com';

const response = await fetch(`${API_URL}/api/projects`);
```

### 4. Serve Static Files

Add static file serving to `cmd/server.go`:

```go
// Serve static files
fs := http.FileServer(http.Dir("./web"))
mux.Handle("/", fs)
```

## Environment Variables Reference

| Variable | Description | Required | Example |
|----------|-------------|----------|---------|
| TURSO_DATABASE_URL | Turso database URL | Yes | `libsql://db-name.turso.io` |
| TURSO_AUTH_TOKEN | Turso authentication token | Yes | `eyJ...` |
| PORT | Server port | No (defaults to 8080) | `8080` |
| GITHUB_TOKEN | GitHub personal access token | No (increases rate limit) | `ghp_...` |

## Local Development

1. Copy `.env.example` to `.env`:
   ```bash
   cp .env.example .env
   ```

2. Fill in your Turso credentials in `.env`

3. Run locally:
   ```bash
   go run cmd/server.go
   ```

## Troubleshooting

### Database Connection Issues
- Verify TURSO_DATABASE_URL and TURSO_AUTH_TOKEN are correct
- Check Turso dashboard for database status
- Ensure auth token hasn't expired

### Build Failures on Render
- Check build logs in Render dashboard
- Verify all dependencies in go.mod
- Ensure Go version compatibility

### CORS Issues
- Verify CORS headers in server.go
- Check browser console for specific errors
- Update allowed origins if needed

## Monitoring

- View logs in Render dashboard
- Monitor database usage in Turso dashboard
- Set up health check endpoint (optional)

## Custom Domain (Optional)

1. In Render dashboard, go to Settings
2. Add custom domain
3. Update DNS records as instructed
4. Update CORS settings if needed
