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

    // Проверяем, выбран ли файл
    if (fileInput.files.length === 0) {
        alert("Please select a file.");
        return;
    }

    const file = fileInput.files[0];
    const reader = new FileReader();

    // Читаем содержимое файла
    reader.onload = function(event) {
        const fileContent = event.target.result;
        bookTextArea.value = fileContent; // Записываем содержимое файла в textarea
    };

    reader.onerror = function() {
        alert("Error reading file.");
    };

    reader.readAsText(file); // Читаем файл как текст
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
function submitChapter(button, bookId) {
    var chapterNameElement = document.getElementById('chapter-name'); 
    var bookTextElement = document.getElementById('chapter-text');
    console.log("Chapter name:", chapterNameElement, "Chapter text:", bookTextElement)

    if (!chapterNameElement || !bookTextElement) {
        console.error('Chapter elements not found');
        return;
    }

    var url = '/library/addbook/' + bookId + '/chapter';

    var payload = {
        title: chapterNameElement.value.trim(),
        text: bookTextElement.value
    };

    var xhr = new XMLHttpRequest();
    xhr.open('POST', url, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status === 201) {
                button.classList.add('clicked');
                try { 
                    if (bookId) { 
                        window.location.href = '/library/book/' + encodeURIComponent(bookId); 
                    } 
                } catch (e) {}
            } else {
                console.error('Failed to save chapter', xhr.status, xhr.responseText);
                alert('Failed to save chapter: ' + xhr.responseText);
            }
        }
    };
    xhr.send(JSON.stringify(payload));
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