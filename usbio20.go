package usbio2

import (
	"fmt"
	"log"
	"bytes"
    "time"
	"encoding/binary"
	"github.com/bearsh/hid"
)

const (
	USB_VENDOR       uint16 = 0x1352
	USB_PRODUCT_ORIG uint16 = 0x0120 // ORIGINAL
	USB_PRODUCT_AKI  uint16 = 0x0121 // AKIZUKI Compatible

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
	CNF_P2_PULLUP_ENABLE  byte = 0x00
	CNF_P2_PULLUP_DISABLE byte = 0x01
	CNF_P2_PULLUP_DEFALUT byte = CNF_P2_PULLUP_ENABLE

	CNF_P1_PIN_DEFAULT  byte = 0x00
	CNF_P2_PIN_DEFAULT  byte = 0x0F
)


type UsbIO2 struct {
	dev *hid.Device
	name string
	seq byte
}


func NewUsbIO2() (*UsbIO2, error) {
	usbio := &UsbIO2{name: ""}
	if err := usbio.openUsbIo2(); err != nil {
		return nil, err
	}
	return usbio, nil
}

func (self *UsbIO2) openUsbIo2() error {
	var detect_devices []hid.DeviceInfo

	self.dev = nil
	self.name = ""

	detect_devices = hid.Enumerate(USB_VENDOR, USB_PRODUCT_ORIG)
	if len(detect_devices) > 0 {
        detect_device := detect_devices[0]
		dev, err := detect_device.Open()
		if err != nil {
			return err
		}
		self.dev = dev
		self.name = "Original USB-IO2.0"
		return nil
	}

	detect_devices = hid.Enumerate(USB_VENDOR, USB_PRODUCT_AKI)
	if len(detect_devices) > 0 {
        detect_device := detect_devices[0]
		dev, err := detect_device.Open()
		if err != nil {
			return err
		}
		self.dev = dev
		self.name = "Akizuki USB-IO2.0"
		return nil
	}

	return fmt.Errorf("Error: openUsbIo2: USB-IO2.0 is not found.")
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
	return self.CreateCommandWithData(requestId, []byte{})
}

func (self *UsbIO2) CreateCommandWithData(requestId byte, data []byte) ([]byte, error) {
	if len(data) > COMMANDSIZE - 2 {
		return nil, fmt.Errorf("len(data) is bigger than %d.",
		                       COMMANDSIZE - 2)
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

	time.Sleep(1 * time.Millisecond)

	read_data, err := self.Read()
	if err != nil {
		return nil, err
	}

	self.seq++

	return read_data, nil
}

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
