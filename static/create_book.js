const TARGET_WIDTH = 168;
const TARGET_HEIGHT = 190;
const MAX_INPUT_MB = 5;           // лимит входного файла
const MAX_INPUT_BYTES = MAX_INPUT_MB * 1024 * 1024;


// Utility function to display messages
function showMessage(message, type) {
    const messageContainer = document.getElementById('message-container');
    if (messageContainer) {
        messageContainer.textContent = message;
        messageContainer.className = `cb-editor-message cb-editor-message--${type}`;
        messageContainer.hidden = false;
        messageContainer.setAttribute('aria-live', type === 'error' ? 'assertive' : 'polite');
    } else {
        console.log(`[${type.toUpperCase()}]: ${message}`);
    }
}

// --- Functions of Image Processing ---

let pendingCoverBlob = null;
let previewObjectURL = null;



// --- LOGIC FOR IMAGE UPLOAD (WITH RESIZING) ---
function handleImageUpload() {
    const fileInput = document.getElementById("image-file");
    const imagePreview = document.getElementById("image-preview");
    const file = fileInput.files[0];

    if (!file) return;

    if (!file.type.startsWith('image/')) {
        showMessage("Please select an image file.", 'error');
        return;
    }

    if (file.size > MAX_INPUT_BYTES) {
        showMessage(
            `Image is too large (${(file.size / 1024 / 1024).toFixed(1)} MB). 
             Maximum allowed size is ${MAX_INPUT_MB} MB.`,
            'error'
        );
        fileInput.value = '';
        return;
    }

    if (previewObjectURL) {
        URL.revokeObjectURL(previewObjectURL);
    }
    previewObjectURL = URL.createObjectURL(file);

    imagePreview.style.backgroundImage = `url(${previewObjectURL})`;
    imagePreview.style.backgroundSize = 'cover';
    imagePreview.style.backgroundPosition = 'center';
    imagePreview.textContent = '';

    const img = new Image();

    img.onload = function () {

        const canvas  = document.createElement('canvas');
        canvas.width  = TARGET_WIDTH;
        canvas.height = TARGET_HEIGHT;

        const ctx = canvas.getContext('2d');

        const srcRatio  = img.naturalWidth / img.naturalHeight;
        const dstRatio  = TARGET_WIDTH / TARGET_HEIGHT;

        let sx, sy, sw, sh;
        if (srcRatio > dstRatio) {
            sh = img.naturalHeight;
            sw = sh * dstRatio;
            sx = (img.naturalWidth - sw) / 2;
            sy = 0;
        } else {
            sw = img.naturalWidth;
            sh = sw / dstRatio;
            sx = 0;
            sy = (img.naturalHeight - sh) / 2;
        }

        ctx.drawImage(img, sx, sy, sw, sh, 0, 0, TARGET_WIDTH, TARGET_HEIGHT);

        canvas.toBlob(
            function (blob) {
                if (!blob) {
                    showMessage("Failed to process image.", 'error');
                    return;
                }

                pendingCoverBlob = blob;

                showMessage(
                    `Cover ready (${(blob.size / 1024).toFixed(0)} KB).`,
                    'success'
                );
            },
            'image/jpeg', 0.85
        );
    };

    img.onerror = function () {
        showMessage("Failed to load image.", 'error');
        URL.revokeObjectURL(previewObjectURL);
        previewObjectURL = null;
    };

    img.src = previewObjectURL;
}

async function releaseBook(button, bookId) {

    const savedId = await saveUpdates(null, bookId);

    if (!savedId) {
        showMessage("Could not publish: book failed to save.", 'error');
        return;
    }

    button.disabled = true;
    button.textContent = 'Publishing...';

    try {
        const response = await fetch(`/library/release/${savedId}`, {
            method: "PUT",
        });

        if (response.ok) {
            showMessage("Book successfully published!", 'success');
            button.textContent = 'Published';
        } else {
            const errorText = await response.text();
            throw new Error(errorText);
        }

    } catch (err) {
        console.error("Publish error:", err);
        showMessage("Error publishing: " + err.message, 'error');
        button.textContent = 'Publish';

    } finally {
        button.disabled = false;
    }
}

async function addChapter(button, bookId) {
    button.disabled = true;
    button.textContent = 'Saving...';

    const savedId = await saveUpdates(null, bookId);

    if (!savedId) {
        showMessage("Could not navigate: book failed to save.", 'error');
        button.disabled = false;
        button.textContent = 'Add Chapter';
        return;
    }

    // Navigate to the chapter page with the confirmed book id
    window.location.href = `/library/addbook/${savedId}/chapter`;
}

// static/create_book.js

// Returns the bookId (existing or newly created) on success, null on failure.
async function saveUpdates(button, bookId) {
    const bookName = document.getElementById('book-name');
    const bookAuthor = document.getElementById('author-name');
    const publicationYear = document.getElementById('publication-year');
    const bookDescription = document.getElementById('book-description');
    const bookForm = bookName?.form;

    if (!bookName || !bookAuthor || !publicationYear || !bookDescription || !bookForm) {
        showMessage("Could not save: book form fields are missing.", 'error');
        return null;
    }

    if (!bookForm.reportValidity()) {
        return null;
    }

    if (!bookName.value.trim() || !bookAuthor.value.trim()) {
        showMessage("Book name and author are required.", 'error');
        return null;
    }

    const formData = new FormData();
    formData.append('name', bookName.value.trim());
    formData.append('author', bookAuthor.value.trim());
    formData.append('publication_year', publicationYear.value.trim());
    formData.append('description', bookDescription.value.trim());

    if (pendingCoverBlob) {
        formData.append('cover_image', pendingCoverBlob, 'cover.jpg');
    }

    const isNew = !bookId || bookId === '0' || bookId === '';
    const method = isNew ? "POST" : "PUT";
    const url = isNew ? `/library/` : `/library/${bookId}`;

    if (button) {
        button.disabled = true;
        button.textContent = 'Saving...';
    }

    try {
        const response = await fetch(url, { method, body: formData });
        const result = await response.json();

        if (!response.ok) throw new Error(result.error || "Server error");

        // Save labels for both new and existing books
        const savedId = isNew ? result.book_id : bookId;
        if (selectedLabels && selectedLabels.some(Boolean)) {
            await saveAllLabels(savedId);
        }

        if (button) button.textContent = 'Saved';

        // Return the id so callers can use it for navigation
        return savedId;

    } catch (err) {
        console.error(err);
        showMessage("Error saving: " + err.message, 'error');
        if (button) button.textContent = 'Save';
        return null;

    } finally {
        if (button) button.disabled = false;
    }
}

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
        showMessage("Could not save: chapter form fields are missing.", 'error');
        return;
    }

    // Set URL and method
    const method = chapterId != 0 ? "PUT" : "POST";
    const url = chapterId != 0
        ? `/library/addbook/${encodeURIComponent(bookId)}/chapter/${encodeURIComponent(chapterId)}`
        : `/library/addbook/${encodeURIComponent(bookId)}/chapter`;

    // Assemble request body
    const payload = {
        title: chapterNameElement.value.trim(),
        text: bookTextElement.value
    };

    const originalButtonText = button?.textContent;
    if (button) {
        button.disabled = true;
        button.textContent = 'Saving...';
    }

    try {
        const response = await fetch(url, {
            method,
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify(payload)
        });

        const responseText = await response.text();
        let result = null;
        try {
            result = responseText ? JSON.parse(responseText) : null;
        } catch {
            // Some error responses are plain text rather than JSON.
        }

        const resultMessage = result && typeof result === 'object'
            ? result.error || result.message
            : null;
        const responseMessage = typeof resultMessage === 'string' && resultMessage
            ? resultMessage
            : responseText;

        if (!response.ok) {
            throw new Error(responseMessage || `Server error (${response.status})`);
        }

        // Create chapter
        if (response.status === 201) {
            const id = result?.chapter_id || chapterId;
            if (!id || id === '0') {
                throw new Error("Chapter was saved, but the server did not return its ID.");
            }
            window.location.href = `/library/addbook/${encodeURIComponent(bookId)}/chapter/${id}`;
            return;
        }

        showMessage("Chapter saved successfully.", 'success');

    } catch (err) {
        console.error("Error of saving chapter:", err);
        const errorMessage = err instanceof Error ? err.message : String(err);
        showMessage(`Error saving chapter: ${errorMessage}`, 'error');
    } finally {
        if (button) {
            button.disabled = false;
            button.textContent = originalButtonText || 'Save';
        }
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
const labelMenuButton = document.getElementById('label-menu-button');

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

function getMenuOptions() {
    return menu ? Array.from(menu.querySelectorAll('[role="option"]')) : [];
}

function setLabelOptionSelected(labelId, isSelected) {
    const option = getMenuOptions().find(item => item.dataset.id === String(labelId));
    if (option) option.setAttribute('aria-selected', String(isSelected));
}

function closeLabelMenu(restoreFocus = false) {
    if (!menu) return;
    menu.classList.remove('open');
    labelMenuButton?.setAttribute('aria-expanded', 'false');
    if (restoreFocus) labelMenuButton?.focus();
}

function openLabelMenu() {
    if (!menu) return;
    menu.classList.add('open');
    labelMenuButton?.setAttribute('aria-expanded', 'true');

    const options = getMenuOptions();
    const selectedOption = options.find(option => option.getAttribute('aria-selected') === 'true');
    (selectedOption || options[0])?.focus();
}

if (menu) {
    getMenuOptions().forEach(option => {
        const id = Number.parseInt(option.dataset.id, 10);
        option.setAttribute('aria-selected', String(Number.isInteger(id) && Boolean(selectedLabels[id])));
    });

    menu.addEventListener('click', event => {
        const target = event.target;
        if (!(target instanceof Element)) return;
        const option = target.closest('[role="option"]');
        if (!option || !menu.contains(option)) return;

        toggleLabelUI(option.dataset.id, option.textContent.trim());
        option.focus();
    });

    menu.addEventListener('keydown', event => {
        const options = getMenuOptions();
        if (!options.length) return;

        const index = options.indexOf(document.activeElement);
        let nextIndex;

        switch (event.key) {
            case 'ArrowDown':
                nextIndex = index < 0 ? 0 : (index + 1) % options.length;
                break;
            case 'ArrowUp':
                nextIndex = index <= 0 ? options.length - 1 : index - 1;
                break;
            case 'Home':
                nextIndex = 0;
                break;
            case 'End':
                nextIndex = options.length - 1;
                break;
            case 'Enter':
            case ' ':
                event.preventDefault();
                if (index >= 0) {
                    const option = options[index];
                    toggleLabelUI(option.dataset.id, option.textContent.trim());
                }
                return;
            case 'Escape':
                event.preventDefault();
                closeLabelMenu(true);
                return;
            case 'Tab':
                closeLabelMenu();
                return;
            default:
                return;
        }

        event.preventDefault();
        options[nextIndex]?.focus();
    });
}

function toggleDropdown(event) {
    if (!menu) return;
    if (event) event.stopPropagation();
    if (menu.classList.contains('open')) {
        closeLabelMenu();
    } else {
        openLabelMenu();
    }
}

function toggleLabelUI(labelId, labelName) {
    if (!labelsList) return;
    
    labelId = parseInt(labelId);
    selectedLabels[0] = true;

    if (!selectedLabels[labelId]) {
        selectedLabels[labelId] = true;
        renderLabel(labelId, labelName);
    } else {
        selectedLabels[labelId] = false;
        const elementToRemove = document.querySelector(
            `.fr-label-item[data-id="${labelId}"]`);
        if (elementToRemove) elementToRemove.remove();
    }
    setLabelOptionSelected(labelId, selectedLabels[labelId]);
}

function renderLabel(id, name) {
    if (!labelsList) return;
    
    const div = document.createElement('div');
    div.className = 'fr-label-item';
    div.setAttribute('data-id', id);
    div.textContent = name;
    labelsList.appendChild(div);
}

async function saveAllLabels(bookId) {
    if (!labelsList) return;
    
    const url = `/library/book/${bookId}/labels`;
    const response = await fetch(url, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ label_ids: selectedLabels })
    });
    return response;
}

if (menu) {
    window.addEventListener('click', event => {
        const target = event.target;
        if (!(target instanceof Node)) return;
        if (!menu.contains(target) && !labelMenuButton?.contains(target)) {
            closeLabelMenu();
        }
    });
}