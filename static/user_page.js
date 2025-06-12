// Handle new book creation
function createNewBook() {
    const bookName = prompt("Enter the name of the new book:");
    if (bookName) {
    const li = document.createElement("li");
    li.className = "book-item";
    li.innerText = bookName;
    createdBooksContainer.appendChild(li);
    user.createdBooks.push(bookName);
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