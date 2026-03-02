// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// HCEProvisionRequest from OpenAPI spec.
type HCEProvisionRequest struct {
	Mtv         map[string]any `json:"mtv,omitempty"`
	Constraints map[string]any `json:"constraints,omitempty"`
	Transaction map[string]any `json:"transaction,omitempty"`
}

// HCEProvisionSuccessfully from OpenAPI spec.
type HCEProvisionSuccessfully struct {
	ResultData  any            `json:"resultData,omitempty"`
	Mtv         map[string]any `json:"mtv,omitempty"`
	Transaction map[string]any `json:"transaction,omitempty"`
}

// CreateHCECardRequest from OpenAPI spec.
type CreateHCECardRequest struct {
	RealProduct map[string]any `json:"real_product,omitempty"`
	Mtv         map[string]any `json:"mtv,omitempty"`
	Transaction map[string]any `json:"transaction,omitempty"`
}

// CreateHCECardSuccesfully from OpenAPI spec.
type CreateHCECardSuccesfully struct {
	ResultData  any            `json:"resultData,omitempty"`
	Mtv         map[string]any `json:"mtv,omitempty"`
	Transaction map[string]any `json:"transaction,omitempty"`
}

// UnprovisionRequest from OpenAPI spec.
type UnprovisionRequest struct {
	Transaction map[string]any `json:"transaction,omitempty"`
}

// UnprovisionSuccesfully from OpenAPI spec.
type UnprovisionSuccesfully struct {
	Response    *ResultData            `json:"response,omitempty"`
	Transaction map[string]any `json:"transaction,omitempty"`
}
