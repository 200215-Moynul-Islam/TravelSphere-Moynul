const searchInput = document.getElementById("search-input");
const regionSelect = document.getElementById("region-select");
const gridContainer = document.getElementById("countries-grid-container");
const filterEmptyState = document.getElementById("filter-empty-state");
let debounceTimer;

async function fetchFilteredCountries() {
  const search = searchInput.value.trim();
  const region = regionSelect.value;

  gridContainer.innerHTML = `
    <div class="empty-explorer-state" style="grid-column: 1 / -1;">
      <h3>Loading...</h3>
      <p>Fetching your next destinations...</p>
    </div>
  `;
  filterEmptyState.classList.add("hidden");

  try {
    const response = await fetch(
      `/api/countries?search=${encodeURIComponent(
        search
      )}&region=${encodeURIComponent(region)}`
    );
    if (!response.ok) {
      throw new Error(`HTTP status error: ${response.status}`);
    }

    const res = await response.json();
    gridContainer.innerHTML = "";

    if (res.success && res.data && res.data.length > 0) {
      filterEmptyState.classList.add("hidden");

      res.data.forEach((country) => {
        const card = document.createElement("a");

        const countryCode = country.code || country.cca3;
        const countryName = country.name?.common || country.name;
        const flagUrl = country.flag || country.flags?.png || "";
        const capitalText = country.capital ? country.capital : "N/A";
        const populationText = country.population
          ? country.population.toLocaleString()
          : "0";
        const currencyText = country.currency || "N/A";
        const languagesText = country.languages || "N/A";

        card.href = `/countries/${countryCode}`;
        card.className = "country-card";
        card.setAttribute("data-name", countryName);
        card.setAttribute("data-capital", capitalText);
        card.setAttribute("data-region", country.region || "");

        card.innerHTML = `
          <div class="card-flag-wrapper">
            <img src="${flagUrl}" alt="Flag of ${countryName}" loading="lazy">
          </div>
          <div class="card-content">
            <h3 class="card-country-name">${countryName}</h3>
            <div class="card-meta-list">
              <div class="card-meta-item">
                <span class="meta-key">Capital:</span>
                <span class="meta-val">${capitalText}</span>
              </div>
              <div class="card-meta-item">
                <span class="meta-key">Population:</span>
                <span class="meta-val">${populationText}</span>
              </div>
              <div class="card-meta-item">
                <span class="meta-key">Currency:</span>
                <span class="meta-val">${currencyText}</span>
              </div>
              <div class="card-meta-item">
                <span class="meta-key">Languages:</span>
                <span class="meta-val">${languagesText}</span>
              </div>
            </div>
          </div>
        `;
        gridContainer.appendChild(card);
      });
    } else {
      filterEmptyState.classList.remove("hidden");
    }
  } catch (err) {
    console.error("Error matching countries:", err);
    gridContainer.innerHTML = `
      <div class="empty-explorer-state" style="grid-column: 1 / -1; color: #ef4444;">
        <h3>Error Loading Destinations</h3>
        <p>An error occurred while updating the country list. Please try again later.</p>
      </div>
    `;
  }
}

searchInput.addEventListener("input", () => {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(fetchFilteredCountries, 300);
});

regionSelect.addEventListener("change", () => {
  fetchFilteredCountries();
});
