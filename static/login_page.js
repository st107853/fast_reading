const loginForm = document.getElementById("login-form");
const registerForm = document.getElementById("register-form");
const formTitle = document.getElementById("form-title");
const toggleText = document.getElementById("toggle-text");
const toggleLink = document.getElementById("toggle-link");

let isLogin = true;

toggleLink.addEventListener("click", () => {
    isLogin = !isLogin;
    if (isLogin) {
    formTitle.textContent = "Log In";
    loginForm.classList.remove("hidden");
    registerForm.classList.add("hidden");
    toggleText.textContent = "Don't have an account?";
    toggleLink.textContent = "Sign Up";
    } else {
    formTitle.textContent = "Sign Up";
    loginForm.classList.add("hidden");
    registerForm.classList.remove("hidden");
    toggleText.textContent = "Already have an account?";
    toggleLink.textContent = "Log In";
    }
});

loginForm.addEventListener("submit", (e) => {
    e.preventDefault();
    const email = document.getElementById("login-email").value;
    const password = document.getElementById("login-password").value;

    var dict = {
        "email": email,
        "password": password
    };

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/auth/login", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(JSON.stringify(dict));
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4 && xhr.status === 200) {
            console.log("Ответ сервера:", xhr.responseText);
            window.location.href = "/api/users/me";  // Переходим на главную страницу
        }
    };
    // Here you can send data to the server
});

registerForm.addEventListener("submit", (e) => {
    e.preventDefault();
    const name = document.getElementById("register-name").value;
    const email = document.getElementById("register-email").value;
    const password = document.getElementById("register-password").value;
    const confirmPassword = document.getElementById("register-confirm-password").value;

    if (password !== confirmPassword) {
    alert("Passwords do not match!");
    return;
    }

    
});