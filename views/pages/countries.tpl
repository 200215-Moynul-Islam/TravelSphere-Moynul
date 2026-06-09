<div class="container">
    <header class="explorer-header">
        <h1 class="explorer-title">Country Explorer</h1>
        <p class="explorer-subtitle">Browse every destination on first load. Search and filter update only the results below — no full page reload.</p>
    </header>

    <div class="filter-bar">
        <div class="filter-group">
            <label class="filter-label" for="search-input">SEARCH</label>
            <input type="text" id="search-input" placeholder="Country or capital..." autocomplete="off">
        </div>
        <div class="filter-group">
            <label class="filter-label" for="region-select">REGION</label>
            <select id="region-select">
                <option value="all">All regions</option>
                <option value="africa">Africa</option>
                <option value="americas">Americas</option>
                <option value="asia">Asia</option>
                <option value="europe">Europe</option>
                <option value="oceania">Oceania</option>
            </select>
        </div>
    </div>

    <div class="countries-grid" id="countries-grid-container">
        {{range .Countries}}
        <a href="/countries/{{.Code}}" class="country-card" data-name="{{.Name}}" data-capital="{{.Capital}}" data-region="{{.Region}}">
            <div class="card-flag-wrapper">
                <img src="{{.Flag}}" alt="Flag of {{.Name}}" loading="lazy">
            </div>
            <div class="card-content">
                <h3 class="card-country-name">{{.Name}}</h3>
                <div class="card-meta-list">
                    <div class="card-meta-item">
                        <span class="meta-key">Capital:</span>
                        <span class="meta-val">{{if .Capital}}{{.Capital}}{{else}}N/A{{end}}</span>
                    </div>
                    <div class="card-meta-item">
                        <span class="meta-key">Population:</span>
                        <span class="meta-val">{{.Population}}</span>
                    </div>
                    <div class="card-meta-item">
                        <span class="meta-key">Currency:</span>
                        <span class="meta-val">{{.Currency}}</span>
                    </div>
                    <div class="card-meta-item">
                        <span class="meta-key">Languages:</span>
                        <span class="meta-val">{{.Languages}}</span>
                    </div>
                </div>
            </div>
        </a>
        {{else}}
        <div class="empty-explorer-state" id="global-empty-state">
            <h3>No countries available</h3>
            <p>We encountered an issue downloading the destination list. Please try again later.</p>
        </div>
        {{end}}
    </div>
    
    <div class="empty-explorer-state hidden" id="filter-empty-state">
        <h3>No matching destinations</h3>
        <p>No countries match your selected search criteria or filter combinations.</p>
    </div>
</div>
