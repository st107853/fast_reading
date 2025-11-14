document.addEventListener('DOMContentLoaded', () => {
    // Query DOM elements here to avoid errors when script runs on pages without them
    const openButtons = document.querySelectorAll('[data-modal-target]');
    const swatches = Array.from(document.querySelectorAll('.fr-swatch'));
    const themeToggles = Array.from(document.querySelectorAll('input[type="checkbox"][id="theme-toggle"], input.fr-theme-toggle'));

    // Swatch persistence key
    const SWATCH_KEY = 'fr_swatch_color';

    const saveSwatch = (color) => {
        try { localStorage.setItem(SWATCH_KEY, color); } catch (e) {}
    };

    const getSavedSwatch = () => {
        try { return localStorage.getItem(SWATCH_KEY); } catch (e) { return null; }
    };

    const applySwatch = (color) => {
        if (!color) return;
        document.documentElement.style.setProperty('--primary-color', color);

        // If there are swatch elements, mark the matching one as selected
        if (swatches.length) {
            swatches.forEach(s => {
                const c = (s.dataset.color || '').toLowerCase();
                s.setAttribute('aria-checked', c === String(color).toLowerCase());
            });
        }
    };

    // Generic modal handling by data attribute
    const openModal = (modalEl) => {
        if (!modalEl) return;
        modalEl.classList.add('active');
        modalEl.setAttribute('aria-hidden', 'false');
        document.body.style.overflow = 'hidden';
        modalEl.focus();
    };

    const closeModal = (modalEl) => {
        if (!modalEl) return;
        modalEl.classList.remove('active');
        modalEl.setAttribute('aria-hidden', 'true');
        document.body.style.overflow = '';
    };

    openButtons.forEach(btn => {
        btn.addEventListener('click', (e) => {
            const selector = btn.getAttribute('data-modal-target');
            if (!selector) return;
            const modalEl = document.querySelector(selector);
            openModal(modalEl);
        });
    });

    // Close buttons inside any modal: uses .fr-modal__close-btn
    document.querySelectorAll('.fr-modal__close-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const modalEl = btn.closest('.fr-modal');
            closeModal(modalEl);
        });
    });

    // Click on backdrop to close
    document.querySelectorAll('.fr-modal').forEach(modalEl => {
        modalEl.addEventListener('click', (e) => {
            if (e.target === modalEl) closeModal(modalEl);
        });
    });

    // Close on ESC for any open modal
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            document.querySelectorAll('.fr-modal.active').forEach(modalEl => closeModal(modalEl));
        }
    });

    // Colour Swatch Selection Logic
    // Apply saved swatch early so pages without controls still show the chosen color
    const initialSavedSwatch = getSavedSwatch();
    if (initialSavedSwatch) applySwatch(initialSavedSwatch);

    if (swatches.length) {
        // If no saved swatch, prefer any swatch already marked in HTML via aria-checked
        if (!initialSavedSwatch) {
            const prechecked = swatches.find(s => s.getAttribute('aria-checked') === 'true');
            if (prechecked && prechecked.dataset.color) applySwatch(prechecked.dataset.color);
        }

        swatches.forEach(swatch => {
            swatch.addEventListener('click', (e) => {
                const target = e.currentTarget;
                // Remove 'selected' state from all swatches
                swatches.forEach(s => s.setAttribute('aria-checked', 'false'));

                // Set 'selected' state on the clicked swatch
                target.setAttribute('aria-checked', 'true');

                // Apply & persist the selected color
                const newColor = target.dataset.color;
                if (newColor) {
                    applySwatch(newColor);
                    saveSwatch(newColor);
                }
            });

            // Handle keyboard navigation for accessibility
            swatch.addEventListener('keydown', (e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                    e.preventDefault();
                    e.currentTarget.click();
                }
            });
        });
    }

    /* Theme (dark mode) handling ------------------------------------------------- */
    const THEME_KEY = 'fr_theme'; // 'dark' | 'light'

    const getSavedTheme = () => {
        try { return localStorage.getItem(THEME_KEY); } catch (e) { return null; }
    };
    const saveTheme = (theme) => { try { localStorage.setItem(THEME_KEY, theme); } catch (e) {} };
    const prefersDark = () => window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;

    const applyTheme = (theme) => {
        const root = document.documentElement; // apply to <html>
        if (theme === 'dark') root.classList.add('fr-theme--dark'); else root.classList.remove('fr-theme--dark');

        // sync any toggle checkboxes on the page
        themeToggles.forEach(cb => {
            cb.checked = (theme === 'dark');
            cb.setAttribute('aria-pressed', theme === 'dark');
        });
    };

    // Initialize theme: saved -> system -> light
    let currentTheme = getSavedTheme();
    if (!currentTheme) currentTheme = prefersDark() ? 'dark' : 'light';
    applyTheme(currentTheme);

    // Listen for changes on theme toggles
    if (themeToggles.length) {
        themeToggles.forEach(cb => cb.addEventListener('change', (e) => {
            const theme = e.target.checked ? 'dark' : 'light';
            applyTheme(theme);
            saveTheme(theme);
        }));
    }

    // React to system preference change if user hasn't explicitly chosen
    const saved = getSavedTheme();
    if (!saved && window.matchMedia) {
        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
            applyTheme(e.matches ? 'dark' : 'light');
        });
    }
});

// Expose login as a global function so inline onclick handlers work across pages.
window.login = function() {
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
};