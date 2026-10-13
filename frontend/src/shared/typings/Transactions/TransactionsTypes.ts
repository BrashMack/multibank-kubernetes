import { PaginatedRequest } from '../App/AppTypes'

export interface Transaction {
  transactionId: string
  amount: string
  currency: 'RUB' | 'USD' | 'EUR' | 'GBP'
  type: 'debit' | 'credit'
  createdAt: string
  transactionInformation: string
}

export interface GetTransactionsByAccountRequest extends PaginatedRequest {
  accountId: string
}

export interface GetTransactionsByAccountResponse {
  transactions: Transaction[]
}

export interface MakeTransactionRequest {
  cardIdFrom: string
  cardIdTo: string
  amount: string
}
