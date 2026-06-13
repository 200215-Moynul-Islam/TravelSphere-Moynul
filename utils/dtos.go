package utils

type CountryDTO struct {
	Codes struct {
		Alpha3 string `json:"alpha_3"`
	} `json:"codes"`
	Population int64 `json:"population"`
	Region string `json:"region"`
	Subregion string `json:"subregion"`
	Names struct {
		Common string `json:"common"`
		Official string `json:"official"`
	} `json:"names"`
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