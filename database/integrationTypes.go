package database

// LutronConfig is the configuration for the Lutron Hub
type LutronConfig struct {
	URL      string `json:"url"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	File     string `json:"file"`
}