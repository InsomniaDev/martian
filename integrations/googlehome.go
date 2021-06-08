package integrations

import (
	"fmt"

	"github.com/evalphobia/google-home-client-go/googlehome"
)

type GoogleHome struct {
	Hostname string ""
	Lang     string
	Accent   string
}

func (gh *GoogleHome) Discover() error {
	dnsEntries := FindCastDNSEntries()
	fmt.Printf("Found %d cast devices\n", len(dnsEntries))
	for i, d := range dnsEntries {
		fmt.Printf("%d) device=%q device_name=%q address=\"%s:%d\" status=%q uuid=%q\n", i+1, d.Device, d.DeviceName, d.AddrV4, d.Port, d.Status, d.UUID)
	}
	return nil
}

func (gh *GoogleHome) Test() {
	cli, err := googlehome.NewClientWithConfig(googlehome.Config{
		// Hostname: "192.168.1.140", // office
		Hostname: "10.10.10.9", // downstairs
		// Hostname: "10.10.10.8", // upstairs
		Lang:   "en",
		Accent: "BB",
	})
	if err != nil {
		panic(err)
	}

	cli.SetVolume(0.3)

	// Speak text on Google Home.
	cli.Notify("Hello Brooke. What you doing?")

	// // Change language
	// cli.SetLang("ja")
	// cli.Notify("こんにちは、グーグル。")

	// // Or set language in Notify()
	// cli.Notify("你好、Google。", "zh")

	// //Play Audio
	// cli.Play("http://127.0.0.1/night.mp3")

	// //Stop Audio
	// cli.StopMedia()

	// //Min of 0.0 Max of 1.0 (Must be of type Float)

	//Get Volume Google Home is running at
	// cli.GetVolume()

	//Kills the running Application (Disconnects from Google Home)
	// defer cli.QuitApp()

}
