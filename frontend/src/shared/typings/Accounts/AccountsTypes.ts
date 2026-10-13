export interface AccountData {
  accountId: string
  number: string
  balance: string
  currency: 'RUB' | 'USD' | 'EUR' | 'GBP'
  bank: string
  type: string
}

export interface NumberCards {
  cardId: string
  number: string
  backgroundImage: string
  type: 'debit' | 'credit'
  cardName: string
}

export interface GetAccountsResponse {
  accounts: AccountData[]
}

export interface AccountIdContext {
  accountId: string
  setAccountId: (accountId: string) => void
}

export interface ExtendedAccountData {
  accountId: string
  number: string
  balance: string
  currency: 'RUB' | 'USD' | 'EUR' | 'GBP'
  bank: string
  type: string
  numberCards: Array<NumberCards>
}
