// AI Search vs Keyword Search Demo - Frontend JavaScript

class SearchApp {
    constructor() {
        this.currentMode = 'keyword';
        this.isSearching = false;
        this.init();
    }

    init() {
        this.bindEvents();
        this.updateModeDisplay();
        this.checkHealth();
    }

    bindEvents() {
        // Search button click
        document.getElementById('searchBtn').addEventListener('click', () => {
            this.performSearch();
        });

        // Enter key in search input
        document.getElementById('searchInput').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.performSearch();
            }
        });

        // Search mode toggle
        document.getElementById('searchModeToggle').addEventListener('change', (e) => {
            this.currentMode = e.target.checked ? 'vector' : 'keyword';
            this.updateModeDisplay();
        });

        // Real-time search (optional - uncomment for live search)
        // document.getElementById('searchInput').addEventListener('input', debounce(() => {
        //     this.performSearch();
        // }, 500));
    }

    updateModeDisplay() {
        const modeText = this.currentMode === 'keyword' ? 'Keyword Search' : 'AI Vector Search';
        const modeIcon = this.currentMode === 'keyword' ? 'bi-keyboard' : 'bi-brain';
        
        document.getElementById('currentMode').innerHTML = 
            `<i class="bi ${modeIcon}"></i> ${modeText}`;
    }

    async checkHealth() {
        try {
            const response = await fetch('/api/health');
            const health = await response.json();
            
            const statusElement = document.getElementById('search-status');
            if (health.status === 'healthy') {
                statusElement.textContent = 'Ready';
                statusElement.className = 'status-ready';
            } else {
                statusElement.textContent = 'Error';
                statusElement.className = 'status-error';
            }
        } catch (error) {
            console.error('Health check failed:', error);
            document.getElementById('search-status').textContent = 'Error';
            document.getElementById('search-status').className = 'status-error';
        }
    }

    async performSearch() {
        const query = document.getElementById('searchInput').value.trim();
        
        if (!query) {
            this.showMessage('Please enter a search query', 'warning');
            return;
        }

        if (this.isSearching) return;

        this.isSearching = true;
        this.showLoading();

        try {
            const response = await fetch('/api/search', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    query: query,
                    mode: this.currentMode,
                    limit: 10
                })
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            this.displayResults(data);

        } catch (error) {
            console.error('Search error:', error);
            this.showMessage('Search failed. Please try again.', 'error');
            this.hideLoading();
        } finally {
            this.isSearching = false;
        }
    }

    displayResults(data) {
        this.hideLoading();

        const resultsContainer = document.getElementById('searchResults');
        const resultsSection = document.getElementById('resultsSection');
        const noResultsSection = document.getElementById('noResultsSection');

        // Update result count and search time
        document.getElementById('resultCount').textContent = `${data.total_results} results`;
        document.getElementById('searchTime').textContent = `${data.search_time_ms.toFixed(1)}ms`;

        if (data.results.length === 0) {
            resultsSection.style.display = 'none';
            noResultsSection.style.display = 'block';
            return;
        }

        // Clear previous results
        resultsContainer.innerHTML = '';

        // Add new results
        data.results.forEach((result, index) => {
            const resultElement = this.createResultElement(result, index + 1);
            resultsContainer.appendChild(resultElement);
        });

        // Show results section
        resultsSection.style.display = 'block';
        noResultsSection.style.display = 'none';

        // Scroll to results
        resultsSection.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }

    createResultElement(result, rank) {
        const div = document.createElement('div');
        div.className = 'list-group-item';

        // Format score for display
        const scoreText = this.currentMode === 'keyword' 
            ? `BM-25: ${result.score.toFixed(3)}`
            : `Similarity: ${(result.score * 100).toFixed(1)}%`;

        // Create tags HTML
        const tagsHtml = result.tags.map(tag => 
            `<span class="badge">${tag}</span>`
        ).join('');

        // Use highlighted content if available, otherwise use regular content
        const content = result.highlighted_content || result.content;
        const truncatedContent = content.length > 300 
            ? content.substring(0, 300) + '...' 
            : content;

        div.innerHTML = `
            <div class="d-flex justify-content-between align-items-start mb-2">
                <h5 class="result-title mb-0">
                    <span class="badge bg-secondary me-2">#${rank}</span>
                    ${result.title}
                </h5>
                <div class="d-flex align-items-center">
                    <span class="badge bg-info category-badge me-2">${result.category}</span>
                    <span class="badge bg-success score-badge">${scoreText}</span>
                </div>
            </div>
            
            <div class="result-content">
                ${truncatedContent}
            </div>
            
            <div class="result-tags">
                ${tagsHtml}
            </div>
        `;

        return div;
    }

    showLoading() {
        document.getElementById('loadingSection').style.display = 'block';
        document.getElementById('resultsSection').style.display = 'none';
        document.getElementById('noResultsSection').style.display = 'none';
    }

    hideLoading() {
        document.getElementById('loadingSection').style.display = 'none';
    }

    showMessage(message, type = 'info') {
        // Create a simple toast notification
        const toast = document.createElement('div');
        toast.className = `alert alert-${type === 'error' ? 'danger' : type} alert-dismissible fade show position-fixed`;
        toast.style.cssText = 'top: 20px; right: 20px; z-index: 1050; min-width: 300px;';
        toast.innerHTML = `
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        `;
        
        document.body.appendChild(toast);
        
        // Auto-remove after 5 seconds
        setTimeout(() => {
            if (toast.parentNode) {
                toast.remove();
            }
        }, 5000);
    }
}

// Utility function for debouncing (used for real-time search)
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Initialize the app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.searchApp = new SearchApp();
});

// Add some example queries for quick testing
const exampleQueries = [
    'python programming',
    'machine learning algorithms',
    'web development',
    'How do I learn artificial intelligence?',
    'What is cloud computing?',
    'database management systems',
    'cybersecurity best practices',
    'mobile app development',
    'blockchain technology',
    'DevOps and CI/CD'
];

// Add example queries to the page (optional)
function addExampleQueries() {
    const searchInfo = document.querySelector('.text-center small');
    if (searchInfo) {
        const examplesDiv = document.createElement('div');
        examplesDiv.className = 'mt-2';
        examplesDiv.innerHTML = `
            <small class="text-muted">
                <strong>Try:</strong> 
                ${exampleQueries.slice(0, 3).map(q => 
                    `<a href="#" class="text-decoration-none example-query">${q}</a>`
                ).join(' • ')}
            </small>
        `;
        searchInfo.appendChild(examplesDiv);

        // Add click handlers for example queries
        document.querySelectorAll('.example-query').forEach(link => {
            link.addEventListener('click', (e) => {
                e.preventDefault();
                document.getElementById('searchInput').value = e.target.textContent;
                window.searchApp.performSearch();
            });
        });
    }
}

// Uncomment the line below to add example queries
// document.addEventListener('DOMContentLoaded', addExampleQueries); 