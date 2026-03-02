// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// CampaignObject from OpenAPI spec.
type CampaignObject struct {
	AgentID             string  `json:"agentId"`
	ProductCode         string  `json:"productCode"`
	CampaignName        string  `json:"campaignName"`
	ExpirationDate      string  `json:"expirationDate,omitempty"`
	StartDaysLate       float64 `json:"startDaysLate,omitempty"`
	EndDaysLate         float64 `json:"endDaysLate,omitempty"`
	StartOnDebitBalance float64 `json:"startOnDebitBalance,omitempty"`
	EndOnDebitBalance   float64 `json:"endOnDebitBalance,omitempty"`
	FileGenActive       bool    `json:"fileGenActive,omitempty"`
	Status              string  `json:"status,omitempty"`
}

// CampaignResponseObject from OpenAPI spec.
type CampaignResponseObject struct {
	ResultData any `json:"resultData,omitempty"`
	Campaign   any `json:"campaign,omitempty"`
}

// CampaignListResponseObject from OpenAPI spec.
type CampaignListResponseObject struct {
	ResultData any              `json:"resultData,omitempty"`
	Campaigns  []map[string]any `json:"campaigns,omitempty"`
	PageNum    string           `json:"page_num,omitempty"`
	TotalPages string           `json:"total_pages,omitempty"`
}

// CampaignAccountsListResponseObject from OpenAPI spec.
type CampaignAccountsListResponseObject struct {
	ResultData any             `json:"resultData,omitempty"`
	Accounts   []AccountObject `json:"accounts,omitempty"`
	PageNum    int             `json:"page_num,omitempty"`
	TotalPages int             `json:"total_pages,omitempty"`
}

// AgentObject from OpenAPI spec.
type AgentObject struct {
	DocID       string `json:"docId"`
	CompanyName string `json:"companyName"`
	TradingName string `json:"tradingName"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	ContactName string `json:"contactName"`
}

// AgentResponseObject from OpenAPI spec.
type AgentResponseObject struct {
	ResultData any `json:"resultData,omitempty"`
	Agent      any `json:"agent,omitempty"`
}
