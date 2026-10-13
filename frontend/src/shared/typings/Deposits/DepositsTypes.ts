export type DepositCurrency = 'RUB' | 'USD' | 'EUR' | 'GBP'

export interface DepositAgreement {
  agreementId: string
  productId: string
  name: string
  balance: string
  currency: DepositCurrency
  interestRate: string
  accountNumber: string
  bank: string
}

export interface GetDepositAgreementsResponse {
  deposits: DepositAgreement[]
}

export interface DepositProduct {
  productId: string
  productName: string
  interestRate: string
  minAmount: string
  maxAmount: string
  currency: DepositCurrency
  termMonths: number
  bank: string
}

export interface GetDepositProductsResponse {
  deposits: DepositProduct[]
}

export interface MakeDepositRequest {
  productId: string
  bank: string
  sourceAccountId: string
  amount: string
  termMonths: number
}

export type MakeDepositResponse = DepositAgreement
