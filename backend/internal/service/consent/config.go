package consent

import "github.com/masharpik/vtb-hack-back/pkg/enums"

type Config struct {
	BankToName map[enums.Bank]enums.Bank
}
