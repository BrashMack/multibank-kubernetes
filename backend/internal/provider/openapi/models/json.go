package models

import (
	"time"

	"github.com/masharpik/vtb-hack-back/pkg/enums"
)

// GetTokenResponse - структура json-ответа от openapi на запрос по получению токена.
type GetTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ClientID    string `json:"client_id"`
	ExpiresIn   int    `json:"expires_in"`
}

type CreateReadConsentRequest struct {
	ClientID           string   `json:"client_id"`
	Permissions        []string `json:"permissions"`
	Reason             string   `json:"reason"`
	RequestingBank     string   `json:"requesting_bank"`
	RequestingBankName string   `json:"requesting_bank_name"`
}

type CreateReadConsentResponse struct {
	RequestID    string              `json:"request_id"`
	ConsentID    string              `json:"consent_id"`
	Status       enums.ConsentStatus `json:"status"`
	Message      string              `json:"message"`
	CreatedAt    string              `json:"created_at"`
	AutoApproved bool                `json:"auto_approved"`
}

type CreateWriteConsentRequest struct {
	RequestingBank string  `json:"requesting_bank"`
	ClientID       string  `json:"client_id"`
	ConsentType    string  `json:"consent_type"`
	DebtorAccount  string  `json:"debtor_account"`
	Amount         float64 `json:"amount"`
	Reference      string  `json:"reference"`
}

type CreateWriteConsentResponse struct {
	RequestID    string              `json:"request_id"`
	ConsentID    string              `json:"consent_id"`
	Status       enums.ConsentStatus `json:"status"`
	ConsentType  string              `json:"consent_type"`
	Message      string              `json:"message"`
	ValidUntil   string              `json:"valid_until"`
	AutoApproved bool                `json:"auto_approved"`
}

type CreateCreateConsentRequest struct {
	RequestingBank         string  `json:"requesting_bank"`
	ClientID               string  `json:"client_id"`
	ReadProductAgreements  bool    `json:"read_product_agreements"`
	OpenProductAgreements  bool    `json:"open_product_agreements"`
	CloseProductAgreements bool    `json:"close_product_agreements"`
	MaxAmount              float64 `json:"max_amount"`
	ValidUntil             string  `json:"valid_until"`
	Reason                 string  `json:"reason"`
}

type CreateCreateConsentResponse struct {
	RequestID    string              `json:"request_id"`
	ConsentID    string              `json:"consent_id"`
	Status       enums.ConsentStatus `json:"status"`
	AutoApproved bool                `json:"auto_approved"`
	Message      string              `json:"message"`
	ValidUntil   string              `json:"valid_until"`
}

type CreateAgreementRequest struct {
	ProductID       string  `json:"product_id"`
	Amount          float64 `json:"amount"`
	TermMonths      int     `json:"term_months"`
	SourceAccountID string  `json:"source_account_id"`
}

type CreateAgreementResponse struct {
	Data struct {
		AgreementID  string              `json:"agreement_id"`
		ProductID    string              `json:"product_id"`
		ProductName  string              `json:"product_name"`
		ProductType  enums.AgreementType `json:"product_type"`
		Amount       float64             `json:"amount"`
		Status       enums.AccountStatus `json:"status"`
		AutoApproved bool                `json:"auto_approved"`
		Message      string              `json:"message"`
		ValidUntil   string              `json:"valid_until"`
	} `json:"data"`
	Meta map[string]any `json:"meta"`
}

type GetCreateConsentResponse struct {
	ConsentID              string              `json:"consent_id"`
	GrantedTo              string              `json:"granted_to"`
	Status                 enums.ConsentStatus `json:"status"`
	ReadProductAgreements  bool                `json:"read_product_agreements"`
	OpenProductAgreements  bool                `json:"open_product_agreements"`
	CloseProductAgreements bool                `json:"close_product_agreements"`
	AllowedProductTypes    []string            `json:"allowed_product_types"`
	MaxAmount              float64             `json:"max_amount"`
	CurrentTotalOpened     float64             `json:"current_total_opened"`
	CreatedAt              string              `json:"created_at"`
	ValidUntil             string              `json:"valid_until"`
}

type GetAgreementResponse struct {
	Data struct {
		AgreementID    string              `json:"agreement_id"`
		ProductID      string              `json:"product_id"`
		ProductName    string              `json:"product_name"`
		ProductType    enums.AgreementType `json:"product_type"`
		InterestRate   float64             `json:"interest_rate"`
		Amount         float64             `json:"amount"`
		Status         enums.AccountStatus `json:"status"`
		StartDate      string              `json:"start_date"`
		EndDate        string              `json:"end_date"`
		AccountNumber  string              `json:"account_number"`
		AccountBalance float64             `json:"account_balance"`
	} `json:"data"`
}

type GetAgreementsResponse struct {
	Data []struct {
		AgreementID   string                `json:"agreement_id"`
		ProductID     string                `json:"product_id"`
		ProductName   string                `json:"product_name"`
		ProductType   enums.AgreementType   `json:"product_type"`
		Amount        float64               `json:"amount"`
		Status        enums.AgreementStatus `json:"status"`
		StartDate     string                `json:"start_date"`
		EndDate       string                `json:"end_date"`
		AccountNumber string                `json:"account_number"`
	} `json:"data"`
	Meta map[string]any `json:"meta"`
}

type GetProductsResponse struct {
	Data struct {
		Product []struct {
			ProductID    string              `json:"productId"`
			ProductType  enums.AgreementType `json:"productType"`
			ProductName  string              `json:"productName"`
			Description  string              `json:"description"`
			InterestRate string              `json:"interestRate"`
			MinAmount    string              `json:"minAmount"`
			MaxAmount    string              `json:"maxAmount"`
			TermMonths   int                 `json:"termMonths"`
		} `json:"product"`
	} `json:"data"`
}

type InstructedAmount struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type DebtorAccount struct {
	SchemeName     string `json:"schemeName"`
	Identification string `json:"identification"`
}

type CreditorAccount struct {
	SchemeName     string `json:"schemeName"`
	Identification string `json:"identification"`
}

type CreditorAccountWithBank struct {
	SchemeName     string `json:"schemeName"`
	Identification string `json:"identification"`
	BankCode       string `json:"bank_code"`
}

type Initiation struct {
	InstructedAmount InstructedAmount `json:"instructedAmount"`
	DebtorAccount    DebtorAccount    `json:"debtorAccount"`
	CreditorAccount  CreditorAccount  `json:"creditorAccount"`
}

type Data struct {
	Initiation Initiation `json:"initiation"`
}

type InitiationWithDescBank struct {
	InstructedAmount InstructedAmount        `json:"instructedAmount"`
	DebtorAccount    DebtorAccount           `json:"debtorAccount"`
	CreditorAccount  CreditorAccountWithBank `json:"creditorAccount"`
}

type DataWithDescBank struct {
	Initiation InitiationWithDescBank `json:"initiation"`
}

type CreatePaymentInRequest struct {
	Data Data `json:"data"`
}

type CreatePaymentOutRequest struct {
	Data DataWithDescBank `json:"data"`
}

type CreatePaymentResponse struct {
	Data struct {
		PaymentID            string              `json:"paymentId"`
		Status               enums.PaymentStatus `json:"status"`
		CreationDateTime     string              `json:"creationDateTime"`
		StatusUpdateDateTime string              `json:"statusUpdateDateTime"`
	}
}

type Links struct {
	Self string `json:"self"`
}

type GetConsentResponse struct {
	Data struct {
		ConsentID            string              `json:"consentId"`
		Status               enums.ConsentStatus `json:"status"`
		CreationDateTime     string              `json:"creationDateTime"`
		StatusUpdateDateTime string              `json:"statusUpdateDateTime"`
		Permissions          []string            `json:"permissions"`
		ExpirationDateTime   string              `json:"expirationDateTime"`
	} `json:"data"`
	Links Links          `json:"links"`
	Meta  map[string]any `json:"meta"`
}

type Account struct {
	AccountID      string              `json:"accountId"`
	Status         enums.AccountStatus `json:"status"`
	AccountType    string              `json:"accountType"`
	AccountSubType enums.AccountType   `json:"accountSubType"`
	Nickname       string              `json:"nickname"`
	Currency       string              `json:"currency"`
	OpeningDate    string              `json:"openingDate"`
	Cards          []struct {
		SchemeName     string `json:"schemeName"`
		Identification string `json:"identification"`
		Name           string `json:"name"`
	} `json:"account"`
}

type GetAccountsResponse struct {
	Data struct {
		Accounts []Account `json:"account"`
	} `json:"data"`
	Links Links          `json:"links"`
	Meta  map[string]any `json:"meta"`
}

type AccountData struct {
	AccountID      string              `json:"accountId"`
	Status         enums.AccountStatus `json:"status"`
	AccountType    string              `json:"accountType"`
	AccountSubType enums.AccountType   `json:"accountSubType"`
	Nickname       string              `json:"nickname"`
	Currency       string              `json:"currency"`
	OpeningDate    string              `json:"openingDate"`
	Description    string              `json:"description"`
}

type GetAccountResponse struct {
	Data struct {
		Accounts []AccountData `json:"account"`
	} `json:"data"`
	Links Links          `json:"links"`
	Meta  map[string]any `json:"meta"`
}

type Card struct {
	CardID         string           `json:"cardId"`
	CardNumber     string           `json:"cardNumber"`
	CardType       enums.CardType   `json:"cardType"`
	CardName       string           `json:"cardName"`
	HolderName     string           `json:"holderName"`
	ExpiryMonth    int              `json:"expiryMonth"`
	ExpiryYear     int              `json:"expiryYear"`
	AccountNumber  string           `json:"accountNumber"`
	AccountBalance string           `json:"accountBalance"`
	Status         enums.CardStatus `json:"status"`
	DailyLimit     string           `json:"dailyLimit"`
	MonthlyLimit   string           `json:"monthlyLimit"`
	IssuedAt       string           `json:"issuedAt"`
}

type GetCardsResponse struct {
	Data struct {
		Cards []Card `json:"cards"`
		Total int    `json:"total"`
	} `json:"data"`
	Meta map[string]any `json:"meta"`
}

type CardData struct {
	CardID         string           `json:"cardId"`
	CardNumber     string           `json:"cardNumber"`
	CardNumberFull string           `json:"cardNumberFull"`
	CardType       enums.CardType   `json:"cardType"`
	CardName       string           `json:"cardName"`
	HolderName     string           `json:"holderName"`
	ExpiryMonth    int              `json:"expiryMonth"`
	ExpiryYear     int              `json:"expiryYear"`
	AccountNumber  string           `json:"accountNumber"`
	AccountBalance string           `json:"accountBalance"`
	Status         enums.CardStatus `json:"status"`
	DailyLimit     string           `json:"dailyLimit"`
	MonthlyLimit   string           `json:"monthlyLimit"`
	IssuedAt       string           `json:"issuedAt"`
}

type GetCardResponse struct {
	CardData CardData       `json:"data"`
	Meta     map[string]any `json:"meta"`
}

type Amount struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type GetAccountBalanceResponse struct {
	Data struct {
		Balance []struct {
			AccountID            string            `json:"accountId"`
			Type                 enums.BalanceType `json:"type"`
			DateTime             string            `json:"dateTime"`
			Amount               Amount            `json:"amount"`
			CreditDebitIndicator string            `json:"creditDebitIndicator"`
		} `json:"balance"`
	} `json:"data"`
}

type GetAccountTransactionsResponse struct {
	Data struct {
		Transaction []Transaction `json:"transaction"`
	} `json:"data"`
	Links Links `json:"links"`
	Meta  struct {
		TotalPages   int `json:"totalPages"`
		TotalRecords int `json:"totalRecords"`
		CurrentPage  int `json:"currentPage"`
		PageSize     int `json:"pageSize"`
	} `json:"meta"`
}

type Transaction struct {
	TransactionID          string    `json:"transactionId"`
	BookingDateTime        time.Time `json:"bookingDateTime"`
	ValueDateTime          time.Time `json:"valueDateTime"`
	Amount                 Amount    `json:"amount"`
	CreditDebitIndicator   string    `json:"creditDebitIndicator"`
	Reference              string    `json:"reference"`
	RemittanceInformation  string    `json:"remittanceInformation"`
	TransactionInformation string    `json:"transactionInformation"`
	CounterpartyAccount    string    `json:"counterpartyAccount"`
	CounterpartyBank       string    `json:"counterpartyBank"`
	Merchant               *struct {
		MerchantID string `json:"merchantId"`
		Name       string `json:"name"`
		MccCode    string `json:"mccCode"`
		Category   string `json:"category"`
		City       string `json:"city"`
		Country    string `json:"country"`
		Address    string `json:"address"`
	} `json:"merchant,omitempty"`
}
