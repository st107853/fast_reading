// Main logic for style and settings handling
document.addEventListener('DOMContentLoaded', () => {
    
    // --- DOM Elements ---
    const openButtons = document.querySelectorAll('[data-modal-target]');
    const swatches = Array.from(document.querySelectorAll('.fr-swatch'));
    const themeToggles = Array.from(document.querySelectorAll('input[type="checkbox"][id="theme-toggle"], input.fr-theme-toggle'));

    // Swatch logic
    const SWATCH_KEY = 'fr_swatch_color';
    const saveSwatch = (color) => { try { localStorage.setItem(SWATCH_KEY, color); } catch (e) {} };
    const getSavedSwatch = () => { try { return localStorage.getItem(SWATCH_KEY); } catch (e) { return null; } };

    const applySwatch = (color) => {
        if (!color) return;
        // Apply the color as a CSS variable
        document.documentElement.style.setProperty('--primary-color', color);

        if (swatches.length) {
            swatches.forEach(s => {
                const c = (s.dataset.color || '').toLowerCase();
                s.setAttribute('aria-checked', c === String(color).toLowerCase());
            });
        }
    };

   // Generic modal handling
    const openModal = (modalEl) => {
        if (!modalEl) return;
        
        // Save the element that was focused before opening the modal
        lastFocusedElement = document.activeElement; 

        modalEl.classList.add('active');
        modalEl.setAttribute('aria-hidden', 'false');
        document.body.style.overflow = 'hidden';
        
        // Move focus to the modal
        modalEl.focus(); 
    };

    const closeModal = (modalEl) => {
        if (!modalEl) return;
        
        // Set aria-hidden="true" and visually hide the modal
        modalEl.setAttribute('aria-hidden', 'true');
        modalEl.classList.remove('active');
        document.body.style.overflow = '';

        // Restore focus to the element that was focused before opening
        if (lastFocusedElement) {
            lastFocusedElement.focus();
            lastFocusedElement = null; // Clear the reference
        } else {
                // Fallback: Focus on the settings button
            const openButton = document.querySelector('[data-modal-target="#settingsModal"]');
            if (openButton) openButton.focus();
        }
    };

    openButtons.forEach(btn => {
        btn.addEventListener('click', (e) => {
            const selector = btn.getAttribute('data-modal-target');
            if (!selector) return;
            const modalEl = document.querySelector(selector);
            openModal(modalEl);
        });
    });

    document.querySelectorAll('.fr-modal__close-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const modalEl = btn.closest('.fr-modal');
            closeModal(modalEl);
        });
    });

    document.querySelectorAll('.fr-modal').forEach(modalEl => {
        modalEl.addEventListener('click', (e) => {
            if (e.target === modalEl) closeModal(modalEl);
        });
    });

    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            document.querySelectorAll('.fr-modal.active').forEach(modalEl => closeModal(modalEl));
        }
    });

    // Swatch Initialization & Events
    const initialSavedSwatch = getSavedSwatch();
    if (initialSavedSwatch) applySwatch(initialSavedSwatch);

    if (swatches.length) {
        if (!initialSavedSwatch) {
            const prechecked = swatches.find(s => s.getAttribute('aria-checked') === 'true');
            if (prechecked && prechecked.dataset.color) applySwatch(prechecked.dataset.color);
        }

        swatches.forEach(swatch => {
            swatch.addEventListener('click', (e) => {
                const target = e.currentTarget;
                swatches.forEach(s => s.setAttribute('aria-checked', 'false'));
                target.setAttribute('aria-checked', 'true');
                const newColor = target.dataset.color;
                if (newColor) {
                    applySwatch(newColor);
                    saveSwatch(newColor);
                }
            });
            swatch.addEventListener('keydown', (e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                    e.preventDefault();
                    e.currentTarget.click();
                }
            });
        });
    }

    /* Theme (dark mode) handling */
    const THEME_KEY = 'fr_theme';
    const getSavedTheme = () => { try { return localStorage.getItem(THEME_KEY); } catch (e) { return null; } };
    const saveTheme = (theme) => { try { localStorage.setItem(THEME_KEY, theme); } catch (e) {} };
    const prefersDark = () => window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;

    const applyTheme = (theme) => {
        const root = document.documentElement;
        if (theme === 'dark') root.classList.add('fr-theme--dark'); else root.classList.remove('fr-theme--dark');

        themeToggles.forEach(cb => {
            cb.checked = (theme === 'dark');
            cb.setAttribute('aria-pressed', theme === 'dark');
        });
    };

    let currentTheme = getSavedTheme();
    if (!currentTheme) currentTheme = prefersDark() ? 'dark' : 'light';
    applyTheme(currentTheme);

    if (themeToggles.length) {
        themeToggles.forEach(cb => cb.addEventListener('change', (e) => {
            const theme = e.target.checked ? 'dark' : 'light';
            applyTheme(theme);
            saveTheme(theme);
        }));
    }

    const saved = getSavedTheme();
    if (!saved && window.matchMedia) {
        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
            applyTheme(e.matches ? 'dark' : 'light');
        });
    }
    
    // LOGIC FOR FONT AND TEXT SIZE
    const FONT_FAMILY_KEY = 'fr_font_family';
    const FONT_SIZE_KEY = 'fr_font_size';
    const FONT_WEIGHT_KEY = 'fr_font_weight';

    // Select elements
    const fontSelect = document.getElementById('fontFamilySelect');
    const weightSelect = document.getElementById('fontWeightSelect');
    const sizeSelect = document.getElementById('fontSizeSelect');

    const saveTextPref = (key, value) => { try { localStorage.setItem(key, value); } catch (e) {} };
    const getSavedTextPref = (key) => { try { return localStorage.getItem(key); } catch (e) { return null; } };
    
    //Loads saved settings and applies them via CSS variables
    const applyTextStyle = () => {
        const family = getSavedTextPref(FONT_FAMILY_KEY) || 'Inter, sans-serif'; 
        const size = getSavedTextPref(FONT_SIZE_KEY) || '1.1rem'; 
        const weight = getSavedTextPref(FONT_WEIGHT_KEY) || '400'; 
        
        // Apply styles via CSS variables on the root <html> element
        const root = document.documentElement.style;
        root.setProperty('--font-family-base', family);
        root.setProperty('--text-font-size', size);
        root.setProperty('--text-font-weight', weight);

        // Update select boxes to reflect the current state
        if (fontSelect) fontSelect.value = family;
        if (sizeSelect) sizeSelect.value = size;
        if (weightSelect) weightSelect.value = weight;
        
    };

    // Load saved settings on startup
    applyTextStyle(); 
    
    // Event listeners for selects
    if (fontSelect) {
        fontSelect.addEventListener('change', (e) => {
            saveTextPref(FONT_FAMILY_KEY, e.target.value);
            applyTextStyle(); 
        });
    }

    if (sizeSelect) {
        sizeSelect.addEventListener('change', (e) => {
            saveTextPref(FONT_SIZE_KEY, e.target.value);
            applyTextStyle(); 
        });
    }
    
    if (weightSelect) {
        weightSelect.addEventListener('change', (e) => {
            saveTextPref(FONT_WEIGHT_KEY, e.target.value);
            applyTextStyle(); 
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

// Scroll synchronization logic
document.addEventListener('DOMContentLoaded', initializeScrollSync);

function initializeScrollSync() {
    const contentArea = document.getElementById('scrollable-content-reading');
    const scrollRange = document.getElementById('scrollRange');
    
    startSync(contentArea, scrollRange);
}

function startSync(scrollableElement, scrollRange) {
    
    // 1. SCROLL SYNC: Text scroll -> Slider update
    scrollableElement.addEventListener('scroll', () => {
        const maxScroll = scrollableElement.scrollHeight - scrollableElement.clientHeight;
        const currentScroll = scrollableElement.scrollTop;

        if (maxScroll > 0) {
            // Calculate scroll percentage (0 to 100)
            const scrollPercentage = (currentScroll / maxScroll) * 100;
            scrollRange.value = scrollPercentage.toFixed(2);
        } else {
            scrollRange.value = 0;
        }
    });

    // 2. SCROLL SYNC: Slider -> Text scroll
    scrollRange.addEventListener('input', () => {
        const sliderValue = parseFloat(scrollRange.value);
        const maxScroll = scrollableElement.scrollHeight - scrollableElement.clientHeight;
        
        if (maxScroll > 0) {
            // Calculate new scrollTop value
            const newScrollTop = (sliderValue / 100) * maxScroll;
            scrollableElement.scrollTop = newScrollTop;
        }
    });
    
    // Set initial value
    scrollableElement.dispatchEvent(new Event('scroll'));
}
