import { Box, Button, Image, Text } from '@gluestack-ui/themed'
import { useNavigation } from '@react-navigation/native'
import { StackNavigationProp } from '@react-navigation/stack'
import * as Clipboard from 'expo-clipboard'
import React, { memo, useCallback, useEffect, useMemo, useState } from 'react'
import { Pressable } from 'react-native'

import useService from '~/entities/useService'
import TransferAction from '~/features/Transfer/ui/TransferAction'
import { RootStackParamList } from '~/shared/typings/App/AppTypes'
import { ExtendedCardData } from '~/shared/typings/Card/CardTypes'

const ExtendedCardView: React.FC<ExtendedCardData> = initialCardData => {
  const navigation = useNavigation<StackNavigationProp<RootStackParamList>>()
  const { getCardById } = useService()
  const [cardDetails, setCardDetails] = useState(initialCardData)
  const [isAccountNumberVisible, setAccountNumberVisible] = useState(false)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const { number, cardNumber, balance, currency, bank, backgroundImage, type, cardId } = cardDetails
  const formattedCardNumber = useMemo(
    () => (cardNumber ? `•••• ${cardNumber.slice(-4)}` : '•••• ••••'),
    [cardNumber],
  )
  const formattedAccountNumber = useMemo(
    () => (number ? `•••• •••• •••• ${number.slice(-4)}` : 'Загрузка данных карты...'),
    [number],
  )
  const displayedAccountNumber = useMemo(
    () => (isAccountNumberVisible && number ? number : formattedAccountNumber),
    [formattedAccountNumber, isAccountNumberVisible, number],
  )
  const onAccountNumberPress = useCallback(() => {
    setAccountNumberVisible(prevState => {
      const nextState = !prevState

      if (nextState && number) {
        Clipboard.setStringAsync(number)
      }

      return nextState
    })
  }, [number])
  useEffect(() => {
    setCardDetails(initialCardData)
  }, [initialCardData])

  const refreshCardDetails = useCallback(async () => {
    if (!cardId) return
    const controller = new AbortController()
    try {
      setIsRefreshing(true)
      const updated = await getCardById(cardId, controller.signal)
      setCardDetails(updated as ExtendedCardData)
    } finally {
      setIsRefreshing(false)
    }
  }, [cardId, getCardById])

  // TODO: отображение лого платёжной системы(по ключевому слову в cardName) на карте(над фоном карты)
  // на карте или рядом с ней отображать dailyLimit и monthlyLimit(может на обороте?)
  // при нажатии на карту, она переворачивается и там показаны полные данные cardNumberFull, holderName, expiryMonth, expiryYear

  return (
    <Box
      flex={1}
      w='100%'
      alignItems='center'
      justifyContent='flex-start'
      px='$4'
      bg='$backgroundDark950'
    >
      <Box w='100%' pb='$4'>
        <Button
          mt='$6'
          px='$4'
          maxWidth='$32'
          borderRadius={16}
          bg='$blue600'
          onPress={() => navigation.navigate('main' as never)}
          shadowColor='#000'
          shadowOpacity={0.3}
          shadowRadius={6}
          elevation={5}
          justifyContent='center'
          alignItems='center'
        >
          <Text color='$white' fontSize='$md' fontWeight='$semibold'>
            На главную
          </Text>
        </Button>
      </Box>
      <Box
        w='90%'
        maxWidth={360}
        borderRadius={20}
        overflow='hidden'
        shadowColor='#000'
        shadowOpacity={0.25}
        shadowRadius={8}
        elevation={6}
        bg='$backgroundDark950'
      >
        <Box position='relative' h={200} justifyContent='center' alignItems='center'>
          {backgroundImage && (
            <Image
              source={{ uri: backgroundImage }}
              alt='Card background'
              w='100%'
              h='100%'
              position='absolute'
            />
          )}
          <Box position='absolute' w='100%' h='100%' bg='rgba(0,0,0,0.8)' />

          <Box
            position='absolute'
            top='$4'
            left='$5'
            right='$5'
            flexDirection='row'
            justifyContent='space-between'
          >
            <Text fontSize='$md' fontWeight='$semibold' color='$white'>
              {bank}
            </Text>
            <Text fontSize='$md' fontWeight='$medium' color='$white'>
              {formattedCardNumber}
            </Text>
          </Box>

          <Box alignItems='center'>
            <Text fontSize='$2xl' fontWeight='$bold' color='$white'>
              {balance} {currency}
            </Text>
            {isRefreshing && (
              <Text fontSize='$xs' color='$textLight400' mt='$1'>
                Обновляем данные…
              </Text>
            )}
          </Box>

          <Box
            position='absolute'
            bottom='$4'
            left='$5'
            right='$5'
            flexDirection='row'
            justifyContent='space-between'
          >
            <Pressable accessibilityRole='button' onPress={onAccountNumberPress}>
              <Text
                fontSize='$md'
                fontWeight='$medium'
                color='$white'
                bg='rgba(61, 61, 61, 0.87)'
                px='$3'
                borderRadius={8}
                selectable
              >
                {displayedAccountNumber}
              </Text>
            </Pressable>
            <Text fontSize='$sm' fontWeight='$medium' color='$white' opacity={0.8}>
              {type}
            </Text>
          </Box>
        </Box>
      </Box>
      <Box w='90%' maxWidth={360} pb='$4'>
        <TransferAction sourceAccountNumber={number} onSuccess={refreshCardDetails} />
      </Box>
    </Box>
  )
}

export default memo(ExtendedCardView)
