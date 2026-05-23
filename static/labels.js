const wait = (ms) => new Promise(resolve => setTimeout(resolve, ms));

// Add event listeners to filter labels
document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.fr-label').forEach(label => {
        label.addEventListener('click', toggleLabel);
    });
});


// Array to store selected label IDs
let selectedLabelIds = [];

// Function to toggle label state
function toggleLabel(event) {
    const labelElement = event.currentTarget;
    const labelId = labelElement.getAttribute('data-id');

    // Toggle 'selected' class
    labelElement.classList.toggle('selected');

    if (labelElement.classList.contains('selected')) {
        // If selected, add ID to array
        selectedLabelIds.push(labelId);
    } else {
        // If deselected, remove ID from array
        selectedLabelIds = selectedLabelIds.filter(id => id !== labelId);
    }
    
    // applyFilters(); 
}


// Main function to apply filters and fetch results
// Filter code:
// 0 - regular
// 1 - continue reading
// 2 - created
// 3 - favourite
async function applyFilters(code) {
    await wait(100); // Small delay to ensure label id updates before fetching

    const keywordInput = document.getElementById('keyword-input');
    const keyword = keywordInput ? keywordInput.value.trim() : '';
    const resultsContainer = getResultsContainer(code);
    
    if (!resultsContainer) {
        console.warn('applyFilters: no results container found, skipping render');
        return;
    }

    resultsContainer.innerHTML = '<p>Searching...</p>';
    
    let queryParams = new URLSearchParams();

    // 1. Keyword
    if (keyword) {
        queryParams.append('keyword', keyword); 
    }

    // 2. Selected labels (passed as comma-separated)
    if (selectedLabelIds.length > 0) {
        queryParams.append('labels', selectedLabelIds.join(','));
    }

    queryParams.append('code', code);

    const url = `/library/filter/?${queryParams.toString()}`;

    try {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Server responded with status ${response.status}`);
        }

        const books = await response.json();
        
        renderResults(books, resultsContainer);

    } catch (error) {
        console.error('Search failed:', error);
        resultsContainer.innerHTML = `<p class="text-red-500">Search error: ${error.message}</p>`;
    }
}

function getResultsContainer(code) {
    if (code === 2) {
        return document.getElementById('created-results-container');
    }
    if (code === 3) {
        return document.getElementById('fav-results-container');
    }
    
    return document.getElementById('results-container');
}

function clearSelectedLabels() {
    selectedLabelIds = [];
    document.querySelectorAll('.fr-label.selected').forEach(label => {
        label.classList.remove('selected');
    });
}

// Function to render results, matching the structure
function renderResults(books, container) {
    if (!container) {
        console.warn('renderResults: no results container found');
        return;
    }
    container.innerHTML = ''; // Clearing container

    if (books.length === 0) {
        container.innerHTML = '<p class="col-span-full text-gray-500">Nothing found, try changing the filters.</p>';
        return;
    }

    books.forEach(book => {
        // Creating fr-card element
        const bookElement = document.createElement('div');
        bookElement.className = 'fr-card bg-white rounded-xl shadow-md p-2 transition-shadow duration-300 hover:shadow-lg'; 
        
        bookElement.innerHTML = `
            <a href="/library/book/${book.id}">
                ${book.cover_path
                    ? `<img src="/covers/${book.cover_path}" alt="book cover" class="book-cover">` 
                    : `<div class="fr-blue-box">${book.name}</div>`
                }
            </a>
            <div class="fr-card__title">${book.name}</div>
            <div class="fr-card__author">${book.author}</div>
        `;
        container.appendChild(bookElement);
    });
}