import { useRoute } from '@react-navigation/native'
import { useQuery } from '@tanstack/react-query'
import React from 'react'

import useService from '~/entities/useService'
import ExtendedCardView from '~/features/ExtendedCard/ExtendedCardView/ui/ExtendedCardView'
import { CardPageRouteProp } from '~/shared/typings/App/AppTypes'
import { ExtendedCardData } from '~/shared/typings/Card/CardTypes'

const CardView: React.FC = () => {
  const { getCardById } = useService()
  //   const navigation = useNavigation<StackNavigationProp<RootStackParamList>>()
  const route = useRoute<CardPageRouteProp>()
  const { cardId } = route.params

  //   useEffect(() => {
  //     if (!cardId) navigation.navigate('main')
  //   }, [cardId, navigation])

  const { data } = useQuery({
    queryKey: ['card', cardId],
    queryFn: ({ signal }) => getCardById(cardId, signal),
    enabled: !!cardId,
  })

  return <ExtendedCardView {...(data as ExtendedCardData)} />
}

export default CardView
