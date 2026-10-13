import { Box, ScrollView, Text } from '@gluestack-ui/themed'
import { useNavigation } from '@react-navigation/native'
import { StackNavigationProp } from '@react-navigation/stack'
import { useQuery } from '@tanstack/react-query'
import React, { useCallback, useEffect, useMemo } from 'react'

import useService from '~/entities/useService'
import Card from '~/features/Card/ui/Card'
import WaitConsentCard from '~/features/Consents/Awaiting/ui/WaitConsentCard'
import RequestConsentCard from '~/features/Consents/RequestConsentCard/ui/RequestConsentCard'
import { ErrorState } from '~/features/ServiceStates/ErrorState'
import SkeletonCard from '~/features/ServiceStates/SkeletonCard'
import { CARDS_COUNT_KEY, DEFAULT_SKELETON_COUNT } from '~/shared/constants/commonConstants'
import { RootStackParamList } from '~/shared/typings/App/AppTypes'
import { getDataFromStore, setDataInStore } from '~/shared/utils/tokenSaving/tokenSaving'
import { isVoid } from '~/shared/utils/validators/isVoid'

interface CardFeedProps {
  onConsetionProcessStart?: () => void
}

const CardFeed: React.FC<CardFeedProps> = () => {
  const { getCards, getAwaitingConsents } = useService()
  const navigation = useNavigation<StackNavigationProp<RootStackParamList>>()

  const onConsetionProcessStart = useCallback(() => {
    navigation.navigate('consention')
  }, [navigation])

  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: ['cards'],
    queryFn: ({ signal }) => getCards(signal),
  })

  const { data: consentCards } = useQuery({
    queryKey: ['consentCards'],
    queryFn: ({ signal }) => getAwaitingConsents(signal),
  })

  const cardsLength = data?.cards?.length
  const skeletonCount = useMemo(() => {
    const storedCount = Number(getDataFromStore(CARDS_COUNT_KEY))
    return !Number.isNaN(storedCount) && storedCount > 0 ? storedCount : DEFAULT_SKELETON_COUNT
  }, [])

  useEffect(() => {
    if (!isLoading && !isVoid(cardsLength)) {
      const nextCount = cardsLength > 0 ? cardsLength : DEFAULT_SKELETON_COUNT
      setDataInStore(CARDS_COUNT_KEY, String(nextCount))
    }
  }, [cardsLength, isLoading])

  const onCardPress = useCallback(
    (cardId: string) => {
      navigation.navigate('card', { cardId })
    },
    [navigation],
  )

  return (
    <Box p='$4'>
      <Text color='$white' fontSize='$2xl' fontWeight='$bold' mb='$4'>
        Карты
      </Text>
      {isError && <ErrorState message='Не удалось загрузить карты' onRetry={refetch} />}
      {!isError && (
        <ScrollView horizontal showsHorizontalScrollIndicator={false}>
          {isLoading &&
            Array.from({ length: skeletonCount }).map((_, index) => (
              <SkeletonCard key={`card-skeleton__${index}`} />
            ))}
          {!isLoading &&
            data?.cards?.map(card => (
              <Card key={`${card.cardId}__${card.cardNumber}`} {...card} onPress={onCardPress} />
            ))}
          {!isLoading &&
            consentCards?.banks?.map(bank => (
              <WaitConsentCard key={bank as unknown as string} data={bank} />
            ))}
          <RequestConsentCard onPress={onConsetionProcessStart} />
        </ScrollView>
      )}
    </Box>
  )
}

export default CardFeed
