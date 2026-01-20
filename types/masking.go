package types

import (
	"fmt"
	"strings"
)

// MaskPAN masks a PAN (Primary Account Number) showing only the last 4 digits.
// Example: "4111111111111111" -> "************1111"
func MaskPAN(pan string) string {
	if len(pan) <= 4 {
		return strings.Repeat("*", len(pan))
	}
	return strings.Repeat("*", len(pan)-4) + pan[len(pan)-4:]
}

// MaskCVV masks a CVV completely.
// Example: "123" -> "***"
func MaskCVV(cvv string) string {
	return strings.Repeat("*", len(cvv))
}

// MaskPIN masks a PIN completely.
// Example: "1234" -> "****"
func MaskPIN(pin string) string {
	return strings.Repeat("*", len(pin))
}

// MaskExpDate masks an expiration date partially.
// Example: "12/25" -> "**/**"
func MaskExpDate(exp string) string {
	return strings.Repeat("*", len(exp))
}

// String returns a string representation of BindCardRequest with sensitive data masked.
// PAN, CVV, and expiration date are masked to prevent accidental logging.
func (r BindCardRequest) String() string {
	return fmt.Sprintf("BindCardRequest{IssuerRequestID:%q, CardID:%q, PAN:%q, CVV:%q, DateExp:%q, AlternativeBindingKey:%q}",
		r.IssuerRequestID,
		r.CardID,
		MaskPAN(r.PAN),
		MaskCVV(r.CVV),
		MaskExpDate(r.DateExp),
		r.AlternativeBindingKey,
	)
}

// String returns a string representation of AssociateAnonymousCardRequest with sensitive data masked.
func (r AssociateAnonymousCardRequest) String() string {
	return fmt.Sprintf("AssociateAnonymousCardRequest{IssuerRequestID:%q, PAN:%q, CVV:%q, DateExp:%q}",
		r.IssuerRequestID,
		MaskPAN(r.PAN),
		MaskCVV(r.CVV),
		MaskExpDate(r.DateExp),
	)
}

// String returns a string representation of FindCardByPANRequest with sensitive data masked.
func (r FindCardByPANRequest) String() string {
	return fmt.Sprintf("FindCardByPANRequest{PAN:%q}", MaskPAN(r.PAN))
}

// String returns a string representation of PinChangeRequest with sensitive data masked.
func (r PinChangeRequest) String() string {
	return fmt.Sprintf("PinChangeRequest{IssuerPINChangeID:%q, NewPin:{IDTransportKey:%q, PinBlock:***}}}",
		r.IssuerPINChangeID,
		r.NewPin.IDTransportKey,
	)
}

// String returns a string representation of PinValidateRequest with sensitive data masked.
func (r PinValidateRequest) String() string {
	return fmt.Sprintf("PinValidateRequest{IssuerPINValidateID:%q, Pin:{IDTransportKey:%q, PinBlock:***}}",
		r.IssuerPINValidateID,
		r.Pin.IDTransportKey,
	)
}

// String returns a string representation of Pin with sensitive data masked.
func (p Pin) String() string {
	return fmt.Sprintf("Pin{IDTransportKey:%q, PinBlock:***, Format:%q}",
		p.IDTransportKey,
		p.Format,
	)
}

// String returns a string representation of InputPin with sensitive data masked.
func (p InputPin) String() string {
	return fmt.Sprintf("InputPin{IDTransportKey:%q, PinBlock:***, Format:%q}",
		p.IDTransportKey,
		p.Format,
	)
}
