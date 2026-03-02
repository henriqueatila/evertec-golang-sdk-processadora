// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// MerchantVanRequest from OpenAPI spec.
type MerchantVanRequest struct {
	AcquirerID         int           `json:"acquirerId"`
	Document           string        `json:"document"`
	DocumentType       string        `json:"documentType"`
	LegalName          string        `json:"legalName"`
	FantasyName        string        `json:"fantasyName,omitempty"`
	MerchantVanCode    string        `json:"merchantVanCode"`
	CPFCNPJSubAcquirer string        `json:"cpfCnpjSubAcquirer,omitempty"`
	Products           []ProductsDTO `json:"products"`
}

// MerchantVanResponse from OpenAPI spec.
type MerchantVanResponse struct {
	ResultData         any           `json:"resultData"`
	BrandErrorList     []any         `json:"brandErrorList,omitempty"`
	LiteralBrandError  string        `json:"literalBrandError,omitempty"`
	Brand              any           `json:"brand,omitempty"`
	AcquirerID         int           `json:"acquirerId,omitempty"`
	Document           string        `json:"document,omitempty"`
	DocumentType       string        `json:"documentType,omitempty"`
	LegalName          string        `json:"legalName,omitempty"`
	FantasyName        string        `json:"fantasyName,omitempty"`
	MerchantVanCode    string        `json:"merchantVanCode,omitempty"`
	CPFCNPJSubAcquirer string        `json:"cpfCnpjSubAcquirer,omitempty"`
	Products           []ProductsDTO `json:"products,omitempty"`
}

// ReducedMerchantVanResponse from OpenAPI spec.
type ReducedMerchantVanResponse struct {
	ResultData        any    `json:"resultData"`
	BrandErrorList    []any  `json:"brandErrorList,omitempty"`
	LiteralBrandError string `json:"literalBrandError,omitempty"`
}
