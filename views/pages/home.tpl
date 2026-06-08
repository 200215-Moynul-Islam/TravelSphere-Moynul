<div class="container">
    <h2 class="section-title">Featured destinations</h2>
    
    <div class="destinations-grid">
        {{range .FeaturedCountries}}
        <a href="/country/{{.Code}}" class="destination-card-link">
            <div class="destination-card">
                <div class="flag-wrapper">
                    <img src="{{.Flag}}" alt="Flag of {{.Name}}" loading="lazy">
                </div>
                <div class="card-info">
                    <h3 class="country-name">{{.Name}}</h3>
                    <p class="country-meta">
                        {{if .Capital}}{{.Capital}}, {{end}}{{.Region}}
                    </p>
                </div>
            </div>
        </a>
        {{else}}
        <div class="empty-state">
            <h3 class="empty-title">No Destinations Found</h3>
            <p class="empty-text">We couldn't load the featured destinations right now. Please try reloading the page later.</p>
        </div>
        {{end}}
    </div>
</div>
