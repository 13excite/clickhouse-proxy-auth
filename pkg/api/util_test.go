package api

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckIPInSubnet(t *testing.T) {
	testCases := []struct {
		name          string
		ipAddr        string
		subnets       []string
		expected      bool
		expectedError error
	}{
		{
			name:          "IP address is in one of the subnets",
			ipAddr:        "192.168.1.10",
			subnets:       []string{"192.168.1.0/24", "10.0.0.0/8"},
			expected:      true,
			expectedError: nil,
		},
		{
			name:          "IP address is not in any of the subnets",
			ipAddr:        "172.16.0.5",
			subnets:       []string{"192.168.1.0/24", "10.0.0.0/8"},
			expected:      false,
			expectedError: nil,
		},
		{
			name:          "Invalid subnet should return an error",
			ipAddr:        "192.168.1.10",
			subnets:       []string{"invalid_subnet"},
			expected:      false,
			expectedError: &net.ParseError{Type: "CIDR address", Text: "invalid_subnet"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.ipAddr, func(t *testing.T) {
			result, err := checkIPInSubnet(tc.ipAddr, tc.subnets)

			assert.Equal(t, tc.expected, result, tc.name)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
