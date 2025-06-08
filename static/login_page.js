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
    xhr.open("POST", "/library/auth/login", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(JSON.stringify(dict));
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4 && (xhr.status === 200 || xhr.status === 201)) {
            window.location.href = "/library/users/me";
        }
    };
    // Here you can send data to the server
});

registerForm.addEventListener("submit", (e) => {
    e.preventDefault();
    const email = document.getElementById("register-email").value;
    const name = document.getElementById("register-name").value;
    const password = document.getElementById("register-password").value;
    const passwordConfirm = document.getElementById("register-confirm-password").value;

    var dict = {
        "email": email,
        "name": name,
        "password": password,
        "passwordConfirm": passwordConfirm
    };

    console.log("Sended to server data:", dict);

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/library/auth/register", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(JSON.stringify(dict));
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4 && (xhr.status === 200 || xhr.status === 201)) {
            xhr.open("POST", "/library/auth/login", true);
            xhr.setRequestHeader('Content-Type', 'application/json');
            xhr.send(JSON.stringify(dict));
            xhr.onreadystatechange = function() {
                if (xhr.readyState === 4 && (xhr.status === 200 || xhr.status === 201)) {
                    window.location.href = "/library/users/me";
                }
            };    
        }
    };
});
