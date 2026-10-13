export interface CardData {
  cardId: string
  cardNumber: string
  cardName: string
  balance: string
  currency: 'RUB' | 'USD' | 'EUR' | 'GBP'
  bank: string
  backgroundImage: string
  type: 'debit' | 'credit'
}

export interface GetCardsResponse {
  cards: CardData[]
}

export interface ExtendedCardData {
  number: string
  cardId: string
  cardNumber: string
  cardNumberFull: string
  cardName: string
  holderName: string
  expiryMonth: number
  expiryYear: number
  dailyLimit: string
  monthlyLimit: string
  balance: string
  currency: 'RUB' | 'USD' | 'EUR' | 'GBP'
  bank: string
  backgroundImage: string
  type: 'debit' | 'credit'
}
