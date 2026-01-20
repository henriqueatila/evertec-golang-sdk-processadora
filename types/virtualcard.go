// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// CreateAccountVirtualCardRequest from OpenAPI spec.
// Used for POST /accounts/{accountId}/virtualcards
type CreateAccountVirtualCardRequest struct {
	IssuerRequestID    string                  `json:"issuerRequestId,omitempty"`
	BirthDate          string                  `json:"birthDate"` // Format: DD/MM/YYYY
	LinkedPhysicalCard string                  `json:"linkedPhysicalCard,omitempty"`
	Constraints        *VirtualCardConstraints `json:"constraints,omitempty"`
}

// CreateVirtualCardRequest from OpenAPI spec.
// Used for POST /cards/{cardId}/virtualcards
type CreateVirtualCardRequest struct {
	IssuerRequestID string                  `json:"issuerRequestId,omitempty"`
	Constraints     *VirtualCardConstraints `json:"constraints,omitempty"`
}

// CreateVirtualCardResponse from OpenAPI spec.
// Reference: https://paysmart-api.gitlab.io/processadora/PT-br/docs/criacao-de-cartoes-virtuais
type CreateVirtualCardResponse struct {
	ResultData      interface{}                     `json:"resultData,omitempty"`
	VirtualCard     *VirtualCardDescriptor          `json:"virtualCard,omitempty"`
	VirtualCardList []ExtendedVirtualCardDescriptor `json:"virtualCardList,omitempty"`
	LinkID          string                          `json:"linkId,omitempty"`
}

// GetVirtualCardResponse from OpenAPI spec.
type GetVirtualCardResponse struct {
	ResultData  interface{} `json:"resultData,omitempty"`
	Virtual     interface{} `json:"virtual,omitempty"`
	Constraints interface{} `json:"constraints,omitempty"`
	Info        interface{} `json:"info,omitempty"`
}

// ListVirtualCardsResponse from OpenAPI spec.
type ListVirtualCardsResponse struct {
	ResultData        interface{}               `json:"resultData,omitempty"`
	VirtualCardsCount int                       `json:"virtual_cards_count,omitempty"`
	VirtualCards      []ComposedVirtualCardInfo `json:"virtual_cards,omitempty"`
}

// VirtualCardInfo from OpenAPI spec.
type VirtualCardInfo struct {
	StatusCode        interface{} `json:"statusCode,omitempty"`
	CardID            string      `json:"cardId,omitempty"`
	StatusDescription string      `json:"statusDescription,omitempty"`
}

// VirtualCardDescriptor from OpenAPI spec.
type VirtualCardDescriptor struct {
	VCardID     string `json:"vCardId,omitempty"`
	VPAN        string `json:"vPan,omitempty"`
	VDateExp    string `json:"vDateExp,omitempty"`
	VCVV        string `json:"vCvv,omitempty"`
	VCardholder string `json:"vCardholder,omitempty"`
}

// VirtualCardDescriptorGet from OpenAPI spec.
type VirtualCardDescriptorGet struct {
	VCardID  string `json:"vCardId,omitempty"`
	VDateExp string `json:"vDateExp,omitempty"`
}

// ExtendedVirtualCardDescriptor from OpenAPI spec.
type ExtendedVirtualCardDescriptor struct {
	VCardID       string `json:"vCardId,omitempty"`
	VPAN          string `json:"vPan,omitempty"`
	VDateExp      string `json:"vDateExp,omitempty"`
	VCVV          string `json:"vCvv,omitempty"`
	VCardholder   string `json:"vCardholder,omitempty"`
	PsProductCode string `json:"psProductCode,omitempty"`
}

// ComposedVirtualCardInfo from OpenAPI spec.
type ComposedVirtualCardInfo struct {
	Virtual     interface{} `json:"virtual,omitempty"`
	Constraints interface{} `json:"constraints,omitempty"`
	Info        interface{} `json:"info,omitempty"`
}

// SimpleCardDescriptor from OpenAPI spec.
type SimpleCardDescriptor struct {
	PsProductCode string `json:"psProductCode"`
	CardID        string `json:"cardId"`
}

// VirtualCardConstraints from OpenAPI spec.
type VirtualCardConstraints struct {
	CurrencyCode        string `json:"currency_code"`
	MaxAmount           string `json:"max_amount"`
	ExpirationTimestamp string `json:"expiration_timestamp"`
}

// CancelVirtualCardRequest from OpenAPI spec.
type CancelVirtualCardRequest struct {
	IssuerRequestID  string       `json:"issuerRequestId,omitempty"`
	CancellationCode int          `json:"cancellationCode"`
	Reason           string       `json:"reason,omitempty"`
	SourceAudit      *SourceAudit `json:"sourceAudit,omitempty"`
}

// CancelVirtualCardResponse from OpenAPI spec.
type CancelVirtualCardResponse struct {
	ResultData           interface{}              `json:"resultData,omitempty"`
	CancelledVirtualCard *ComposedVirtualCardInfo `json:"cancelledVirtualCard,omitempty"`
}

// ModifyVirtualCardResponse from OpenAPI spec.
type ModifyVirtualCardResponse struct {
	ResultData  interface{} `json:"resultData,omitempty"`
	Virtual     interface{} `json:"virtual,omitempty"`
	Constraints interface{} `json:"constraints,omitempty"`
	Info        interface{} `json:"info,omitempty"`
}
