package omron

import (
	"os"

	"github.com/thomasf/bpchart/pkg/omron"
	"github.com/thomasf/lg"
	"github.com/truveris/gousb/usb"
)

var (
	cmdInit = []byte{0x02, 0x08, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x18}
	cmdData = []byte{0x02, 0x08, 0x01, 0x00, 0x02, 0xac, 0x28, 0x00, 0x8f}
	cmdDone = []byte{0x02, 0x08, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x07}
	cmdFail = []byte{0x02, 0x08, 0x0f, 0x0f, 0x0f, 0x0f, 0x00, 0x00, 0x08}

// unsigned char rawdata[64];
// unsigned char payload[40*70];
)

var usbdevice *usb.Device

func Open() error {
	ctx := usb.NewContext()
	// ctx.Debug(10)
	devs, err := ctx.ListDevices(func(desc *usb.Descriptor) bool {
		if desc.Vendor.String() == "0590" {
			return true
		}
		return false
	})
	if err != nil {
		for _, v := range devs {
			err := v.Close()
			if err != nil {
				lg.Errorln(err)
			}
		}
		return err
	}
	usbdevice = devs[0]
	return nil
}

func Close() error {
	return usbdevice.Close()
}

// func buildCRC(data []byte) int {
// 	crc := int(0);
// 	int len = data[1];

// 	while(--len)
// 	{
// 		crc ^= data[len];
// 	}

// 	return crc;
// }

func Read(bank int) ([]omron.Entry, error) {
	var entries []omron.Entry

	const (
		config      = 0x1
		iface       = 0x0
		setup       = 0x0
		endpointOut = 2
		endpointIn  = 0x81
	)
	epOut, err := usbdevice.OpenEndpoint(config, iface, setup, endpointOut|uint8(usb.ENDPOINT_DIR_OUT))
	// dev, err := usbdevice.Open()
	if err != nil {
		return entries, err
	}

	epIn, err := usbdevice.OpenEndpoint(config, iface, setup, endpointIn|uint8(usb.ENDPOINT_DIR_IN))
	if err != nil {
		return entries, err
	}

	n, err := epOut.Write(cmdInit)
	if err != nil {
		return entries, err
	}
	lg.Infoln("wrote", n, "bytes")
	// dev.Configs()
	var readBytes []byte
	n, err = epIn.Read(readBytes)
	if err != nil {
		return entries, err
	}
	lg.Infoln("read", n, "bytes")

	os.Exit(1)
	addr := (0x02AC)

	for i := 0; i < 70; i++ {
		// if (abort) {}

		cmdData[4] = byte(addr >> 8)
		cmdData[5] = byte(addr & 0xFF)

	}

	os.Exit(1)
	return entries, nil
}
