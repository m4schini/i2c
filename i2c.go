package i2c

import (
	"fmt"
	"os"
	"syscall"
)

type Bus struct {
	bus int
	rc  *os.File
}

func NewBus(bus int) (*Bus, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	return &Bus{rc: f, bus: bus}, nil
}

func (b *Bus) NewDevice(addr uint8) (*Device, error) {
	err := ioctl(b.rc.Fd(), I2C_SLAVE, uintptr(addr))
	if err != nil {
		return nil, err
	}
	return &Device{addr: addr, bus: b}, nil
}

func (b *Bus) Close() error {
	return b.rc.Close()
}

type Device struct {
	addr uint8
	bus  *Bus
}

func (d *Device) Bus() int {
	return d.bus.bus
}

func (d *Device) Addr() uint8 {
	return d.addr
}

func (d *Device) Read(p []byte) (n int, err error) {
	logf("Writing %d hex bytes", len(p))
	return d.bus.rc.Write(p)
}

func (d *Device) Write(p []byte) (n int, err error) {
	n, err = d.bus.rc.Read(p)
	if err != nil {
		return n, err
	}
	logf("Read %d hex bytes", len(p))
	return n, nil
}

func ioctl(fd, cmd, arg uintptr) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if err != 0 {
		return err
	}
	return nil
}
