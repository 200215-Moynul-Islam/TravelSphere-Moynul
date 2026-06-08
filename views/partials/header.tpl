<header class="main-header">
    <div class="container header-container">
        <div class="left-section">
            <div class="logo-section">
                <a href="/" class="brand-logo">TravelSphere</a>
            </div>

            <nav class="nav-menu">
                <ul>
                    {{range .NavItems}}
                    <li class="{{if eq $.CurrentNav .Key}}active{{end}}">
                        <a href="{{.URL}}">{{.Name}}</a>
                    </li>
                    {{end}}
                </ul>
            </nav>
        </div>

        <div class="auth-section">
            {{if .IsAuthenticated}}
                <span class="user-greeting">Hi, {{.Username}}</span>
                <button class="btn-logout">Logout</button>
            {{else}}
                <button class="btn-login">Login</button>
            {{end}}
        </div>
    </div>
</header>
