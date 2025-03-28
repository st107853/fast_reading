function submitForm() {
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

    console.log("Создан объект:", dict);

    console.log("будет передан:", dict);
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/library/", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(JSON.stringify(dict));
    console.log("Передан:", JSON.stringify(dict));
    xhr.onreadystatechange = function() {
        if (xhr.readyState == 4 && xhr.status == 200) {
            console.log("Ответ сервера:", xhr.responseText);
            alert("Книга успешно добавлена!");
            window.location.href = "/";  // Переходим на главную страницу
        }
    };
}

document.addEventListener("DOMContentLoaded", function() {
    document.getElementById("submitBtn").addEventListener("click", submitForm);
});