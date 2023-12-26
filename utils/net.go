package utils

import (
	"log"
	"net"
)

func Get_mtu( iface string ) int {
	handle , err := net.InterfaceByName(iface)
	if( err != nil ) {
		log.Fatal(err)
	}

	return handle.MTU
}