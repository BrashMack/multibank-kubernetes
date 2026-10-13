export type ConsentCardData = {
  bank: string
}

export interface GetAwaitingConsentResponse {
  banks: ConsentCardData[]
}
