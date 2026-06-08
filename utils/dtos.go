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
	Flags struct {
		Png string `json:"png"`
	} `json:"flags"`
	Capital []string `json:"capital"`
	Currencies map[string]struct {
		Name string `json:"name"`
	} `json:"currencies"`
	Languages map[string]string `json:"languages"`
}
