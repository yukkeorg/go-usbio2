package main

import (
	"log"
	"github.com/yukkeorg/go-usbio2"
)

func main() {
	usbio, err := usbio2.NewUsbIO2()
	if err != nil {
		log.Printf("USB-IO2.0 not detected.")
		return
	}
	defer usbio.Cleanup()

	config, err := usbio.GetConfig()
	if err != nil {
		log.Printf("Can't read device config.")
		return
	}

	log.Printf("Detected device: %s", usbio.GetDeviceName())
	log.Printf("Pullup Enabled: %s", bool2str(config.EnablePullUp, "ENABLE", "DISABLE"))
	log.Printf("Port1 Pin Setup")
	for i, b := range config.Port1 {
		log.Printf("  Pin%d: %s", i, bool2str(b, "INPUT", "OUTPUT"))
	}
	log.Printf("Port2 Pin Setup")
	for i, b := range config.Port2 {
		log.Printf("  Pin%d: %s", i, bool2str(b, "INPUT", "OUTPUT"))
	}
}

func bool2str(b bool, true_str string, false_str string) string {
	if(b) {
		return true_str
	} else {
		return false_str
	}
}
