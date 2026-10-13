import { Box, Button, Text } from '@gluestack-ui/themed'
import React, { memo, useEffect, useMemo, useState } from 'react'
import { FlatList, TouchableOpacity } from 'react-native'

import { CURRENCY_SYMBOL } from '~/shared/constants/commonConstants'
import { AccountData } from '~/shared/typings/Accounts/AccountsTypes'

export interface QRTransaction {
  primaryId: string
  secondaryId: string
  deltaNum: number
  delta: string
  accountIdTo: string
  amount: string
}

export interface QRPaymentCardAutocompleteProps {
  accounts: AccountData[]
  accountIdTo: string
  amount: string
  onContinueClick?: (transaction: QRTransaction) => void
  currency: string
  errorMsg: string
  isLoading?: boolean
}

const parseBalance = (balance: string) => parseFloat(balance.replace(/[,\s]/g, '')) || 0

const QRPaymentCardAutocomplete: React.FC<QRPaymentCardAutocompleteProps> = ({
  accounts,
  accountIdTo,
  amount,
  onContinueClick,
  currency = 'RUB',
  errorMsg,
  isLoading = false,
}) => {
  const numericAmount = useMemo(() => parseBalance(amount), [amount])

  const eligibleAccounts = useMemo(
    () =>
      accounts.filter(
        acc => acc.type !== 'loan' && acc.currency === currency && acc.number !== accountIdTo,
      ),
    [accounts, currency, accountIdTo],
  )

  // DEFAULT primary = richest
  const defaultPrimary = useMemo(
    () =>
      [...eligibleAccounts].sort((a, b) => parseBalance(b.balance) - parseBalance(a.balance))[0] ||
      null,
    [eligibleAccounts],
  )

  const [selectedPrimary, setSelectedPrimary] = useState<string>(defaultPrimary?.number || '')
  const [selectedSecondary, setSelectedSecondary] = useState<string>('')

  const [drawerType, setDrawerType] = useState<'primary' | 'secondary' | null>(null)

  // recalc when accounts changed
  useEffect(() => {
    ;(() => {
      setSelectedPrimary(defaultPrimary?.number || '')
      setSelectedSecondary('')
    })()
  }, [defaultPrimary])

  const primaryAccount = useMemo(
    () => eligibleAccounts.find(a => a.number === selectedPrimary) || null,
    [selectedPrimary, eligibleAccounts],
  )

  const primaryBalance = primaryAccount ? parseBalance(primaryAccount.balance) : 0

  const isEnough = primaryBalance >= numericAmount

  const secondaryCandidates = useMemo(
    () => eligibleAccounts.filter(a => a.number !== selectedPrimary),
    [eligibleAccounts, selectedPrimary],
  )

  const secondaryAccount = useMemo(
    () => secondaryCandidates.find(a => a.number === selectedSecondary) || null,
    [selectedSecondary, secondaryCandidates],
  )

  const totalBalance =
    primaryBalance + (secondaryAccount ? parseBalance(secondaryAccount.balance) : 0)

  const isEnoughWithSecondary = !isEnough && totalBalance >= numericAmount

  const deltaNum = Number(primaryAccount?.balance) - Number(amount)
  const [integerPart, decimals = ''] = String(Math.abs(deltaNum)).split('.')
  const delta = `${integerPart}.${decimals.padEnd(2, '0')}`

  const getAccountLabel = (num?: string) => {
    if (!num) return 'Выберите счёт'
    const acc = accounts.find(a => a.number === num)
    if (!acc) return num
    return `${acc.number} — ${acc.balance} ${CURRENCY_SYMBOL[acc.currency] ?? acc.currency}`
  }

  const handleSelectAccount = (acc: AccountData) => {
    if (drawerType === 'primary') {
      setSelectedPrimary(acc.number)
      setSelectedSecondary('') // reset secondary on change of primary
    }

    if (drawerType === 'secondary') {
      setSelectedSecondary(acc.number)
    }

    setDrawerType(null)
  }

  const renderAccountItem = ({ item }: { item: AccountData }) => {
    const isSelected =
      (drawerType === 'primary' && item.number === selectedPrimary) ||
      (drawerType === 'secondary' && item.number === selectedSecondary)

    return (
      <TouchableOpacity
        onPress={() => handleSelectAccount(item)}
        style={{ paddingVertical: 12, paddingHorizontal: 16 }}
        activeOpacity={0.7}
      >
        <Box flexDirection='row' justifyContent='space-between' alignItems='center'>
          <Box>
            <Text fontSize={15} fontWeight='$medium' color='$textLight100'>
              {item.number}
            </Text>
            <Text fontSize={13} color='$textLight400' mt={4}>
              {item.balance} {item.currency}
            </Text>
          </Box>

          {isSelected && (
            <Text fontSize={13} color='$blue600'>
              Выбрано
            </Text>
          )}
        </Box>
      </TouchableOpacity>
    )
  }

  return (
    <Box flex={1} gap={20} px='$4' pt='$4' backgroundColor='$backgroundDark950'>
      <Box bg='$backgroundDark900' borderRadius={12} p={14} mb={8}>
        <Text fontSize={16} color='$textLight400' mb={4}>
          Счёт получателя:
        </Text>
        <Text fontSize={17} fontWeight='$semibold' color='$textLight100' mb={8}>
          {accountIdTo || 'Не указан'}
        </Text>

        <Text fontSize={16} color='$textLight400' mb={4}>
          Сумма перевода:
        </Text>
        <Text fontSize={17} fontWeight='$semibold' color='$textLight100'>
          {amount} {currency}
        </Text>
      </Box>

      <Box>
        <Text fontSize={16} mb={6} color='$textLight200'>
          Выберите основной счет для для оплаты
        </Text>

        <TouchableOpacity onPress={() => setDrawerType('primary')}>
          <Box
            bg='$backgroundDark800'
            borderColor='$backgroundDark700'
            borderWidth={1}
            borderRadius={12}
            px={14}
            py={12}
          >
            <Text color='$textLight100'>{getAccountLabel(selectedPrimary)}</Text>
          </Box>
        </TouchableOpacity>
      </Box>

      {!isEnough && (
        <Box>
          <Text fontSize={16} mb={6} color='$textLight200'>
            {`Выберите счет для перевода недостающих ${delta || ''} ${CURRENCY_SYMBOL[currency] ?? currency}`}
          </Text>

          <TouchableOpacity onPress={() => setDrawerType('secondary')}>
            <Box
              bg='$backgroundDark800'
              borderColor='$backgroundDark700'
              borderWidth={1}
              borderRadius={12}
              px={14}
              py={12}
            >
              <Text color='$textLight100'>{getAccountLabel(selectedSecondary)}</Text>
            </Box>
          </TouchableOpacity>
        </Box>
      )}

      {/* ERRORS */}
      {!isEnough && !isEnoughWithSecondary && (
        <Text color='$red400' fontSize={14}>
          Недостаточно средств для совершения платежа.
        </Text>
      )}

      {(isEnough || isEnoughWithSecondary) && !!errorMsg && (
        <Text color='$red400' fontSize={14}>
          {errorMsg}
        </Text>
      )}

      {/* BUTTON */}
      <Button
        size='md'
        isDisabled={(!isEnough && !isEnoughWithSecondary) || isLoading}
        mt={10}
        bg={!isEnough && !isEnoughWithSecondary ? '$backgroundDark700' : '$blue600'}
        borderRadius={16}
        justifyContent='center'
        alignItems='center'
        onPress={() => {
          onContinueClick?.({
            secondaryId: selectedSecondary,
            primaryId: selectedPrimary,
            deltaNum,
            delta,
            accountIdTo,
            amount,
          })
        }}
      >
        <Text color='$white' fontSize='$md' fontWeight='$semibold'>
          Продолжить
        </Text>
      </Button>

      {/* DRAWER */}
      {drawerType && (
        <Box position='absolute' top={0} left={0} right={0} bottom={0} justifyContent='flex-end'>
          <TouchableOpacity
            style={{ flex: 1, backgroundColor: 'rgba(0,0,0,0.45)' }}
            onPress={() => setDrawerType(null)}
          />

          <Box
            bg='$backgroundDark950'
            borderTopLeftRadius={20}
            borderTopRightRadius={20}
            maxHeight='60%'
            p={12}
          >
            <Box flexDirection='row' justifyContent='space-between' alignItems='center' mb={8}>
              <Text fontSize={16} fontWeight='$semibold' color='$textLight100'>
                {drawerType === 'primary'
                  ? 'Выберите основной счёт'
                  : 'Выберите счёт для доп. списания'}
              </Text>

              <TouchableOpacity onPress={() => setDrawerType(null)} activeOpacity={0.7}>
                <Text color='$textLight400'>Отмена</Text>
              </TouchableOpacity>
            </Box>

            <FlatList
              data={drawerType === 'primary' ? eligibleAccounts : secondaryCandidates}
              keyExtractor={item => `${item.accountId}__${item.number}`}
              renderItem={renderAccountItem}
              showsVerticalScrollIndicator={false}
              ItemSeparatorComponent={() => <Box height={1} bg='$backgroundDark850' my={6} />}
            />
          </Box>
        </Box>
      )}
    </Box>
  )
}

export default memo(QRPaymentCardAutocomplete)
