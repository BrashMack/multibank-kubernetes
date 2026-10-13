import { Box, HStack, Text, VStack } from '@gluestack-ui/themed'
import React, { memo, useCallback, useMemo } from 'react'
import { Pressable } from 'react-native'

import { CURRENCY_SYMBOL } from '~/shared/constants/commonConstants'

interface AccountCardProps {
  accountId: string
  number: string
  balance: string
  currency: 'RUB' | 'USD' | 'EUR' | 'GBP'
  bank: string
  type: string
  onPress?: (id: string) => void
}

const AccountCard: React.FC<AccountCardProps> = ({
  accountId,
  number,
  balance,
  currency,
  bank,
  type,
  onPress,
}) => {
  const maskedNumber = useMemo(() => (number ? `•• ${number.slice(-4)}` : ''), [number])
  const symbol = useMemo(() => CURRENCY_SYMBOL[currency] ?? '', [currency])
  const handlePress = useCallback(() => onPress?.(accountId), [onPress, accountId])

  return (
    <Pressable accessibilityRole='button' onPress={handlePress}>
      <Box
        w={180}
        h={110}
        borderRadius={12}
        overflow='hidden'
        m='$2'
        bg='$backgroundDark700'
        p='$3'
        shadowColor='#000'
        shadowOpacity={0.3}
        shadowRadius={6}
        elevation={4}
        borderWidth={1}
        borderColor='$backgroundDark600'
      >
        <VStack flex={1} justifyContent='space-between'>
          <VStack space='xs'>
            <Text fontSize='$sm' fontWeight='$semibold' color='$textLight0'>
              {bank}
            </Text>
            <Text fontSize='$xs' color='$textLight300'>
              {type}
            </Text>
          </VStack>

          <VStack space='xs'>
            <Text fontSize='$xs' color='$textLight300'>
              {maskedNumber}
            </Text>

            <HStack alignItems='flex-end' space='xs'>
              <Text fontSize='$lg' fontWeight='$bold' color='$textLight0'>
                {balance}
              </Text>
              <Text fontSize='$lg' color='$textLight0' fontWeight='$bold'>
                {symbol}
              </Text>
            </HStack>
          </VStack>
        </VStack>
      </Box>
    </Pressable>
  )
}

export default memo(AccountCard)
