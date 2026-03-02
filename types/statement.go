// Code generated from OpenAPI spec. DO NOT EDIT.
package types

// Statement from OpenAPI spec.
type Statement struct {
	AccountID                               string           `json:"accountId,omitempty"`
	PaymentDue                              string           `json:"paymentDue,omitempty"`
	PaymentDueActual                        string           `json:"paymentDueActual,omitempty"`
	CloseDate                               string           `json:"closeDate,omitempty"`
	Balance                                 any      `json:"balance,omitempty"`
	MinimumPaymentDue                       any      `json:"minimumPaymentDue,omitempty"`
	TransactionsList                        []StatementEntry `json:"transactionsList,omitempty"`
	AdditionalTransactionsList              []StatementEntry `json:"additionalTransactionsList,omitempty"`
	OutOfStatementPeriodTransactions        []StatementEntry `json:"outOfStatementPeriodTransactions,omitempty"`
	FirstTransactionInserted                string           `json:"firstTransactionInserted,omitempty"`
	LastTransactionInserted                 string           `json:"lastTransactionInserted,omitempty"`
	OpeningDateTime                         string           `json:"openingDateTime,omitempty"`
	StatementClosingDateTime                string           `json:"statementClosingDateTime,omitempty"`
	ChargesInNextStatementForMinimumPayment any      `json:"chargesInNextStatementForMinimumPayment,omitempty"`
	StatementID                             string           `json:"statementId,omitempty"`
	LastStatementID                         string           `json:"lastStatementId,omitempty"`
	InterestCalculationMemory               []DailyFeeEntry  `json:"interestCalculationMemory,omitempty"`
	LastUpdated                             string           `json:"lastUpdated,omitempty"`
	PaymentsAndCredits                      any      `json:"paymentsAndCredits,omitempty"`
	PurchasesAndDebits                      any      `json:"purchasesAndDebits,omitempty"`
	CET                                     any      `json:"CET,omitempty"`
	ResultData                              any      `json:"resultData"`
	QueryDate                               string           `json:"query_date,omitempty"`
	SourceAudit                             *SourceAudit     `json:"sourceAudit,omitempty"`
}

// StatementList from OpenAPI spec.
type StatementList struct {
	ResultData    any `json:"resultData,omitempty"`
	StatementList []Statement `json:"statementList,omitempty"`
	StartingAfter string      `json:"startingAfter,omitempty"`
}

// StatementEntry from OpenAPI spec.
type StatementEntry struct {
	TransactionID            string      `json:"transactionId,omitempty"`
	CardID                   string      `json:"cardId,omitempty"`
	LastFourDigits           string      `json:"last_four_digits"`
	InternationalTransaction bool        `json:"internationalTransaction"`
	TransactionDate          any `json:"transactionDate"`
	TransactionDescription   string      `json:"transactionDescription"`
	Amount                   Amount      `json:"amount"`
	AmountDollar             *Amount     `json:"amountDollar,omitempty"`
	DebitOrCredit            string      `json:"debit_or_credit,omitempty"`
}

// StatementTransactionItem from OpenAPI spec.
type StatementTransactionItem struct {
	TransactionID            string      `json:"transactionId,omitempty"`
	CardID                   string      `json:"cardId,omitempty"`
	LastFourDigits           string      `json:"last_four_digits,omitempty"`
	InternationalTransaction bool        `json:"internationalTransaction,omitempty"`
	TransactionDate          any `json:"transactionDate,omitempty"`
	TransactionDescription   string      `json:"transactionDescription,omitempty"`
	TransactionType          string      `json:"transactionType,omitempty"`
	Amount                   any `json:"amount,omitempty"`
	AmountDollar             any `json:"amountDollar,omitempty"`
	DebitOrCredit            string      `json:"debit_or_credit,omitempty"`
}

// OpenStatement from OpenAPI spec.
type OpenStatement struct {
	AccountID                    string           `json:"accountId"`
	PaymentDue                   string           `json:"paymentDue"`
	CloseDate                    string           `json:"closeDate,omitempty"`
	PreviousBalance              any      `json:"previousBalance,omitempty"`
	Balance                      any      `json:"balance"`
	CreditLimit                  any      `json:"creditLimit"`
	WithdrawalCreditLimit        any      `json:"withdrawalCreditLimit"`
	CurrentCreditLimit           any      `json:"currentCreditLimit"`
	CurrentWithdrawalCreditLimit any      `json:"currentWithdrawalCreditLimit"`
	TransactionsList             []StatementEntry `json:"transactionsList"`
	FirstTransactionInserted     string           `json:"firstTransactionInserted,omitempty"`
	OpeningDateTime              string           `json:"openingDateTime,omitempty"`
	QueryDate                    string           `json:"query_date,omitempty"`
	SourceAudit                  *SourceAudit     `json:"sourceAudit,omitempty"`
}

// FutureStatement from OpenAPI spec.
type FutureStatement struct {
	AccountID        string           `json:"accountId"`
	PaymentDue       string           `json:"paymentDue"`
	Balance          any      `json:"balance"`
	TransactionsList []StatementEntry `json:"transactionsList"`
	QueryDate        string           `json:"query_date,omitempty"`
	SourceAudit      *SourceAudit     `json:"sourceAudit,omitempty"`
}

// FutureStatementList from OpenAPI spec.
type FutureStatementList struct {
	ResultData    any       `json:"resultData,omitempty"`
	StatementList []FutureStatement `json:"statementList,omitempty"`
}
