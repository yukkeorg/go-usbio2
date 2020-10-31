package usbio2


type UsbIO2Config struct {
	// "Disable" Port2 pullup settings
	// bit0 on: Disable, bit0 off: Enable
	EnablePullUp bool

	// Pin Settings
	// true: input, false: output
	Port1 []bool
	Port2 []bool
}

func NewUsbIO2Config() *UsbIO2Config {
	return &UsbIO2Config{
		false,
		make([]bool, 8, 8),
		make([]bool, 8, 8),
	}
}

func (self *UsbIO2Config) FromBytes(buf []byte) {
	if (buf[1] & 0x01) == 0x01 {
		self.EnablePullUp = false
	} else {
		self.EnablePullUp = true
	}

	self.Port1 = byte_to_boolarray(buf[4])
	self.Port2 = byte_to_boolarray(buf[5])
}

func (self *UsbIO2Config) ToBytes() []byte {
	buf := make([]byte, COMMANDSIZE - 2)
	if self.EnablePullUp {
		buf[1] = 0x00
	} else {
		buf[1] = 0x01
	}

	buf[4] = boolarray_to_byte(self.Port1)
	buf[5] = boolarray_to_byte(self.Port2)

	return buf
}
