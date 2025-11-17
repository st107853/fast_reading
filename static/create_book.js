function submitForm(button, bookId) {

    var bookNameElement = document.getElementById('book-name');
    var bookAuthorElement = document.getElementById('author-name');
    var releaseDateElement = document.getElementById('publication-year');
    var bookTextElement = document.getElementById('book-description');

    if (!bookNameElement || !bookAuthorElement || !releaseDateElement) {
        console.error("Ошибка: один из элементов не найден!");
        return;
    }

    var dict = {
        "name": bookNameElement.value.trim(),
        "author": bookAuthorElement.value.trim(),
        "release_date": releaseDateElement.value.trim()+"-01-02T15:04:05Z",
        "description": bookTextElement.value.trim()
    };

    var xhr = new XMLHttpRequest();

    if (bookId) {
        xhr.open("PUT", "/library/" + bookId, true);
    } else {
        xhr.open("POST", "/library/", true);
    }
    xhr.setRequestHeader('Content-Type', 'application/json');
    // allow cookies to be set on the response (same-origin)
    xhr.withCredentials = true;
    xhr.send(JSON.stringify(dict));
    xhr.onreadystatechange = function() {
        console.log("Ответ сервера:", xhr.responseText);
        // Ожидайте 200 (OK) для PUT или 201 (Created) для POST
        if (xhr.readyState === 4 && (xhr.status === 200 || xhr.status === 201)) {
            console.log("Успех:", xhr.responseText);
            button.classList.add("clicked");
        } else if (xhr.readyState === 4) {
            // Обработка ошибок
            console.error("Ошибка обновления:", xhr.responseText);
        }
    };
}

function updateText() {
    const fileInput = document.getElementById("file");
    const bookTextArea = document.getElementById("bookText");

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
            addChapterAnchor.setAttribute('href', '/library/book/' + encodeURIComponent(bookId) + '/chapter');
        }
    } catch (e) {
        console.error('failed to update add chapter link', e);
    }
});

// Submit chapter for the current book
function submitChapter(button, bookId) {
    var chapterNameElement = document.getElementById('chapter-name');
    var bookTextElement = document.getElementById('bookText');

    if (!chapterNameElement || !bookTextElement) {
        console.error('Chapter elements not found');
        return;
    }

    var url = '/library/addbook/' + bookId + '/chapter';

    var payload = {
        title: chapterNameElement.value.trim(),
        text: bookTextElement.innerText || bookTextElement.textContent || ''
    };

    var xhr = new XMLHttpRequest();
    xhr.open('POST', url, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
            if (xhr.status === 201) {
                // optional: visually indicate success
                button.classList.add('clicked');
                try { if (bookId) { window.location.href = '/library/book/' + encodeURIComponent(bookId); } } catch (e) {}
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