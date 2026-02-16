// Fetch and display projects from the API
async function loadProjects() {
    try {
        const response = await fetch('http://localhost:8080/api/projects');
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const projects = await response.json();
        displayProjects(projects);
    } catch (error) {
        console.error('Error fetching projects:', error);
        displayError();
    }
}

// Display projects in the cards wrapper
function displayProjects(projects) {
    const cardsWrapper = document.querySelector('.cards-wrapper');
    
    // Clear existing cards
    cardsWrapper.innerHTML = '';
    
    // Check if there are no projects
    if (!projects || projects.length === 0) {
        cardsWrapper.innerHTML = '<p style="text-align: center; color: #8b949e; padding: 40px;">No projects found. Be the first to add one!</p>';
        return;
    }
    
    // Create cards for each project
    projects.forEach(project => {
        const card = createProjectCard(project);
        cardsWrapper.appendChild(card);
    });
}

// Create a card element for a project
function createProjectCard(project) {
    const card = document.createElement('div');
    card.className = 'card';
    
    // Format time ago (you can enhance this)
    const timeAgo = getTimeAgo(project.created_at);
    
    card.innerHTML = `
        <div class="card-header">
            <div>
                <h3>
                    <a href="${project.github_url}" target="_blank" rel="noopener noreferrer">
                        ${project.owner_name} / ${project.name}
                    </a>
                </h3>
                <p>${project.description || 'No description available'}</p>
                <div class="meta">
                    ${project.language ? `lang: ${project.language}` : 'lang: Unknown'}
                    <span>⭐ ${project.stars}</span>
                    <span>last activity: ${timeAgo}</span>
                </div>
            </div>
        </div>
    `;
    
    return card;
}

// Simple time ago function
function getTimeAgo(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const seconds = Math.floor((now - date) / 1000);
    
    if (seconds < 60) return 'just now';
    if (seconds < 3600) return `${Math.floor(seconds / 60)} minutes ago`;
    if (seconds < 86400) return `${Math.floor(seconds / 3600)} hours ago`;
    if (seconds < 2592000) return `${Math.floor(seconds / 86400)} days ago`;
    if (seconds < 31536000) return `${Math.floor(seconds / 2592000)} months ago`;
    return `${Math.floor(seconds / 31536000)} years ago`;
}

// Display error message
function displayError() {
    const cardsWrapper = document.querySelector('.cards-wrapper');
    cardsWrapper.innerHTML = `
        <p style="text-align: center; color: #f85149; padding: 40px;">
            Failed to load projects. Please try again later.
        </p>
    `;
}

// Load projects when page loads
document.addEventListener('DOMContentLoaded', loadProjects);
