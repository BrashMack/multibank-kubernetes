import { StackScreenProps } from '@react-navigation/stack'

export type RootStackParamList = {
  main: undefined
  login: undefined
  signup: undefined
  consention: undefined
  card: { cardId: string }
  account: { accountId: string }
  qrPayment: { accountIdTo: string; amount: string }
  showcase: undefined
}

export type CardPageRouteProp = StackScreenProps<RootStackParamList, 'card'>

export type AccountPageRouteProp = StackScreenProps<RootStackParamList, 'account'>

export type QRPaymentPageRouteProp = StackScreenProps<RootStackParamList, 'qrPayment'>

export interface PaginatedRequest {
  page: number
  limit: number
}
