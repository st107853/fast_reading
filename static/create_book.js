const CHAPTER_TEXT_KEY = 'chapterEditor_text';
const COVER_IMAGE_KEY = 'chapterEditor_coverImage';
const TARGET_WIDTH = 168;
const TARGET_HEIGHT = 190;

// Фиктивные данные для демонстрации отправки на сервер
const MOCK_BOOK_ID = 42; 
const MOCK_API_URL = `/books/${MOCK_BOOK_ID}`; // Эмулируем ваш роут PUT /books/:book_id

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

// --- ФУНКЦИИ КОНВЕРТАЦИИ ДЛЯ ОТПРАВКИ ---

/**
 * Конвертирует Data URL (Base64) в объект Blob.
 * Используется для подготовки изображения обложки к отправке через FormData.
 * @param {string} dataurl - Data URL изображения.
 * @returns {Blob} Объект Blob.
 */
function dataURLtoBlob(dataurl) {
    const parts = dataurl.split(',');
    
    // 1. Усиленная проверка: Data URL должен содержать две части: метаданные и данные.
    if (parts.length !== 2) {
        throw new Error("Data URL is improperly formatted (missing base64 data part).");
    }

    // Получаем MIME-тип из первой части (e.g., image/jpeg)
    const mimeMatch = parts[0].match(/:(.*?);/);
    if (!mimeMatch || mimeMatch.length < 2) {
        // Если MIME-тип не найден, используем по умолчанию image/jpeg
        console.warn("MIME type not found in Data URL header. Defaulting to 'image/jpeg'.");
        var mime = 'image/jpeg';
    } else {
        var mime = mimeMatch[1];
    }
    
    const base64Data = parts[1];

    // 2. Декодируем Base64-строку. Это наиболее вероятное место сбоя, 
    // если данные в localStorage повреждены.
    const bstr = atob(base64Data);
    let n = bstr.length;
    const u8arr = new Uint8Array(n);

    // Преобразуем декодированные символы в Uint8Array
    while (n--) {
        u8arr[n] = bstr.charCodeAt(n);
    }
    return new Blob([u8arr], { type: mime });
}

// --- ЛОГИКА ЗАГРУЗКИ ИЗОБРАЖЕНИЯ (С РЕСАЙЗОМ) ---
function handleImageUpload() {
    const fileInput = document.getElementById("image-file");
    const imagePreview = document.getElementById("image-preview");
    const file = fileInput.files[0];

    if (!file || !file.type.startsWith('image/')) {
        showMessage("Пожалуйста, выберите файл изображения.", 'error');
        return;
    }

    const reader = new FileReader();

    reader.onload = function(event) {
        const img = new Image();
        img.onload = function() {
            // Создаем холст для ресайза
            const canvas = document.createElement('canvas');
            const ctx = canvas.getContext('2d');

            canvas.width = TARGET_WIDTH;
            canvas.height = TARGET_HEIGHT;

            // Рисуем изображение на холсте, сжимая его
            ctx.drawImage(img, 0, 0, TARGET_WIDTH, TARGET_HEIGHT);

            // Получаем ресайзнутое изображение в формате DataURL
            const resizedImageUrl = canvas.toDataURL('image/jpeg', 0.8); // Сжатие качества 0.8

            // 1. Обновление превью
            imagePreview.style.backgroundImage = `url(${resizedImageUrl})`;
            imagePreview.textContent = "";

            // 2. Сохранение в ГЛОБАЛЬНУЮ ПЕРЕМЕННУЮ (ИСПРАВЛЕНИЕ)
            window.tempCoverImageBase64 = resizedImageUrl;
            
            showMessage(`Обложка успешно загружена, сжата до ${TARGET_WIDTH}x${TARGET_HEIGHT}px.`, 'success');
        };
        
        img.onerror = function() {
            showMessage("Ошибка загрузки изображения в память.", 'error');
        };

        img.src = event.target.result;
    };

    reader.onerror = function() {
        showMessage("Ошибка чтения файла изображения.", 'error');
    };
    
    reader.readAsDataURL(file);
}

// --- ЛОГИКА ОТПРАВКИ ОБНОВЛЕНИЙ ---
async function saveUpdates(button, bookId) {
    
    // ИСПРАВЛЕНИЕ: Получаем Data URL из ГЛОБАЛЬНОЙ ПЕРЕМЕННОЙ
    const savedImageURL = window.tempCoverImageBase64; 

    const bookName = document.getElementById('book-name');
    const bookAuthor = document.getElementById('author-name');
    const releaseDate = document.getElementById('publication-year');
    const bookText = document.getElementById('book-description');

    if (!bookName || !bookAuthor || !releaseDate || !bookText) {
        showMessage("Ошибка: Не удалось найти один из обязательных элементов формы.", 'error');
        return;
    }
    
    const formData = new FormData();
    formData.append('name', bookName.value.trim());
    formData.append('author', bookAuthor.value.trim());
    // Добавляем дату в формате ISO 8601, как вы указали
    formData.append('release_date', releaseDate.value.trim() + "-01-02T00:00:00Z"); 
    formData.append('description', bookText.value.trim());

    
    // 2. Добавление изображения обложки как файла (если доступно)
    if (savedImageURL) {
        try {
            // Преобразуем Data URL в Blob
            const imageBlob = dataURLtoBlob(savedImageURL);
            
            let extension = imageBlob.type.split('/')[1];
            if (extension === 'jpeg') extension = 'jpg'; 
            
            // Добавляем Blob в FormData с именем 'cover_image' (как ожидает Go-контроллер)
            formData.append('cover_image', imageBlob, `cover_${bookId || Date.now()}.${extension}`);
            
            showMessage("Подготовка: Изображение обложки добавлено в форму.", 'info');
        } catch (e) {
            console.error("Blob conversion error:", e);
            showMessage(`Ошибка при преобразовании изображения обложки: ${e.message}. Отправка без обложки.`, 'error');
        }
    } else {
        // Теперь это сообщение будет появляться только если пользователь не загружал обложку
        showMessage("Подготовка: Изображение обложки не найдено. Отправка только текста.", 'info');
    }

    const method = bookId ? "PUT" : "POST";
    const url = bookId ? `/library/${bookId}` : `/library/`;
    
    // Mocking the fetch call for demonstration purposes in Canvas
    button.disabled = true;
    button.textContent = 'Отправка...';

   try {
        const response = await fetch(url, {
            method,
            body: formData
        });
        if (response.status === 201) {
            const result = await response.json();
            button.classList.add("clicked");
            showMessage("Книга успешно создана!", 'success');
            // Перенаправление на страницу создания главы
            window.location.href = `/library/addbook/${result.book_id}`; 
        } 
        else if (response.status === 200) {
            button.classList.add("clicked");
            showMessage("Книга успешно обновлена!", 'success');
        } 
        else {
            const errorText = await response.text();
            throw new Error(errorText);
        }
    } catch (err) {
        console.error("Error during fetch:", err);
        showMessage("Ошибка при отправке данных на сервер: " + err.message, 'error');
    } finally {
        button.classList.remove("clicked");
    }
    
    button.disabled = false;
    button.textContent = 'Saved';
}

// Инициализация превью при загрузке
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
        alert("Ошибка: не найдены элементы формы главы", !chapterNameElement, !bookTextElement);
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