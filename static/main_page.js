function updateAction() {
    var name = document.getElementById('name').value;
    document.getElementById('bookForm').action += name;
}

function login() {
    const loginForm = document.getElementById("login-form");
    function getCookie(name) {
        let matches = document.cookie.match(new RegExp(
            "(?:^|; )" + name.replace(/([\.$?*|{}\(\)\[\]\\\/\+^])/g, '\\$1') + "=([^;]*)"
        ));
        return matches ? decodeURIComponent(matches[1]) : undefined;
    }
    if (getCookie('logged_in') === 'true') {
        window.location.href = '/library/users/me';
    } else {
        window.location.href = '/library/auth/login';
    }
}