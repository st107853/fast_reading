document.addEventListener('DOMContentLoaded', () => {
    const favouriteButton = document.getElementById('favBtn');
    const favouriteStatus = document.getElementById('favStatus');

    if (!favouriteButton) return;

    const isLoggedIn = favouriteButton.dataset.loggedIn === 'true';
    const bookId = Number(favouriteButton.dataset.bookId);

    if (!Number.isInteger(bookId) || bookId <= 0) {
        favouriteButton.disabled = true;
        console.error('Invalid book ID for favourite toggle.');
        return;
    }

    favouriteButton.addEventListener('change', async () => {
        const previousState = !favouriteButton.checked;

        if (!isLoggedIn) {
            favouriteButton.checked = previousState;
            if (favouriteStatus) {
                favouriteStatus.textContent = 'Please log in to manage favourites.';
            }
            return;
        }

        favouriteButton.disabled = true;

        try {
            const response = await fetch(`/library/book/${bookId}/favourite`, {
                method: 'POST'
            });

            if (!response.ok) {
                throw new Error(`Favourite request failed with status ${response.status}`);
            }

        } catch (error) {
            favouriteButton.checked = previousState;
            if (favouriteStatus) {
                favouriteStatus.textContent = 'Unable to update favourites.';
            }
            console.error(error);
        } finally {
            favouriteButton.disabled = false;
            if (favouriteStatus) {
                favouriteStatus.textContent = '';
            }
        }
    });
});
