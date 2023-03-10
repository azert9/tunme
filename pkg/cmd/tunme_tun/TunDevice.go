//go:build app_tun

package tunme_tun

import (
	"fmt"
	"golang.org/x/sys/unix"
	"io"
)

// #include <string.h>
// #include <unistd.h>
// #include <fcntl.h>
// #include <linux/if.h>
// #include <linux/if_tun.h>
// #include <sys/ioctl.h>
//
// #define TUN_CLONE_DEVICE "/dev/net/tun"
//
// struct tun_device {
//     int fd;
//     char *name;
//     char name_buff[IFNAMSIZ];
// };
//
// struct tun_device create_tun_device(_GoString_ interface_name)
// {
//     struct tun_device tun_device;
//     memset(&tun_device, 0, sizeof(tun_device));
//     tun_device.name = tun_device.name_buff;
//
//     tun_device.fd = open(TUN_CLONE_DEVICE, O_RDWR);
//     if (tun_device.fd < 0)
//         return tun_device;
//
//     struct ifreq params;
//     memset(&params, 0, sizeof(params));
//     params.ifr_flags = IFF_TUN | IFF_NO_PI;
//     if (_GoStringLen(interface_name) != 0)
//         strncpy(params.ifr_name, _GoStringPtr(interface_name), IFNAMSIZ);
//
//     if (ioctl(tun_device.fd, TUNSETIFF, &params) == -1)
//     {
//         close(tun_device.fd);
//         tun_device.fd = -1;
//         return tun_device;
//     }
//
//     strcpy(tun_device.name_buff, params.ifr_name);
//
//     return tun_device;
// }
import "C"

type TunDevice interface {
	io.ReadWriteCloser
	Name() string
}

type _tunDevice struct {
	_fd   int
	_name string
}

func _newTunDeviceWithName(name string) (TunDevice, error) {

	tunDeviceInfo := C.create_tun_device(name)
	if tunDeviceInfo.fd < 0 {
		return nil, fmt.Errorf("failed to create TUN device") // TODO: Add details
	}

	return _tunDevice{
		_fd:   int(tunDeviceInfo.fd),
		_name: C.GoString(tunDeviceInfo.name),
	}, nil
}

func NewTunDevice() (TunDevice, error) {
	return _newTunDeviceWithName("")
}

// TODO: use
func NewTunDeviceWithName(name string) (TunDevice, error) {
	return _newTunDeviceWithName(name)
}

func (device _tunDevice) Close() error {
	return unix.Close(device._fd)
}

func (device _tunDevice) Read(data []byte) (int, error) {
	return unix.Read(device._fd, data)
}

func (device _tunDevice) Write(p []byte) (int, error) {
	return unix.Write(device._fd, p)
}

func (device _tunDevice) Name() string {
	return device._name
}

// TODO: remove
//import (
//	"bytes"
//	"fmt"
//	"golang.org/x/sys/unix"
//	"io"
//	"syscall"
//	"unsafe"
//)
//
//// #include <string.h>
//// #include <unistd.h>
//// #include <fcntl.h>
//// #include <linux/if.h>
//// #include <linux/if_tun.h>
//// #include <sys/ioctl.h>
////
//// #define TUN_CLONE_DEVICE "/dev/net/tun"
////
//// struct tun_device {
////     int fd;
////     char *name;
////     char name_buff[IFNAMSIZ];
//// };
////
//// struct tun_device create_tun_device(_GoString_ interface_name)
//// {
////     struct tun_device tun_device;
////     memset(&tun_device, 0, sizeof(tun_device));
////     tun_device.name = tun_device.name_buff;
////
////     tun_device.fd = open(TUN_CLONE_DEVICE, O_RDWR);
////     if (tun_device.fd < 0)
////         return tun_device;
////
////     struct ifreq params;
////     memset(&params, 0, sizeof(params));
////     params.ifr_flags = IFF_TUN | IFF_NO_PI;
////     if (_GoStringLen(interface_name) != 0)
////         strncpy(params.ifr_name, _GoStringPtr(interface_name), IFNAMSIZ);
////
////     if (ioctl(tun_device.fd, TUNSETIFF, &params) == 0)
////     {
////         close(tun_device.fd);
////         tun_device.fd = -1;
////         return tun_device;
////     }
////
////     strcpy(tun_device.name_buff, params.ifr_name);
////
////     return tun_device;
//// }
//import "C"
//
//type TunDevice interface {
//	io.ReadWriteCloser
//	Name() string
//}
//
//type _tunDevice struct {
//	_fd   int
//	_name string
//}
//
//// TODO: make more platform-agnostic
//type _ifreqStruct struct {
//	IfrnName [16]byte // IFNAMSIZ == 16 on my system
//	Fd       int
//	Name     uintptr
//}
//
//func _newTunDeviceWithName(name string) (TunDevice, error) {
//
//	// TODO: initialize structure
//
//	cloneDevice := "/dev/net/tun"
//
//	success := false
//
//	fd, err := unix.Open(cloneDevice, syscall.O_RDWR, 0)
//	if err != nil {
//		return nil, fmt.Errorf("creating tun device: %w", err)
//	}
//	defer func() {
//		if !success {
//			unix.Close(fd)
//		}
//	}()
//
//	// TODO: hold this struct in memory after the function returns?
//	var params _tunDeviceStruct
//	params.Name = uintptr(unsafe.Pointer(&params.NameBuff[0]))
//
//	// performing ioctl TUNSETIFF (code 202)
//	if err := unix.IoctlSetInt(fd, 202, int(uintptr(unsafe.Pointer(&params)))); err != nil {
//		return nil, fmt.Errorf("creating tun device: %w", err)
//	}
//
//	finalNameSize := bytes.IndexByte(params.NameBuff[:], 0)
//	if finalNameSize == -1 {
//		finalNameSize = len(params.NameBuff)
//	}
//	finalName := make([]byte, finalNameSize)
//	copy(finalName, params.NameBuff[:])
//
//	success = true
//
//	return _tunDevice{
//		_fd:   fd,
//		_name: string(finalName),
//	}, nil
//}
//
//func NewTunDevice() (TunDevice, error) {
//	return _newTunDeviceWithName("")
//}
//
//func CreateTunDeviceWithName(name string) (TunDevice, error) {
//	return _newTunDeviceWithName(name)
//}
//
//func (device _tunDevice) Close() error {
//	return unix.Close(device._fd)
//}
//
//func (device _tunDevice) Read(data []byte) (int, error) {
//	return unix.Read(device._fd, data)
//}
//
//func (device _tunDevice) Write(p []byte) (int, error) {
//	return unix.Write(device._fd, p)
//}
//
//func (device _tunDevice) Name() string {
//	return device._name
//}
