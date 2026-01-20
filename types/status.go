// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// HealthCheckResponse from OpenAPI spec.
type HealthCheckResponse struct {
	ResultData interface{}                `json:"resultData,omitempty"`
	Systems    []AuthorizerHealthResponse `json:"systems,omitempty"`
	Datetime   string                     `json:"datetime,omitempty"`
}

// AuthorizerHealthResponse from OpenAPI spec.
type AuthorizerHealthResponse struct {
	Health       *HealthStatus `json:"health,omitempty"`
	AuthorizerID string        `json:"authorizerId,omitempty"`
}

// DataprepStatus - Cursor contendo status detalhado de um arquivo de dataprep
type DataprepStatus struct {
	Filename            string `json:"filename"`
	IssuerName          string `json:"issuerName,omitempty"`
	Registers           int    `json:"registers,omitempty"`
	CreationTimestamp   string `json:"creationTimestamp,omitempty"`
	LastUpdateTimestamp string `json:"lastUpdateTimestamp,omitempty"`
	MatStart            string `json:"matStart,omitempty"`
	MatEnd              string `json:"matEnd,omitempty"`
	MatExpectedSteps    int    `json:"matExpectedSteps,omitempty"`
	MatCurrentStep      int    `json:"matCurrentStep,omitempty"`
	MatStatus           string `json:"matStatus,omitempty"`
	EmbStart            string `json:"embStart,omitempty"`
	EmbEnd              string `json:"embEnd,omitempty"`
	EmbExpectedSteps    int    `json:"embExpectedSteps,omitempty"`
	EmbCurrentStep      int    `json:"embCurrentStep,omitempty"`
	EmbStatus           string `json:"embStatus,omitempty"`
	DtpStart            string `json:"dtpStart,omitempty"`
	DtpEnd              string `json:"dtpEnd,omitempty"`
	DtpExpectedSteps    int    `json:"dtpExpectedSteps,omitempty"`
	DtpCurrentStep      int    `json:"dtpCurrentStep,omitempty"`
	DtpStatus           string `json:"dtpStatus,omitempty"`
	SftpStart           string `json:"sftpStart,omitempty"`
	SftpEnd             string `json:"sftpEnd,omitempty"`
	SftpExpectedSteps   int    `json:"sftpExpectedSteps,omitempty"`
	SftpCurrentStep     int    `json:"sftpCurrentStep,omitempty"`
	SftpStatus          string `json:"sftpStatus,omitempty"`
	CourierName         string `json:"courierName,omitempty"`
	BureauxName         string `json:"bureauxName,omitempty"`
}

// DataprepStatusResult from OpenAPI spec.
type DataprepStatusResult struct {
	ResultData *ResultData     `json:"resultData,omitempty"`
	File       *DataprepStatus `json:"file,omitempty"`
}

// DataprepStatusListResult from OpenAPI spec.
type DataprepStatusListResult struct {
	ResultData *ResultData      `json:"resultData,omitempty"`
	HasMore    bool             `json:"hasMore,omitempty"`
	Files      []DataprepStatus `json:"files,omitempty"`
}
