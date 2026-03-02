// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// TravelNotice from OpenAPI spec.
type TravelNotice struct {
	Cards          []Card `json:"cards,omitempty"`
	TravelNoticeID string `json:"travelNoticeId,omitempty"`
	CountryCodes   string `json:"countryCodes,omitempty"`
	BeginDate      string `json:"beginDate,omitempty"`
	EndDate        string `json:"endDate,omitempty"`
}

// NewTravelNotice from OpenAPI spec.
type NewTravelNotice struct {
	IssuerRequestID string `json:"issuerRequestId,omitempty"`
	Cards           []Card `json:"cards"`
	CountryCodes    string `json:"countryCodes"`
	BeginDate       string `json:"beginDate"`
	EndDate         string `json:"endDate"`
}

// TravelNoticeCreatedSuccessfully from OpenAPI spec.
type TravelNoticeCreatedSuccessfully struct {
	ResultData   any  `json:"resultData"`
	TravelNotice TravelNotice `json:"travelNotice"`
}

// UpdateTravelNotice from OpenAPI spec.
type UpdateTravelNotice struct {
	IssuerRequestID string `json:"issuerRequestId,omitempty"`
	Cards           []Card `json:"cards,omitempty"`
	CountryCodes    string `json:"countryCodes,omitempty"`
	BeginDate       string `json:"beginDate,omitempty"`
	EndDate         string `json:"endDate,omitempty"`
}
