const TARGET_WIDTH = 168;
const TARGET_HEIGHT = 190;


// Utility function to display messages
function showMessage(message, type) {
    const messageContainer = document.getElementById('message-container');
    if (messageContainer) {
        messageContainer.textContent = message;
        messageContainer.className = 'p-2 rounded text-sm font-medium transition-colors duration-300 ' + 
                                        (type === 'success' ? 'bg-green-100 text-green-700' : 
                                        type === 'error' ? 'bg-red-100 text-red-700' : 'bg-blue-100 text-blue-700');
    } else {
        console.log(`[${type.toUpperCase()}]: ${message}`);
    }
}

// --- Functions of Image Processing ---

/**
 * Converts a Data URL (Base64) to a Blob object.
 * Used to prepare the cover image for submission via FormData.
 * @param {string} dataurl - Data URL of the image.
 * @returns {Blob} The resulting Blob object.
 */
function dataURLtoBlob(dataurl) {
    const parts = dataurl.split(',');
    
    // Data URL should have two parts: header and data
    if (parts.length !== 2) {
        throw new Error("Data URL is improperly formatted (missing base64 data part).");
    }

    // Get the MIME type from the first part (e.g., image/jpeg)
    const mimeMatch = parts[0].match(/:(.*?);/);
    if (!mimeMatch || mimeMatch.length < 2) {
        // If MIME type is not found, use the default image/jpeg
        console.warn("MIME type not found in Data URL header. Defaulting to 'image/jpeg'.");
        var mime = 'image/jpeg';
    } else {
        var mime = mimeMatch[1];
    }
    
    const base64Data = parts[1];

    // Decoding Base64 string to binary data
    const bstr = atob(base64Data);
    let n = bstr.length;
    const u8arr = new Uint8Array(n);

    // Convert the decoded characters to a Uint8Array
    while (n--) {
        u8arr[n] = bstr.charCodeAt(n);
    }
    return new Blob([u8arr], { type: mime });
}

// --- LOGIC FOR IMAGE UPLOAD (WITH RESIZING) ---
function handleImageUpload() {
    const fileInput = document.getElementById("image-file");
    const imagePreview = document.getElementById("image-preview");
    const file = fileInput.files[0];

    if (!file || !file.type.startsWith('image/')) {
        showMessage("Please select an image file.", 'error');
        return;
    }

    const reader = new FileReader();

    reader.onload = function(event) {
        const img = new Image();
        img.onload = function() {
            // Create a canvas to resize the image
            const canvas = document.createElement('canvas');
            const ctx = canvas.getContext('2d');

            canvas.width = TARGET_WIDTH;
            canvas.height = TARGET_HEIGHT;

            // Draw the image on the canvas, resizing it
            ctx.drawImage(img, 0, 0, TARGET_WIDTH, TARGET_HEIGHT);

            // Get the resized image in DataURL format
            const resizedImageUrl = canvas.toDataURL('image/jpeg', 0.8); // Quality compression 0.8

            // Update the preview
            imagePreview.style.backgroundImage = `url(${resizedImageUrl})`;
            imagePreview.textContent = "";

            // Save in GLOBAL VARIABLE (FIX)
            window.tempCoverImageBase64 = resizedImageUrl;
            
            showMessage(`Cover successfully uploaded and resized to ${TARGET_WIDTH}x${TARGET_HEIGHT}px.`, 'success');
        };
        
        img.onerror = function() {
            showMessage("Error loading image into memory.", 'error');
        };

        img.src = event.target.result;
    };

    reader.onerror = function() {
        showMessage("Error reading image file.", 'error');
    };
    
    reader.readAsDataURL(file);
}

async function releaseBook(button, bookId) {
    const bookName = document.getElementById('book-name');
    const bookAuthor = document.getElementById('author-name');
    const releaseDate = document.getElementById('publication-year');
    const bookText = document.getElementById('book-description');

    if (!bookName || !bookAuthor || !releaseDate || !bookText) {
        showMessage("Error: Could not find one of the required form elements.", 'error');
        return;
    }

    const method = "PUT";
    const url = `/library/release/${bookId}`;
    
    // Mocking the fetch call for demonstration purposes in Canvas
    button.disabled = true;
    button.textContent = 'Publishing...';

   try {
        const response = await fetch(url, {
            method,
        });
        if (response.status === 200) {
            // button.classList.add("clicked");
            showMessage("Book successfully published!", 'success');
        } 
        else {
            const errorText = await response.text();
            throw new Error(errorText);
        }
    } catch (err) {
        console.error("Error during fetch:", err);
        showMessage("Error: " + err.message, 'error');
    } finally {
        //button.classList.remove("clicked");
    }
    
    button.disabled = false;
    button.textContent = 'Publish';

    
}

async function saveUpdates(button, bookId) {
    const savedImageURL = window.tempCoverImageBase64; 
    const bookName = document.getElementById('book-name');
    const bookAuthor = document.getElementById('author-name');
    const releaseDate = document.getElementById('publication-year');
    const bookText = document.getElementById('book-description');

    if (!bookName || !bookAuthor || !releaseDate || !bookText) {
        showMessage("Please fill in all required fields.", 'error');
        return;
    }
    
    const formData = new FormData();
    formData.append('name', bookName.value.trim());
    formData.append('author', bookAuthor.value.trim());
    formData.append('release_date', releaseDate.value.trim() + "-01-02T00:00:00Z"); 
    formData.append('description', bookText.value.trim());

    if (savedImageURL) {
        try {
            const imageBlob = dataURLtoBlob(savedImageURL);
            let extension = imageBlob.type.split('/')[1] || 'jpg';
            formData.append('cover_image', imageBlob, `cover.${extension}`);
        } catch (e) {
            console.error("Blob error:", e);
        }
    }

    const method = bookId ? "PUT" : "POST";
    const url = bookId ? `/library/${bookId}` : `/library/`;
    
    button.disabled = true;
    button.textContent = 'Отправка...';

    try {
        const response = await fetch(url, { method, body: formData });
        const result = await response.json();

        if (!response.ok) throw new Error(result.error || "Server error");

        // Creating new book
        if (response.status === 201) {
            const newBookId = result.book_id;

            if (selectedLabels && selectedLabels.length > 0) {
                await saveAllLabels(newBookId); 
            }

            window.location.href = `/library/addbook/${newBookId}`; 
        } 
        
        // Updating existing book
        else if (response.status === 200) {
            if (selectedLabels && selectedLabels.length > 0) {
                await saveAllLabels(bookId);
            }
            button.textContent = 'Saved';
        }

    } catch (err) {
        console.error(err);
        showMessage("Error: " + err.message, 'error');
        button.textContent = 'Error';
    } finally {
        button.disabled = false;
    }
}

// Initialize image preview if there's a temporary cover image stored (e.g., after page reload)
document.addEventListener('DOMContentLoaded', () => {
    const imagePreview = document.getElementById("image-preview");
    if (window.tempCoverImageBase64) {
        imagePreview.style.backgroundImage = `url(${window.tempCoverImageBase64})`;
    }
});

function updateText() {
    const fileInput = document.getElementById("file");
    const bookTextArea = document.getElementById("scrollable-content-reading");
    const file = fileInput.files[0];

    if (fileInput.files.length === 0) {
        alert("Please select a file.");
        return;
    }

    const reader = new FileReader();

    
    reader.readAsText(file);

    // Read file content
    reader.onload = function(event) {
        const fileContent = event.target.result;
        bookTextArea.value = fileContent; // Update textarea with file content
    };

    reader.onerror = function() {
        alert("Error reading file.");
    };
}

// Update the "Add chapter" anchor to point to the current book (if cookie exists)
document.addEventListener('DOMContentLoaded', function () {
    try {
        var cookies = document.cookie.split(';').map(c => c.trim());
        var bookId = null;
        for (var i = 0; i < cookies.length; i++) {
            if (cookies[i].startsWith('book_id=')) {
                bookId = cookies[i].substring('book_id='.length);
                break;
            }
        }
        var addChapterAnchor = document.querySelector('a[href="/library/addbook/' + encodeURIComponent(bookId) + '/chapter"]');
        if (addChapterAnchor && bookId) {
            addChapterAnchor.setAttribute('href', '/library/addbook/' + encodeURIComponent(bookId) + '/chapter');
        }
    } catch (e) {
        console.error('failed to update add chapter link', e);
    }
});

// Submit chapter for the current book
async function submitChapter(button, bookId, chapterId) {
    const chapterNameElement = document.getElementById('chapter-name');
    const bookTextElement = document.getElementById('scrollable-content-reading');

    if (!chapterNameElement || !bookTextElement) {
        alert("Error: Could not find one of the required chapter form elements.", !chapterNameElement, !bookTextElement);
        return;
    }

    // Set URL and method
    const method = chapterId ? "PUT" : "POST";
    const url = chapterId
        ? `/library/addbook/${encodeURIComponent(bookId)}/chapter/${encodeURIComponent(chapterId)}`
        : `/library/addbook/${encodeURIComponent(bookId)}/chapter`;

    // Assemble request body
    const payload = {
        title: chapterNameElement.value.trim(),
        text: bookTextElement.value
    };

    try {
        const response = await fetch(url, {
            method,
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify(payload)
        });

        // Create chapter
        if (!chapterId && response.status === 201) {
            button.classList.add('clicked');
            window.location.href = `/library/addbook/${encodeURIComponent(bookId)}`;
            return;
        }

        // Update chapter
        if (chapterId && response.status === 200) {
            button.classList.add('clicked');
            return;
        }

        // Errors
        const errorText = await response.text();
        throw new Error(errorText || "Unknown error");

    } catch (err) {
        console.error("Error of saving chapter:", err);
    }
}


// Handle book deletion
function deleteBook(id) {
    if (!id) {
        return;
    }

    var xhr = new XMLHttpRequest();
    xhr.open("DELETE", "/library/" + id, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4 && (xhr.status === 200 || xhr.status === 201)) {
            window.location.href = "/library/users/me";
        }
    };
    xhr.send();
}

// Handle book deletion
function deleteChapter(id) {
    if (!id) {
        return;
    }

    var xhr = new XMLHttpRequest();
    xhr.open("DELETE", "/library/chapter/" + id, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4 && (xhr.status === 200 || xhr.status === 201)) {
            window.location.href = "/library/users/me";
        }
    };
    xhr.send();
}

const menu = document.getElementById("dropdownMenu");
const labelsList = document.getElementById('labelsList');

// 17 existing labels + 1 for "no label" (id=0) which is always selected
let selectedLabels = Array(18).fill(false);

if (labelsList) {
    const items = labelsList.querySelectorAll('.fr-label-item');

    items.forEach(item => {
        const id = parseInt(item.getAttribute('data-id'));
        if (!isNaN(id) && id >= 0 && id < selectedLabels.length) {
            selectedLabels[id] = true;
        }
    });
}

function toggleDropdown(event) {
    if (event) event.stopPropagation();
    menu.classList.toggle("open");
}

function toggleLabelUI(labelId, labelName) {
    labelId = parseInt(labelId);
    selectedLabels[0] = true;


    if ( !selectedLabels[labelId] ) {
        selectedLabels[labelId] = true;
        renderLabel(labelId, labelName);
    } else {
        selectedLabels[labelId] = false;
        const elementToRemove = document.querySelector(`.fr-label-item[data-id="${labelId}"]`);
        if (elementToRemove) elementToRemove.remove();
    }
}

function renderLabel(id, name) {
    const div = document.createElement('div');
    div.className = 'fr-label-item';
    div.setAttribute('data-id', id);
    div.textContent = name;
    labelsList.appendChild(div);
}

async function saveAllLabels(bookId) {
    const url = `/library/book/${bookId}/labels`;
    const response = await fetch(url, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ label_ids: selectedLabels })
    });
    return response;
}

// Clouse dropdown when clicking outside
window.onclick = function(event) {
    if (!event.target.matches('.fr-btn-with-icon') && !event.target.closest('.fr-btn-with-icon')) {
        if (menu.classList.contains('open')) {
            menu.classList.remove('open');
        }
    }
}