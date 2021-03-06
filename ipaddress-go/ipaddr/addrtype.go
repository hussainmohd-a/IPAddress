package ipaddr

type addrType string

const (
	zeroType addrType = ""     // no segments
	ipv4Type addrType = "IPv4" // ipv4 segments
	ipv6Type addrType = "IPv6" // ipv6 segments
	macType  addrType = "MAC"  // mac segments
)

func (a addrType) isNil() bool {
	return a == zeroType
}

func (a addrType) isIPv4() bool {
	return a == ipv4Type
}

func (a addrType) isIPv6() bool {
	return a == ipv6Type
}

func (a addrType) isIP() bool {
	return a.isIPv4() || a.isIPv6()
}

func (a addrType) isMAC() bool {
	return a == macType
}

func (a addrType) getCreator() (creator ParsedAddressCreator) {
	if a.isIPv6() {
		creator = DefaultIPv6Network.GetIPv6AddressCreator()
	} else if a.isIPv4() {
		creator = DefaultIPv4Network.GetIPv4AddressCreator()
	} else if a.isMAC() {
		creator = DefaultMACNetwork.GetMACAddressCreator()
	}
	return
}