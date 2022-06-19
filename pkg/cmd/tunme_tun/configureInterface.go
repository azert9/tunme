package tunme_tun

// #include <string.h>
// #include <unistd.h>
// #include <stdbool.h>
// #include <errno.h>
// #include <sys/ioctl.h>
// #include <arpa/inet.h>
// #include <net/if.h>
//
// const char *configure_interface(_GoString_ iface, _GoString_ addr, _GoString_ mask)
// {
//     bool success = false;
//
//     struct ifreq ifr;
//
//     int fd = socket(PF_INET, SOCK_DGRAM, IPPROTO_IP);
//     if (fd < 0)
//         goto cleanup;
//
//     strncpy(ifr.ifr_name, _GoStringPtr(iface), IFNAMSIZ);
//
//     // address
//
//     ifr.ifr_addr.sa_family = AF_INET;
//     inet_pton(AF_INET, _GoStringPtr(addr), ifr.ifr_addr.sa_data + 2);
//     if (ioctl(fd, SIOCSIFADDR, &ifr) < 0)
//         goto cleanup;
//
//     // mask
//
//     inet_pton(AF_INET, _GoStringPtr(mask), ifr.ifr_addr.sa_data + 2);
//     if(ioctl(fd, SIOCSIFNETMASK, &ifr) < 0)
//         goto cleanup;
//
//     // flags
//
//     ioctl(fd, SIOCGIFFLAGS, &ifr);
//     ifr.ifr_flags |= (IFF_UP | IFF_RUNNING);
//     if (ioctl(fd, SIOCSIFFLAGS, &ifr) < 0)
//         goto cleanup;
//
//     //
//
//     success = true;
//
// cleanup:
//     const char *ret = NULL;
//     if (!success)
//         ret = strerror(errno);
//
//     close(fd);
//
//     return ret;
// }
import "C"

import (
	"fmt"
	"net"
)

func configureInterface(iface string, address string) error {

	addr, subnet, err := net.ParseCIDR(address)
	if err != nil {
		return err
	}

	mask := net.IPv4bcast.Mask(subnet.Mask)

	errStr := C.configure_interface(iface+"\000", addr.String()+"\000", mask.String()+"\000")
	if errStr != nil {
		return fmt.Errorf("failed to configure interface: %s", C.GoString(errStr))
	}

	return nil
}
