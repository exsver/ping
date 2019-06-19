package ping

var IPv4DestinationUnreachableCode = map[int]string{
	//0-5 RFC-792
	//6-12 RFC-1122
	//13-15 RFC-1812
	0:        "net unreachable",
	1:        "host unreachable",
	2:        "protocol unreachable",
	3:        "port unreachable",
	4:        "fragmentation needed and DF set",
	5:        "source route failed",
	6:        "destination network unknown",
	7:        "destination host unknown",
	8:        "source host isolated",
	9:        "communication with destination network administratively prohibited",
	10:       "communication with destination host administratively prohibited",
	11:       "network unreachable for type of service",
	12:       "host unreachable for type of service",
	13:       "communication administratively prohibited",
	14:       "host precedence violation",
	15:       "precedence cutoff in effect",
	16 - 255: "unknown code", // 8 bits in ICMP Header
}

var IPv4TimeExceededCode = map[int]string{
	0:       "time to live exceeded in transit",
	1:       "fragment reassembly time exceeded",
	2 - 255: "unknown code", // 8 bits in ICMP Header
}

var IPv6DestinationUnreachableCode = map[int]string{
	//0-6 RFC-4443
	0:       "no route to destination",
	1:       "communication with destination administratively prohibited",
	2:       "beyond scope of source address",
	3:       "address unreachable",
	4:       "port unreachable",
	5:       "source address failed ingress/egress policy",
	6:       "reject route to destination",
	7 - 255: "unknown code", // 8 bits in ICMPv6 Header
}

var IPv6TimeExceededCode = map[int]string{
	//0-1 RFC-4443
	0:       "hop limit exceeded in transit",
	1:       "fragment reassembly time exceeded",
	2 - 255: "unknown code", // 8 bits in ICMPv6 Header
}

var IPv6ParameterProblemCode = map[int]string{
	//0-2 RFC-4443
	//3 RFC-7112
	0:       "erroneous header field encountered",
	1:       "unrecognized Next Header type encountered",
	2:       "unrecognized IPv6 option encountered",
	3:       "IPv6 First Fragment has incomplete IPv6 Header Chain",
	4 - 255: "unknown code", // 8 bits in ICMPv6 Header
}
