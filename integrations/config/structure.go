package config

// Config is the overall configuration
type Config struct {
	RedisData   Redis        `yaml:"redis"`
	LutronData  LutronConfig `yaml:"lutron"`
	MyQData     MyQ          `yaml:"myq"`
	HarmonyData Harmony      `yaml:"harmony"`
	KasaData    Kasa         `yaml:"kasa"`
	SleepIqData SleepIq      `yaml:"sleepiq"`
	Life360Data Life360      `yaml:"life360"`
	ZwaveData   Zwave        `yaml:"zwave"`
	MenuConfig  []Menu       `yaml:"menu"`
}

// Redis is the configuration for Redis
type Redis struct {
	URL  string `yaml:"url"`
	Port string `yaml:"port"`
}

// LutronConfig is the configuration for the Lutron Hub
type LutronConfig struct {
	URL      string `yaml:"url"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	File     string `yaml:"file"`
}

// MyQ is the configuration for the MyQ Hub
type MyQ struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Harmony is the configuration for the Harmony Hub
type Harmony struct {
	IPAddress string `yaml:"ip_address"`
}

// Kasa is the configuration for the Kasa integration
type Kasa struct {
	Devices []string `yaml:"devices"`
}

// SleepIq is the configuration for the sleep number bed
type SleepIq struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Life360 is the configuration for life60
type Life360 struct {
	Username            string `yaml:"username"`
	Password            string `yaml:"password"`
	AuthenticationToken string `yaml:"authenticationToken"`
}

// Zwave is the configuration for life60
type Zwave struct {
	URL string `yaml:"url"`
}

// Menu is the setup for the configuration to display on the Web UI
type Menu struct {
	Index    int      `yaml:"index"`
	AreaName string   `yaml:"areaName"`
	Lutron   []int    `yaml:"lutron"`
	Kasa     []string `yaml:"kasa"`
	Harmony  []struct {
		Activities []struct {
			State    string `yaml:"state"`
			Activity string `yaml:"activity"`
		} `yaml:"activities"`
	} `yaml:"harmony"`
	Custom []struct {
		Type    string   `yaml:"type"`
		Name    string   `yaml:"name"`
		Devices []string `yaml:"devices"`
	} `yaml:"custom"`
}
