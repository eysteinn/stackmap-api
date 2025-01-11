package utils

import (
	"fmt"
	"net"

	"golang.org/x/exp/rand"
)

// Return the first mac address as string
func GetMacAddr() (string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	/*if len(ifas) == 0 {
	}*/
	//var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			//as = append(as, a)
			return a, nil
		}
	}
	//return ifas[0].HardwareAddr.String(), nil
	return "", fmt.Errorf("unable to find mac address for secret key")

	//return as, nil*/
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
