package types

import (
	"strings"
	"testing"
)

func TestMaskPAN(t *testing.T) {
	tests := []struct {
		name     string
		pan      string
		expected string
	}{
		{
			name:     "full PAN",
			pan:      "4111111111111111",
			expected: "************1111",
		},
		{
			name:     "short PAN",
			pan:      "1234",
			expected: "****",
		},
		{
			name:     "very short PAN",
			pan:      "12",
			expected: "**",
		},
		{
			name:     "empty PAN",
			pan:      "",
			expected: "",
		},
		{
			name:     "15 digit PAN",
			pan:      "378282246310005",
			expected: "***********0005",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskPAN(tt.pan)
			if result != tt.expected {
				t.Errorf("MaskPAN(%q) = %q, want %q", tt.pan, result, tt.expected)
			}
		})
	}
}

func TestMaskCVV(t *testing.T) {
	tests := []struct {
		name     string
		cvv      string
		expected string
	}{
		{
			name:     "3 digit CVV",
			cvv:      "123",
			expected: "***",
		},
		{
			name:     "4 digit CVV",
			cvv:      "1234",
			expected: "****",
		},
		{
			name:     "empty CVV",
			cvv:      "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskCVV(tt.cvv)
			if result != tt.expected {
				t.Errorf("MaskCVV(%q) = %q, want %q", tt.cvv, result, tt.expected)
			}
		})
	}
}

func TestMaskPIN(t *testing.T) {
	tests := []struct {
		name     string
		pin      string
		expected string
	}{
		{
			name:     "4 digit PIN",
			pin:      "1234",
			expected: "****",
		},
		{
			name:     "6 digit PIN",
			pin:      "123456",
			expected: "******",
		},
		{
			name:     "empty PIN",
			pin:      "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskPIN(tt.pin)
			if result != tt.expected {
				t.Errorf("MaskPIN(%q) = %q, want %q", tt.pin, result, tt.expected)
			}
		})
	}
}

func TestMaskExpDate(t *testing.T) {
	tests := []struct {
		name     string
		exp      string
		expected string
	}{
		{
			name:     "MM/YY format",
			exp:      "12/25",
			expected: "*****",
		},
		{
			name:     "MMYY format",
			exp:      "1225",
			expected: "****",
		},
		{
			name:     "empty",
			exp:      "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskExpDate(tt.exp)
			if result != tt.expected {
				t.Errorf("MaskExpDate(%q) = %q, want %q", tt.exp, result, tt.expected)
			}
		})
	}
}

func TestBindCardRequestString(t *testing.T) {
	req := BindCardRequest{
		IssuerRequestID:       "req-123",
		CardID:                "card-456",
		PAN:                   "4111111111111111",
		CVV:                   "123",
		DateExp:               "12/25",
		AlternativeBindingKey: "alt-key",
	}

	result := req.String()

	// Verify PAN is masked
	if strings.Contains(result, "4111111111111111") {
		t.Error("String() should not contain unmasked PAN")
	}
	if !strings.Contains(result, "************1111") {
		t.Error("String() should contain masked PAN")
	}

	// Verify CVV is masked
	if strings.Contains(result, `CVV:"123"`) {
		t.Error("String() should not contain unmasked CVV")
	}
	if !strings.Contains(result, `CVV:"***"`) {
		t.Error("String() should contain masked CVV")
	}

	// Verify expiration date is masked
	if strings.Contains(result, `DateExp:"12/25"`) {
		t.Error("String() should not contain unmasked expiration date")
	}

	// Verify non-sensitive fields are present
	if !strings.Contains(result, "req-123") {
		t.Error("String() should contain IssuerRequestID")
	}
	if !strings.Contains(result, "card-456") {
		t.Error("String() should contain CardID")
	}
}

func TestAssociateAnonymousCardRequestString(t *testing.T) {
	req := AssociateAnonymousCardRequest{
		IssuerRequestID: "req-123",
		PAN:             "4111111111111111",
		CVV:             "123",
		DateExp:         "12/25",
	}

	result := req.String()

	if strings.Contains(result, "4111111111111111") {
		t.Error("String() should not contain unmasked PAN")
	}
	if strings.Contains(result, `CVV:"123"`) {
		t.Error("String() should not contain unmasked CVV")
	}
}

func TestFindCardByPANRequestString(t *testing.T) {
	req := FindCardByPANRequest{
		PAN: "4111111111111111",
	}

	result := req.String()

	if strings.Contains(result, "4111111111111111") {
		t.Error("String() should not contain unmasked PAN")
	}
	if !strings.Contains(result, "************1111") {
		t.Error("String() should contain masked PAN")
	}
}

func TestPinString(t *testing.T) {
	pin := Pin{
		IDTransportKey: "key-123",
		PinBlock:       "encrypted-pin-block",
		Format:         "ISO-0",
	}

	result := pin.String()

	if strings.Contains(result, "encrypted-pin-block") {
		t.Error("String() should not contain PinBlock value")
	}
	if !strings.Contains(result, "PinBlock:***") {
		t.Error("String() should contain masked PinBlock")
	}
	if !strings.Contains(result, "key-123") {
		t.Error("String() should contain IDTransportKey")
	}
}

func TestInputPinString(t *testing.T) {
	pin := InputPin{
		IDTransportKey: "key-123",
		PinBlock:       "encrypted-pin-block",
		Format:         "ISO-0",
	}

	result := pin.String()

	if strings.Contains(result, "encrypted-pin-block") {
		t.Error("String() should not contain PinBlock value")
	}
	if !strings.Contains(result, "PinBlock:***") {
		t.Error("String() should contain masked PinBlock")
	}
}

func TestPinChangeRequestString(t *testing.T) {
	req := PinChangeRequest{
		IssuerPINChangeID: "pin-change-123",
		NewPin: Pin{
			IDTransportKey: "key-456",
			PinBlock:       "secret-pin-block",
			Format:         "ISO-0",
		},
	}

	result := req.String()

	// Verify PinBlock is not exposed
	if strings.Contains(result, "secret-pin-block") {
		t.Error("String() should not contain PinBlock value")
	}
	// Verify masked output
	if !strings.Contains(result, "PinBlock:***") {
		t.Error("String() should contain masked PinBlock")
	}
	// Verify non-sensitive fields are present
	if !strings.Contains(result, "pin-change-123") {
		t.Error("String() should contain IssuerPINChangeID")
	}
	if !strings.Contains(result, "key-456") {
		t.Error("String() should contain IDTransportKey")
	}
}

func TestPinValidateRequestString(t *testing.T) {
	req := PinValidateRequest{
		IssuerPINValidateID: "pin-validate-123",
		Pin: InputPin{
			IDTransportKey: "key-789",
			PinBlock:       "secret-pin-block",
			Format:         "ISO-1",
		},
	}

	result := req.String()

	// Verify PinBlock is not exposed
	if strings.Contains(result, "secret-pin-block") {
		t.Error("String() should not contain PinBlock value")
	}
	// Verify masked output
	if !strings.Contains(result, "PinBlock:***") {
		t.Error("String() should contain masked PinBlock")
	}
	// Verify non-sensitive fields are present
	if !strings.Contains(result, "pin-validate-123") {
		t.Error("String() should contain IssuerPINValidateID")
	}
	if !strings.Contains(result, "key-789") {
		t.Error("String() should contain IDTransportKey")
	}
}
