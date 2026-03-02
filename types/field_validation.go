package types

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

// Field length validation errors
var (
	ErrFieldTooLong  = errors.New("field exceeds maximum length")
	ErrFieldTooShort = errors.New("field is below minimum length")
)

// Max field lengths from API specification
// Reference: https://paysmart-api.gitlab.io/processadora/PT-br/docs/formatacao-campos
const (
	MaxIssuerRequestID          = 50
	MaxIssuerCardID             = 50
	MaxMotherName               = 40
	MaxFullName                 = 50
	MaxIdentityDocumentNumber   = 14
	MaxNationality              = 24
	MaxPhoneNumber              = 14
	MaxEmail                    = 60
	MaxEmbossingName            = 24
	MaxRecipient                = 60
	MaxAddressLine1             = 60
	MaxAddressLine2             = 5
	MaxAddressLine3             = 30
	MaxReference                = 60
	MaxNeighborhood             = 30
	MaxCity                     = 30
	MaxState                    = 2
	MaxZipcode                  = 9
	MaxCountry                  = 24
	MaxExtraData                = 30
	MaxProcessID                = 20
	MaxOtherIdentityDocNumber   = 10
	MaxOtherIdentityDocIssuedBy = 3
)

// ValidateMaxLength validates that a string field does not exceed the maximum length.
func ValidateMaxLength(fieldName, value string, maxLen int) error {
	if utf8.RuneCountInString(value) > maxLen {
		return fmt.Errorf("%s: %w (max %d characters, got %d)",
			fieldName, ErrFieldTooLong, maxLen, utf8.RuneCountInString(value))
	}
	return nil
}

// ValidateRangeLength validates that a string field is within the specified range.
func ValidateRangeLength(fieldName, value string, minLen, maxLen int) error {
	length := utf8.RuneCountInString(value)
	if length < minLen {
		return fmt.Errorf("%s: %w (min %d characters, got %d)",
			fieldName, ErrFieldTooShort, minLen, length)
	}
	if length > maxLen {
		return fmt.Errorf("%s: %w (max %d characters, got %d)",
			fieldName, ErrFieldTooLong, maxLen, length)
	}
	return nil
}

// Validate validates the Address fields including length constraints.
func (a Address) Validate() error {
	if err := ValidateMaxLength("AddressLine1", a.AddressLine1, MaxAddressLine1); err != nil {
		return err
	}
	if err := ValidateMaxLength("AddressLine2", a.AddressLine2, MaxAddressLine2); err != nil {
		return err
	}
	if a.AddressLine3 != "" {
		if err := ValidateMaxLength("AddressLine3", a.AddressLine3, MaxAddressLine3); err != nil {
			return err
		}
	}
	if err := ValidateMaxLength("Neighborhood", a.Neighborhood, MaxNeighborhood); err != nil {
		return err
	}
	if err := ValidateMaxLength("City", a.City, MaxCity); err != nil {
		return err
	}
	if err := ValidateMaxLength("State", a.State, MaxState); err != nil {
		return err
	}
	if err := ValidateMaxLength("Zipcode", a.Zipcode, MaxZipcode); err != nil {
		return err
	}
	if a.Country != "" {
		if err := ValidateMaxLength("Country", a.Country, MaxCountry); err != nil {
			return err
		}
	}
	return nil
}

// Validate validates the AddressWithRecipient fields including length constraints.
func (a AddressWithRecipient) Validate() error {
	if a.Recipient != "" {
		if err := ValidateMaxLength("Recipient", a.Recipient, MaxRecipient); err != nil {
			return err
		}
	}
	if err := ValidateMaxLength("AddressLine1", a.AddressLine1, MaxAddressLine1); err != nil {
		return err
	}
	if err := ValidateMaxLength("AddressLine2", a.AddressLine2, MaxAddressLine2); err != nil {
		return err
	}
	if err := ValidateMaxLength("City", a.City, MaxCity); err != nil {
		return err
	}
	if err := ValidateMaxLength("State", a.State, MaxState); err != nil {
		return err
	}
	if err := ValidateMaxLength("Zipcode", a.Zipcode, MaxZipcode); err != nil {
		return err
	}
	return nil
}

// Validate validates the ContactInformation fields including length constraints.
func (c ContactInformation) Validate() error {
	if err := ValidateMaxLength("PersonalPhoneNumber1", c.PersonalPhoneNumber1, MaxPhoneNumber); err != nil {
		return err
	}
	if c.PersonalPhoneNumber2 != "" {
		if err := ValidateMaxLength("PersonalPhoneNumber2", c.PersonalPhoneNumber2, MaxPhoneNumber); err != nil {
			return err
		}
	}
	if c.ComercialPhoneNumber1 != "" {
		if err := ValidateMaxLength("ComercialPhoneNumber1", c.ComercialPhoneNumber1, MaxPhoneNumber); err != nil {
			return err
		}
	}
	if err := ValidateMaxLength("Email", c.Email, MaxEmail); err != nil {
		return err
	}
	return nil
}

// Validate validates the CardEmbossing fields including length constraints.
func (c CardEmbossing) Validate() error {
	return ValidateMaxLength("EmbossingName", c.EmbossingName, MaxEmbossingName)
}

// Validate validates the CardholderData fields including length constraints.
func (c CardholderData) Validate() error {
	if err := ValidateMaxLength("FullName", c.FullName, MaxFullName); err != nil {
		return err
	}
	if err := ValidateMaxLength("IdentityDocumentNumber", c.IdentityDocumentNumber, MaxIdentityDocumentNumber); err != nil {
		return err
	}
	if err := ValidateMaxLength("Nationality", c.Nationality, MaxNationality); err != nil {
		return err
	}
	if c.Gender != "" {
		if !Gender(c.Gender).IsValid() {
			return fmt.Errorf("Gender: invalid value %s", c.Gender)
		}
	}
	if c.CivilStatus != "" {
		if !CivilStatus(c.CivilStatus).IsValid() {
			return fmt.Errorf("CivilStatus: invalid value %s", c.CivilStatus)
		}
	}
	if c.ContactInformation != nil {
		if err := c.ContactInformation.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate validates the PersonalIdentityDocumentInfo fields including length constraints.
func (p PersonalIdentityDocumentInfo) Validate() error {
	if err := ValidateMaxLength("IdentityDocumentNumber", p.IdentityDocumentNumber, MaxOtherIdentityDocNumber); err != nil {
		return err
	}
	if err := ValidateMaxLength("State", p.State, MaxState); err != nil {
		return err
	}
	if err := ValidateMaxLength("IssuedBy", p.IssuedBy, MaxOtherIdentityDocIssuedBy); err != nil {
		return err
	}
	return nil
}

// Validate validates the AccountOwnerData fields including length constraints.
func (a AccountOwnerData) Validate() error {
	if err := ValidateMaxLength("FullName", a.FullName, MaxFullName); err != nil {
		return err
	}
	if err := ValidateMaxLength("IdentityDocumentNumber", a.IdentityDocumentNumber, MaxIdentityDocumentNumber); err != nil {
		return err
	}
	if a.MotherName != "" {
		if err := ValidateMaxLength("MotherName", a.MotherName, MaxMotherName); err != nil {
			return err
		}
	}
	if err := a.ContactInformation.Validate(); err != nil {
		return err
	}
	if a.OtherIdentityDocumentNumber != nil {
		if err := a.OtherIdentityDocumentNumber.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate validates the SourceAudit fields including length constraints.
// Calling Validate on a nil *SourceAudit pointer is safe and returns nil.
func (s *SourceAudit) Validate() error {
	if s == nil {
		return nil
	}
	if s.OperatorID != "" {
		if err := ValidateMaxLength("OperatorID", s.OperatorID, MaxIssuerRequestID); err != nil {
			return err
		}
	}
	if err := ValidateMaxLength("ProcessID", s.ProcessID, MaxProcessID); err != nil {
		return err
	}
	return nil
}

// Validate validates the NewCardRequest fields including length constraints.
func (n NewCardRequest) Validate() error {
	if n.IssuerRequestID != "" {
		if err := ValidateMaxLength("IssuerRequestID", n.IssuerRequestID, MaxIssuerRequestID); err != nil {
			return err
		}
	}
	if n.IssuerCardID != "" {
		if err := ValidateMaxLength("IssuerCardID", n.IssuerCardID, MaxIssuerCardID); err != nil {
			return err
		}
	}
	if n.ExtraData != "" {
		if err := ValidateMaxLength("ExtraData", n.ExtraData, MaxExtraData); err != nil {
			return err
		}
	}
	if n.CustomizedTrackingID != "" {
		if err := ValidateMaxLength("CustomizedTrackingID", n.CustomizedTrackingID, MaxIssuerRequestID); err != nil {
			return err
		}
	}
	if n.Cardholder != nil {
		// Handle any to CardholderData conversion
		if holderData, ok := n.Cardholder.(CardholderData); ok {
			if err := holderData.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Validate validates the PinChangeRequest fields including length constraints.
func (p PinChangeRequest) Validate() error {
	if p.IssuerPINChangeID != "" {
		if err := ValidateMaxLength("IssuerPINChangeID", p.IssuerPINChangeID, MaxIssuerRequestID); err != nil {
			return err
		}
	}
	return p.SourceAudit.Validate()
}

// Validate validates the NewAccountRequest fields including length constraints.
func (n NewAccountRequest) Validate() error {
	if n.IssuerRequestID != "" {
		if err := ValidateMaxLength("IssuerRequestID", n.IssuerRequestID, MaxIssuerRequestID); err != nil {
			return err
		}
	}
	if n.IssuerAccountID != "" {
		if err := ValidateMaxLength("IssuerAccountID", n.IssuerAccountID, MaxIssuerRequestID); err != nil {
			return err
		}
	}
	return n.SourceAudit.Validate()
}

// Validate validates the CardReissueRequest fields including length constraints.
func (c CardReissueRequest) Validate() error {
	if c.IssuerCardID != "" {
		if err := ValidateMaxLength("IssuerCardID", c.IssuerCardID, MaxIssuerCardID); err != nil {
			return err
		}
	}
	if c.IssuerCardReissueID != "" {
		if err := ValidateMaxLength("IssuerCardReissueID", c.IssuerCardReissueID, MaxIssuerRequestID); err != nil {
			return err
		}
	}
	if c.ExtraData != "" {
		if err := ValidateMaxLength("ExtraData", c.ExtraData, MaxExtraData); err != nil {
			return err
		}
	}
	if c.CustomizedTrackingID != "" {
		if err := ValidateMaxLength("CustomizedTrackingID", c.CustomizedTrackingID, MaxIssuerRequestID); err != nil {
			return err
		}
	}
	if c.EmbossingName != "" {
		if err := ValidateMaxLength("EmbossingName", c.EmbossingName, MaxEmbossingName); err != nil {
			return err
		}
	}
	return c.SourceAudit.Validate()
}

// Validate validates the CardBlockRequest fields including length constraints.
func (c CardBlockRequest) Validate() error {
	if c.IssuerCardBlockID != "" {
		if err := ValidateMaxLength("IssuerCardBlockID", c.IssuerCardBlockID, MaxIssuerRequestID); err != nil {
			return err
		}
	}
	return c.SourceAudit.Validate()
}

// Validate validates the CardUnblockRequest fields including length constraints.
func (c CardUnblockRequest) Validate() error {
	if c.IssuerCardUnblockID != "" {
		if err := ValidateMaxLength("IssuerCardUnblockID", c.IssuerCardUnblockID, MaxIssuerRequestID); err != nil {
			return err
		}
	}
	return c.SourceAudit.Validate()
}

// Validate validates the CancelAccountRequest fields including length constraints.
func (c CancelAccountRequest) Validate() error {
	return c.SourceAudit.Validate()
}

// Validate validates the BlockAccountRequest fields including length constraints.
func (b BlockAccountRequest) Validate() error {
	return nil
}

// Validate validates the UnblockAccountRequest fields including length constraints.
func (u UnblockAccountRequest) Validate() error {
	return nil
}
