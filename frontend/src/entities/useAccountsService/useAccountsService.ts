import { useCallback, useMemo } from 'react'

import { axiosInstance } from '~/shared'
import { ACCESS_TOKEN_KEY } from '~/shared/constants/authConstants'
import { ExtendedAccountData, GetAccountsResponse } from '~/shared/typings/Accounts/AccountsTypes'
import { getDataFromStore } from '~/shared/utils/tokenSaving/tokenSaving'

const useAccountsService = () => {
  const getAccounts = useCallback(async (signal: AbortSignal): Promise<GetAccountsResponse> => {
    const accessToken = getDataFromStore(ACCESS_TOKEN_KEY)

    if (!accessToken) throw new Error('No access token found')
    const response = await axiosInstance.get(`/accounts`, {
      signal,
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    })

    return response.data
  }, [])

  const getAccountById = useCallback(
    async (accountId: string, signal: AbortSignal): Promise<ExtendedAccountData> => {
      const accessToken = getDataFromStore(ACCESS_TOKEN_KEY)
      if (!accessToken) throw new Error('No access token found')
      const response = await axiosInstance.get<ExtendedAccountData>(`/account/${accountId}`, {
        signal,
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      })

      return response.data
    },
    [],
  )

  return useMemo(() => ({ getAccounts, getAccountById }), [getAccounts, getAccountById])
}

export default useAccountsService
