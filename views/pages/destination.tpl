<div class="container">
    <header class="country-header-card">
        <div class="flag-wrapper">
            <img src="{{.CountryDetails.Flag}}" alt="Flag of {{.CountryDetails.Name}}" class="flag-img">
        </div>
        <div class="profile-details">
            <span class="region-badge">{{.CountryDetails.Region}}</span>
            <h1 class="country-title">{{.CountryDetails.Name}}</h1>
            
            <div class="meta-grid">
                <div class="meta-item">
                    <span class="meta-label">CAPITAL</span>
                    <span class="meta-value">{{if .CountryDetails.Capital}}{{.CountryDetails.Capital}}{{else}}N/A{{end}}</span>
                </div>
                <div class="meta-item">
                    <span class="meta-label">POPULATION</span>
                    <span class="meta-value">{{if .CountryDetails.Population}}{{.CountryDetails.Population}}{{else}}N/A{{end}}</span>
                </div>
                <div class="meta-item">
                    <span class="meta-label">CURRENCY</span>
                    <span class="meta-value">{{if .CountryDetails.Currency}}{{.CountryDetails.Currency}}{{else}}N/A{{end}}</span>
                </div>
                <div class="meta-item">
                    <span class="meta-label">LANGUAGES</span>
                    <span class="meta-value">{{if .CountryDetails.Languages}}{{.CountryDetails.Languages}}{{else}}N/A{{end}}</span>
                </div>
            </div>
        </div>
    </header>

    <div class="action-strip">
        <button class="wishlist-add-btn" id="add-to-wishlist-trigger" data-code="{{.CountryDetails.Code}}">Add to Wishlist</button>
    </div>

    <section class="panel-card attractions-sidebar">
        <h2 class="panel-heading">Attractions & landmarks</h2>
        <div class="attraction-rows-stack">
            {{range .Attractions}}
            <div class="detail-attraction-row">
                <div class="row-main-content">
                    <span class="landmark-title">{{.Name}}</span>
                    <span class="landmark-kinds">
                        {{range $index, $tag := .Tags}}
                            {{if $index}}, {{end}}{{$tag}}
                        {{end}}
                    </span>
                </div>
            </div>
            {{else}}
            <div class="empty-state-fallback">
                <p>No historical or cultural attractions cataloged around the capital region yet.</p>
            </div>
            {{end}}
        </div>
    </section>
</div>
