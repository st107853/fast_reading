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

// Toggle between favourite and created books
document.addEventListener('DOMContentLoaded', function () {
    try {
        var tabButtons = document.querySelectorAll('[data-target]');
        if (!tabButtons || tabButtons.length === 0) return;

        var panels = function () { return document.querySelectorAll('#favBooks, #createdBooks'); };

        tabButtons.forEach(function (btn) {
            btn.addEventListener('click', function (ev) {
                ev.preventDefault();
                var target = btn.getAttribute('data-target');
                if (!target) return;

                // hide all panels
                panels().forEach(function (p) { p.classList.add('fr-hidden'); });

                // show the requested panel
                var show = document.querySelector(target);
                if (show) show.classList.remove('fr-hidden');

                // update active state on buttons
                tabButtons.forEach(function (b) { b.classList.remove('fr-btn--chosed'); b.setAttribute('aria-pressed', 'false'); });
                btn.classList.add('fr-btn--chosed');
                btn.setAttribute('aria-pressed', 'true');
            });
        });
    } catch (e) {
        console.error('user_page toggle init error', e);
    }
});