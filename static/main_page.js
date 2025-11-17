function updateAction() {
  const nameEl = document.getElementById('name');
  const form = document.getElementById('bookForm');
  const name = nameEl ? nameEl.value : '';
  if (form) form.action = form.action.replace(/\/+$/, '') + '/' + encodeURIComponent(name);
}

function initScrollSync() {
  const cardList = document.getElementById("cardList");
  const range = document.getElementById("customRange1");

  if (!cardList || !range) return;

  // Function to recalculate and set range max
  const syncRangeMax = () => {
    const max = Math.max(0, cardList.scrollWidth - cardList.clientWidth);
    range.max = String(max); // Ensure it's a string for the attribute
  };

  // Initial sync after a brief delay to allow layout to settle
  setTimeout(syncRangeMax, 100);

  // Sync when range input changes
  range.addEventListener("input", () => {
    cardList.scrollLeft = Number(range.value) || 0;
  });

  // Sync when cardList scrolls
  cardList.addEventListener("scroll", () => {
    range.value = cardList.scrollLeft;
  });

  // Recalculate max on window resize
  window.addEventListener("resize", syncRangeMax);

  // Also try to recalculate when images load
  const images = cardList.querySelectorAll('img');
  images.forEach(img => {
    img.addEventListener('load', syncRangeMax);
  });
}

// Wait for DOM to be ready
document.addEventListener('DOMContentLoaded', initScrollSync);
// Also call immediately in case DOMContentLoaded already fired
if (document.readyState === 'loading') {
  // DOM still loading, wait for event
} else {
  // DOM already loaded
  initScrollSync();
}
