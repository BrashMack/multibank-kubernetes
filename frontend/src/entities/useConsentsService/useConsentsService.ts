import { useCallback, useMemo } from 'react'

import { axiosInstance } from '~/shared'
import { ACCESS_TOKEN_KEY } from '~/shared/constants/authConstants'
import {
  GetActiveConsentsResponse,
  RequestConsentsResponse,
} from '~/shared/typings/Consents/ConsetnsTypes'
import { GetAwaitingConsentResponse } from '~/shared/typings/Consents/WaitConsentCardTypes'
import { getDataFromStore } from '~/shared/utils/tokenSaving/tokenSaving'

const useConsentsService = () => {
  const requestConsent = useCallback(
    async (bank: string, type: string, signal?: AbortSignal): Promise<RequestConsentsResponse> => {
      const accessToken = getDataFromStore(ACCESS_TOKEN_KEY)
      if (!accessToken) throw new Error('No access token found')
      const response = await axiosInstance.post(
        `/consent`,
        {
          bank,
          type,
        },
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

  const getActiveConsents = useCallback(
    async (signal?: AbortSignal): Promise<GetActiveConsentsResponse> => {
      const accessToken = getDataFromStore(ACCESS_TOKEN_KEY)
      if (!accessToken) throw new Error('No access token found')
      const response = await axiosInstance.get(
        `/consents/active`,

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

  const getAwaitingConsents = useCallback(
    async (signal?: AbortSignal): Promise<GetAwaitingConsentResponse> => {
      const accessToken = getDataFromStore(ACCESS_TOKEN_KEY)
      if (!accessToken) throw new Error('No access token found')
      const response = await axiosInstance.get(`/consents/awaiting`, {
        signal,
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      })
      return response.data
    },
    [],
  )

  return useMemo(
    () => ({ requestConsent, getActiveConsents, getAwaitingConsents }),
    [requestConsent, getActiveConsents, getAwaitingConsents],
  )
}

export default useConsentsService
