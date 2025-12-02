function initScrollSync() {
  const cardList = document.getElementById("cardList");
  const range = document.getElementById("customRange1");

  if (!cardList || !range) return;

  // Function to recalculate and set range max
  const syncRangeMax = () => {
    const max = Math.max(0, cardList.scrollWidth - cardList.clientWidth);
    range.max = String(max); // Ensure it's a string for the attribute
  };

  // Initial sync after a brief delay to allow layout to settle
  setTimeout(syncRangeMax, 100);

  // Sync when range input changes
  range.addEventListener("input", () => {
    cardList.scrollLeft = Number(range.value) || 0;
  });

  // Sync when cardList scrolls
  cardList.addEventListener("scroll", () => {
    range.value = cardList.scrollLeft;
  });

  // Recalculate max on window resize
  window.addEventListener("resize", syncRangeMax);

  // Also try to recalculate when images load
  const images = cardList.querySelectorAll('img');
  images.forEach(img => {
    img.addEventListener('load', syncRangeMax);
  });
}

// Wait for DOM to be ready
document.addEventListener('DOMContentLoaded', initScrollSync);
// Also call immediately in case DOMContentLoaded already fired
if (document.readyState === 'loading') {
  // DOM still loading, wait for event
} else {
  // DOM already loaded
  initScrollSync();
}

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
async function applyFilters() {
    const keyword = document.getElementById('keyword-input').value.trim();
    const resultsContainer = document.getElementById('results-container');
    
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

    const url = `/library/filter/?${queryParams.toString()}`;

    try {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Server responded with status ${response.status}`);
        }

        const books = await response.json();
        
        renderResults(books);

    } catch (error) {
        console.error('Search failed:', error);
        resultsContainer.innerHTML = `<p class="text-red-500">Search error: ${error.message}</p>`;
    }
}

// Function to render results, matching the structure
function renderResults(books) {
    const container = document.getElementById('results-container');
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
                <!-- Using placeholder image as cover is not provided in DTO -->
                <div class="fr-blue-box">
                    ${book.name.substring(0, 10)}
                </div>
            </a>
            <div class="fr-card__title">${book.name}</div>
            <div class="fr-card__author">${book.author}</div>
        `;
        container.appendChild(bookElement);
    });
}