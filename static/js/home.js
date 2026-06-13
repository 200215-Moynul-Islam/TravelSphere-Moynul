const searchInput = document.getElementById("hero-country-search");
const resultsDropdown = document.getElementById("hero-search-results");
let debounceTimer;

searchInput.addEventListener("input", () => {
  clearTimeout(debounceTimer);
  const query = searchInput.value.trim();

  if (query.length === 0) {
    resultsDropdown.classList.add("hidden");
    return;
  }

  debounceTimer = setTimeout(async () => {
    resultsDropdown.innerHTML = "";

    try {
      const response = await fetch(
        `/api/countries?search=${encodeURIComponent(query)}`,
      );
      if (!response.ok) {
        throw new Error(`HTTP status error: ${response.status}`);
      }

      const res = await response.json();

      if (res.success && res.data && res.data.length > 0) {
        res.data.forEach((country) => {
          const li = document.createElement("li");
          const capitalText = country.capital ? ` - ${country.capital}` : "";
          li.textContent = `${country.name}${capitalText}`;

          li.addEventListener("click", () => {
            window.location.href = `/countries/${country.code}`;
          });
          resultsDropdown.appendChild(li);
        });
        resultsDropdown.classList.remove("hidden");
      } else {
        resultsDropdown.innerHTML =
          '<li class="no-match">No matching destinations found</li>';
        resultsDropdown.classList.remove("hidden");
      }
    } catch (err) {
      console.error("Error matching countries:", err);
      resultsDropdown.innerHTML =
        '<li class="no-match" style="color: #ef4444;">An error occurred while searching</li>';
      resultsDropdown.classList.remove("hidden");
    }
  }, 300);
});

document.addEventListener("click", (e) => {
  if (!searchInput.contains(e.target) && !resultsDropdown.contains(e.target)) {
    resultsDropdown.classList.add("hidden");
  }
});
