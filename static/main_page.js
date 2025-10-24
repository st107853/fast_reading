const cardList = document.getElementById("cardList");
const range = document.getElementById("customRange1");

function updateAction() {
    var name = document.getElementById('name').value;
    document.getElementById('bookForm').action += name;
}

range.max = cardList.scrollWidth - cardList.clientWidth;

range.addEventListener("input", () => {
  cardList.scrollLeft = range.value;

});

cardList.addEventListener("scroll", () => {
  range.value = cardList.scrollLeft;
});
