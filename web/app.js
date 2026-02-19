// API Configuration - automatically uses production or local URL
const API_URL = (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1')
    ? 'http://localhost:8080' 
    : window.location.origin;

// Fetch and display projects from the API
async function loadProjects() {
    try {
        const response = await fetch(`${API_URL}/api/projects`);
        
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
    
    // Make card clickable
    card.addEventListener('click', () => {
        window.open(project.github_url, '_blank', 'noopener,noreferrer');
    });
    
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
                    <span style="display: inline-flex; align-items: center; gap: 4px;">
                        <svg style="width: 14px; height: 14px; color: #8b949e;" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M11.48 3.499a.562.562 0 011.04 0l2.125 5.111a.563.563 0 00.475.345l5.518.442c.499.04.701.663.321.988l-4.204 3.602a.563.563 0 00-.182.557l1.285 5.385a.562.562 0 01-.84.61l-4.725-2.885a.563.563 0 00-.586 0L6.982 20.54a.562.562 0 01-.84-.61l1.285-5.386a.562.562 0 00-.182-.557l-4.204-3.602a.563.563 0 01.321-.988l5.518-.442a.563.563 0 00.475-.345L11.48 3.5z"/>
                        </svg>
                        ${project.stars}
                    </span>
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

// Modal functionality
const modal = document.getElementById('addProjectModal');
const addProjectBtn = document.querySelector('.add-project');
const cancelBtn = document.getElementById('cancelButton');
const submitBtn = document.getElementById('submitButton');
const urlInput = document.getElementById('githubUrlInput');
const modalError = document.getElementById('modalError');
const modalPreview = document.getElementById('modalPreview');

let previewedProject = null;

// Show modal
addProjectBtn.addEventListener('click', () => {
    modal.classList.add('show');
});

// Hide modal
function closeModal() {
    modal.classList.remove('show');
    urlInput.value = '';
    modalError.classList.remove('show');
    modalPreview.classList.remove('show');
    modalPreview.innerHTML = '';
    submitBtn.disabled = true;
    previewedProject = null;
}

cancelBtn.addEventListener('click', closeModal);

// Close modal when clicking outside
modal.addEventListener('click', (e) => {
    if (e.target === modal) {
        closeModal();
    }
});

// Fetch GitHub preview when URL is entered
urlInput.addEventListener('input', async () => {
    const url = urlInput.value.trim();
    
    // Reset state
    modalError.classList.remove('show');
    modalPreview.classList.remove('show');
    submitBtn.disabled = true;
    previewedProject = null;
    
    if (!url) {
        return;
    }
    
    // Validate GitHub URL format
    const githubUrlPattern = /^https?:\/\/github\.com\/([^\/]+)\/([^\/]+)\/?$/;
    const match = url.match(githubUrlPattern);
    
    if (!match) {
        modalError.textContent = 'Please enter a valid GitHub repository URL (e.g., https://github.com/username/repository)';
        modalError.classList.add('show');
        return;
    }
    
    const owner = match[1];
    const repo = match[2];
    
    // Fetch from GitHub API
    try {
        const response = await fetch(`https://api.github.com/repos/${owner}/${repo}`);
        
        if (!response.ok) {
            if (response.status === 404) {
                modalError.textContent = 'Repository not found. Please check the URL.';
            } else {
                modalError.textContent = `GitHub API error: ${response.status}`;
            }
            modalError.classList.add('show');
            return;
        }
        
        const repoData = await response.json();
        
        // Store previewed project
        previewedProject = {
            github_url: url
        };
        
        // Show preview
        modalPreview.innerHTML = createProjectCard({
            name: repoData.name,
            owner_name: repoData.owner.login,
            description: repoData.description,
            language: repoData.language,
            stars: repoData.stargazers_count,
            created_at: repoData.created_at,
            github_url: url
        }).outerHTML;
        
        modalPreview.classList.add('show');
        submitBtn.disabled = false;
        
    } catch (error) {
        console.error('Error fetching repository:', error);
        modalError.textContent = 'Failed to fetch repository information. Please try again.';
        modalError.classList.add('show');
    }
});

// Submit project
submitBtn.addEventListener('click', async () => {
    if (!previewedProject) {
        return;
    }
    
    // Disable button during submission
    submitBtn.disabled = true;
    submitBtn.textContent = 'Adding...';
    
    try {
        const response = await fetch(`${API_URL}/api/projects`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(previewedProject)
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        // Success - close modal and reload projects
        closeModal();
        await loadProjects();
        
    } catch (error) {
        console.error('Error adding project:', error);
        modalError.textContent = 'Failed to add project. Please try again.';
        modalError.classList.add('show');
        submitBtn.disabled = false;
        submitBtn.textContent = 'Add Project';
    }
});
