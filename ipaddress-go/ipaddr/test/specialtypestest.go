package test

import (
	"github.com/seancfoley/ipaddress/ipaddress-go/ipaddr"
	"math/big"
	"net"
	"strconv"
)

type specialTypesTester struct {
	testBase
}

var (
	//hostOptionsSpecial            = new(ipaddr.HostNameParametersBuilder).AllowEmpty(true).ParseEmptyStrAs(ipaddr.LoopbackOption).GetIPAddressParametersBuilder().AllowEmpty(false).SetRangeParameters(ipaddr.WildcardOnly).AllowAll(true).GetParentBuilder().ToParams()
	//addressOptionsSpecial         = new(ipaddr.IPAddressStringParametersBuilder).Set(hostOptionsSpecial.GetIPAddressParameters()).AllowEmpty(true).ParseEmptyStrAs(ipaddr.LoopbackOption).ToParams()
	hostOptionsSpecial            = new(ipaddr.HostNameParametersBuilder).AllowEmpty(true).GetIPAddressParametersBuilder().AllowEmpty(true).ParseEmptyStrAs(ipaddr.LoopbackOption).SetRangeParameters(ipaddr.WildcardOnly).AllowAll(true).GetParentBuilder().ToParams()
	addressOptionsSpecial         = new(ipaddr.IPAddressStringParametersBuilder).Set(hostOptionsSpecial.GetIPAddressParameters()).AllowEmpty(true).ParseEmptyStrAs(ipaddr.LoopbackOption).ToParams()
	macOptionsSpecial             = new(ipaddr.MACAddressStringParametersBuilder).Set(macAddressOptions).AllowEmpty(true).SetRangeParameters(ipaddr.WildcardOnly).AllowAll(true).ToParams()
	emptyAddressOptions           = new(ipaddr.HostNameParametersBuilder).Set(hostOptions).GetIPAddressParametersBuilder().AllowEmpty(true).ParseEmptyStrAs(ipaddr.LoopbackOption).GetParentBuilder().ToParams()
	emptyAddressNoLoopbackOptions = new(ipaddr.HostNameParametersBuilder).Set(emptyAddressOptions).GetIPAddressParametersBuilder().ParseEmptyStrAs(ipaddr.NoAddressOption).GetParentBuilder().ToParams()
)

func (t specialTypesTester) run() {
	allSingleHex := "0x00000000-0xffffffff"
	allSingleOctal := "000000000000-037777777777"

	t.testIPv4Strings("*", true, "*.*.*.*", "*.*.*.*", "%.%.%.%", "000-255.000-255.000-255.000-255", "*.*.*.*.in-addr.arpa", allSingleHex, allSingleOctal)
	t.testIPv4Strings("***.***.***.***", true, "*.*.*.*", "*.*.*.*", "%.%.%.%", "000-255.000-255.000-255.000-255", "*.*.*.*.in-addr.arpa", allSingleHex, allSingleOctal)
	t.testIPv4Strings("*.*", false, "*.*.*.*", "*.*.*.*", "%.%.%.%", "000-255.000-255.000-255.000-255", "*.*.*.*.in-addr.arpa", allSingleHex, allSingleOctal)
	t.testIPv4Strings("*.*/16", false, "*.*.0.0/16", "*.*.*.*", "%.%.%.%", "000-255.000-255.000.000/16", "*.*.*.*.in-addr.arpa", allSingleHex, allSingleOctal)
	t.testIPv4Strings("*.*/16", true, "*.*.0.0/16", "*.*.*.*", "%.%.%.%", "000-255.000-255.000.000/16", "*.*.*.*.in-addr.arpa", allSingleHex, allSingleOctal)
	t.testIPv4Strings("*/16", true, "*.*.0.0/16", "*.*.*.*", "%.%.%.%", "000-255.000-255.000.000/16", "*.*.*.*.in-addr.arpa", allSingleHex, allSingleOctal)
	t.testIPv4Strings("*/255.255.0.0", false, "*.*.0.0/16", "*.*.*.*", "%.%.%.%", "000-255.000-255.000.000/16", "*.*.*.*.in-addr.arpa", allSingleHex, allSingleOctal)
	t.testIPv4Strings("*/255.255.0.0", true, "*.*.0.0/16", "*.*.*.*", "%.%.%.%", "000-255.000-255.000.000/16", "*.*.*.*.in-addr.arpa", allSingleHex, allSingleOctal)
	t.testIPv4Strings("", false, "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.000.000.001", "1.0.0.127.in-addr.arpa", "0x7f000001", "017700000001")
	t.testIPv4Strings("", true, "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.000.000.001", "1.0.0.127.in-addr.arpa", "0x7f000001", "017700000001")

	base85All := "00000000000000000000" + ipaddr.ExtendedDigitsRangeSeparatorStr + "=r54lj&NUUO~Hi%c2ym0"
	//base85AllPrefixed := base85All + "/16"
	//base85AllPrefixed64 := base85All + "/64"
	base8516 := "00000000000000000000" + ipaddr.ExtendedDigitsRangeSeparatorStr + "=q{+M|w0(OeO5^EGP660" + "/16"
	base8564 := "00000000000000000000" + ipaddr.ExtendedDigitsRangeSeparatorStr + "=r54lj&NUTUTif>jH#O0" + "/64"
	allSingleHexIPv6 := "0x00000000000000000000000000000000-0xffffffffffffffffffffffffffffffff"
	allSingleOctalIPv6 := "00000000000000000000000000000000000000000000-03777777777777777777777777777777777777777777"

	t.testIPv6Strings("*", true, "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "%:%:%:%:%:%:%:%", "0000-ffff:0000-ffff:0000-ffff:0000-ffff:0000-ffff:0000-ffff:0000-ffff:0000-ffff", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*.*.*.*", "*:*:*:*:*:*:*.*.*.*", "*:*:*:*:*:*:*.*.*.*", "*:*:*:*:*:*:*.*.*.*", "*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.ip6.arpa", "*-*-*-*-*-*-*-*.ipv6-literal.net", base85All, allSingleHexIPv6, allSingleOctalIPv6)
	t.testIPv6Strings("*:*", false, "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "%:%:%:%:%:%:%:%", "0000-ffff:0000-ffff:0000-ffff:0000-ffff:0000-ffff:0000-ffff:0000-ffff:0000-ffff", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*.*.*.*", "*:*:*:*:*:*:*.*.*.*", "*:*:*:*:*:*:*.*.*.*", "*:*:*:*:*:*:*.*.*.*", "*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.ip6.arpa", "*-*-*-*-*-*-*-*.ipv6-literal.net", base85All, allSingleHexIPv6, allSingleOctalIPv6)
	t.testIPv6Strings("*:*", true, "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "%:%:%:%:%:%:%:%", "0000-ffff:0000-ffff:0000-ffff:0000-ffff:0000-ffff:0000-ffff:0000-ffff:0000-ffff", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*.*.*.*", "*:*:*:*:*:*:*.*.*.*", "*:*:*:*:*:*:*.*.*.*", "*:*:*:*:*:*:*.*.*.*", "*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.ip6.arpa", "*-*-*-*-*-*-*-*.ipv6-literal.net", base85All, allSingleHexIPv6, allSingleOctalIPv6)

	t.testIPv6Strings("*/16", true,
		"*:0:0:0:0:0:0:0/16",
		"*:*:*:*:*:*:*:*",
		"*:*:*:*:*:*:*:*",
		"%:%:%:%:%:%:%:%",
		"0000-ffff:0000:0000:0000:0000:0000:0000:0000/16",
		"*::/16",
		"*::/16",
		"*::/16",
		"*:*:*:*:*:*:*:*",
		"*::0.0.0.0/16",
		"*::0.0.0.0/16",
		"*::/16",
		"*::/16",
		"*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.ip6.arpa",
		"*-0-0-0-0-0-0-0.ipv6-lit˜eral.net/16",
		base8516, allSingleHexIPv6, allSingleOctalIPv6)
	t.testIPv6Strings("*:*/16", false,
		"*:0:0:0:0:0:0:0/16", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "%:%:%:%:%:%:%:%", "0000-ffff:0000:0000:0000:0000:0000:0000:0000/16", "*::/16", "*::/16", "*::/16", "*:*:*:*:*:*:*:*", "*::0.0.0.0/16", "*::0.0.0.0/16", "*::/16", "*::/16", "*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.ip6.arpa", "*-0-0-0-0-0-0-0.ipv6-literal.net/16", base8516, allSingleHexIPv6, allSingleOctalIPv6)
	t.testIPv6Strings("*:*/16", true, "*:0:0:0:0:0:0:0/16", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "%:%:%:%:%:%:%:%", "0000-ffff:0000:0000:0000:0000:0000:0000:0000/16", "*::/16", "*::/16", "*::/16", "*:*:*:*:*:*:*:*", "*::0.0.0.0/16", "*::0.0.0.0/16", "*::/16", "*::/16", "*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.ip6.arpa", "*-0-0-0-0-0-0-0.ipv6-literal.net/16", base8516, allSingleHexIPv6, allSingleOctalIPv6)
	t.testIPv6Strings("*/64", false, "*:*:*:*:0:0:0:0/64", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "%:%:%:%:%:%:%:%", "0000-ffff:0000-ffff:0000-ffff:0000-ffff:0000:0000:0000:0000/64", "*:*:*:*::/64", "*:*:*:*::/64", "*:*:*:*::/64", "*:*:*:*:*:*:*:*", "*:*:*:*::0.0.0.0/64", "*:*:*:*::0.0.0.0/64", "*:*:*:*::/64", "*:*:*:*::/64", "*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.ip6.arpa", "*-*-*-*-0-0-0-0.ipv6-literal.net/64", base8564, allSingleHexIPv6, allSingleOctalIPv6)
	t.testIPv6Strings("*/64", true, "*:*:*:*:0:0:0:0/64", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "%:%:%:%:%:%:%:%", "0000-ffff:0000-ffff:0000-ffff:0000-ffff:0000:0000:0000:0000/64", "*:*:*:*::/64", "*:*:*:*::/64", "*:*:*:*::/64", "*:*:*:*:*:*:*:*", "*:*:*:*::0.0.0.0/64", "*:*:*:*::0.0.0.0/64", "*:*:*:*::/64", "*:*:*:*::/64", "*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.ip6.arpa", "*-*-*-*-0-0-0-0.ipv6-literal.net/64", base8564, allSingleHexIPv6, allSingleOctalIPv6)
	t.testIPv6Strings("*:*/64", false, "*:*:*:*:0:0:0:0/64", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "%:%:%:%:%:%:%:%", "0000-ffff:0000-ffff:0000-ffff:0000-ffff:0000:0000:0000:0000/64", "*:*:*:*::/64", "*:*:*:*::/64", "*:*:*:*::/64", "*:*:*:*:*:*:*:*", "*:*:*:*::0.0.0.0/64", "*:*:*:*::0.0.0.0/64", "*:*:*:*::/64", "*:*:*:*::/64", "*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.ip6.arpa", "*-*-*-*-0-0-0-0.ipv6-literal.net/64", base8564, allSingleHexIPv6, allSingleOctalIPv6)
	t.testIPv6Strings("*:*/64", true, "*:*:*:*:0:0:0:0/64", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "%:%:%:%:%:%:%:%", "0000-ffff:0000-ffff:0000-ffff:0000-ffff:0000:0000:0000:0000/64", "*:*:*:*::/64", "*:*:*:*::/64", "*:*:*:*::/64", "*:*:*:*:*:*:*:*", "*:*:*:*::0.0.0.0/64", "*:*:*:*::0.0.0.0/64", "*:*:*:*::/64", "*:*:*:*::/64", "*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.ip6.arpa", "*-*-*-*-0-0-0-0.ipv6-literal.net/64", base8564, allSingleHexIPv6, allSingleOctalIPv6)
	t.testIPv6Strings("*/ffff::", false, "*:0:0:0:0:0:0:0/16", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "%:%:%:%:%:%:%:%", "0000-ffff:0000:0000:0000:0000:0000:0000:0000/16", "*::/16", "*::/16", "*::/16", "*:*:*:*:*:*:*:*", "*::0.0.0.0/16", "*::0.0.0.0/16", "*::/16", "*::/16", "*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.ip6.arpa", "*-0-0-0-0-0-0-0.ipv6-literal.net/16", base8516, allSingleHexIPv6, allSingleOctalIPv6)
	t.testIPv6Strings("*/ffff::", true, "*:0:0:0:0:0:0:0/16", "*:*:*:*:*:*:*:*", "*:*:*:*:*:*:*:*", "%:%:%:%:%:%:%:%", "0000-ffff:0000:0000:0000:0000:0000:0000:0000/16", "*::/16", "*::/16", "*::/16", "*:*:*:*:*:*:*:*", "*::0.0.0.0/16", "*::0.0.0.0/16", "*::/16", "*::/16", "*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.*.ip6.arpa", "*-0-0-0-0-0-0-0.ipv6-literal.net/16", base8516, allSingleHexIPv6, allSingleOctalIPv6)

	t.testIPv6Strings("", true, "0:0:0:0:0:0:0:1", "0:0:0:0:0:0:0:1", "::1", "0:0:0:0:0:0:0:1", "0000:0000:0000:0000:0000:0000:0000:0001", "::1", "::1", "::1", "::1", "::0.0.0.1", "::0.0.0.1", "::0.0.0.1", "::0.0.0.1", "1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa", "0-0-0-0-0-0-0-1.ipv6-literal.net", "00000000000000000001", "0x00000000000000000000000000000001", "00000000000000000000000000000000000000000001")

	t.testInvalidValues()

	t.testValidity()

	t.testEmptyValues()

	t.testAllValues()
	t.testAllValuesVersioned(ipaddr.IPv4, getCount(255, 4))
	t.testAllValuesVersioned(ipaddr.IPv6, getCount(0xffff, 8))
	t.testAllMACValues(getCount(0xff, 6), getCount(0xff, 8))

	addressEmpty := t.createParamsHost("", emptyAddressOptions)
	t.hostLabelsHostTest(addressEmpty, []string{"127", "0", "0", "1"})
	addressEmpty2 := t.createParamsHost("", emptyAddressNoLoopbackOptions)
	t.hostLabelsHostTest(addressEmpty2, []string{})
	hostEmpty := t.createParamsHost("", hostOptionsSpecial)
	//t.hostLabelsHostTest(hostEmpty, []string{})
	t.hostLabelsHostTest(hostEmpty, []string{"127", "0", "0", "1"})

	t.testEmptyIsSelf()
	t.testSelf("localhost", true)
	t.testSelf("127.0.0.1", true)
	t.testSelf("::1", true)
	t.testSelf("[::1]", true)
	t.testSelf("*", false)
	t.testSelf("sean.com", false)
	t.testSelf("1.2.3.4", false)
	t.testSelf("::", false)
	t.testSelf("[::]", false)
	t.testSelf("[1:2:3:4:1:2:3:4]", false)
	t.testSelf("1:2:3:4:1:2:3:4", false)

	t.testEmptyLoopback()
	t.testLoopback("127.0.0.1", true)
	t.testLoopback("::1", true)
	t.testLoopback("*", false)
	t.testLoopback("1.2.3.4", false)
	t.testLoopback("::", false)
	t.testLoopback("1:2:3:4:1:2:3:4", false)

	t.testNils()
}

func (t specialTypesTester) testIPv4Strings(addr string, explicit bool, normalizedString, normalizedWildcardString, sqlString, fullString, reverseDNSString, singleHex, singleOctal string) {
	w := t.createParamsAddress(addr, addressOptionsSpecial)
	var ipAddr *ipaddr.IPAddress
	if explicit {
		ipAddr = w.GetVersionedAddress(ipaddr.IPv4)
	} else {
		ipAddr = w.GetAddress()
	}
	t.testStrings(w, ipAddr, normalizedString, normalizedWildcardString, normalizedWildcardString, sqlString, fullString,
		normalizedString, normalizedString, normalizedWildcardString, normalizedString, normalizedWildcardString, reverseDNSString, normalizedString,
		singleHex, singleOctal)
}

func (t specialTypesTester) testIPv6Strings(addr string,
	explicit bool,
	normalizedString,
	normalizedWildcardString,
	canonicalWildcardString,
	sqlString,
	fullString,
	compressedString,
	canonicalString,
	subnetString,
	compressedWildcardString,
	mixedStringNoCompressMixed,
	mixedStringNoCompressHost,
	mixedStringCompressCoveredHost,
	mixedString,
	reverseDNSString,
	uncHostString,
	base85String,
	singleHex,
	singleOctal string) {
	w := t.createParamsAddress(addr, addressOptionsSpecial)
	var ipAddr *ipaddr.IPAddress
	if explicit {
		ipAddr = w.GetVersionedAddress(ipaddr.IPv6)
	} else {
		ipAddr = w.GetAddress()
	}
	t.testBase.testIPv6Strings(w,
		ipAddr,
		normalizedString,
		normalizedWildcardString,
		canonicalWildcardString,
		sqlString,
		fullString,
		compressedString,
		canonicalString,
		subnetString,
		compressedWildcardString,
		mixedStringNoCompressMixed,
		mixedStringNoCompressHost,
		mixedStringCompressCoveredHost,
		mixedString,
		reverseDNSString,
		uncHostString,
		base85String,
		singleHex,
		singleOctal)
}

func (t specialTypesTester) testEmptyValues() {
	//zeroHostOptions := new(ipaddr.HostNameParametersBuilder).ParseEmptyStrAs(ipaddr.LoopbackOption).ToParams()
	zeroHostOptions := new(ipaddr.HostNameParametersBuilder).GetIPAddressParametersBuilder().ParseEmptyStrAs(ipaddr.LoopbackOption).GetParentBuilder().ToParams()
	zeroAddrOptions := new(ipaddr.IPAddressStringParametersBuilder).ParseEmptyStrAs(ipaddr.LoopbackOption).ToParams()
	t.testEmptyValuesOpts(hostOptionsSpecial, addressOptionsSpecial)

	zeroHostOptions = new(ipaddr.HostNameParametersBuilder).GetIPAddressParametersBuilder().ParseEmptyStrAs(ipaddr.ZeroAddressOption).GetParentBuilder().ToParams()
	zeroAddrOptions = new(ipaddr.IPAddressStringParametersBuilder).ParseEmptyStrAs(ipaddr.ZeroAddressOption).ToParams()
	t.testEmptyValuesOpts(zeroHostOptions, zeroAddrOptions)

	zeroHostOptions = new(ipaddr.HostNameParametersBuilder).GetIPAddressParametersBuilder().ParseEmptyStrAs(ipaddr.NoAddressOption).GetParentBuilder().ToParams()
	zeroAddrOptions = new(ipaddr.IPAddressStringParametersBuilder).ParseEmptyStrAs(ipaddr.NoAddressOption).ToParams()
	t.testEmptyValuesOpts(zeroHostOptions, zeroAddrOptions)
}

func (t specialTypesTester) testEmptyValuesOpts(hp ipaddr.HostNameParameters, sp ipaddr.IPAddressStringParameters) {
	hostEmpty := t.createParamsHost("", hp)
	addressEmpty := t.createParamsAddress("", sp)

	// preferredVersion := new(ipaddr.IPAddressStringParametersBuilder).ToParams().GetPreferredVersion()
	preferredAddressVersion := sp.GetPreferredVersion()
	preferredHostVersion := hp.GetPreferredVersion()

	//var addr, addr2 net.IP
	var addr net.IP
	if preferredAddressVersion.IsIPv6() {
		if sp.EmptyStrParsedAs() == ipaddr.LoopbackOption {
			addr = net.IPv6loopback
		} else if sp.EmptyStrParsedAs() == ipaddr.ZeroAddressOption {
			addr = net.IPv6zero
		}
	} else {
		if sp.EmptyStrParsedAs() == ipaddr.LoopbackOption {
			addr = net.IPv4(127, 0, 0, 1)
		} else if sp.EmptyStrParsedAs() == ipaddr.ZeroAddressOption {
			addr = net.IPv4(0, 0, 0, 0)
		}
	}

	//if preferredAddressVersion != preferredHostVersion || hp.EmptyStrParsedAs() != sp.EmptyStrParsedAs() {
	//	t.addFailure(newFailure("failure: precondition to test is that options have same preferred version and emptry string options", addressEmpty))
	//}
	if preferredAddressVersion != preferredHostVersion {
		t.addFailure(newFailure("failure: precondition to test is that options have same preferred version", addressEmpty))
	}
	if addr == nil {
		// empty string not parsed as an address
		if addressEmpty.GetAddress() != nil {
			t.addFailure(newFailure("no match "+addressEmpty.GetAddress().String(), addressEmpty))
		}
		addr, err := addressEmpty.ToAddress()
		if addr != nil {
			t.addFailure(newFailure("no match "+addr.String(), addressEmpty))
		}
		if err != nil {
			t.addFailure(newFailure("no match "+err.Error(), addressEmpty))
		}
		if hostEmpty.AsAddress() != nil {
			t.addFailure(newHostFailure("host "+hostEmpty.String()+" treated as address "+hostEmpty.AsAddress().String(), hostEmpty))
			//t.addFailure(newHostFailure("no match "+hostEmpty.AsAddress().String(), hostEmpty))
		}
		return
	}
	address := ipaddr.FromIP(addr)
	//xxxx so we created the inetaddress, then creating an address from that, and expected equality with the address created directly
	//xxxx so here we could get the equivalent net.IP, and create an address from that
	////try {
	//	InetAddress addr = InetAddress.getByName("");
	//	InetAddress addr2 = InetAddress.getByName(null);
	//
	//	var  params ipaddr.IPAddressStringFormatParameters
	//	if preferredVersion.IsIPv6() {
	//		params = addressOptionsSpecial.GetIPv6Parameters()
	//	} else {
	//		params = addressOptionsSpecial.GetIPv4Parameters()
	//	}

	//IPAddressStringFormatParameters params = addr instanceof Inet6Address ? ADDRESS_OPTIONS_SPECIAL.getIPv6Parameters() : ADDRESS_OPTIONS_SPECIAL.getIPv4Parameters();
	//IPAddressNetwork<?, ?, ?, ?, ?> network = params.getNetwork();
	//IPAddress address = network.getAddressCreator().createAddress(addr.getAddress());

	//IPAddressStringFormatParameters params2 = addr2 instanceof Inet6Address ? ADDRESS_OPTIONS_SPECIAL.getIPv6Parameters() : ADDRESS_OPTIONS_SPECIAL.getIPv4Parameters();
	//IPAddressNetwork<?, ?, ?, ?, ?> network2 = params2.getNetwork();
	//IPAddress address2 = network2.getAddressCreator().createAddress(addr2.getAddress());

	if !addressEmpty.GetAddress().Equal(address) {
		t.addFailure(newFailure("no match "+addr.String(), addressEmpty))
		//} else if(!addressEmpty.GetAddress().Equal(address2)) {
		//	t.addFailure(newFailure("no match " + addr2, addressEmpty));
	} else if addressEmpty.GetAddress().Compare(address) != 0 {
		t.addFailure(newFailure("no match "+addr.String(), addressEmpty))
		//} else if(addressEmpty.GetAddress().Compare(address2) != 0) {
		//	t.addFailure(newFailure("no match " + addr2, addressEmpty));
	} else if addressEmpty.GetAddress().GetCount().Cmp(bigOneConst()) != 0 {
		t.addFailure(newFailure("no count match "+addr.String(), addressEmpty))
	} else {
		addressEmpty = hostEmpty.AsAddressString() //note that hostEmpty allows empty strings and they resolve to loopbacks, but they are not treated as addresses
		//if addressEmpty == nil {
		//	t.addFailure(newFailure("host "+hostEmpty.String()+" treated as address "+addressEmpty.String(), addressEmpty))
		//} else {
		//addressEmpty = t.createParamsHost("", emptyAddressOptions).AsAddressString()
		if addressEmpty == nil {
			t.addFailure(newFailure("host "+hostEmpty.String()+" treated as address "+addressEmpty.String(), addressEmpty))
			//t.addFailure(newFailure("no match ", addressEmpty))
			//} else if(!addressEmpty.GetAddress().Equal(address2)) {
			//	t.addFailure(newFailure("no match " + addr2, addressEmpty));
		} else if !addressEmpty.GetAddress().Equal(address) {
			t.addFailure(newFailure("no match "+addressEmpty.GetAddress().String()+" with "+address.String(), addressEmpty))
			//} else if(!addressEmpty.GetAddress().Equal(address2)) {
			//	t.addFailure(newFailure("no match " + addr2, addressEmpty));
		} else if addressEmpty.GetAddress().Compare(address) != 0 {
			t.addFailure(newFailure("no match "+addr.String(), addressEmpty))
			//} else if(addressEmpty.GetAddress().Compare(address2) != 0) {
			//	t.addFailure(newFailure("no match " + addr2, addressEmpty));
		} else if addressEmpty.GetAddress().GetCount().Cmp(bigOneConst()) != 0 {
			t.addFailure(newFailure("no count match "+addr.String(), addressEmpty))
		} else {
			addressEmptyValue := hostEmpty.GetAddress()
			if !addressEmptyValue.Equal(address) {
				t.addFailure(newFailure("no match "+addr.String(), addressEmpty))
				//} else if(!addressEmptyValue.Equal(address2)) {
				//	t.addFailure(newFailure("no match " + addr2, addressEmpty));
			} else if addressEmptyValue.Compare(address) != 0 {
				t.addFailure(newFailure("no match "+addr.String(), addressEmpty))
				//} else if(addressEmptyValue.Compare(address2) != 0) {
				//	t.addFailure(newFailure("no match " + addr2, addressEmpty));
			} else if addressEmptyValue.GetCount().Cmp(bigOneConst()) != 0 {
				t.addFailure(newFailure("no count match "+addr.String(), addressEmpty))
			}
		}
		//}
	}
	//} catch(UnknownHostException e) {
	//	addFailure(new Failure("unexpected unknown host", addressEmpty));
	//}
	t.incrementTestCount()
}

func (t specialTypesTester) testInvalidValues() {
	// invalid mask
	addressAll := t.createParamsAddress("*/f0ff::", addressOptionsSpecial)
	//try {
	_, err := addressAll.ToAddress()
	if err == nil {
		t.addFailure(newFailure("unexpectedly valid", addressAll))
	} else {
		//} catch(IncompatibleAddressException e) {
		// valid mask
		addressAll = t.createParamsAddress("*/fff0::", addressOptionsSpecial)
		//try {
		if addressAll.GetAddress() == nil {
			t.addFailure(newFailure("unexpectedly invalid", addressAll))
		} else {
			//ambiguous
			addressAll = t.createParamsAddress("*", addressOptionsSpecial)
			if addressAll.GetAddress() != nil {
				t.addFailure(newFailure("unexpectedly invalid", addressAll))
			} else {
				//ambiguous
				addressAll = t.createParamsAddress("*/16", addressOptionsSpecial)
				if addressAll.GetAddress() != nil {
					t.addFailure(newFailure("unexpectedly invalid", addressAll))
				}
				//unambiguous similar addresses tested with testStrings()
			}
		}
	}
	//	} catch(IncompatibleAddressException e2) {
	//	t.addFailure(new Failure("unexpectedly valid", addressAll));
	//	}
	//} catch(AddressStringException e) {
	//t.addFailure(new Failure("unexpected exception: " + e, addressAll));
	//}
}

func (t specialTypesTester) testValidity() {
	hostEmpty := t.createHost("")
	hostAll := t.createHost("*")
	hostAllIPv4 := t.createHost("*.*.*.*")
	hostAllIPv6 := t.createHost("*:*:*:*:*:*:*:*")
	addressEmpty := t.createAddress("")
	addressAll := t.createAddress("*")
	macEmpty := t.createMACAddress("")
	macAll := t.createMACAddress("*")

	if hostEmpty.IsValid() {
		t.addFailure(newHostFailure("unexpectedly valid", hostEmpty))
	} else if hostAll.IsValid() {
		t.addFailure(newHostFailure("unexpectedly valid", hostAll))
	} else if hostAllIPv4.IsValid() {
		t.addFailure(newHostFailure("unexpectedly valid", hostAllIPv4))
	} else if hostAllIPv6.IsValid() {
		t.addFailure(newHostFailure("unexpectedly valid", hostAllIPv6))
	} else if addressEmpty.IsValid() {
		t.addFailure(newFailure("unexpectedly valid", addressEmpty))
	} else if addressAll.IsValid() {
		t.addFailure(newFailure("unexpectedly valid", addressAll))
	} else if macEmpty.IsValid() {
		t.addFailure(newMACFailure("unexpectedly valid", macEmpty))
	} else if macAll.IsValid() {
		t.addFailure(newMACFailure("unexpectedly valid", macAll))
	} else if hostAll.GetAddress() != nil {
		t.addFailure(newHostFailure("unexpectedly valid", hostAll))
	} else if hostEmpty.GetAddress() != nil {
		t.addFailure(newHostFailure("unexpectedly valid", hostEmpty))
	} else {
		hostEmpty = t.createParamsHost("", hostOptionsSpecial)
		hostAll = t.createParamsHost("*", hostOptionsSpecial)
		hostAllIPv4 = t.createParamsHost("*.*.*.*", hostOptionsSpecial)
		hostAllIPv6 = t.createParamsHost("*:*:*:*:*:*:*:*", hostOptionsSpecial)
		addressEmpty = t.createParamsAddress("", addressOptionsSpecial)
		addressAll = t.createParamsAddress("*", addressOptionsSpecial)
		macEmpty = t.createMACParamsAddress("", macOptionsSpecial)
		macAll = t.createMACParamsAddress("*", macOptionsSpecial)
		if !hostEmpty.IsValid() {
			t.addFailure(newHostFailure("unexpectedly invalid", hostEmpty))
		} else if !hostAll.IsValid() {
			t.addFailure(newHostFailure("unexpectedly invalid", hostAll))
		} else if !hostAllIPv4.IsValid() {
			t.addFailure(newHostFailure("unexpectedly invalid", hostAllIPv4))
		} else if !hostAllIPv6.IsValid() {
			t.addFailure(newHostFailure("unexpectedly invalid", hostAllIPv6))
		} else if !addressEmpty.IsValid() {
			t.addFailure(newFailure("unexpectedly invalid", addressEmpty))
		} else if !addressAll.IsValid() {
			t.addFailure(newFailure("unexpectedly invalid", addressAll))
		} else if !macEmpty.IsValid() {
			t.addFailure(newMACFailure("unexpectedly invalid", macEmpty))
		} else if !macAll.IsValid() {
			t.addFailure(newMACFailure("unexpectedly invalid", macAll))
		} else if hostEmpty.GetAddress() == nil { //loopback
			t.addFailure(newHostFailure("unexpectedly invalid", hostEmpty))
		} else if hostAll.GetAddress() != nil {
			t.addFailure(newHostFailure("unexpectedly invalid", hostAll))
		} else {
			//With empty strings, if we wish to allow them, there are two options,
			//we can either treat them as host names and we defer to the validation options for host names, as done above,
			//or we treat than as addresses and use the address options to control behaviour, as we do here.
			hostEmpty = t.createParamsHost("", emptyAddressOptions)
			if !hostEmpty.IsValid() {
				t.addFailure(newHostFailure("unexpectedly invalid", hostEmpty))
			} else if hostEmpty.GetAddress() == nil { //loopback
				t.addFailure(newHostFailure("unexpectedly invalid", hostEmpty))
			} else {
				addressAll = t.createParamsAddress("*.*/64", addressOptionsSpecial) // invalid prefix
				if addressAll.IsValid() {
					t.addFailure(newFailure("unexpectedly valid: "+addressAll.String(), addressAll))
				}
			}
		}
	}
	t.incrementTestCount()
}

func (t specialTypesTester) testAllMACValues(count1, count2 *big.Int) {
	macAll := t.createMACParamsAddress("*", macOptionsSpecial).GetAddress()
	macAll2 := t.createMACParamsAddress("*:*:*:*:*:*:*", macOptionsSpecial).GetAddress()
	address1Str := "*:*:*:*:*:*"
	address2Str := "*:*:*:*:*:*:*:*"
	mac1 := t.createMACParamsAddress(address1Str, macOptionsSpecial).GetAddress()
	mac2 := t.createMACParamsAddress(address2Str, macOptionsSpecial).GetAddress()
	if !macAll.Equal(mac1) {
		t.addFailure(newSegmentSeriesFailure("no match "+macAll.String(), mac1))
	} else if !macAll2.Equal(mac2) {
		t.addFailure(newSegmentSeriesFailure("no match "+macAll2.String(), mac2))
	} else if macAll.Compare(mac1) != 0 {
		t.addFailure(newSegmentSeriesFailure("no match "+macAll.String(), mac1))
	} else if macAll2.Compare(mac2) != 0 {
		t.addFailure(newSegmentSeriesFailure("no match "+macAll2.String(), mac2))
	} else if macAll.GetCount().Cmp(count1) != 0 {
		t.addFailure(newSegmentSeriesFailure("no count match ", macAll))
	} else if macAll2.GetCount().Cmp(count2) != 0 {
		t.addFailure(newSegmentSeriesFailure("no count match ", macAll2))
	}
	t.incrementTestCount()
}

func (t specialTypesTester) testAllValuesVersioned(version ipaddr.IPVersion, count *big.Int) {
	hostAll := t.createParamsHost("*", hostOptionsSpecial)
	addressAllStr := t.createParamsAddress("*", addressOptionsSpecial)
	addressAll := addressAllStr.GetVersionedAddress(version)
	var address2Str = "*.*.*.*"
	if !version.IsIPv4() {
		address2Str = "*:*:*:*:*:*:*:*"
	}
	address := t.createParamsAddress(address2Str, addressOptionsSpecial).GetAddress()
	if !addressAll.Equal(address) {
		t.addFailure(newIPAddrFailure("no match "+address.String(), addressAll))
	} else if addressAll.Compare(address) != 0 {
		t.addFailure(newIPAddrFailure("no match "+address.String(), addressAll))
	} else if addressAll.GetCount().Cmp(count) != 0 {
		// x := t.createParamsAddress("*", addressOptionsSpecial).GetVersionedAddress(version);
		//x.getCount();
		t.addFailure(newIPAddrFailure("no count match ", addressAll))
	} else {
		str := hostAll.AsAddressString()
		addressAll = str.GetVersionedAddress(version)
		if !addressAll.Equal(address) {
			t.addFailure(newIPAddrFailure("no match "+address.String(), addressAll))
		} else if addressAll.Compare(address) != 0 {
			t.addFailure(newIPAddrFailure("no match "+address.String(), addressAll))
		} else if addressAll.GetCount().Cmp(count) != 0 {
			t.addFailure(newIPAddrFailure("no count match ", addressAll))
		}
	}
	t.incrementTestCount()
}

func (t specialTypesTester) testAllValues() {
	hostAll := t.createParamsHost("*", hostOptionsSpecial)
	addressAll := t.createParamsAddress("*", addressOptionsSpecial)
	macAll := t.createMACParamsAddress("*", macOptionsSpecial)
	if addressAll.GetAddress() != nil {
		t.addFailure(newFailure("non null", addressAll))
	} else if hostAll.AsAddress() != nil {
		t.addFailure(newHostFailure("non null", hostAll))
	} else if hostAll.GetAddress() != nil {
		t.addFailure(newHostFailure("non null", hostAll))
	} else if macAll.GetAddress() == nil {
		t.addFailure(newMACFailure("null", macAll))
	}
	t.incrementTestCount()
}

func (t specialTypesTester) testEmptyIsSelf() {
	w := t.createParamsHost("", hostOptionsSpecial)
	if !w.IsSelf() {
		t.addFailure(newHostFailure("failed: isSelf is "+strconv.FormatBool(w.IsSelf()), w))
	}
	w2 := t.createParamsHost("", emptyAddressOptions)
	if !w2.IsSelf() {
		t.addFailure(newHostFailure("failed: isSelf is "+strconv.FormatBool(w2.IsSelf()), w2))
	}
	t.incrementTestCount()
}

func (t specialTypesTester) testSelf(host string, isSelf bool) {
	w := t.createParamsHost(host, hostOptionsSpecial)
	if isSelf != w.IsSelf() {
		t.addFailure(newHostFailure("failed: isSelf is "+strconv.FormatBool(isSelf), w))
	}
	t.incrementTestCount()
}

func (t specialTypesTester) testEmptyLoopback() {
	w := t.createParamsHost("", hostOptionsSpecial)
	if !w.IsLoopback() {
		t.addFailure(newHostFailure("failed: isSelf is "+strconv.FormatBool(w.IsSelf()), w))
	}
	addressEmptyValue := w.GetAddress()
	if !addressEmptyValue.IsLoopback() {
		t.addFailure(newIPAddrFailure("failed: isSelf is "+strconv.FormatBool(addressEmptyValue.IsLoopback()), addressEmptyValue))
	}
	w2 := t.createParamsHost("", emptyAddressOptions)
	if !w2.IsLoopback() {
		t.addFailure(newHostFailure("failed: isSelf is "+strconv.FormatBool(w2.IsSelf()), w2))
	}
	t.incrementTestCount()
}

func (t specialTypesTester) testLoopback(host string, isSelf bool) {
	w := t.createParamsHost(host, hostOptionsSpecial)
	if isSelf != w.IsLoopback() {
		t.addFailure(newHostFailure("failed: isSelf is "+strconv.FormatBool(isSelf), w))
	}
	w2 := t.createParamsAddress(host, addressOptionsSpecial)
	if isSelf != w2.IsLoopback() {
		t.addFailure(newHostFailure("failed: isSelf is "+strconv.FormatBool(isSelf), w))
	}
	t.incrementTestCount()
}

func (t specialTypesTester) testNils() {
	var ipRangesIPv4 []*ipaddr.IPAddressSeqRange
	//var ipv4Addresses []*ipaddr.IPv4Address
	ipv4Addr1 := ipaddr.NewIPAddressString("1.2.3.3").GetAddress().ToIPv4Address()
	ipv4Addr2 := ipaddr.NewIPAddressString("2.2.3.4-5").GetAddress().ToIPv4Address()

	ipRangesIPv4 = append(ipRangesIPv4, nil)
	ipRangesIPv4 = append(ipRangesIPv4, &ipaddr.IPAddressSeqRange{})
	ipRangesIPv4 = append(ipRangesIPv4, ipaddr.NewIPv4SeqRange(nil, nil).ToIPAddressSeqRange())
	//fmt.Println(" it is " + ipaddr.NewIPv4SeqRange(nil, nil).ToIPAddressSeqRange().String())
	ipRangesIPv4 = append(ipRangesIPv4, (&ipaddr.IPv4AddressSeqRange{}).ToIPAddressSeqRange())
	ipRangesIPv4 = append(ipRangesIPv4, ipaddr.NewIPv4SeqRange(&ipaddr.IPv4Address{}, nil).ToIPAddressSeqRange())
	ipRangesIPv4 = append(ipRangesIPv4, ipaddr.NewIPv4SeqRange(ipv4Addr1, nil).ToIPAddressSeqRange())
	ipRangesIPv4 = append(ipRangesIPv4, ipaddr.NewIPv4SeqRange(nil, ipv4Addr2).ToIPAddressSeqRange())
	ipRangesIPv4 = append(ipRangesIPv4, ipaddr.NewIPv4SeqRange(ipv4Addr1, ipv4Addr2).ToIPAddressSeqRange())

	for i := range ipRangesIPv4 {
		range1 := ipRangesIPv4[i]
		//fmt.Printf("range %d using fmt is %v\n", i+1, range1)
		//fmt.Printf("range %d using Stringer is "+range1.String()+"\n\n", i+1)
		for j := i; j < len(ipRangesIPv4); j++ {
			range2 := ipRangesIPv4[j]
			if i == j {
				if range1.Compare(range2) != 0 {
					t.addFailure(newSeqRangeFailure("comparison of "+range1.String()+" with "+range2.String()+" yields "+strconv.Itoa(range1.Compare(range2)), range1))
				} else if range2.Compare(range1) != 0 {
					t.addFailure(newSeqRangeFailure("comparison of "+range2.String()+" with "+range1.String()+" yields "+strconv.Itoa(range2.Compare(range1)), range1))
				} else if !range1.Equal(range2) {
					t.addFailure(newSeqRangeFailure(range1.String()+" and "+range2.String()+" not equal", range1))
				} else if !range2.Equal(range1) {
					t.addFailure(newSeqRangeFailure(range2.String()+" and "+range1.String()+" not equal", range1))
				}
			} else {
				if c := range1.Compare(range2); c > 0 {
					t.addFailure(newSeqRangeFailure("comparison of "+range1.String()+" with "+range2.String()+" yields "+strconv.Itoa(range1.Compare(range2)), range1))
				} else if c == 0 && !range1.Equal(range2) {
					t.addFailure(newSeqRangeFailure(range1.String()+" and "+range2.String()+" not equal", range1))
				} else if c2 := range2.Compare(range1); c2 < 0 {
					t.addFailure(newSeqRangeFailure("comparison of "+range2.String()+" with "+range1.String()+" yields "+strconv.Itoa(range2.Compare(range1)), range1))
				} else if c2 == 0 && (!range2.Equal(range1) || c != 0) {
					t.addFailure(newSeqRangeFailure(range2.String()+" and "+range1.String()+" not equal", range1))
				}
			}
		}
	}

	ipv6Addr1 := ipaddr.NewIPAddressString("1:2:3:3::").GetAddress().ToIPv6Address()
	ipv6Addr2 := ipaddr.NewIPAddressString("2:2:3:4-5::").GetAddress().ToIPv6Address()

	var ipRangesIPv6 []*ipaddr.IPAddressSeqRange

	ipRangesIPv6 = append(ipRangesIPv6, nil)
	ipRangesIPv6 = append(ipRangesIPv6, &ipaddr.IPAddressSeqRange{})
	ipRangesIPv6 = append(ipRangesIPv6, ipaddr.NewIPv6SeqRange(nil, nil).ToIPAddressSeqRange())
	ipRangesIPv6 = append(ipRangesIPv6, (&ipaddr.IPv6AddressSeqRange{}).ToIPAddressSeqRange())
	ipRangesIPv6 = append(ipRangesIPv6, ipaddr.NewIPv6SeqRange(ipv6Addr1, nil).ToIPAddressSeqRange())
	ipRangesIPv6 = append(ipRangesIPv6, ipaddr.NewIPv6SeqRange(nil, ipv6Addr2).ToIPAddressSeqRange())
	ipRangesIPv6 = append(ipRangesIPv6, ipaddr.NewIPv6SeqRange(ipv6Addr1, ipv6Addr2).ToIPAddressSeqRange())

	for i := range ipRangesIPv6 {
		range1 := ipRangesIPv6[i]
		//fmt.Printf("range %d using fmt is %v\n", i+1, range1)
		//fmt.Printf("range %d using Stringer is "+range1.String()+"\n\n", i+1)
		for j := i; j < len(ipRangesIPv6); j++ {
			range2 := ipRangesIPv6[j]
			if i == j {
				if range1.Compare(range2) != 0 {
					t.addFailure(newSeqRangeFailure("comparison of "+range1.String()+" with "+range2.String()+" yields "+strconv.Itoa(range1.Compare(range2)), range1))
				} else if range2.Compare(range1) != 0 {
					t.addFailure(newSeqRangeFailure("comparison of "+range2.String()+" with "+range1.String()+" yields "+strconv.Itoa(range2.Compare(range1)), range1))
				} else if !range1.Equal(range2) {
					t.addFailure(newSeqRangeFailure(range1.String()+" and "+range2.String()+" not equal", range1))
				} else if !range2.Equal(range1) {
					t.addFailure(newSeqRangeFailure(range2.String()+" and "+range1.String()+" not equal", range1))
				}
			} else {
				if c := range1.Compare(range2); c > 0 {
					t.addFailure(newSeqRangeFailure("comparison of "+range1.String()+" with "+range2.String()+" yields "+strconv.Itoa(range1.Compare(range2)), range1))
				} else if c == 0 && !range1.Equal(range2) {
					t.addFailure(newSeqRangeFailure(range1.String()+" and "+range2.String()+" not equal", range1))
				} else if c2 := range2.Compare(range1); c2 < 0 {
					t.addFailure(newSeqRangeFailure("comparison of "+range2.String()+" with "+range1.String()+" yields "+strconv.Itoa(range2.Compare(range1)), range1))
				} else if c2 == 0 && (!range2.Equal(range1) || c != 0) {
					t.addFailure(newSeqRangeFailure(range2.String()+" and "+range1.String()+" not equal", range1))
				}
			}
		}
	}

	for _, range1 := range ipRangesIPv4 {
		for _, range2 := range ipRangesIPv6 {
			// the nils and the blank ranges
			c1 := range1.Compare(range2)
			c2 := range2.Compare(range1)
			if range1 == nil {
				if range2 == nil {
					if c1 != 0 || c2 != 0 {
						t.addFailure(newSeqRangeFailure("comparison of "+range1.String()+" with "+range2.String()+" yields "+strconv.Itoa(range1.Compare(range2)), range1))
					}
				} else if c1 >= 0 {
					t.addFailure(newSeqRangeFailure("comparison of "+range1.String()+" with "+range2.String()+" yields "+strconv.Itoa(range1.Compare(range2)), range1))
				}
			} else if range2 == nil {
				if c1 <= 0 || c2 >= 0 {
					t.addFailure(newSeqRangeFailure("comparison of "+range1.String()+" with "+range2.String()+" yields "+strconv.Itoa(range1.Compare(range2)), range1))
				}
			} else if range1.GetByteCount() == 0 {
				if range2.GetByteCount() == 0 {
					if c1 != 0 || c2 != 0 {
						t.addFailure(newSeqRangeFailure("comparison of "+range1.String()+" with "+range2.String()+" yields "+strconv.Itoa(range1.Compare(range2)), range1))
					}
				} else {
					if c1 >= 0 || c2 <= 0 {
						t.addFailure(newSeqRangeFailure("comparison of "+range1.String()+" with "+range2.String()+" yields "+strconv.Itoa(range1.Compare(range2)), range1))
					}
				}
			} else if range2.GetByteCount() == 0 {
				if c1 <= 0 || c2 >= 0 {
					t.addFailure(newSeqRangeFailure("comparison of "+range1.String()+" with "+range2.String()+" yields "+strconv.Itoa(range1.Compare(range2)), range1))
				}
			} else if c1 >= 0 {
				t.addFailure(newSeqRangeFailure("comparison of "+range1.String()+" with "+range2.String()+" yields "+strconv.Itoa(range1.Compare(range2)), range1))
			} else if c2 <= 0 {
				t.addFailure(newSeqRangeFailure("comparison of "+range2.String()+" with "+range1.String()+" yields "+strconv.Itoa(range2.Compare(range1)), range1))
			}
		}
	}

	nil1 := ipaddr.CountComparator.CompareSeries(nil, nil)
	nil2 := ipaddr.CountComparator.CompareRanges(nil, nil)
	nil3 := ipaddr.CountComparator.CompareAddresses(nil, nil)
	nil4 := ipaddr.CountComparator.CompareDivisions(nil, nil)
	nil5 := ipaddr.CountComparator.CompareAddressSections(nil, nil)
	nil6 := ipaddr.CountComparator.CompareSegments(nil, nil)
	nil7 := ipaddr.CountComparator.Compare(nil, nil)
	if nil1 != 0 || nil2 != 0 || nil3 != 0 || nil4 != 0 || nil5 != 0 || nil6 != 0 || nil7 != 0 {
		t.addFailure(newSegmentSeriesFailure("comparison of nils yields non-zero", nil))
	}

	ipv4Section1 := ipv4Addr1.GetSection()
	ipv6Section1 := ipv6Addr1.GetSection()

	ipv4Range1 := ipRangesIPv4[len(ipRangesIPv4)-1]
	ipv6Range1 := ipRangesIPv6[len(ipRangesIPv6)-1]

	ipv4Segment1 := ipv4Section1.GetSegment(0)
	ipv6Segment1 := ipv6Section1.GetSegment(0)
	ipDivision := ipaddr.NewDivision(11, 8, 0)

	nil1 = ipaddr.CountComparator.CompareSeries(ipv4Addr1, nil)
	nil11 := ipaddr.CountComparator.CompareSeries(ipv6Addr1, nil)
	nil2 = ipaddr.CountComparator.CompareRanges(ipv4Range1, nil)
	nil21 := ipaddr.CountComparator.CompareRanges(ipv6Range1, nil)
	nil3 = ipaddr.CountComparator.CompareAddresses(ipv4Addr1, nil)
	nil4 = ipaddr.CountComparator.CompareDivisions(ipv4Segment1, nil)
	nil400 := ipaddr.CountComparator.CompareDivisions(ipv6Segment1, nil)
	nil40 := ipaddr.CountComparator.CompareDivisions(ipDivision, nil)
	nil41 := ipaddr.CountComparator.CompareSeries(ipv4Section1, nil)
	nil42 := ipaddr.CountComparator.CompareSeries(ipv6Section1, nil)
	nil5 = ipaddr.CountComparator.CompareAddressSections(ipv4Section1, nil)
	nil51 := ipaddr.CountComparator.CompareAddressSections(ipv6Section1, nil)
	nil6 = ipaddr.CountComparator.CompareSegments(ipv4Segment1, nil)
	nil60 := ipaddr.CountComparator.CompareSegments(ipv6Segment1, nil)
	nil7 = ipaddr.CountComparator.Compare(ipv4Addr1, nil)
	if nil1 <= 0 || nil11 <= 0 || nil2 <= 0 || nil21 <= 0 || nil3 <= 0 || nil4 <= 0 || nil400 <= 0 || nil40 <= 0 || nil41 <= 0 || nil42 <= 0 || nil5 <= 0 || nil51 <= 0 || nil6 <= 0 || nil60 <= 0 || nil7 <= 0 {
		t.addFailure(newSegmentSeriesFailure("comparison of nils yields negative", nil))
	}

	nil1 = ipaddr.CountComparator.CompareSeries(nil, ipv4Addr1)
	nil11 = ipaddr.CountComparator.CompareSeries(nil, ipv6Addr1)
	nil2 = ipaddr.CountComparator.CompareRanges(nil, ipv4Range1)
	nil21 = ipaddr.CountComparator.CompareRanges(nil, ipv6Range1)
	nil3 = ipaddr.CountComparator.CompareAddresses(nil, ipv4Addr1)
	nil4 = ipaddr.CountComparator.CompareDivisions(nil, ipv4Segment1)
	nil400 = ipaddr.CountComparator.CompareDivisions(nil, ipv6Segment1)
	nil40 = ipaddr.CountComparator.CompareDivisions(nil, ipDivision)
	nil41 = ipaddr.CountComparator.CompareSeries(nil, ipv4Section1)
	nil42 = ipaddr.CountComparator.CompareSeries(nil, ipv6Section1)
	nil5 = ipaddr.CountComparator.CompareAddressSections(nil, ipv4Section1)
	nil51 = ipaddr.CountComparator.CompareAddressSections(nil, ipv6Section1)
	nil6 = ipaddr.CountComparator.CompareSegments(nil, ipv4Segment1)
	nil60 = ipaddr.CountComparator.CompareSegments(nil, ipv6Segment1)
	nil7 = ipaddr.CountComparator.Compare(nil, ipv4Addr1)
	if nil1 >= 0 || nil11 >= 0 || nil2 >= 0 || nil21 >= 0 || nil3 >= 0 || nil4 >= 0 || nil400 >= 0 || nil40 >= 0 || nil41 >= 0 || nil42 >= 0 || nil5 >= 0 || nil51 >= 0 || nil6 >= 0 || nil60 >= 0 || nil7 >= 0 {
		t.addFailure(newSegmentSeriesFailure("comparison of nils yields positive", nil))
	}

	noIPV6Error := func(sect *ipaddr.IPv6AddressSection) *ipaddr.IPAddress {
		ipv6addrx, _ := ipaddr.NewIPv6Address(sect)
		return ipv6addrx.ToIPAddress()
	}

	var ipAddressesIPv6 []*ipaddr.IPAddress

	ipAddressesIPv6 = append(ipAddressesIPv6, nil)
	ipAddressesIPv6 = append(ipAddressesIPv6, &ipaddr.IPAddress{})
	ipAddressesIPv6 = append(ipAddressesIPv6, (&ipaddr.IPv6Address{}).ToIPAddress())
	ipAddressesIPv6 = append(ipAddressesIPv6, (&ipaddr.IPv6AddressSeqRange{}).GetLowerIPAddress())
	ipAddressesIPv6 = append(ipAddressesIPv6, noIPV6Error(nil))
	ipAddressesIPv6 = append(ipAddressesIPv6, noIPV6Error(ipv6Section1))

	for i := range ipAddressesIPv6 {
		range1 := ipAddressesIPv6[i]
		//fmt.Printf("range %d using fmt is %v\n", i+1, range1)
		//fmt.Printf("range %d using Stringer is "+range1.String()+"\n\n", i+1)
		for j := i; j < len(ipAddressesIPv6); j++ {
			range2 := ipAddressesIPv6[j]
			if i == j {
				if range1.Compare(range2) != 0 {
					t.addFailure(newSegmentSeriesFailure("comparison of "+range1.String()+" with "+range2.String()+" yields "+strconv.Itoa(range1.Compare(range2)), range1))
				} else if range2.Compare(range1) != 0 {
					t.addFailure(newSegmentSeriesFailure("comparison of "+range2.String()+" with "+range1.String()+" yields "+strconv.Itoa(range2.Compare(range1)), range1))
				} else if !range1.Equal(range2) {
					t.addFailure(newSegmentSeriesFailure(range1.String()+" and "+range2.String()+" not equal", range1))
				} else if !range2.Equal(range1) {
					t.addFailure(newSegmentSeriesFailure(range2.String()+" and "+range1.String()+" not equal", range1))
				}
			} else {
				if c := range1.Compare(range2); c > 0 {
					t.addFailure(newSegmentSeriesFailure("comparison of "+range1.String()+" with "+range2.String()+" yields "+strconv.Itoa(range1.Compare(range2)), range1))
				} else if c == 0 && !range1.Equal(range2) {
					t.addFailure(newSegmentSeriesFailure(range1.String()+" and "+range2.String()+" not equal", range1))
				} else if c2 := range2.Compare(range1); c2 < 0 {
					t.addFailure(newSegmentSeriesFailure("comparison of "+range2.String()+" with "+range1.String()+" yields "+strconv.Itoa(range2.Compare(range1)), range1))
				} else if c2 == 0 && (!range2.Equal(range1) || c != 0) {
					t.addFailure(newSegmentSeriesFailure(range2.String()+" and "+range1.String()+" not equal", range1))
				}
			}
		}
	}

	noIPV4Error := func(sect *ipaddr.IPv4AddressSection) *ipaddr.IPAddress {
		ipv4addrx, _ := ipaddr.NewIPv4Address(sect)
		return ipv4addrx.ToIPAddress()
	}

	var ipAddressesIPv4 []*ipaddr.IPAddress

	ipAddressesIPv4 = append(ipAddressesIPv4, nil)
	ipAddressesIPv4 = append(ipAddressesIPv4, &ipaddr.IPAddress{})
	ipAddressesIPv4 = append(ipAddressesIPv4, (&ipaddr.IPv4Address{}).ToIPAddress())
	ipAddressesIPv4 = append(ipAddressesIPv4, (&ipaddr.IPv4AddressSeqRange{}).GetLowerIPAddress())
	ipAddressesIPv4 = append(ipAddressesIPv4, noIPV4Error(nil))
	ipAddressesIPv4 = append(ipAddressesIPv4, noIPV4Error(ipv4Section1))

	for i := range ipAddressesIPv4 {
		range1 := ipAddressesIPv4[i]
		//fmt.Printf("range %d using fmt is %v\n", i+1, range1)
		//fmt.Printf("range %d using Stringer is "+range1.String()+"\n\n", i+1)
		for j := i; j < len(ipAddressesIPv4); j++ {
			range2 := ipAddressesIPv4[j]
			if i == j {
				if range1.Compare(range2) != 0 {
					t.addFailure(newSegmentSeriesFailure("comparison of "+range1.String()+" with "+range2.String()+" yields "+strconv.Itoa(range1.Compare(range2)), range1))
				} else if range2.Compare(range1) != 0 {
					t.addFailure(newSegmentSeriesFailure("comparison of "+range2.String()+" with "+range1.String()+" yields "+strconv.Itoa(range2.Compare(range1)), range1))
				} else if !range1.Equal(range2) {
					t.addFailure(newSegmentSeriesFailure(range1.String()+" and "+range2.String()+" not equal", range1))
				} else if !range2.Equal(range1) {
					t.addFailure(newSegmentSeriesFailure(range2.String()+" and "+range1.String()+" not equal", range1))
				}
			} else {
				if c := range1.Compare(range2); c > 0 {
					t.addFailure(newSegmentSeriesFailure("comparison of "+range1.String()+" with "+range2.String()+" yields "+strconv.Itoa(range1.Compare(range2)), range1))
				} else if c == 0 && !range1.Equal(range2) {
					t.addFailure(newSegmentSeriesFailure(range1.String()+" and "+range2.String()+" not equal", range1))
				} else if c2 := range2.Compare(range1); c2 < 0 {
					t.addFailure(newSegmentSeriesFailure("comparison of "+range2.String()+" with "+range1.String()+" yields "+strconv.Itoa(range2.Compare(range1)), range1))
				} else if c2 == 0 && (!range2.Equal(range1) || c != 0) {
					t.addFailure(newSegmentSeriesFailure(range2.String()+" and "+range1.String()+" not equal", range1))
				}
			}
		}
	}

	//TODO for addresses, when supplying new sections, or nil segment arrays, or nil whatever, I guess we should get the zero address
	/*

		Copied this code over to know how the cmoparisons should shake out in here:

				      if divSeries1, ok := one.(AddressDivisionSeries); ok {
				      		if divSeries2, ok := two.(AddressDivisionSeries); ok {
				      			return comp.CompareSeries(divSeries1, divSeries2)
				      		} else {
				      			return 1
				      		}
				      	} else if div1, ok := one.(DivisionType); ok {
				      		if div2, ok := two.(DivisionType); ok {
				      			return comp.CompareDivisions(div1, div2)
				      		} else {
				      			return -1
				      		}
				      	} else if rng1, ok := one.(IPAddressSeqRangeType); ok {
				      		if rng2, ok := two.(IPAddressSeqRangeType); ok {
				      			return comp.CompareRanges(rng1, rng2)
				      		} else if _, ok := two.(AddressDivisionSeries); ok {
				      			return -1
				      		}
				      		return 1
				      	}
				      	// we've covered all known address items for one, so check two
				      	if _, ok := two.(AddressDivisionSeries); ok {
				      		return -1
				      	} else if _, ok := two.(DivisionType); ok {
				      		return 1
				      	} else if _, ok := two.(IPAddressSeqRangeType); ok {
				      		return -1
				      	}

				so ranges are bigger than the others

				   then

				     if addrSeries1, ok := one.(AddressType); ok {
						if addrSeries2, ok := two.(AddressType); ok {
							return comp.CompareAddresses(addrSeries1, addrSeries2)
						}
						return 1
					} else if _, ok := two.(AddressType); ok {
						return -1
					}
					if addrSection1, ok := one.(AddressSectionType); ok {
						if addrSection2, ok := two.(AddressSectionType); ok {
							return comp.CompareAddressSections(addrSection1, addrSection2)
						}
					}

				addresses are bigger than sections

				   then

				      const (
				ipv6sectype          = 7
				ipv6v4groupingtype   = 6
				ipv4sectype          = 5
				ipsectype            = 4
				macsectype           = 3
				sectype              = 1
				largegroupingtype    = -2
				standardgroupingtype = -3
			)

			const (
				ipv6segtype     = 6
				ipv4segtype     = 5
				ipsegtype       = 4
				macsegtype      = 3
				segtype         = 1
				largedivtype    = -2
				standarddivtype = 0
			)

			const (
				ipv6rangetype = 2
				ipv4rangetype = 1
				iprangetype   = 0
			)

			the above ordering orders the versions/types of each, so ipv6 > ipv4 > ip

	*/
}

func getCount(segmentMax, segmentCount uint64) *big.Int {
	segMax := new(big.Int).SetUint64(segmentMax + 1)
	return segMax.Exp(segMax, new(big.Int).SetUint64(segmentCount), nil)
	//return segMax.pow(segmentCount);
}