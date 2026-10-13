import { Box, Button, Text } from '@gluestack-ui/themed'
import { useNavigation } from '@react-navigation/native'
import { StackNavigationProp } from '@react-navigation/stack'
import { useMutation, useQuery } from '@tanstack/react-query'
import React, { useCallback, useMemo, useState } from 'react'

import useService from '~/entities/useService'
import BankChoiceScroll from '~/features/Consents/BankChoiceScroll/ui/BankChoiceScroll'
import { BANK_CHOOSE } from '~/shared/constants/consentsConstants'
import { RootStackParamList } from '~/shared/typings/App/AppTypes'

const BankChoosing: React.FC = () => {
  const [selectedBank, setSelectedBank] = useState<string | null>(null)
  const navigation = useNavigation<StackNavigationProp<RootStackParamList>>()
  const [errorMsg, setErrorMsg] = useState<string>('')
  const { requestConsent, getActiveConsents } = useService()

  const { data, isLoading, isRefetching } = useQuery({
    queryKey: ['active-consents'],
    queryFn: ({ signal }) => getActiveConsents(signal),
  })

  const availableBanks = useMemo(() => {
    return data?.awaitingBanks?.length || data?.confirmedBanks?.length
      ? BANK_CHOOSE.filter(
          bank =>
            !(data.awaitingBanks.includes(bank.value) || data.confirmedBanks.includes(bank.value)),
        )
      : BANK_CHOOSE
  }, [data])

  const mutation = useMutation({
    mutationFn: () => requestConsent(selectedBank as string, 'read'),
    onSuccess: () => {
      navigation.navigate('main')
    },
    onError: () => {
      setErrorMsg('Произошла ошибка')
      setTimeout(() => setErrorMsg(''), 2000)
    },
  })

  const navigateToMain = useCallback(() => {
    navigation.navigate('main')
  }, [navigation])

  return (
    <Box flex={1} justifyContent='center' alignItems='center' bg='$backgroundDark950' p='$4'>
      <Box p='$2' alignSelf='flex-start'>
        <Button
          borderRadius={12}
          shadowColor='#000'
          shadowOpacity={0.3}
          shadowRadius={6}
          elevation={6}
          onPress={navigateToMain}
        >
          <Text color='$white'>На главную</Text>
        </Button>
      </Box>
      <BankChoiceScroll
        onBankSelect={bank => setSelectedBank(bank)}
        availableBanks={availableBanks}
        isLoading={isLoading}
        isRefetching={isRefetching}
        onPressClick={() => {
          if (availableBanks?.length) mutation.mutateAsync()
          else navigation.navigate('main')
        }}
        errorMsg={errorMsg}
      />
    </Box>
  )
}

export default BankChoosing
