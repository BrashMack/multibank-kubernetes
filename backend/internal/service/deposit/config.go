package deposit

import "github.com/masharpik/vtb-hack-back/pkg/enums"

type PriorityConfig struct {
	ProductID    string
	ProductName  string
	Bank         enums.Bank
	Currency     string
	InterestRate string
	MinAmount    string
	MaxAmount    string
	TermMonths   int
}

type Config struct {
	DefaultCurrency string `default:"RUB"`
	Banks           []enums.Bank
	BankToName      map[enums.Bank]enums.Bank
	Priority        PriorityConfig
}
