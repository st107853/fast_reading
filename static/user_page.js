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