import { useMemo } from 'react'

import useAccountsService from './useAccountsService/useAccountsService'
import useAuthService from './useAuthService/useAuthService'
import useCardService from './useCardService/useCardService'
import useConsentsService from './useConsentsService/useConsentsService'
import useDepositService from './useDepositService/useDepositService'
import useTransactionsService from './useTransactionsService/useTransactionsService'
// Сюда импортим сервисы с запросами к бэкенду
const useService = () => {
  const cardsService = useCardService()
  const consentService = useConsentsService()
  const authService = useAuthService()
  const accountsService = useAccountsService()
  const transactionsService = useTransactionsService()
  const depositService = useDepositService()

  return useMemo(
    () => ({
      ...cardsService,
      ...consentService,
      ...authService,
      ...accountsService,
      ...depositService,
      ...transactionsService,
    }),
    [
      cardsService,
      consentService,
      authService,
      accountsService,
      depositService,
      transactionsService,
    ],
  )
}

export default useService
