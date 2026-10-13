import { GetAccountsResponse } from '~/shared/typings/Accounts/AccountsTypes'

export const ACCOUNT_MOCK: GetAccountsResponse = {
  accounts: [
    {
      number: '1234 5678 9012 3456',
      balance: '150,000.00',
      currency: 'USD',
      bank: 'VTB',
      type: 'Checking',
      numberCards: [],
    },
    {
      number: '1234 5678 9012 3457',
      balance: '150,000.00',
      currency: 'GBP',
      bank: 'Sberbank',
      type: 'Checking',
      numberCards: [],
    },
    {
      number: '6554 3210 9876 5432',
      balance: '950.00',
      currency: 'RUB',
      bank: 'Alfabank',
      type: 'Card',
      numberCards: [],
    },
    {
      number: '6554 3210 9876 5431',
      balance: '950.00',
      currency: 'RUB',
      bank: 'Alfabank',
      type: 'loan',
      numberCards: [],
    },
    {
      number: '6554 3210 9876 5488',
      balance: '50.00',
      currency: 'RUB',
      bank: 'Alfabank',
      type: 'card',
      numberCards: [],
    },
  ],
}
