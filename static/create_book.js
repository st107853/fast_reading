async function saveUpdates(button, bookId) {
    const bookName = document.getElementById('book-name');
    const bookAuthor = document.getElementById('author-name');
    const releaseDate = document.getElementById('publication-year');
    const bookText = document.getElementById('book-description');

    if (!bookName || !bookAuthor || !releaseDate || !bookText) {
        alert("Error: Cannot find one of the required elements.");
        return;
    }

    const data = {
        name: bookName.value.trim(),
        author: bookAuthor.value.trim(),
        release_date: releaseDate.value.trim() + "-01-02T00:00:00Z",
        description: bookText.value.trim()
    };

    const method = bookId ? "PUT" : "POST";
    const url = bookId ? `/library/${bookId}` : `/library/`;

    try {
        const response = await fetch(url, {
            method,
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify(data)
        });

        if (response.status === 201) {
            const result = await response.json();
            button.classList.add("clicked");
            window.location.href = `/library/addbook/${result.book_id}`;
        } 
        else if (response.status === 200) {
            button.classList.add("clicked");
        } 
        else {
            const errorText = await response.text();
            throw new Error(errorText);
        }

    } catch (err) {
        console.error("Error:", err);
    }
}



function updateText() {
    const fileInput = document.getElementById("file");
    const bookTextArea = document.getElementById("chapter-text");

    // Check if a file is selected
    if (fileInput.files.length === 0) {
        alert("Please select a file.");
        return;
    }

    const file = fileInput.files[0];
    const reader = new FileReader();

    // Read the file content
    reader.onload = function(event) {
        const fileContent = event.target.result;
        bookTextArea.value = fileContent; // Write the file content to the textarea
    };

    reader.onerror = function() {
        alert("Error reading file.");
    };

    reader.readAsText(file); // Read the file as text
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
    const bookTextElement = document.getElementById('chapter-text');

    if (!chapterNameElement || !bookTextElement) {
        alert("Ошибка: не найдены элементы формы главы");
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
            alert("Chapter successfully created");
            window.location.href = `/library/addbook/${encodeURIComponent(bookId)}`;
            return;
        }

        // Update chapter
        if (chapterId && response.status === 200) {
            button.classList.add('clicked');
            alert("Chapter successfully updated");
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
    console.log("Sended to server data:", id);
    if (!id) {
        alert("Please select a book to delete.");
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
    console.log("Sended to server data:", id);
    if (!id) {
        alert("Please select a chapter to delete.");
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

function toggleDropdown(event) {
    // stops form submission & page reload
    if (event) event.preventDefault();

    const menu = document.getElementById("dropdownMenu");
    menu.classList.toggle("open");
}

async function submitNewLabel(bookId, labelId) {
    const url = `/library/addbook/${bookId}/label/${labelId}`;

    try {
        const response = await fetch(url, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include"
        });

        if (response.status === 201) {
            return;
        }

        const responseText = await response.text();
        let errorData = {};
        try {
            errorData = JSON.parse(responseText);
        } catch (e) {
            errorData = { error: responseText };
        }

        throw new Error(errorData.error || `Server error (Status: ${response.status})`);

    } catch (err) {
        console.error("Error adding label:", err);
        alert("Failed to add label: " + err.message);
    }
}