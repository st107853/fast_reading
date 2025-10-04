const modal = document.getElementById('settingsModal');
const openBtn = document.getElementById('openSettingsBtn');
const closeBtn = document.getElementById('closeSettingsBtn');
const swatches = document.querySelectorAll('.swatch');
const body = document.body;
const loginForm = document.getElementById("login-form");
const themeToggle = document.getElementById('theme-toggle');
const cardList = document.getElementById("cardList");
const range = document.getElementById("customRange1");

function updateAction() {
    var name = document.getElementById('name').value;
    document.getElementById('bookForm').action += name;
}

function login() {
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

range.max = cardList.scrollWidth - cardList.clientWidth;

range.addEventListener("input", () => {
  cardList.scrollLeft = range.value;

});

cardList.addEventListener("scroll", () => {
  range.value = cardList.scrollLeft;
});

document.addEventListener('DOMContentLoaded', () => {
    // Modal Open/Close Logic
    openBtn.addEventListener('click', () => {
        modal.classList.add('active');
        body.style.overflow = 'hidden'; // Prevent scrolling the background
    });

    const closeModal = () => {
        modal.classList.remove('active');
        body.style.overflow = '';
    };

    closeBtn.addEventListener('click', closeModal);

    // Close when clicking outside the modal content
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            closeModal();
        }
    });

    // Close on Escape key
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && modal.classList.contains('active')) {
            closeModal();
        }
    });

    // Colour Swatch Selection Logic
    swatches.forEach(swatch => {
        swatch.addEventListener('click', (e) => {
            // Remove 'selected' state from all swatches
            swatches.forEach(s => s.setAttribute('aria-checked', 'false'));

            // Set 'selected' state on the clicked swatch
            e.currentTarget.setAttribute('aria-checked', 'true');

            // Apply the selected color to the CSS variable
            const newColor = e.currentTarget.dataset.color;
            document.documentElement.style.setProperty('--primary-color', newColor);
            
            console.log(`New primary color set to: ${newColor}`);
        });
        
        // Handle keyboard navigation for accessibility
        swatch.addEventListener('keydown', (e) => {
            if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                e.currentTarget.click();
            }
        });
    });

    // Theme Toggle Logic
    themeToggle.addEventListener('change', (e) => {
        if (e.target.checked) {
            body.classList.add('dark-theme');
            console.log('Switched to Dark Theme');
        } else {
            body.classList.remove('dark-theme');
            console.log('Switched to Light Theme');
        }
   });
}); 