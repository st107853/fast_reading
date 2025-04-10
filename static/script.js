function submitForm(button) {

    console.log("Функция submitForm вызвана!");

    var bookNameElement = document.getElementById('bookName');
    var bookAuthorElement = document.getElementById('bookAuthor');
    var releaseDateElement = document.getElementById('releaseDate');
    var bookTextElement = document.getElementById('bookText');

    if (!bookNameElement || !bookAuthorElement || !releaseDateElement || !bookTextElement) {
        console.error("Ошибка: один из элементов не найден!");
        return;
    }

    var dict = {
        "name": bookNameElement.value.trim(),
        "author": bookAuthorElement.value.trim(),
        "date": releaseDateElement.value.trim()+"T15:04:05Z",
        "text": bookTextElement.value.trim()
    };

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/library/", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(JSON.stringify(dict));
    console.log("Передан:", JSON.stringify(dict));
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4 && xhr.status === 201) {
            console.log("Ответ сервера:", xhr.responseText);
            button.classList.add("clicked"); // Если всё ок, кнопка меняет цвет
            // window.location.href = "/library/";  // Переходим на главную страницу
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

