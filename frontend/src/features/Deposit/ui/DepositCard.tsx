import { Box, Text } from '@gluestack-ui/themed'
import React, { memo, useMemo } from 'react'

import { CURRENCY_SYMBOL } from '~/shared/constants/commonConstants'
import { DepositAgreement } from '~/shared/typings/Deposits/DepositsTypes'

const DepositCard: React.FC<DepositAgreement> = ({
  // agreementId,
  // productId,
  name,
  balance,
  currency,
  interestRate,
  accountNumber,
  bank,
}) => {
  const symbol = useMemo(() => CURRENCY_SYMBOL[currency] ?? '', [currency])

  const maskedAccount = useMemo(
    () => (accountNumber ? `•• ${accountNumber.slice(-4)}` : ''),
    [accountNumber],
  )

  const formattedRate = useMemo(() => {
    const numericRate = Number(interestRate)
    if (Number.isNaN(numericRate)) {
      return `${interestRate}%`
    }
    return `${numericRate.toFixed(1)}%`
  }, [interestRate])

  return (
    <Box
      minWidth={190}
      maxWidth={270}
      h={135}
      borderRadius={16}
      overflow='hidden'
      mr='$4'
      p='$3'
      bg='$backgroundDark700'
      borderWidth={1}
      borderColor='$backgroundDark600'
      shadowColor='#000'
      shadowOpacity={0.3}
      shadowRadius={8}
      elevation={6}
    >
      <Box flex={1} flexDirection='column' justifyContent='flex-start' alignItems='flex-start'>
        <Text fontSize='$md' color='$textLight0' fontWeight='$semibold' mb='$1'>
          {name}
        </Text>
        <Text fontSize='$xs' color='$textLight300' fontWeight='$semibold'>
          Ставка {formattedRate}
        </Text>
      </Box>
      <Box flex={1} mt='$2' justifyContent='flex-end' alignItems='flex-start' top='$1.5'>
        <Box
          flexDirection='row'
          justifyContent='center'
          alignItems='baseline'
          borderColor='$backgroundDark500'
          borderBottomWidth={2}
        >
          <Text fontSize='$xl' fontWeight='$bold' color='$textLight0'>
            {balance}
          </Text>
          <Text fontSize='$lg' color='$textLight0' fontWeight='$bold'>
            {` ${symbol}`}
          </Text>
        </Box>
      </Box>
      <Box flex={1} flexDirection='row' justifyContent='space-between' alignItems='flex-end'>
        <Text fontSize='$xs' color='$textLight300'>
          {bank}
        </Text>
        <Text fontSize='$xs' color='$textLight300'>
          {maskedAccount}
        </Text>
      </Box>
    </Box>
  )
}

export default memo(DepositCard)
