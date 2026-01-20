// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// ParseQrCodeParams from OpenAPI spec.
type ParseQrCodeParams struct {
	QRCode        string `json:"qrCode"`
	PsProductCode string `json:"psProductCode,omitempty"`
}

// ParseQrCodeSuccess from OpenAPI spec.
type ParseQrCodeSuccess struct {
	PayloadFormatIndicator     string                 `json:"payloadFormatIndicator,omitempty"`
	PointOfInitiationMethod    string                 `json:"pointOfInitiationMethod,omitempty"`
	MerchantAccountInformation *MerchantAccountInfo   `json:"merchantAccountInformation,omitempty"`
	MerchantCategoryCode       string                 `json:"merchantCategoryCode,omitempty"`
	TransactionCurrency        string                 `json:"transactionCurrency,omitempty"`
	TransactionAmount          string                 `json:"transactionAmount,omitempty"`
	CountryCode                string                 `json:"countryCode,omitempty"`
	MerchantName               string                 `json:"merchantName,omitempty"`
	MerchantCity               string                 `json:"merchantCity,omitempty"`
	POStalCode                 string                 `json:"postalCode,omitempty"`
	AdditionalDataInformation  *AdditionalDataInfo    `json:"additionalDataInformation,omitempty"`
	UnreservedTemplates        string                 `json:"unreservedTemplates,omitempty"`
	TransactionInformation     *QrCodeTransactionInfo `json:"transactionInformation,omitempty"`
	Crc                        string                 `json:"crc,omitempty"`
}

// PaymentCard from OpenAPI spec.
type PaymentCard struct {
	Status          string           `json:"status,omitempty"`
	Error           *QRCODEError     `json:"error,omitempty"`
	TransactionInfo *TransactionInfo `json:"transactionInfo,omitempty"`
}

// PaymentCardParams from OpenAPI spec.
type PaymentCardParams struct {
	KeyID               string `json:"keyId"`
	CipheredInformation string `json:"cipheredInformation"`
	QRCode              string `json:"qrCode"`
	ReceiveCallback     bool   `json:"receiveCallback,omitempty"`
}

// SimplePaymentCardParams from OpenAPI spec.
type SimplePaymentCardParams struct {
	VCardID string `json:"vCardId"`
	QRCode  string `json:"qrCode"`
}

// QRCODEError from OpenAPI spec.
type QRCODEError struct {
	Code                    string `json:"code,omitempty"`
	IssuerAuthorizationCode string `json:"issuerAuthorizationCode,omitempty"`
	Description             string `json:"description,omitempty"`
}

// InternationalData from OpenAPI spec.
type InternationalData struct {
	DollarAmount   *Amount `json:"dollar_amount,omitempty"`
	OriginalAmount *Amount `json:"original_amount,omitempty"`
	DollarRealRate string  `json:"dollar_real_rate,omitempty"`
	Spread         string  `json:"spread,omitempty"`
}
