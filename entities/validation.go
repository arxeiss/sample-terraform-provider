package entities

import (
	"errors"
	"net"
	"regexp"
	"strconv"
)

var (
	nameRegex    = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	ipRegex      = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)
	ipRangeRegex = regexp.MustCompile(`^([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})/([0-9]{1,2})$`)
)

type Validatable interface {
	Validate() error
}

func ValidateName(name string) error {
	if len(name) < 2 || len(name) > 50 {
		return errors.New("name must be between 2 a 50 chars length")
	}
	if !nameRegex.MatchString(name) {
		return errors.New("name can contains only alphanumerical characters, dash and underscore")
	}

	return nil
}

func ValidateDisplayName(dn *string) error {
	if dn == nil {
		return nil
	}
	if len(*dn) < 2 || len(*dn) > 250 {
		return errors.New("display_name must be between 2 a 250 chars length")
	}

	return nil
}

func ValidIP(ip string) bool {
	if !ipRegex.MatchString(ip) {
		return false
	}
	return net.ParseIP(ip) != nil
}

func ValidIPRange(ipRange string) bool {
	parts := ipRangeRegex.FindStringSubmatch(ipRange)
	if len(parts) != 3 {
		return false
	}
	if net.ParseIP(parts[1]) == nil {
		return false
	}
	mask, _ := strconv.Atoi(parts[2])

	return mask < 32
}
