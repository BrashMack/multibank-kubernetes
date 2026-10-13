import { Box, Button, Spinner, Text, VStack } from '@gluestack-ui/themed'
import React, { memo, useCallback, useState } from 'react'
import { FlatList, TouchableOpacity } from 'react-native'

import { BANK_CHOOSE } from '~/shared/constants/consentsConstants'

interface BankChoiceScrollProps {
  onBankSelect?: (bankId: string) => void
  onPressClick?: () => void
  errorMsg?: string
  availableBanks?: (typeof BANK_CHOOSE)[0][]
  isLoading?: boolean
  isRefetching?: boolean
}

const getBankColor = (bankLabel: string, selected: boolean) => {
  switch (bankLabel) {
    case 'AwesomeBank':
      return selected ? '$red800' : '$red500'
    case 'SmartBank':
      return selected ? '$green800' : '$green500'
    case 'VirtualBank':
      return selected ? '$blue800' : '$blue500'
    default:
      return selected ? '$primary800' : '$primary500'
  }
}

const BankChoiceScroll: React.FC<BankChoiceScrollProps> = ({
  onBankSelect,
  availableBanks = [],
  isLoading = false,
  isRefetching = false,
  onPressClick,
  errorMsg = '',
}) => {
  const [selectedBank, setSelectedBank] = useState<string | null>(null)

  const handleBankSelect = useCallback(
    (bank: string) => {
      setSelectedBank(bank)
      onBankSelect?.(bank)
    },
    [onBankSelect],
  )

  return (
    <VStack flex={1} justifyContent='center' alignItems='center' w='100%' py='$6' gap='$8'>
      <VStack w='100%' alignItems='center' gap='$12' flexShrink={1}>
        <Text fontWeight='$medium' color='$white' fontSize='$2xl'>
          {availableBanks.length ? 'Выберите банк' : 'Нет доступных банков'}
        </Text>

        {isLoading && !isRefetching ? (
          <Spinner color='$primary500' size={80} />
        ) : (
          <FlatList
            horizontal
            data={availableBanks}
            keyExtractor={item => item.value}
            showsHorizontalScrollIndicator={false}
            contentContainerStyle={{ paddingHorizontal: 8 }}
            renderItem={({ item }) => {
              const isSelected = selectedBank === item.value
              const bgColor = getBankColor(item.label, isSelected)
              return (
                <TouchableOpacity
                  onPress={() => handleBankSelect(item.value)}
                  style={{ marginHorizontal: 8 }}
                >
                  <Box
                    bg={bgColor}
                    borderRadius={12}
                    p='$3'
                    minWidth={100}
                    justifyContent='center'
                    alignItems='center'
                    shadowColor='#000'
                    shadowOpacity={0.1}
                    shadowRadius={4}
                  >
                    <Text color='$white' fontWeight='$medium'>
                      {item.label}
                    </Text>
                  </Box>
                </TouchableOpacity>
              )
            }}
          />
        )}
      </VStack>

      <Box w='100%' alignItems='center'>
        <Button
          onPress={onPressClick}
          w='60%'
          borderRadius={12}
          shadowColor='#000'
          shadowOpacity={0.2}
          shadowRadius={6}
        >
          <Text color='$white' fontWeight='$bold' fontSize='$md'>
            Продолжить
          </Text>
        </Button>
        {!!errorMsg && (
          <Text color='$red500' fontSize='$md' mt='$2'>
            {errorMsg}
          </Text>
        )}
      </Box>
    </VStack>
  )
}

export default memo(BankChoiceScroll)
