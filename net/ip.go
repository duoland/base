package net

import (
	"math/big"
	"net"
)

// IPv4ToInt64 convert ip address to int64
func IPv4ToInt64(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}
