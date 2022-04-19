package hubitat

type HubitatData struct {
	Config Config
}

type Config struct {
	AccessKey  string `json:"accessKey"`
	HubitatUrl string `json:"url"`
}
