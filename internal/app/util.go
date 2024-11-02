package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
)

func HashIp(ip string) (string, error) {
	parsedIp := net.ParseIP(ip)
	if parsedIp == nil {
		return "", fmt.Errorf("invalid IP address: %s", ip)
	}

	hasher := sha256.New()

	_, err := hasher.Write(parsedIp)
	if err != nil {
		return "", err
	}

	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}
