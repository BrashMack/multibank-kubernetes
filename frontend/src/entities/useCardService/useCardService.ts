import { useCallback, useMemo } from 'react'

import { axiosInstance } from '~/shared'
import { ACCESS_TOKEN_KEY } from '~/shared/constants/authConstants'
import { ExtendedCardData, GetCardsResponse } from '~/shared/typings/Card/CardTypes'
import { getDataFromStore } from '~/shared/utils/tokenSaving/tokenSaving'

const useCardService = () => {
  const getCards = useCallback(async (signal: AbortSignal): Promise<GetCardsResponse> => {
    const accessToken = getDataFromStore(ACCESS_TOKEN_KEY)

    if (!accessToken) throw new Error('No access token found')
    const response = await axiosInstance.get(`/cards`, {
      signal,
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    })

    return response.data
  }, [])

  const getCardById = useCallback(
    async (cardId: string, signal: AbortSignal): Promise<ExtendedCardData> => {
      const accessToken = getDataFromStore(ACCESS_TOKEN_KEY)
      if (!accessToken) throw new Error('No access token found')
      const response = await axiosInstance.get(`/card/${cardId}`, {
        signal,
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      })

      return response.data
    },
    [],
  )
  return useMemo(() => ({ getCards, getCardById }), [getCards, getCardById])
}

export default useCardService
