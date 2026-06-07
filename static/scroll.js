// Scroll synchronization logic
document.addEventListener('DOMContentLoaded', initializeScrollSync);

function initializeScrollSync() {
    const contentArea = document.getElementById('scrollable-content-reading');
    const scrollRange = document.getElementById('scrollRange');
    
    startSync(contentArea, scrollRange);
}

function startSync(scrollableElement, scrollRange) {
    
    // SCROLL SYNC: Text scroll -> Slider update
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

    // SCROLL SYNC: Slider -> Text scroll
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
