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
    log.Printf("Pullup Enabled: %v", config.EnablePullUp)
    log.Printf("Port1 Pin Setup (value in 'ture' is a INPUT)")
    for i, b := range config.Port1 {
        log.Printf("  Pin%d: %v", i, b)
    }
    log.Printf("Port2 Pin Setup (value in 'ture' is a INPUT)")
    for i, b := range config.Port2 {
        log.Printf("  Pin%d: %v", i, b)
    }
}
