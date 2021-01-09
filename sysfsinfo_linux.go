package hid

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type SysfsInformation struct {
	Serial string
	Manufacturer string
	Product string
}

var reSysBus = regexp.MustCompile(`^/sys/bus/usb/devices/usb(\d+)$`)
var reSysBusDevice = regexp.MustCompile(`^/sys/bus/usb/devices/usb(\d+)/(\d+)-(\d+)$`)


func readFileNoErr(file string) string {
	b, _ := ioutil.ReadFile(file)
	return strings.TrimSuffix(string(b), "\n")
}

func readFileIntNoErr(file string) int {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return 0
	}
	num, _ := strconv.Atoi(strings.TrimSuffix(string(b), "\n"))
	return num
}

func readDeviceProperties(bus int, dev int) SysfsInformation {
	ret := SysfsInformation{}
	found := false
	filepath.Walk("/sys/bus/usb/devices/", func(busPath string, info os.FileInfo, err error) error {

		if ! reSysBus.MatchString(busPath) || found {
			return nil
		}
		filepath.Walk(busPath + "/", func(path string, info os.FileInfo, err error) error {

			if ! reSysBusDevice.MatchString(path) || found {
				return nil
			}
			devId := readFileIntNoErr(filepath.Join(path, "devnum"))
			busId := readFileIntNoErr(filepath.Join(path, "busnum"))
			if devId != dev || busId != bus {
				return nil
			}

			found = true
			ret.Serial = readFileNoErr(filepath.Join(path, "serial"))
			ret.Manufacturer = readFileNoErr(filepath.Join(path, "manufacturer"))
			ret.Product = readFileNoErr(filepath.Join(path, "product"))

			return nil
		})
		return nil
	})

	return ret
}


