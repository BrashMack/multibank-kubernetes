import { Box, HStack, Text, VStack } from '@gluestack-ui/themed'
import React, { memo, useMemo } from 'react'

import { CURRENCY_SYMBOL } from '~/shared/constants/commonConstants'
import { Transaction } from '~/shared/typings/Transactions/TransactionsTypes'

//eslint-disable-next-line @typescript-eslint/no-empty-object-type
export interface TransactionItemProps extends Transaction {}

const TransactionItem: React.FC<TransactionItemProps> = ({
  amount,
  currency,
  type,
  createdAt,
  transactionInformation,
}) => {
  const isIncome = useMemo(() => type === 'credit', [type])

  const symbol = useMemo(() => currency && (CURRENCY_SYMBOL[currency] ?? ''), [currency])
  const formattedAmount = useMemo(
    () => `${isIncome ? '+' : '-'}${amount} ${symbol}`,
    [isIncome, amount, symbol],
  )

  const formattedDate = useMemo(() => new Date(createdAt).toLocaleString(), [createdAt])

  const amountColor = useMemo(() => (isIncome ? '$success600' : '$error600'), [isIncome])

  return (
    <Box
      bg='$backgroundDark800'
      borderRadius={12}
      p='$3'
      m='$2'
      shadowColor='#000'
      shadowOpacity={0.25}
      shadowRadius={6}
      elevation={3}
      borderWidth={1}
      borderColor='$backgroundDark700'
    >
      <HStack justifyContent='space-between' alignItems='center'>
        <VStack space='xs' flex={1} mr='$2'>
          <Text fontSize='$sm' fontWeight='$semibold' color='$textLight0'>
            {transactionInformation}
          </Text>

          <Text fontSize='$xs' color='$textLight400'>
            {formattedDate}
          </Text>
        </VStack>

        <Text fontSize='$md' fontWeight='$bold' color={amountColor}>
          {formattedAmount}
        </Text>
      </HStack>
    </Box>
  )
}

export default memo(TransactionItem)
