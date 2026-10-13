import { Box, Button, ScrollView, Text } from '@gluestack-ui/themed'
import { useNavigation } from '@react-navigation/native'
import { StackNavigationProp } from '@react-navigation/stack'
import { useQueryClient } from '@tanstack/react-query'
import * as Clipboard from 'expo-clipboard'
import React, { memo, useCallback, useEffect, useMemo, useState } from 'react'
import { Pressable } from 'react-native'

import useService from '~/entities/useService'
import Card from '~/features/Card/ui/Card'
import TransferAction from '~/features/Transfer/ui/TransferAction'
import { ExtendedAccountData } from '~/shared/typings/Accounts/AccountsTypes'
import { RootStackParamList } from '~/shared/typings/App/AppTypes'

const ExtendedAccountView: React.FC<ExtendedAccountData> = initialAccountData => {
  const navigation = useNavigation<StackNavigationProp<RootStackParamList>>()
  const { getAccountById } = useService()
  const queryClient = useQueryClient()
  const [accountDetails, setAccountDetails] = useState(initialAccountData)
  const [isAccountNumberVisible, setAccountNumberVisible] = useState(false)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const { number, balance, currency, bank, type, numberCards, accountId } = accountDetails
  const formattedNumber = useMemo(
    () => (number ? `•••• •••• •••• ${number.slice(-4)}` : 'Загрузка данных счёта...'),
    [number],
  )
  const displayedAccountNumber = useMemo(
    () => (isAccountNumberVisible && number ? number : formattedNumber),
    [formattedNumber, isAccountNumberVisible, number],
  )

  // TODO: Сделать нормальное отображение загрузки данных, пока они не загрузились

  const onCardPress = useCallback(
    (cardId: string) => {
      navigation.navigate('card', { cardId })
    },
    [navigation],
  )

  const onNumberPress = useCallback(() => {
    setAccountNumberVisible(prevState => {
      const nextState = !prevState

      if (nextState && number) {
        Clipboard.setStringAsync(number)
      }

      return nextState
    })
  }, [number])

  useEffect(() => {
    setAccountDetails(initialAccountData)
  }, [initialAccountData])

  const refreshAccountDetails = useCallback(async () => {
    if (!accountId) return
    const controller = new AbortController()
    try {
      setIsRefreshing(true)
      const updated = await getAccountById(accountId, controller.signal)
      setAccountDetails(updated)
    } finally {
      setIsRefreshing(false)
    }
  }, [accountId, getAccountById])

  const handleTransferSuccess = useCallback(async () => {
    await refreshAccountDetails()
    if (accountId) {
      await queryClient.invalidateQueries({ queryKey: ['transactions', accountId] })
    }
  }, [accountId, queryClient, refreshAccountDetails])

  return (
    <Box bg='$backgroundDark950'>
      <Box w='100%' alignItems='center' justifyContent='flex-start' px='$4'>
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
          <Box
            position='relative'
            h={200}
            justifyContent='center'
            alignItems='center'
            bg='$backgroundDark900'
          >
            <Box position='absolute' top='$4' left='$5' right='$5'>
              <Text fontSize='$md' fontWeight='$semibold' color='$white'>
                {bank}
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
              <Pressable accessibilityRole='button' onPress={onNumberPress}>
                <Text
                  fontSize='$md'
                  fontWeight='$medium'
                  color='$white'
                  px='$3'
                  bg='rgba(61, 61, 61, 0.87)'
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
          <TransferAction sourceAccountNumber={number} onSuccess={handleTransferSuccess} />
        </Box>
      </Box>
      <Box w='$full' px='$4' pb='$4'>
        {numberCards?.length > 0 && (
          <>
            <Text color='$white' fontSize='$2xl' px='$4' fontWeight='$bold' mb='$4'>
              Карты привязанные к счёту
            </Text>
            <ScrollView horizontal showsHorizontalScrollIndicator={false}>
              {numberCards.map(card => (
                <Card
                  key={card.cardId}
                  cardId={card.cardId}
                  cardNumber={card.number}
                  // cardName={card.cardName}
                  backgroundImage={card.backgroundImage}
                  type={card.type}
                  onPress={onCardPress}
                />
              ))}
            </ScrollView>
          </>
        )}
      </Box>
    </Box>
  )
}

export default memo(ExtendedAccountView)
