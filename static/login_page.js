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
    loginForm.classList.remove("fr-hidden");
    registerForm.classList.add("fr-hidden");
    toggleText.textContent = "Don't have an account?";
    toggleLink.textContent = "Sign Up";
    } else {
    formTitle.textContent = "Sign Up";
    loginForm.classList.add("fr-hidden");
    registerForm.classList.remove("fr-hidden");
    toggleText.textContent = "Already have an account?";
    toggleLink.textContent = "Log In";
    }
});

loginForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    
    const email = document.getElementById("login-email").value;
    const password = document.getElementById("login-password").value;

    try {
        const response = await fetch("/library/auth/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, password })
        });

        if (!response.ok) {
            const data = await response.json();
            alert(data.error || "Invalid email or password");
            return;
        }

        window.location.href = "/library/users/me";

    } catch (err) {
        alert("Something went wrong, please try again");
        console.error(err);
    }
});

registerForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    
    const email = document.getElementById("register-email").value;
    const name = document.getElementById("register-name").value;
    const password = document.getElementById("register-password").value;
    const passwordConfirm = document.getElementById("register-confirm-password").value;

    if (password !== passwordConfirm) {
        alert("Passwords do not match");
        return;
    }

    try {
        const registerResponse = await fetch("/library/auth/register", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, name, password, passwordConfirm })
        });

        if (!registerResponse.ok) {
            const data = await registerResponse.json();
            alert(data.error || "Registration failed");
            return;
        }

        const loginResponse = await fetch("/library/auth/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, password })
        });

        if (!loginResponse.ok) {
            const data = await loginResponse.json();
            alert(data.error || "Login after registration failed");
            return;
        }

        window.location.href = "/library/users/me";

    } catch (err) {
        alert("Something went wrong, please try again");
        console.error(err);
    }
});
