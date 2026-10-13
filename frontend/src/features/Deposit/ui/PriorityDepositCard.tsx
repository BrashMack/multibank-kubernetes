import { Box, Text } from '@gluestack-ui/themed'
import React, { memo, useCallback, useMemo } from 'react'
import { Pressable } from 'react-native'

import { CURRENCY_SYMBOL } from '~/shared/constants/commonConstants'
import { DepositProduct } from '~/shared/typings/Deposits/DepositsTypes'

interface PriorityDepositCardProps extends DepositProduct {
  onPress?: (product: DepositProduct) => void
}

const PriorityDepositCard: React.FC<PriorityDepositCardProps> = ({
  productId,
  productName,
  interestRate,
  minAmount,
  maxAmount,
  currency,
  termMonths,
  bank,
  onPress,
}) => {
  const min = useMemo(() => Number(minAmount), [minAmount])
  const max = useMemo(() => Number(maxAmount), [maxAmount])
  const expectedDeposit = useMemo(() => {
    if (Number.isNaN(min) || Number.isNaN(max)) {
      return null
    }
    return (min + (max ? max : min * 10)) / 2
  }, [min, max])

  const formattedRate = useMemo(() => {
    if (!interestRate) return ''
    const numericRate = Number(interestRate)
    if (Number.isNaN(numericRate)) return `${interestRate}%`
    return `${numericRate.toFixed(1)}%`
  }, [interestRate])

  const symbol = useMemo(() => (currency ? (CURRENCY_SYMBOL[currency] ?? '') : ''), [currency])

  const incomePreview = useMemo(() => {
    if (!interestRate || !expectedDeposit) return null
    const numericRate = Number(interestRate)
    if (Number.isNaN(numericRate)) return null
    const income = (numericRate / 100 / 12) * termMonths * expectedDeposit
    return income.toFixed(0)
  }, [expectedDeposit, interestRate, termMonths])

  const handlePress = useCallback(() => {
    onPress?.({
      productId,
      productName,
      interestRate,
      minAmount,
      maxAmount,
      currency,
      termMonths,
      bank,
    })
  }, [
    bank,
    currency,
    interestRate,
    maxAmount,
    minAmount,
    onPress,
    productId,
    productName,
    termMonths,
  ])

  return (
    <Box maxWidth={550}>
      <Pressable onPress={handlePress}>
        <Box
          w='100%'
          h={140}
          borderRadius={18}
          borderWidth={1}
          borderColor='$backgroundDark600'
          bg='$backgroundDark900'
          p='$3.5'
          mb='$4'
          shadowColor='#000'
          shadowOpacity={0.35}
          shadowRadius={10}
          elevation={8}
        >
          <Box flex={1} flexDirection='row' justifyContent='space-between' alignItems='flex-start'>
            <Box
              bg='$green500'
              pb='$0.5'
              mr={20}
              borderRadius={8}
              shadowColor='#000'
              shadowOpacity={0.35}
              shadowRadius={8}
              elevation={4}
              flex={1}
              maxWidth={176}
              minWidth={109}
            >
              <Text color='$backgroundDark950' fontWeight='$semibold' fontSize='$sm' px='$2'>
                Выгодное предложение
              </Text>
            </Box>
            <Box
              //flex={1}
              flexDirection='column'
              justifyContent='flex-start'
              alignItems='flex-end'
              bottom='$0.5'
            >
              <Text color='$textLight0' fontWeight='$bold' fontSize='$md' mb='$1'>
                {productName}
              </Text>
              <Text color='$textLight300' fontWeight='$medium' fontSize='$sm'>
                {bank}
              </Text>
            </Box>
          </Box>
          <Box flex={2} flexDirection='column' justifyContent='flex-end' alignItems='flex-start'>
            <Text color='$green500' fontWeight='$semibold' mb='$1'>
              Ставка {formattedRate}
            </Text>
            <Text color='$textLight300' fontWeight='$light' fontSize='$sm' mb='$0.5'>
              От {min.toFixed(0)} {symbol} {maxAmount && `до ${max.toFixed(0)} ${symbol} `}на{' '}
              {termMonths} мес.
            </Text>
            {incomePreview && expectedDeposit && (
              <Text color='$textLight0' fontWeight='$medium' fontSize='$md'>
                {incomePreview} {symbol} — доход с {expectedDeposit.toFixed(0)} {symbol} за{' '}
                {termMonths} мес.
              </Text>
            )}
          </Box>
        </Box>
      </Pressable>
    </Box>
  )
}

export default memo(PriorityDepositCard)
