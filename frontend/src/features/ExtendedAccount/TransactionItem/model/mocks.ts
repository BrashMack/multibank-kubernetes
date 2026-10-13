import { Transaction } from '~/shared/typings/Transactions/TransactionsTypes'

export const mockTransactions: Transaction[] = [
  {
    transactionId: 'tx_001',
    amount: '1250.50',
    currency: 'RUB',
    type: 'debit',
    createdAt: '2025-01-15T10:24:00Z',
    transactionInformation: 'Оплата в магазине Магнит',
  },
  {
    transactionId: 'tx_002',
    amount: '5600',
    currency: 'USD',
    type: 'credit',
    createdAt: '2025-01-14T18:10:00Z',
    transactionInformation: 'Зарплата',
  },
  {
    transactionId: 'tx_003',
    amount: '249',
    currency: 'EUR',
    type: 'debit',
    createdAt: '2025-01-13T08:50:00Z',
    transactionInformation: 'Подписка Spotify',
  },
]
