package types

import (
	"strings"
	"testing"
)

func TestValidateMaxLength(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		value     string
		maxLen    int
		wantErr   bool
	}{
		{
			name:      "valid length",
			fieldName: "testField",
			value:     "hello",
			maxLen:    10,
			wantErr:   false,
		},
		{
			name:      "exactly max length",
			fieldName: "testField",
			value:     "hello",
			maxLen:    5,
			wantErr:   false,
		},
		{
			name:      "exceeds max length",
			fieldName: "testField",
			value:     "hello world",
			maxLen:    5,
			wantErr:   true,
		},
		{
			name:      "empty value is valid",
			fieldName: "testField",
			value:     "",
			maxLen:    5,
			wantErr:   false,
		},
		{
			name:      "unicode characters count correctly",
			fieldName: "testField",
			value:     "日本語",
			maxLen:    2,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMaxLength(tt.fieldName, tt.value, tt.maxLen)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMaxLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRangeLength(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		value     string
		minLen    int
		maxLen    int
		wantErr   bool
	}{
		{
			name:      "valid range",
			fieldName: "testField",
			value:     "hello",
			minLen:    3,
			maxLen:    10,
			wantErr:   false,
		},
		{
			name:      "below min length",
			fieldName: "testField",
			value:     "hi",
			minLen:    3,
			maxLen:    10,
			wantErr:   true,
		},
		{
			name:      "above max length",
			fieldName: "testField",
			value:     "hello world",
			minLen:    3,
			maxLen:    5,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRangeLength(tt.fieldName, tt.value, tt.minLen, tt.maxLen)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRangeLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddress_Validate(t *testing.T) {
	tests := []struct {
		name    string
		address Address
		wantErr bool
	}{
		{
			name: "valid address",
			address: Address{
				AddressLine1: "Rua Example",
				AddressLine2: "123",
				Neighborhood: "Bairro",
				City:         "São Paulo",
				State:        "SP",
				Zipcode:      "01234567",
			},
			wantErr: false,
		},
		{
			name: "addressLine1 too long",
			address: Address{
				AddressLine1: strings.Repeat("a", MaxAddressLine1+1),
				AddressLine2: "123",
				Neighborhood: "Bairro",
				City:         "São Paulo",
				State:        "SP",
				Zipcode:      "01234567",
			},
			wantErr: true,
		},
		{
			name: "zipcode too long",
			address: Address{
				AddressLine1: "Rua Example",
				AddressLine2: "123",
				Neighborhood: "Bairro",
				City:         "São Paulo",
				State:        "SP",
				Zipcode:      strings.Repeat("0", MaxZipcode+1),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.address.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Address.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContactInformation_Validate(t *testing.T) {
	tests := []struct {
		name    string
		contact ContactInformation
		wantErr bool
	}{
		{
			name: "valid contact",
			contact: ContactInformation{
				PersonalPhoneNumber1: "+5511999999999",
				Email:                "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "email too long",
			contact: ContactInformation{
				PersonalPhoneNumber1: "+5511999999999",
				Email:                strings.Repeat("a", MaxEmail+1) + "@example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.contact.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ContactInformation.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCardEmbossing_Validate(t *testing.T) {
	tests := []struct {
		name      string
		embossing CardEmbossing
		wantErr   bool
	}{
		{
			name: "valid embossing name",
			embossing: CardEmbossing{
				EmbossingName: "JOAO DA SILVA",
			},
			wantErr: false,
		},
		{
			name: "embossing name too long",
			embossing: CardEmbossing{
				EmbossingName: strings.Repeat("A", MaxEmbossingName+1),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.embossing.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CardEmbossing.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCardholderData_Validate(t *testing.T) {
	tests := []struct {
		name    string
		holder  CardholderData
		wantErr bool
	}{
		{
			name: "valid holder",
			holder: CardholderData{
				CardholderType:         "main",
				FullName:               "João da Silva",
				CardData:               CardEmbossing{EmbossingName: "JOAO SILVA"},
				IDentityDocumentNumber: "12345678901",
				BirthDate:              "01/01/1990",
				Nationality:            "BRA",
				ContactInformation:     &ContactInformation{},
			},
			wantErr: false,
		},
		{
			name: "fullName too long",
			holder: CardholderData{
				CardholderType:         "main",
				FullName:               strings.Repeat("A", MaxFullName+1),
				CardData:               CardEmbossing{},
				IDentityDocumentNumber: "12345678901",
				BirthDate:              "01/01/1990",
				Nationality:            "BRA",
			},
			wantErr: true,
		},
		{
			name: "invalid gender",
			holder: CardholderData{
				CardholderType:         "main",
				FullName:               "João da Silva",
				CardData:               CardEmbossing{},
				IDentityDocumentNumber: "12345678901",
				BirthDate:              "01/01/1990",
				Nationality:            "BRA",
				Gender:                 "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.holder.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CardholderData.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSourceAudit_Validate(t *testing.T) {
	tests := []struct {
		name    string
		audit   SourceAudit
		wantErr bool
	}{
		{
			name: "valid audit",
			audit: SourceAudit{
				OperatorID: "op123",
				ProcessID:  "proc456",
			},
			wantErr: false,
		},
		{
			name: "processId too long",
			audit: SourceAudit{
				ProcessID: strings.Repeat("p", MaxProcessID+1),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.audit.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SourceAudit.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewCardRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request NewCardRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: NewCardRequest{
				IssuerRequestID: "req123",
				IssuerCardID:    "card123",
			},
			wantErr: false,
		},
		{
			name: "issuerRequestID too long",
			request: NewCardRequest{
				IssuerRequestID: strings.Repeat("r", MaxIssuerRequestID+1),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCardRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAccountRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request NewAccountRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: NewAccountRequest{
				IssuerRequestID: "req123",
				IssuerAccountID: "acc123",
				SourceAudit:     &SourceAudit{},
			},
			wantErr: false,
		},
		{
			name: "issuerRequestID too long",
			request: NewAccountRequest{
				IssuerRequestID: strings.Repeat("r", MaxIssuerRequestID+1),
				SourceAudit:     &SourceAudit{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccountRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
