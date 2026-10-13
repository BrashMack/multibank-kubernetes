import { useCallback, useMemo } from 'react'

import { axiosInstance } from '~/shared'
import { ACCESS_TOKEN_KEY } from '~/shared/constants/authConstants'
import {
  DepositProduct,
  GetDepositAgreementsResponse,
  GetDepositProductsResponse,
  MakeDepositRequest,
  MakeDepositResponse,
} from '~/shared/typings/Deposits/DepositsTypes'
import { getDataFromStore } from '~/shared/utils/tokenSaving/tokenSaving'

const useDepositService = () => {
  const getDepositAgreements = useCallback(
    async (signal: AbortSignal): Promise<GetDepositAgreementsResponse> => {
      const accessToken = getDataFromStore(ACCESS_TOKEN_KEY)
      if (!accessToken) throw new Error('No access token found')

      const response = await axiosInstance.get<GetDepositAgreementsResponse>(
        '/deposits/agreements',
        {
          signal,
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        },
      )

      return response.data
    },
    [],
  )

  const getPriorityDeposit = useCallback(async (signal: AbortSignal): Promise<DepositProduct> => {
    const accessToken = getDataFromStore(ACCESS_TOKEN_KEY)
    if (!accessToken) throw new Error('No access token found')

    const response = await axiosInstance.get<DepositProduct>('/deposit/priority', {
      signal,
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    })

    return response.data
  }, [])

  const getDepositProducts = useCallback(
    async (signal: AbortSignal): Promise<GetDepositProductsResponse> => {
      const accessToken = getDataFromStore(ACCESS_TOKEN_KEY)
      if (!accessToken) throw new Error('No access token found')

      const response = await axiosInstance.get<GetDepositProductsResponse>('/deposits', {
        signal,
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      })

      return response.data
    },
    [],
  )

  const makeDeposit = useCallback(
    async (
      { productId, bank, sourceAccountId, amount, termMonths }: MakeDepositRequest,
      signal?: AbortSignal,
    ): Promise<MakeDepositResponse> => {
      const accessToken = getDataFromStore(ACCESS_TOKEN_KEY)
      if (!accessToken) throw new Error('No access token found')
      const response = await axiosInstance.post<MakeDepositResponse>(
        `/deposits/agreement`,
        { productId, bank, amount, termMonths, sourceAccountId: sourceAccountId },
        {
          signal,
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        },
      )

      return response.data
    },
    [],
  )

  return useMemo(
    () => ({ getDepositAgreements, getPriorityDeposit, getDepositProducts, makeDeposit }),
    [getDepositAgreements, getPriorityDeposit, getDepositProducts, makeDeposit],
  )
}

export default useDepositService
