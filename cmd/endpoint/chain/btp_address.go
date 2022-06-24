package chain

import (
	"fmt"
	"strings"
)

type BTPAddress string

func (a BTPAddress) Protocol() string {
	s := string(a)
	if i := strings.Index(s, "://"); i > 0 {
		return s[:i]
	}
	return ""
}
func (a BTPAddress) NetworkAddress() string {
	if a.Protocol() != "" {
		ss := strings.Split(string(a), "/")
		if len(ss) > 2 {
			return ss[2]
		}
	}
	return ""
}
func (a BTPAddress) network() (string, string) {
	if s := a.NetworkAddress(); s != "" {
		ss := strings.Split(s, ".")
		if len(ss) > 1 {
			return ss[0], ss[1]
		} else {
			return "", ss[0]
		}
	}
	return "", ""
}
func (a BTPAddress) BlockChain() string {
	_, v := a.network()
	return v
}
func (a BTPAddress) NetworkID() string {
	n, _ := a.network()
	return n
}
func (a BTPAddress) ContractAddress() string {
	if a.Protocol() != "" {
		ss := strings.Split(string(a), "/")
		if len(ss) > 3 {
			return ss[3]
		}
	}
	return ""
}

func (a BTPAddress) String() string {
	return string(a)
}

func (a *BTPAddress) Set(v string) error {
	*a = BTPAddress(v)
	return nil
}

func (a BTPAddress) Type() string {
	return "BtpAddress"
}

func ValidateBtpAddress(ba BTPAddress) error {
	switch p := ba.Protocol(); p {
	case "btp":
	default:
		return fmt.Errorf("not supported protocol:%s", p)
	}
	switch v := ba.BlockChain(); v {
	case "icon":
	case "iconee":
	default:
		return fmt.Errorf("not supported blockchain:%s", v)
	}
	if len(ba.ContractAddress()) < 1 {
		return fmt.Errorf("empty contract address")
	}
	return nil
}
