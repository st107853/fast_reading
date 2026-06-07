const wait = (ms) => new Promise(resolve => setTimeout(resolve, ms));

// Read the filter code from the page once at startup.
// The HTML sets this via data-filter-code on the main container.
function getPageFilterCode() {
    const container = document.querySelector('[data-filter-code]');
    return container ? container.getAttribute('data-filter-code') : '0';
}

let selectedLabelIds = [];

document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.fr-label, .fr-label--interactive').forEach(label => {
        label.addEventListener('click', toggleLabel);
    });
});

function toggleLabel(event) {
    const labelElement = event.currentTarget;
    const labelId = labelElement.getAttribute('data-id');

    labelElement.classList.toggle('selected');

    if (labelElement.classList.contains('selected')) {
        selectedLabelIds.push(labelId);
    } else {
        selectedLabelIds = selectedLabelIds.filter(id => id !== labelId);
    }
}

// code parameter is now optional — falls back to the page's own filter code
async function applyFilters(code) {
    await wait(100);

    const filterCode = (code !== undefined) ? code : getPageFilterCode();
    const keywordInput = document.getElementById('keyword-input');
    const keyword = keywordInput ? keywordInput.value.trim() : '';
    const resultsContainer = getResultsContainer(filterCode);

    if (!resultsContainer) {
        console.warn('applyFilters: no results container found');
        return;
    }

    resultsContainer.innerHTML = '<p>Searching...</p>';

    const queryParams = new URLSearchParams();
    if (keyword) queryParams.append('keyword', keyword);
    if (selectedLabelIds.length > 0) queryParams.append('labels', selectedLabelIds.join(','));
    queryParams.append('code', filterCode);

    try {
        const response = await fetch(`/library/filter/?${queryParams.toString()}`);
        if (!response.ok) throw new Error(`Server responded with status ${response.status}`);

        const books = await response.json();
        renderResults(books, resultsContainer);

    } catch (error) {
        console.error('Search failed:', error);
        resultsContainer.innerHTML = `<p>Search error: ${error.message}</p>`;
    }
}

function getResultsContainer(code) {
    const map = {
        '2': 'created-results-container',
        '3': 'fav-results-container',
    };
    const id = map[String(code)] || 'results-container';
    return document.getElementById(id);
}

function clearSelectedLabels() {
    selectedLabelIds = [];
    document.querySelectorAll('.fr-label.selected, .fr-label--interactive.selected').forEach(label => {
        label.classList.remove('selected');
    });
}

function renderResults(books, container) {
    if (!container) return;
    container.innerHTML = '';

    if (books.length === 0) {
        container.innerHTML = '<p>Nothing found, try changing the filters.</p>';
        return;
    }

    books.forEach(book => {
        const bookElement = document.createElement('div');
        bookElement.className = 'fr-card';

        bookElement.innerHTML = `
            <a href="/library/book/${book.id}">
                ${book.cover_path
                    ? `<img src="/covers/${book.cover_path}" alt="book cover" class="fr-card__img">`
                    : `<div class="fr-blue-box">${book.name}</div>`
                }
            </a>
            <div class="fr-card__title">${book.name}</div>
            <div class="fr-card__author">${book.author}</div>
        `;
        container.appendChild(bookElement);
    });
}