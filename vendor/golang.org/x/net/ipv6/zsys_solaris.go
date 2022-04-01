// Code generated by cmd/cgo -godefs; DO NOT EDIT.
// cgo -godefs defs_solaris.go

package ipv6

const (
	sizeofSockaddrStorage = 0x100
	sizeofSockaddrInet6   = 0x20
	sizeofInet6Pktinfo    = 0x14
	sizeofIPv6Mtuinfo     = 0x24

	sizeofIPv6Mreq       = 0x14
	sizeofGroupReq       = 0x104
	sizeofGroupSourceReq = 0x204

	sizeofICMPv6Filter = 0x20
)

type sockaddrStorage struct {
	Family     uint16
	X_ss_pad1  [6]int8
	X_ss_align float64
	X_ss_pad2  [240]int8
}

type sockaddrInet6 struct {
	Family         uint16
	Port           uint16
	Flowinfo       uint32
	Addr           [16]byte /* in6_addr */
	Scope_id       uint32
	X__sin6_src_id uint32
}

type inet6Pktinfo struct {
	Addr    [16]byte /* in6_addr */
	Ifindex uint32
}

type ipv6Mtuinfo struct {
	Addr sockaddrInet6
	Mtu  uint32
}

type ipv6Mreq struct {
	Multiaddr [16]byte /* in6_addr */
	Interface uint32
}

type groupReq struct {
	Interface uint32
	Pad_cgo_0 [256]byte
}

type groupSourceReq struct {
	Interface uint32
	Pad_cgo_0 [256]byte
	Pad_cgo_1 [256]byte
}

type icmpv6Filter struct {
	X__icmp6_filt [8]uint32
}
