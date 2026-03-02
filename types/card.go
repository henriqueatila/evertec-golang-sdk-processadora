// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// Card from OpenAPI spec.
type Card struct {
	CardID             string         `json:"cardId"`
	AccountID          string         `json:"accountId"`
	IssuerCardID       string         `json:"issuerCardId,omitempty"`
	IssuerCardHolderID string         `json:"issuerCardHolderId,omitempty"`
	Bin                string         `json:"bin"`
	Last4Digits        string         `json:"last4Digits"`
	ExpirationDate     string         `json:"expirationDate"`
	AllowContactless   bool           `json:"allowContactless,omitempty"`
	LinkedPhysicalCard string         `json:"linkedPhysicalCard,omitempty"`
	Physical           bool           `json:"physical,omitempty"`
	Bound              bool           `json:"bound,omitempty"`
	PsProductCode      string         `json:"psProductCode,omitempty"`
	EmbossingFileName  string         `json:"embossingFileName,omitempty"`
	CourierID          *CourierID     `json:"CourierID,omitempty"`
	CardholderData     CardholderData `json:"cardholderData"`
	Profile            string         `json:"profile"`
	Product            string         `json:"product"`
	Status             CardStatus     `json:"status"`
	BlockInformation   any            `json:"blockInformation,omitempty"`
	IssuanceDate       string         `json:"issuanceDate"`
	TrackingID         string         `json:"trackingId,omitempty"`
}

// CardStatus from OpenAPI spec.
// CardStatusInfo - status container
type CardStatusInfo struct {
	Status        string `json:"status,omitempty"`
	StatusDetails string `json:"statusDetails,omitempty"`
}

// CardEmbossing from OpenAPI spec.
type CardEmbossing struct {
	EmbossingName string `json:"embossingName"`
}

// CardListResult from OpenAPI spec.
type CardListResult struct {
	HasMore bool   `json:"hasMore,omitempty"`
	Cards   []Card `json:"cards,omitempty"`
}

// CardSingleListResult from OpenAPI spec.
type CardSingleListResult struct {
	Cards []Card `json:"cards,omitempty"`
}

// CardQueryResult from OpenAPI spec.
type CardQueryResult struct {
	BIN            string `json:"BIN,omitempty"`
	LastFourDigits string `json:"last_four_digits,omitempty"`
	CardID         string `json:"cardId,omitempty"`
}

// CardholderData from OpenAPI spec.
type CardholderData struct {
	CardholderType              string                        `json:"cardholderType"`
	FullName                    string                        `json:"fullName"`
	CardData                    CardEmbossing                 `json:"cardData"`
	IdentityDocumentNumber      string                        `json:"identityDocumentNumber"`
	OtherIdentityDocumentNumber *PersonalIdentityDocumentInfo `json:"otherIdentityDocumentNumber,omitempty"`
	BirthDate                   string                        `json:"birthDate"`
	Nationality                 string                        `json:"nationality"`
	Gender                      string                        `json:"gender,omitempty"`
	CivilStatus                 string                        `json:"civilStatus,omitempty"`
	ContactInformation          *ContactInformation           `json:"contactInformation,omitempty"`
}

// Cardholder from OpenAPI spec.
type Cardholder struct {
	CardholderID   string          `json:"cardholderId,omitempty"`
	CardholderData *CardholderData `json:"cardholderData,omitempty"`
}

// CardholderListResult from OpenAPI spec.
type CardholderListResult struct {
	HasMore     bool         `json:"hasMore,omitempty"`
	Cardholders []Cardholder `json:"cardholders,omitempty"`
}

// CardBlockRequest from OpenAPI spec.
type CardBlockRequest struct {
	IssuerCardBlockID string       `json:"issuerCardBlockId,omitempty"`
	BlockCode         int          `json:"blockCode"`
	Reason            string       `json:"reason,omitempty"`
	ApplyToSingle     bool         `json:"applyToSingle,omitempty"`
	SourceAudit       *SourceAudit `json:"sourceAudit,omitempty"`
}

// CardBlockResult from OpenAPI spec.
type CardBlockResult struct {
	ResultData        any      `json:"resultData"`
	BlockedCard       Card     `json:"blocked_card"`
	OtherChangedCards []string `json:"otherChangedCards,omitempty"`
}

// CardBlockInformation from OpenAPI spec.
type CardBlockInformation struct {
	BlockCodeTransitionsAsList []string `json:"blockCodeTransitionsAsList,omitempty"`
	BlockCode                  int      `json:"blockCode"`
	RequestedBy                string   `json:"requestedBy"`
	Description                string   `json:"description,omitempty"`
	BlockTime                  string   `json:"blockTime"`
}

// CardUnblockRequest from OpenAPI spec.
type CardUnblockRequest struct {
	IssuerCardUnblockID string       `json:"issuerCardUnblockId,omitempty"`
	UnblockCode         int          `json:"unblockCode"`
	Reason              string       `json:"reason,omitempty"`
	ApplyToSingle       bool         `json:"applyToSingle,omitempty"`
	SourceAudit         *SourceAudit `json:"sourceAudit,omitempty"`
}

// CardUnblockResult from OpenAPI spec.
type CardUnblockResult struct {
	ResultData        any      `json:"resultData"`
	UnblockedCard     Card     `json:"unblocked_card"`
	OtherChangedCards []string `json:"otherChangedCards,omitempty"`
}

// CardReissueRequest from OpenAPI spec.
type CardReissueRequest struct {
	IssuerCardID          string       `json:"issuerCardId,omitempty"`
	ExtraData             string       `json:"extraData,omitempty"`
	DeliveryKitCode       float64      `json:"deliveryKitCode,omitempty"`
	IssuerCardReissueID   string       `json:"issuerCardReissueId,omitempty"`
	AllowContactless      bool         `json:"allowContactless,omitempty"`
	CourierID             *CourierID   `json:"CourierID,omitempty"`
	CustomizedTrackingID  string       `json:"customizedTrackingId,omitempty"`
	ReturnAddress         string       `json:"returnAddress,omitempty"`
	EmbossingName         string       `json:"embossingName,omitempty"`
	Reason                string       `json:"reason,omitempty"`
	CardDeliveryAddress   any          `json:"cardDeliveryAddress,omitempty"`
	AlternativeBindingKey string       `json:"alternativeBindingKey,omitempty"`
	SourceAudit           *SourceAudit `json:"sourceAudit,omitempty"`
	BureauxID             *BureauxID   `json:"BureauxID,omitempty"`
}

// CardReissueResult from OpenAPI spec.
type CardReissueResult struct {
	ResultData        any                    `json:"resultData"`
	NewCardID         string                 `json:"newCardId"`
	NewCardIDList     []SimpleCardDescriptor `json:"newCardIdList,omitempty"`
	LinkID            string                 `json:"linkId,omitempty"`
	OtherChangedCards []string               `json:"otherChangedCards,omitempty"`
}

// CardBlockAndReissueRequest from OpenAPI spec.
type CardBlockAndReissueRequest struct {
	IssuerCardBlockID     string       `json:"issuerCardBlockId,omitempty"`
	IssuerCardID          string       `json:"issuerCardId,omitempty"`
	DeliveryKitCode       float64      `json:"deliveryKitCode,omitempty"`
	ExtraData             string       `json:"extraData,omitempty"`
	AllowContactless      bool         `json:"allowContactless,omitempty"`
	CourierID             *CourierID   `json:"CourierID,omitempty"`
	CustomizedTrackingID  string       `json:"customizedTrackingId,omitempty"`
	ReturnAddress         string       `json:"returnAddress,omitempty"`
	EmbossingName         string       `json:"embossingName,omitempty"`
	BlockCode             int          `json:"blockCode"`
	Reason                string       `json:"reason,omitempty"`
	CardDeliveryAddress   any          `json:"cardDeliveryAddress,omitempty"`
	AlternativeBindingKey string       `json:"alternativeBindingKey,omitempty"`
	SourceAudit           *SourceAudit `json:"sourceAudit,omitempty"`
	BureauxID             *BureauxID   `json:"BureauxID,omitempty"`
}

// CardBlockAndReissueResult from OpenAPI spec.
type CardBlockAndReissueResult struct {
	ResultData        any                    `json:"resultData"`
	BlockedCard       Card                   `json:"blocked_card"`
	NewCardID         string                 `json:"newCardId"`
	NewCardIDList     []SimpleCardDescriptor `json:"newCardIdList,omitempty"`
	LinkID            string                 `json:"linkId,omitempty"`
	OtherChangedCards []string               `json:"otherChangedCards,omitempty"`
}

// NewCardRequest from OpenAPI spec.
type NewCardRequest struct {
	IssuerRequestID      string       `json:"issuerRequestId,omitempty"`
	PsProductCode        string       `json:"psProductCode,omitempty"`
	IssuerCardID         string       `json:"issuerCardId,omitempty"`
	DeliveryKitCode      float64      `json:"deliveryKitCode,omitempty"`
	ExtraData            string       `json:"extraData,omitempty"`
	InhibitPhysicalCard  bool         `json:"inhibitPhysicalCard,omitempty"`
	AllowContactless     bool         `json:"allowContactless,omitempty"`
	CourierID            *CourierID   `json:"CourierID,omitempty"`
	CustomizedTrackingID string       `json:"customizedTrackingId,omitempty"`
	ReturnAddress        string       `json:"returnAddress,omitempty"`
	Cardholder           any          `json:"cardholder,omitempty"`
	CardDeliveryAddress  any          `json:"cardDeliveryAddress,omitempty"`
	SourceAudit          *SourceAudit `json:"sourceAudit,omitempty"`
	BureauxID            *BureauxID   `json:"BureauxID,omitempty"`
}

// NewCardrequestCreatedSuccessfully from OpenAPI spec.
// Reference: https://paysmart-api.gitlab.io/processadora/PT-br/docs/criacao-de-cartoes-fisicos
type NewCardrequestCreatedSuccessfully struct {
	ResultData any                    `json:"resultData,omitempty"`
	CardID     string                 `json:"cardId"`
	AccountID  string                 `json:"accountId,omitempty"`
	CardIDList []SimpleCardDescriptor `json:"cardIdList,omitempty"`
	LinkID     string                 `json:"linkId,omitempty"`
}

// NewAnonymousCardRequest from OpenAPI spec.
type NewAnonymousCardRequest struct {
	IssuerRequestID       string       `json:"issuerRequestId,omitempty"`
	IssuerCardID          string       `json:"issuerCardId,omitempty"`
	PsProductCode         string       `json:"psProductCode"`
	AlternativeBindingKey string       `json:"alternativeBindingKey,omitempty"`
	DeliveryKitCode       float64      `json:"deliveryKitCode,omitempty"`
	ExtraData             string       `json:"extraData,omitempty"`
	AllowContactless      bool         `json:"allowContactless,omitempty"`
	CourierID             *CourierID   `json:"CourierID,omitempty"`
	CustomizedTrackingID  string       `json:"customizedTrackingId,omitempty"`
	ReturnAddress         string       `json:"returnAddress,omitempty"`
	CardDeliveryAddress   any          `json:"cardDeliveryAddress"`
	EmbossingName         string       `json:"embossingName,omitempty"`
	SourceAudit           *SourceAudit `json:"sourceAudit,omitempty"`
	BureauxID             *BureauxID   `json:"BureauxID,omitempty"`
}

// NewAnonymousCardRequestCreatedSuccessfully from OpenAPI spec.
type NewAnonymousCardRequestCreatedSuccessfully struct {
	ResultData    any                    `json:"resultData,omitempty"`
	CardID        string                 `json:"cardId"`
	AccountID     string                 `json:"accountId,omitempty"`
	NewCardIDList []SimpleCardDescriptor `json:"newCardIdList,omitempty"`
	LinkID        string                 `json:"linkId,omitempty"`
}

// BindCardRequest from OpenAPI spec.
type BindCardRequest struct {
	IssuerRequestID       string `json:"issuerRequestId,omitempty"`
	CardID                string `json:"cardId"`
	PAN                   string `json:"pan"`
	CVV                   string `json:"cvv"`
	DateExp               string `json:"dateExp"`
	AlternativeBindingKey string `json:"alternativeBindingKey,omitempty"`
}

// BindCardResult from OpenAPI spec.
type BindCardResult struct {
	ResultData any      `json:"resultData"`
	Account    Account  `json:"account"`
	BoundCards []string `json:"boundCards,omitempty"`
}

// BindAnonymousCardholderData from OpenAPI spec.
type BindAnonymousCardholderData struct {
	CardholderType              string                        `json:"cardholderType"`
	FullName                    string                        `json:"fullName"`
	CardData                    CardEmbossing                 `json:"cardData"`
	IdentityDocumentNumber      string                        `json:"identityDocumentNumber"`
	OtherIdentityDocumentNumber *PersonalIdentityDocumentInfo `json:"otherIdentityDocumentNumber,omitempty"`
	BirthDate                   string                        `json:"birthDate"`
	Nationality                 string                        `json:"nationality"`
	Gender                      string                        `json:"gender,omitempty"`
	CivilStatus                 string                        `json:"civilStatus,omitempty"`
	ContactInformation          *ContactInformation           `json:"contactInformation,omitempty"`
	MotherName                  string                        `json:"motherName,omitempty"`
	OccupationInfo              *OccupationInfo               `json:"occupationInfo,omitempty"`
}

// AssociateAnonymousCardRequest from OpenAPI spec.
type AssociateAnonymousCardRequest struct {
	IssuerRequestID string       `json:"issuerRequestId,omitempty"`
	PAN             string       `json:"PAN"`
	CVV             string       `json:"CVV"`
	DateExp         string       `json:"dateExp"`
	Cardholder      any          `json:"cardholder"`
	SourceAudit     *SourceAudit `json:"sourceAudit,omitempty"`
}

// FindCardByPANRequest from OpenAPI spec.
type FindCardByPANRequest struct {
	PAN string `json:"PAN"`
}

// PinChangeRequest from OpenAPI spec.
type PinChangeRequest struct {
	IssuerPINChangeID string       `json:"issuerPINChangeId,omitempty"`
	NewPin            Pin          `json:"newPin"`
	SourceAudit       *SourceAudit `json:"sourceAudit,omitempty"`
}

// PinChangeCreatedSuccessfully from OpenAPI spec.
type PinChangeCreatedSuccessfully struct {
	ResultData                        any      `json:"resultData,omitempty"`
	UpdatingLinkedPhysicalForVirtuals bool     `json:"updatingLinkedPhysicalForVirtuals,omitempty"`
	OtherChangedCards                 []string `json:"otherChangedCards,omitempty"`
}

// Pin from OpenAPI spec.
type Pin struct {
	IDTransportKey string `json:"idTransportKey"`
	PinBlock       string `json:"pinBlock"`
	Format         string `json:"format,omitempty"`
}

// InputPin from OpenAPI spec.
type InputPin struct {
	IDTransportKey string `json:"idTransportKey"`
	PinBlock       string `json:"pinBlock"`
	Format         string `json:"format,omitempty"`
}

// PinValidateRequest from OpenAPI spec.
type PinValidateRequest struct {
	IssuerPINValidateID string       `json:"issuerPINValidateId,omitempty"`
	Pin                 InputPin     `json:"pin"`
	SourceAudit         *SourceAudit `json:"sourceAudit,omitempty"`
}

// PinValidateSuccessfully from OpenAPI spec.
type PinValidateSuccessfully struct {
	ResultData any `json:"resultData,omitempty"`
}

// FunctionalitiesInformation from OpenAPI spec.
type FunctionalitiesInformation struct {
	Description   string `json:"description,omitempty"`
	Code          string `json:"code"`
	International bool   `json:"international"`
	Active        bool   `json:"active"`
}

// GetFunctionBlockingResult from OpenAPI spec.
type GetFunctionBlockingResult struct {
	ResultData    *ResultData    `json:"resultData,omitempty"`
	Functionality map[string]any `json:"functionality,omitempty"`
}

// UpdateFunctionBlockingRequest from OpenAPI spec.
type UpdateFunctionBlockingRequest struct {
	Functionalities []FunctionalitiesInformation `json:"functionalities"`
}

// UpdateFunctionBlockingResult from OpenAPI spec.
type UpdateFunctionBlockingResult struct {
	ResultData    *ResultData    `json:"resultData,omitempty"`
	Functionality map[string]any `json:"functionality,omitempty"`
}

// UpdateCardRequest from OpenAPI spec.
type UpdateCardRequest struct {
	IssuerRequestID    string `json:"issuerRequestId,omitempty"`
	IssuerCardID       string `json:"issuerCardId,omitempty"`
	AllowContactless   bool   `json:"allowContactless,omitempty"`
	LinkedPhysicalCard string `json:"linkedPhysicalCard,omitempty"`
	ApplyToSingle      bool   `json:"applyToSingle,omitempty"`
	Cardholder         any    `json:"cardholder,omitempty"`
}
