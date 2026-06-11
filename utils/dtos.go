package utils

type CountryDTO struct {
	Cca3 string `json:"cca3"`
	Population int64 `json:"population"`
	Region string `json:"region"`
	Subregion string `json:"subregion"`
	Name struct {
		Common string `json:"common"`
		Official string `json:"official"`
	} `json:"name"`
	Flag struct {
		Png string `json:"url_png"`
	} `json:"flag"`
	Capitals []struct {
		Name string `json:"name"`
	} `json:"capitals"`
	Currencies []struct {
		Code string `json:"code"`
		Name string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Languages []struct {
		Name string `json:"name"`
	} `json:"languages"`
}