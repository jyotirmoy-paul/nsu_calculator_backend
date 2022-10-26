package utils

import "os"

const AddressEnvVariableName = "address"
const DefaultAddress = "0.0.0.0:30031"

// Gets the address from environment variable OR uses the default address
func GetAddress() string {
	if addr := os.Getenv(AddressEnvVariableName); len(addr) != 0 {
		return addr
	}

	return DefaultAddress
}
