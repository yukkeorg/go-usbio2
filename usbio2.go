package usbio2

import (
    "fmt"
    "log"
    "bytes"
    "encoding/binary"
    "github.com/GeertJohan/go.hid"
)

const (
    USB_VENDOR       uint16 = 0x1352
    USB_PRODUCT_ORIG uint16 = 0x0120
    USB_PRODUCT_AKI  uint16 = 0x0121

    COMMANDSIZE  int = 64

    // Command
    CMD_WRITEREAD     byte = 0x20
    CMD_READWRITE     byte = 0x21  // for original
    CMD_REPEATWRITE   byte = 0x22  // for original

    CMD_READFROMFLASH byte = 0xF0
    CMD_WRITETOFLASH  byte = 0xF1

    CMD_READCONFIG    byte = 0xF8
    CMD_WEIRWCONFIG   byte = 0xF9

    // Config
    CNF_P2_PULLUP_DEFALUT byte = 0x00
    CNF_P2_PULLUP_ENABLE  byte = 0x00
    CNF_P2_PULLUP_DISABLE byte = 0x01

    CNF_P1_PIN_DEFAULT  byte = 0x00
    CNF_P2_PIN_DEFAULT  byte = 0x0F
)


/*
*/

type UsbIO2 struct {
    dev *hid.Device
    name string
    seq byte
}


func NewUsbIO2() (*UsbIO2, error) {
    usbio := &UsbIO2{name: ""}
    if err := usbio.openUsbIo2(); err != nil {
        return nil,err
    }
    return usbio, nil
}


func (self *UsbIO2) openUsbIo2() error {
    var err error

    self.dev, err = hid.Open(USB_VENDOR, USB_PRODUCT_ORIG, "")
    if err == nil {
        self.name = "Original USB-IO2"
        return nil
    }

    self.dev, err = hid.Open(USB_VENDOR, USB_PRODUCT_AKI, "")
    if err == nil {
        self.name = "Akizuki USB-IO2"
        return nil
    }

    return fmt.Errorf("Error: openUsbIo2: USB-IO2 is not found.")
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


func (self *UsbIO2) CreateCommand(requestId byte) ([]byte, error) {
    return self.CreateCommandWithData(requestId, make([]byte, 0))
}


func (self *UsbIO2) CreateCommandWithData(requestId byte, data []byte) ([]byte, error) {
    if len(data) > COMMANDSIZE - 2 {
        return nil, fmt.Errorf("len(data) is bigger than %d.", COMMANDSIZE - 2)
    }

    command := make([]byte, COMMANDSIZE)
    command[0] = requestId
    copy(command[1:], data)
    command[COMMANDSIZE - 1] = self.seq

    return command, nil
}


func (self *UsbIO2) Write(data []byte) error {
    write_size, err := self.dev.Write(data)
    if err != nil {
        return err
    }

    if write_size != COMMANDSIZE {
        log.Printf("write: WARN: Strainge write data size : %d", write_size)
    }

    return nil
}


func (self *UsbIO2) Read() ([]byte, error) {
    read_data := make([]byte, COMMANDSIZE)
    read_size, err := self.dev.Read(read_data)
    if err != nil {
        return nil, err
    }

    if read_size != COMMANDSIZE {
        log.Printf("read: WARN: Strainge read data size: %d", read_size)
    }
    if read_data[COMMANDSIZE - 1] != self.seq {
        log.Printf("read: WARN: Don't match sequence number")
    }

    return read_data, nil
}

func (self *UsbIO2) WriteRead(command []byte) ([]byte, error) {
    err := self.Write(command)
    if err != nil {
        return nil, err
    }

    read_data, err2 := self.Read()
    if err != nil {
        return nil, err2
    }

    return read_data, nil
}


// Utility Function
func (self *UsbIO2) GetPortStatus() (uint16, error) {
    var err error

    command, _ := self.CreateCommand(CMD_WRITEREAD)

    read_data, err := self.WriteRead(command)
    if err != nil {
        return 0, err
    }

    var ret uint16
    buf := bytes.NewReader(read_data[1:3])
    err = binary.Read(buf, binary.LittleEndian, &ret)
    if err != nil {
        return 0, fmt.Errorf("binary.Read: %s", err)
    }

    self.seq++

    return ret, nil
}


func (self *UsbIO2) GetConfig() (*UsbIO2Config, error) {
    command, _ := self.CreateCommand(CMD_READCONFIG)

    read_data, err := self.WriteRead(command)
    if err != nil {
        return nil, err
    }

    config := NewUsbIO2Config()
    config.FromBytes(read_data[1:63])

    self.seq++

    return config, nil

}



/*
*/
type UsbIO2Config struct {
    // "Disable" Port2 pullup settings
    // bit0 on: Disable, bit0 off: Enable
    EnablePullUp bool

    // Pin Settings
    // true: input, false: output
    PinPort1 [8]bool
    PinPort2 [8]bool
}

func NewUsbIO2Config() *UsbIO2Config {
    return &UsbIO2Config{}
}

func (self *UsbIO2Config) FromBytes(buf []byte) {
    if (buf[1] & 0x01) == 0x01 {
        self.EnablePullUp = false
    } else {
        self.EnablePullUp = true
    }

    self.PinPort1 = byte_to_boolarray(buf[4])
    self.PinPort2 = byte_to_boolarray(buf[5])
}


func (self *UsbIO2Config) ToBytes() []byte {
    buf := make([]byte, COMMANDSIZE - 2)
    if self.EnablePullUp {
        buf[1] = 0x00
    } else {
        buf[1] = 0x01
    }

    buf[4] = boolarray_to_byte(self.PinPort1)
    buf[5] = boolarray_to_byte(self.PinPort2)

    return buf
}


func boolarray_to_byte(ba [8]bool) byte {
    var b byte = 0
    for i, p := range ba {
        if p {
            b |= (1 << uint(i))
        }
    }
    return b
}


func byte_to_boolarray(b byte) [8]bool {
    var ba [8]bool
    for i := 0; b != 0; i++ {
        if (b & 0x01) == 0x01 {
            ba[i] = true
        } else {
            ba[i] = false
        }
        b = b >> 1
    }
    return ba
}
