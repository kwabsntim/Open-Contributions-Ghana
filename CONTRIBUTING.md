# Contributing to Contributions Ghana

Thanks for your interest in contributing.

## Guidelines

- Keep changes small and focused
- Write clear commit messages
- Test your changes before submitting a PR
- Be respectful and constructive in discussions
- Follow the existing code style
- Update documentation when adding features

## Development Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/kwabsntim/Open-source-Ghana.git
   cd Open-source-Ghana
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Run the application:

   ```bash
   go run cmd/server.go
   ```

4. Access the application at `http://localhost:8080`

## Project Structure

- `cmd/` - Application entry points
- `internal/` - Core business logic (not importable by external projects)
- `web/` - Frontend files (HTML, CSS, JavaScript)

## Making Changes

### Backend (Go)

- Database operations go in `internal/repository.go`
- Business logic goes in `internal/service.go`
- HTTP handlers go in `internal/handlers.go`
- Data models go in `internal/models.go`

### Frontend

- HTML structure in `web/index.html`
- Styling in `web/style.css`
- JavaScript logic in `web/app.js`

### Testing

Before submitting:

1. Test the server runs without errors
2. Test the API endpoints with curl or browser
3. Test the frontend displays correctly
4. Check for console errors in browser DevTools

## Reporting Issues

Open an issue describing:

- The problem or feature request
- Expected behavior
- Steps to reproduce (for bugs)
- Screenshots (if applicable)

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to your fork (`git push origin feature/amazing-feature`)
5. Open a Pull Request with a clear description

## Code Style

- Go: Follow standard Go conventions (use `gofmt`)
- JavaScript: Use modern ES6+ syntax
- CSS: Keep selectors organized and commented
- Use meaningful variable and function names

## Questions?

Feel free to open an issue for questions or clarifications.
