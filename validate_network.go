package validationcontext

import (
	"fmt"
	"net"
)

// ValidateIPAddress checks if the value is a valid IP address (IPv4 or IPv6).
func (vc *ValidationContext) ValidateIPAddress(value, field, errMsg string) {
	if ip := net.ParseIP(value); ip == nil {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、有効なIPアドレスを指定してください。", field))
	}
}

// ValidateIPv4 checks if the value is a valid IPv4 address.
func (vc *ValidationContext) ValidateIPv4(value, field, errMsg string) {
	ip := net.ParseIP(value)
	if ip == nil || ip.To4() == nil {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、有効なIPv4アドレスを指定してください。", field))
	}
}

// ValidateIPv6 checks if the value is a valid IPv6 address.
func (vc *ValidationContext) ValidateIPv6(value, field, errMsg string) {
	ip := net.ParseIP(value)
	if ip == nil || ip.To4() != nil {
		vc.addErrorWithMessage(field, errMsg, fmt.Sprintf("%sには、有効なIPv6アドレスを指定してください。", field))
	}
}
