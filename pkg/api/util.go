package api

import (
	"net"
)

// checkIpInSubnet checks contains ip in subnet
func checkIpInSubnet(ipAddr string, subnets []string) (bool, error) {
	// iterate by subnets array and check
	// does subnet contain addr or not
	for _, subnet := range subnets {
		_, subnetParse, err := net.ParseCIDR(subnet)
		if err != nil {
			return false, err
		}
		ipAddrParse := net.ParseIP(ipAddr)
		if subnetParse.Contains(ipAddrParse) {
			return true, nil
		} // end if contains
	} // end for

	return false, nil
}
