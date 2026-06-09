package constants

const (
	NavHome= "home"
	NavCountries = "countries"
	NavWishlist = "wishlist"
	NavDashboard = "dashboard"
)

type NavItem struct {
	Name string
	URL string
	Key string
}

var NavigationItems = []NavItem{
	{Name: "Home", URL: "/", Key: NavHome},
	{Name: "Countries", URL: "/countries", Key: NavCountries},
	{Name: "Wishlist", URL: "/wishlist", Key: NavWishlist},
	{Name: "Dashboard", URL: "/dashboard", Key: NavDashboard},
}

var FeaturedCountryCodes = []string{"USA", "FRA", "JPN", "AUS", "BRA", "BGD"}
