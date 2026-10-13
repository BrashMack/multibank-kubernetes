import { useCallback, useMemo } from 'react'

import { axiosInstance } from '~/shared'
import { ACCESS_TOKEN_KEY } from '~/shared/constants/authConstants'
import {
  GetTransactionsByAccountRequest,
  GetTransactionsByAccountResponse,
  MakeTransactionRequest,
} from '~/shared/typings/Transactions/TransactionsTypes'
import { getDataFromStore } from '~/shared/utils/tokenSaving/tokenSaving'

const useTransactionsService = () => {
  const getTransactionsByAccount = useCallback(
    async (
      { accountId, page, limit }: GetTransactionsByAccountRequest,
      signal?: AbortSignal,
    ): Promise<GetTransactionsByAccountResponse> => {
      const accessToken = getDataFromStore(ACCESS_TOKEN_KEY)
      if (!accessToken) throw new Error('No access token found')
      const response = await axiosInstance.get(
        `/account/${accountId}/transactions`,

        {
          signal,
          params: { limit, page },
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        },
      )

      return response.data
    },
    [],
  )

  const makeTransaction = useCallback(
    async ({ cardIdFrom, cardIdTo, amount }: MakeTransactionRequest, signal?: AbortSignal) => {
      const accessToken = getDataFromStore(ACCESS_TOKEN_KEY)
      if (!accessToken) throw new Error('No access token found')
      await axiosInstance.post(
        `/card/payment`,
        { cardIdFrom, cardIdTo, amount },
        {
          signal,
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        },
      )
    },
    [],
  )

  return useMemo(
    () => ({ getTransactionsByAccount, makeTransaction }),
    [getTransactionsByAccount, makeTransaction],
  )
}

export default useTransactionsService
