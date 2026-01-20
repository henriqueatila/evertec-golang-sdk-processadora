package types

import "fmt"

// ErrorCode represents PaySmart API error codes.
// Documentation: https://paysmart-api.gitlab.io/processadora/PT-br/docs/codigos-erro-api
type ErrorCode int

// Error code constants from PaySmart API.
const (
	// Success
	ErrorCodeSuccess ErrorCode = 0

	// Communication errors
	ErrorCodeCommunication ErrorCode = 6

	// Card errors (14-24)
	ErrorCodeInvalidCard                  ErrorCode = 14
	ErrorCodeCardRequestedEmissionReject  ErrorCode = 15
	ErrorCodeCardInEmission               ErrorCode = 16
	ErrorCodeInvalidBlockingCode          ErrorCode = 17
	ErrorCodeBlockingCodeUnlockMismatch   ErrorCode = 18
	ErrorCodeBlockingCodeLockMismatch     ErrorCode = 19
	ErrorCodeUnlockNotPermitted           ErrorCode = 20
	ErrorCodeCardAlreadyUnlocked          ErrorCode = 21
	ErrorCodeCardBlockedHigherSeverity    ErrorCode = 22
	ErrorCodeCardStatusPreventsReissuance ErrorCode = 23
	ErrorCodePendingReissuanceExists      ErrorCode = 24

	// Password/CVV errors (62-75)
	ErrorCodeInvalidPasswordData ErrorCode = 62
	ErrorCodePasswordLocked      ErrorCode = 70
	ErrorCodeInvalidPassword     ErrorCode = 71
	ErrorCodeInvalidPasswordCVV  ErrorCode = 75

	// System errors
	ErrorCodeSystemUnavailable ErrorCode = 91

	// Validation errors (104-170)
	ErrorCodeInvalidData                    ErrorCode = 104
	ErrorCodeAccountOrCardBlocked           ErrorCode = 109
	ErrorCodeInvalidAccount                 ErrorCode = 110
	ErrorCodeMandatoryDataMissing           ErrorCode = 115
	ErrorCodeInvalidFieldFormat             ErrorCode = 120
	ErrorCodeCardAlreadyAssociated          ErrorCode = 121
	ErrorCodeAnonymousAccountNoNewCards     ErrorCode = 122
	ErrorCodeAccountCanceled                ErrorCode = 123
	ErrorCodeValidationFailures             ErrorCode = 124
	ErrorCodeTransactionNotFound            ErrorCode = 125
	ErrorCodeInvalidTokenizationResponse    ErrorCode = 127
	ErrorCodeNoVirtualPANAvailable          ErrorCode = 128
	ErrorCodeCryptographicPANFailed         ErrorCode = 129
	ErrorCodeVirtualCardAnonymousDisallowed ErrorCode = 130
	ErrorCodeAnonymousCardPrintingOnly      ErrorCode = 131
	ErrorCodeQRCodeTokenFailed              ErrorCode = 132
	ErrorCodeFlagAPIKeyError                ErrorCode = 133
	ErrorCodeInvalidFlagQRCodeResponse      ErrorCode = 134
	ErrorCodeInvalidPINFormat               ErrorCode = 135
	ErrorCodeVirtualCardDataUnavailable     ErrorCode = 136
	ErrorCodeQRCodeEncryptionFailed         ErrorCode = 137
	ErrorCodeDuplicateQRCodeTransaction     ErrorCode = 147
	ErrorCodeAcquirerNotConfiguredQR        ErrorCode = 148
	ErrorCodeAccountClosed                  ErrorCode = 149
	ErrorCodeCardCanceled                   ErrorCode = 150
	ErrorCodeDebitConfirmationFailed        ErrorCode = 151
	ErrorCodeEmbossingFileNoStatus          ErrorCode = 157
	ErrorCodeCardNoEmbossingFile            ErrorCode = 158
	ErrorCodePhysicalCardChildProhibited    ErrorCode = 159
	ErrorCodeAccountChildProductProhibited  ErrorCode = 160
	ErrorCodeVirtualCardParentProductError  ErrorCode = 162
	ErrorCodeCardholderAssociationRequired  ErrorCode = 163
	ErrorCodeAltBindingKeyMultipleCards     ErrorCode = 164
	ErrorCodeAltBindingKeySecurityMismatch  ErrorCode = 165
	ErrorCodeAltBindingKeyOrDataRequired    ErrorCode = 166
	ErrorCodeAltBindingKeyInUse             ErrorCode = 167
	ErrorCodeNewAltBindingKeyRequired       ErrorCode = 168
	ErrorCodeLinkedPhysicalCardInvalid      ErrorCode = 169
	ErrorCodePhysicalCardUpdateInapplicable ErrorCode = 170

	// API errors (804-813)
	ErrorCodeProductNotRegistered     ErrorCode = 804
	ErrorCodeIdempotentRequestRunning ErrorCode = 808
	ErrorCodeAccountAlreadyExists     ErrorCode = 809
	ErrorCodeInvalidJSONStructure     ErrorCode = 810
	ErrorCodeAPIKeyInvalid            ErrorCode = 811
	ErrorCodeTransactionAPIError      ErrorCode = 812
	ErrorCodeIdempotencyKeyMissing    ErrorCode = 813

	// Internal errors (977-999)
	ErrorCodePANHashFailed             ErrorCode = 977
	ErrorCodeInvoiceExpirationInvalid  ErrorCode = 979
	ErrorCodeAccountCreationBalanceErr ErrorCode = 980
	ErrorCodeQRCodeEloConfigNotFound   ErrorCode = 981
	ErrorCodeNoDisputesFound           ErrorCode = 982
	ErrorCodeDisputeNotInResubmission  ErrorCode = 983
	ErrorCodeDisputeNotFound           ErrorCode = 984
	ErrorCodeDisputeInitialTxMissing   ErrorCode = 985
	ErrorCodeInclusiveTxAlreadyExists  ErrorCode = 986
	ErrorCodePartialValueExceedsTotal  ErrorCode = 987
	ErrorCodeDisputeReasonNotFound     ErrorCode = 988
	ErrorCodeTransactionParamsNotFound ErrorCode = 989
	ErrorCodeDisputedTxNotSettled      ErrorCode = 990
	ErrorCodeDisputedTxCanceled        ErrorCode = 991
	ErrorCodeDisputeAlreadyExists      ErrorCode = 992
	ErrorCodeTransactionAPIConfigError ErrorCode = 993
	ErrorCodeBufferServiceConfigError  ErrorCode = 994
	ErrorCodeTokenizationConfigError   ErrorCode = 995
	ErrorCodeIssuerConfigNotFound      ErrorCode = 996
	ErrorCodePINValidationConfigError  ErrorCode = 997
	ErrorCodeInternalWriteError        ErrorCode = 998
	ErrorCodeInternalError             ErrorCode = 999
)

// errorMessages maps error codes to their Portuguese descriptions.
// Messages are exactly as documented at:
// https://paysmart-api.gitlab.io/processadora/PT-br/docs/codigos-erro-api
var errorMessages = map[ErrorCode]string{
	ErrorCodeSuccess:                        "Comando concluído com sucesso",
	ErrorCodeCommunication:                  "Erro de comunicação",
	ErrorCodeInvalidCard:                    "Cartão inválido",
	ErrorCodeCardRequestedEmissionReject:    "Cartão requisitado, mas emissão rejeitada",
	ErrorCodeCardInEmission:                 "Cartão em processo de emissão / sem dados para cifragem",
	ErrorCodeInvalidBlockingCode:            "Código de bloqueio inválido",
	ErrorCodeBlockingCodeUnlockMismatch:     "Código de bloqueio não corresponde a uma operação de desbloqueio",
	ErrorCodeBlockingCodeLockMismatch:       "Código de bloqueio não corresponde a uma operação de bloqueio",
	ErrorCodeUnlockNotPermitted:             "Desbloqueio não permitido para este cartão",
	ErrorCodeCardAlreadyUnlocked:            "Cartão já desbloqueado",
	ErrorCodeCardBlockedHigherSeverity:      "Cartão já bloqueado por um motivo mais crítico ao informado",
	ErrorCodeCardStatusPreventsReissuance:   "Status atual do cartão não permite reemissão",
	ErrorCodePendingReissuanceExists:        "Já existe uma outra solicitação de reemissão pendente para este cartão",
	ErrorCodeInvalidPasswordData:            "Dados inválidos de senha para validação",
	ErrorCodePasswordLocked:                 "Senha está bloqueada",
	ErrorCodeInvalidPassword:                "Senha inválida",
	ErrorCodeInvalidPasswordCVV:             "Senha ou CVV inválido",
	ErrorCodeSystemUnavailable:              "Sistema indisponível",
	ErrorCodeInvalidData:                    "Dados inválidos",
	ErrorCodeAccountOrCardBlocked:           "Conta / cartão bloqueado",
	ErrorCodeInvalidAccount:                 "Conta inválida",
	ErrorCodeMandatoryDataMissing:           "Faltam dados mandatórios",
	ErrorCodeInvalidFieldFormat:             "Formato inválido para o campo",
	ErrorCodeCardAlreadyAssociated:          "Cartão já associado a um portador",
	ErrorCodeAnonymousAccountNoNewCards:     "Conta vinculada a um cartão anônimo não pode receber novos cartões",
	ErrorCodeAccountCanceled:                "Conta cancelada",
	ErrorCodeValidationFailures:             "Validações falharam",
	ErrorCodeTransactionNotFound:            "Transação não encontrada",
	ErrorCodeInvalidTokenizationResponse:    "Resposta inválida do serviço de tokenização",
	ErrorCodeNoVirtualPANAvailable:          "Falha na obtenção de um PAN virtual disponível",
	ErrorCodeCryptographicPANFailed:         "Não foi possível fazer a segurança criptográfica do PAN",
	ErrorCodeVirtualCardAnonymousDisallowed: "Geração de cartões virtuais não permitida para cartões anônimos",
	ErrorCodeAnonymousCardPrintingOnly:      "Cartões só podem ser anônimos para impressão. Se não for um cartão para impressão deve haver um cardholder informado",
	ErrorCodeQRCodeTokenFailed:              "Não foi possível obter o token de autenticação para validação de QRCode",
	ErrorCodeFlagAPIKeyError:                "Erro na API de obtenção da chave da bandeira",
	ErrorCodeInvalidFlagQRCodeResponse:      "Resposta inválida da bandeira na obtenção do resultado da transação com QR Code",
	ErrorCodeInvalidPINFormat:               "Formato inválido para a entrada de PIN. ISO-1 mandatório",
	ErrorCodeVirtualCardDataUnavailable:     "Não foi possível obter os dados do cartão virtual",
	ErrorCodeQRCodeEncryptionFailed:         "Não foi possível cifrar os dados para transação de QRCode",
	ErrorCodeDuplicateQRCodeTransaction:     "Transação de QRCode duplicada",
	ErrorCodeAcquirerNotConfiguredQR:        "Adquirente não configurado para aceitar transações de crédito com QRCode",
	ErrorCodeAccountClosed:                  "Conta fechada",
	ErrorCodeCardCanceled:                   "Cartão cancelado",
	ErrorCodeDebitConfirmationFailed:        "Não foi possível confirmar a transação de débito",
	ErrorCodeEmbossingFileNoStatus:          "Arquivo de embossing não possui status",
	ErrorCodeCardNoEmbossingFile:            "Cartão não possui um arquivo de embossing",
	ErrorCodePhysicalCardChildProhibited:    "Não é possível emitir cartões físicos com produtos filhos",
	ErrorCodeAccountChildProductProhibited:  "Não é possível criar ou modificar contas para produtos filhos",
	ErrorCodeVirtualCardParentProductError:  "Cartões virtuais criados abaixo de cartões já existentes não podem utilizar produtos pais",
	ErrorCodeCardholderAssociationRequired:  "Associação de portadores a cartões criados sem cardholder tem que ser aplicada a todos os cartões vinculados",
	ErrorCodeAltBindingKeyMultipleCards:     "O campo alternativeBindingKey informado se refere a mais de um cartão",
	ErrorCodeAltBindingKeySecurityMismatch:  "O campo alternativeBindingKey e os dados de seguranca do cartao nao conferem",
	ErrorCodeAltBindingKeyOrDataRequired:    "Deve ser enviado o alternativeBindingKey ou os dados do cartao (PAN/CVV/dateExp completos)",
	ErrorCodeAltBindingKeyInUse:             "Chave alternativa de binding ja utilizada",
	ErrorCodeNewAltBindingKeyRequired:       "O envio de uma nova chave alternativa de binding é obrigatório para reemissão de cartões anônimos que já possuam uma chave alternativa anterior",
	ErrorCodeLinkedPhysicalCardInvalid:      "Cartão físico vinculado não existente, não pertence à conta ou não é físico",
	ErrorCodePhysicalCardUpdateInapplicable: "Atualização de cartão físico vinculado não faz sentido para cartão físico",
	ErrorCodeProductNotRegistered:           "Produto não cadastrado para este emissor",
	ErrorCodeIdempotentRequestRunning:       "Outra requisição idempotente já está rodando",
	ErrorCodeAccountAlreadyExists:           "Já existe uma conta / solicitação de conta para este documento / produto / solicitante",
	ErrorCodeInvalidJSONStructure:           "Não foi possível encontrar a estrutura JSON esperada no corpo da requisição",
	ErrorCodeAPIKeyInvalid:                  "API-Key não enviada ou inválida",
	ErrorCodeTransactionAPIError:            "Erro ao obter transações da API interna de transações",
	ErrorCodeIdempotencyKeyMissing:          "Campo IdempotencyKey ausente no cabeçalho",
	ErrorCodePANHashFailed:                  "Não foi possível calcular um hash de comparação para o PAN",
	ErrorCodeInvoiceExpirationInvalid:       "Data de vencimento da fatura deve ser maior que 0 e menor que 29",
	ErrorCodeAccountCreationBalanceErr:      "Erro na criação da conta no processo de saldos e limites",
	ErrorCodeQRCodeEloConfigNotFound:        "Não foi possível encontrar configurações do serviço de QR Code Elo",
	ErrorCodeNoDisputesFound:                "Nenhuma disputa encontrada com os parâmetros de busca utilizados",
	ErrorCodeDisputeNotInResubmission:       "Disputa não está na etapa de reapresentação",
	ErrorCodeDisputeNotFound:                "Disputa não encontrada",
	ErrorCodeDisputeInitialTxMissing:        "Transação inclusiva inicial (TE-05 ou TE-06) não existe para essa disputa",
	ErrorCodeInclusiveTxAlreadyExists:       "Transação inclusiva já existe para essa transação",
	ErrorCodePartialValueExceedsTotal:       "Valor parcial não pode ser maior ou igual ao valor total da transação",
	ErrorCodeDisputeReasonNotFound:          "Configurações para a razão da disputa não existente",
	ErrorCodeTransactionParamsNotFound:      "Não foi possível encontrar alguns parâmetros da transação",
	ErrorCodeDisputedTxNotSettled:           "Transação a ser disputada não foi confirmada na liquidação, tente novamente nos próximos dias",
	ErrorCodeDisputedTxCanceled:             "Transação a ser disputada está cancelada",
	ErrorCodeDisputeAlreadyExists:           "Disputa já existe para essa transação",
	ErrorCodeTransactionAPIConfigError:      "Não foi possível encontrar configurações da API de Transações para o emissor representado pela API-Key informada",
	ErrorCodeBufferServiceConfigError:       "Não foi possível encontrar configurações do Colchão para o emissor representado pela API-Key informada",
	ErrorCodeTokenizationConfigError:        "Não foi possível encontrar as configurações do serviço de tokenização",
	ErrorCodeIssuerConfigNotFound:           "Não foi possível encontrar configurações para o emissor representado pela API-Key informada",
	ErrorCodePINValidationConfigError:       "Não foi possível encontrar as configurações para validar o PIN / CVV para esse cartão",
	ErrorCodeInternalWriteError:             "Erro interno de gravação",
	ErrorCodeInternalError:                  "Erro interno paySmart",
}

// Message returns the human-readable description in Portuguese.
// Returns empty string if code is unknown.
func (e ErrorCode) Message() string {
	if msg, ok := errorMessages[e]; ok {
		return msg
	}
	return ""
}

// Error implements the error interface.
func (e ErrorCode) Error() string {
	if msg := e.Message(); msg != "" {
		return fmt.Sprintf("paySmart error %d: %s", e, msg)
	}
	return fmt.Sprintf("paySmart error %d", e)
}

// IsSuccess returns true if the code indicates success (0).
func (e ErrorCode) IsSuccess() bool {
	return e == ErrorCodeSuccess
}
