# Open Contributions Ghana

Discover and contribute to Ghanaian open-source projects.

Open Contributions Ghana is a simple directory that lists Ghana-based GitHub projects that are actively looking for contributors. It helps developers find meaningful projects to work on and helps maintainers get visibility and contributors.

---

## What It Does

- Submit Ghanaian GitHub repositories
- Display project details with real-time GitHub API integration
- Browse projects by language
- View repository metadata (stars, description, last activity)
- Link directly to GitHub for contributions
- Dark-themed, responsive interface

Open Source Ghana is not a replacement for GitHub. It is a discovery layer built on top of GitHub.

---

## Why This Exists

It was difficult to discover Ghanaian open-source projects that are actively looking for contributors.

There are communities and meetups, but no simple, searchable directory focused on:

- Discoverability
- Contribution intent
- Local ecosystem growth

Open Source Ghana is an experiment to validate whether a Ghana-focused open-source discovery platform is useful.

---

## Tech Stack

**Backend:**
- Go (Golang) with standard library net/http
- SQLite database for project storage
- GitHub API integration for repository metadata

**Frontend:**
- HTML5
- Pure CSS (no frameworks)
- Vanilla JavaScript with async/await
- Quicksand font from Google Fonts

---

## How To Run Locally

1. Clone the repository

   ```bash
   git clone https://github.com/kwabsntim/Open-source-Ghana.git
   cd Open-source-Ghana
   ```

2. Install dependencies

   ```bash
   go mod tidy
   ```

3. Run the application

   ```bash
   go run cmd/server.go
   ```

4. Open your browser

   ```
   http://localhost:8080
   ```

The application will create a SQLite database (`my.db`) automatically on first run.

---

## Project Structure

```
Open source Ghana/
│
├── cmd/
│   └── server.go           # Main application entry point
│
├── internal/
│   ├── database.go         # Database initialization
│   ├── handlers.go         # HTTP request handlers
│   ├── interface.go        # Repository interface
│   ├── models.go           # Data structures
│   ├── repository.go       # Database operations
│   └── service.go          # Business logic
│
├── web/
│   ├── index.html          # Frontend HTML
│   ├── style.css           # Styling
│   └── app.js              # JavaScript logic
│
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

---

## How To Contribute

Contributions are welcome.

You can help by:

- Improving UI/UX
- Adding better filtering and search
- Improving GitHub API integration
- Adding tests
- Improving documentation
- Fixing bugs
- Adding category assignment logic
- Implementing authentication

### Steps:

1. Fork the repository

2. Create a new branch

   ```bash
   git checkout -b feature/your-feature-name
   ```

3. Make your changes

4. Commit clearly

   ```bash
   git commit -m "Add filtering by language"
   ```

5. Push and open a Pull Request

Please keep changes focused and well-documented.

---

## API Endpoints

**GET /api/projects**
- Returns all projects from the database
- Response: JSON array of project objects

**POST /api/projects**
- Accepts GitHub repository URL
- Fetches metadata from GitHub API
- Stores project in database
- Request body: `{"github_url": "https://github.com/owner/repo"}`

---

## Features Implemented

- GitHub URL parsing and validation
- Real-time repository metadata fetching
- SQLite persistence
- CORS-enabled API
- Modal-based project submission
- GitHub API preview before submission
- Responsive card layout
- Sticky navigation and sidebar
- Smooth scrolling
- Dark theme design

---

## Future Ideas

- User authentication
- Project categories and tagging
- Trending projects based on activity
- Activity-based sorting
- Contributor highlights
- Search functionality
- Multi-language support
- Admin dashboard
- Email notifications for new projects

---

## Database Schema

**Projects Table:**
- id (INTEGER, PRIMARY KEY)
- name (TEXT)
- description (TEXT)
- github_url (TEXT, UNIQUE)
- owner_name (TEXT)
- owner_avatar (TEXT)
- language (TEXT)
- stars (INTEGER)
- category (TEXT)
- created_at (DATETIME)

---

## License

This project is licensed under the MIT License. See the LICENSE file for details.

---

## Acknowledgments

Built as a community initiative to strengthen Ghana's open-source ecosystem.

If you're a Ghanaian developer or maintainer, consider adding your project to help others discover it.
