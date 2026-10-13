export type RequestConsentsResponse = {
  auto_approve: boolean
}

export type GetActiveConsentsResponse = {
  confirmedBanks: string[]
  awaitingBanks: string[]
}
