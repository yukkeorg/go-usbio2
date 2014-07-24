package main

import (
    "log"
    "github.com/yukkeorg/go-usbio2"
)

func main() {
    uio, err := usbio2.NewUsbIO2()
    if err != nil {
        return
    }

    defer uio.Cleanup()

    config, err := uio.GetConfig()
    if err != nil {
        return
    }

    log.Printf("Detected: %s", uio.GetDeviceName())
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
