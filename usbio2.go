package usbio2

import (
    "fmt"
    "log"
    "bytes"
    "encoding/binary"
    "github.com/GeertJohan/go.hid"
)

const (
    USBIO2_VENDOR       uint16 = 0x1352
    USBIO2_PRODUCT_ORIG uint16 = 0x0120
    USBIO2_PRODUCT_AKI  uint16 = 0x0121

    USBIO2_COMMANDSIZE  int = 64
)

type UsbIO2 struct {
    dev *hid.Device
    name string
    seq byte
}

func NewUsbIO2() (*UsbIO2, error) {
    usbio := &UsbIO2{name: ""}
    err := usbio.openUsbIo2()
    if err != nil {
        return nil,err
    }
    return usbio, nil
}


func (self *UsbIO2) GetDeviceName() string {
    return self.name
}

func (self *UsbIO2) Cleanup() {
    if self.dev != nil {
        self.dev.Close()
    } else {
        log.Printf("self.dev is nil")
    }
}

func (self *UsbIO2) GetPortStatus() (uint16, error) {
    var err error

    command := self.createCommand(0x20)
    err = self.write(command)
    if err != nil {
        return 0, err
    }

    read_data, err := self.read()
    if err != nil {
        return 0, err
    }

    var ret uint16
    buf := bytes.NewReader(read_data[1:3])
    err = binary.Read(buf, binary.LittleEndian, &ret)
    if err != nil {
        return 0, fmt.Errorf("binary.Read: %s", err)
    }

    self.incrementSequence()

    return ret, nil
}

// ----

func (self *UsbIO2) openUsbIo2() error {
    var err error

    self.dev, err = hid.Open(USBIO2_VENDOR, USBIO2_PRODUCT_ORIG, "")
    if err == nil {
        self.name = "Original USB-IO2"
        return nil
    }

    self.dev, err = hid.Open(USBIO2_VENDOR, USBIO2_PRODUCT_AKI, "")
    if err == nil {
        self.name = "Akizuki USB-IO2"
        return nil
    }

    return fmt.Errorf("Error: openUsbIo2: USB-IO2 is not found.")
}


func (self *UsbIO2) createCommand(requestId byte) []byte {
    command := make([]byte, USBIO2_COMMANDSIZE)
    command[0] = requestId
    command[USBIO2_COMMANDSIZE - 1] = self.seq
    return command
}

func (self *UsbIO2) write(data []byte) error {
    write_size, err := self.dev.Write(data)
    if err != nil {
        return err
    }

    if write_size != USBIO2_COMMANDSIZE {
        log.Printf("write: WARN: Strainge write data size : %d", write_size)
    }

    return nil
}

func (self *UsbIO2) read() (data []byte, err error) {
    read_data := make([]byte, USBIO2_COMMANDSIZE)
    read_size, err := self.dev.Read(read_data)
    if err != nil {
        return nil, err
    }

    if read_size != USBIO2_COMMANDSIZE {
        log.Printf("read: WARN: Strainge read data size: %d", read_size)
    }
    if read_data[USBIO2_COMMANDSIZE - 1] != self.seq {
        log.Printf("read: WARN: Don't match sequence number")
    }

    return read_data, nil
}

func (self *UsbIO2) incrementSequence() {
    self.seq++
}
