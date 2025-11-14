const cardList = document.getElementById("cardList");
const range = document.getElementById("customRange1");

function updateAction() {
  const nameEl = document.getElementById('name');
  const form = document.getElementById('bookForm');
  const name = nameEl ? nameEl.value : '';
  if (form) form.action = form.action.replace(/\/+$/, '') + '/' + encodeURIComponent(name);
}

// Only wire scroll syncing if both elements exist
if (cardList && range) {
  // compute a safe max (guard against NaN)
  const max = Math.max(0, cardList.scrollWidth - cardList.clientWidth);
  range.max = max;

  range.addEventListener("input", () => {
    cardList.scrollLeft = Number(range.value) || 0;
  });

  cardList.addEventListener("scroll", () => {
    range.value = cardList.scrollLeft;
  });
}
