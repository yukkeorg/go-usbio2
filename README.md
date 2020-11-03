go-usbio2
=========


Dependencies
------------

- go.hid (https://github.com/GeertJohan/go.hid)


Install
-------

``` sh
export GOPATH=/path/to/go-develop-directory
go get github.com/yukkeorg/go-usbio2
```

Usage
-----

``` golang
import "fmt"
import "github.com/yukkeorg/usbio2"

func main() {
	usbio, err := usbio2.NewUsbIO2()
	if err != nil {
		return
	}
	deffer usbio.Cleanup()

	fmt.Printf("Device Name: %s\n", usbio.GetDeviceName())

	r, err := usbio.WriteRead([]byte{})
	if err != nil {
		return
	}

	fmt.Printf("ReadData : %v\n", r)
}
```

License
-------

MIT
