package gelfexporter

import (
	"net"
	"strings"
)

func ResolveEndpoint(endpoint string) (string, error) {
	var err error
	var host = endpoint
	var port = ""

	if strings.LastIndexByte(endpoint, ':') != -1 {
		host, port, err = net.SplitHostPort(endpoint)

		if err != nil {
			return "", err
		}
	}

	ips, err := net.LookupIP(host)

	if err != nil || ips == nil || len(ips) == 0 {
		return "", err
	}

	if port != "" {
		return net.JoinHostPort(ips[0].String(), port), nil
	}

	return ips[0].String(), nil
}
