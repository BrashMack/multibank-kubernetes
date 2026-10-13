/* eslint-disable react-hooks/set-state-in-effect */
import {
  Box,
  Button,
  FormControl,
  FormControlError,
  FormControlErrorText,
  FormControlLabel,
  FormControlLabelText,
  Input,
  InputField,
  Modal,
  ModalBackdrop,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  Spinner,
  Text,
  VStack,
} from '@gluestack-ui/themed'
import { useMutation, useQuery } from '@tanstack/react-query'
import axios from 'axios'
import React, { memo, useCallback, useEffect, useMemo, useState } from 'react'
import { Alert, Pressable, ScrollView } from 'react-native'

import useService from '~/entities/useService'
import { CURRENCY_SYMBOL } from '~/shared/constants/commonConstants'
import { AccountData } from '~/shared/typings/Accounts/AccountsTypes'
import {
  DepositAgreement,
  DepositProduct,
  MakeDepositRequest,
} from '~/shared/typings/Deposits/DepositsTypes'

interface DepositCreateModalProps {
  product: DepositProduct | null
  isOpen: boolean
  isPriority: boolean
  onClose: () => void
  onSuccess?: (deposit: DepositAgreement) => void
}

const parseAmount = (value?: string) => {
  if (!value) return 0
  const normalized = value.replace(/\s/g, '').replace(',', '.')
  const numeric = Number(normalized)
  return Number.isNaN(numeric) ? 0 : numeric
}

const DepositCreateModal: React.FC<DepositCreateModalProps> = ({
  product,
  isOpen,
  isPriority,
  onClose,
  onSuccess,
}) => {
  const { getAccounts, makeDeposit } = useService()
  const [amountValue, setAmountValue] = useState('')
  const [selectedAccountId, setSelectedAccountId] = useState('')
  const [isAccountListOpen, setAccountListOpen] = useState(false)
  const [amountError, setAmountError] = useState('')
  const [accountError, setAccountError] = useState('')
  const [generalError, setGeneralError] = useState('')

  const {
    data: accountsData,
    isLoading: isAccountsLoading,
    isRefetching: isAccountsRefetching,
  } = useQuery({
    queryKey: ['accounts'],
    queryFn: ({ signal }) => getAccounts(signal),
    enabled: isOpen,
  })

  const minAmount = parseAmount(product?.minAmount)
  const maxAmount = parseAmount(product?.maxAmount)
  const productSymbol = product ? (CURRENCY_SYMBOL[product.currency] ?? product.currency) : ''

  const eligibleAccounts = useMemo(() => {
    const accounts = accountsData?.accounts ?? []
    return accounts.filter(
      account =>
        parseAmount(account.balance) >= minAmount &&
        (!product || account.currency === product?.currency) &&
        (!product || account.bank === product?.bank),
    )
  }, [accountsData, minAmount, product])

  useEffect(() => {
    if (isOpen && product) {
      setAmountValue(product.minAmount ?? '')
      setSelectedAccountId('')
      setAccountError('')
      setAmountError('')
      setGeneralError('')
      setAccountListOpen(false)
    }
  }, [isOpen, product])

  useEffect(() => {
    if (!selectedAccountId && eligibleAccounts.length > 0) {
      setSelectedAccountId(eligibleAccounts[0].accountId)
    }
  }, [eligibleAccounts, selectedAccountId])

  const selectedAccount = eligibleAccounts.find(acc => acc.accountId === selectedAccountId)

  const formattedRate = useMemo(() => {
    if (!product?.interestRate) return ''
    const numeric = Number(product.interestRate)
    if (Number.isNaN(numeric)) {
      return `${product.interestRate}%`
    }
    return `${numeric.toFixed(1)}%`
  }, [product])

  const toggleAccountList = useCallback(() => {
    if (!eligibleAccounts.length) return
    setAccountListOpen(prev => !prev)
  }, [eligibleAccounts.length])

  const handleAccountSelect = useCallback((accountId: string) => {
    setSelectedAccountId(accountId)
    setAccountError('')
    setAccountListOpen(false)
  }, [])

  const closeModal = useCallback(() => {
    setAccountListOpen(false)
    setAccountError('')
    setAmountError('')
    setGeneralError('')
    onClose()
  }, [onClose])

  const { mutate: submitDeposit, isPending: isSubmitting } = useMutation({
    mutationFn: (payload: MakeDepositRequest) => makeDeposit(payload),
    onSuccess: deposit => {
      Alert.alert('Вклад открыт', 'Заявка успешно создана.')
      onSuccess?.(deposit)
      closeModal()
    },
    onError: error => {
      if (axios.isAxiosError(error)) {
        const backendMessage =
          typeof error.response?.data === 'string'
            ? error.response?.data
            : error.response?.data?.message
        if (backendMessage) {
          setGeneralError(backendMessage)
          return
        }
      }

      setGeneralError('Не удалось открыть вклад. Попробуйте позже.')
    },
  })

  const handleSubmit = useCallback(() => {
    if (!product) return
    setAmountError('')
    setAccountError('')
    setGeneralError('')

    if (!selectedAccount) {
      setAccountError('Нет подходящих счетов для пополнения')
      return
    }

    const normalizedAmount = amountValue.replace(',', '.').trim()
    const amountFormatRegex = /^\d+(\.\d{0,2})?$/
    const parsedAmount = Number(normalizedAmount)

    if (
      !normalizedAmount ||
      Number.isNaN(parsedAmount) ||
      !amountFormatRegex.test(normalizedAmount)
    ) {
      setAmountError('Введите корректную сумму пополнения')
      return
    }

    if (parsedAmount < minAmount) {
      setAmountError(`Минимальная сумма — ${product.minAmount} ${productSymbol}`)
      return
    }

    if (product.maxAmount && parsedAmount > maxAmount) {
      setAmountError(`Максимальная сумма — ${product.maxAmount} ${productSymbol}`)
      return
    }

    if (parsedAmount > parseAmount(selectedAccount.balance)) {
      setAmountError('Недостаточно средств на счёте')
      return
    }

    const [integerPart, decimals = ''] = normalizedAmount.split('.')
    const formattedAmount = `${integerPart}.${decimals.padEnd(2, '0')}`

    submitDeposit({
      productId: product.productId,
      bank: product.bank,
      sourceAccountId: selectedAccount.number,
      amount: formattedAmount,
      termMonths: product.termMonths,
    })
  }, [amountValue, maxAmount, minAmount, product, selectedAccount, productSymbol, submitDeposit])

  const renderAccountItem = (account: AccountData) => (
    <Pressable
      key={account.accountId}
      onPress={() => handleAccountSelect(account.accountId)}
      style={{
        paddingVertical: 12,
        paddingHorizontal: 16,
        borderBottomWidth: 1,
        borderBottomColor: '#2e2e2e',
        backgroundColor: account.accountId === selectedAccountId ? '#1f1f1f' : 'transparent',
      }}
    >
      <Text fontWeight='$semibold' color='$textLight100'>
        {account.bank} •• {account.number.slice(-4)}
      </Text>
      <Text fontSize='$sm' color='$textLight300'>
        {account.balance} {CURRENCY_SYMBOL[account.currency] ?? account.currency}
      </Text>
    </Pressable>
  )

  const isAccountsLoadingState = isAccountsLoading || isAccountsRefetching
  const disabledSubmit =
    !product ||
    !selectedAccountId ||
    !eligibleAccounts.length ||
    isSubmitting ||
    isAccountsLoadingState

  return (
    <Modal isOpen={isOpen && Boolean(product)} onClose={closeModal}>
      <ModalBackdrop />
      <ModalContent bg='$backgroundDark900'>
        <ModalHeader justifyContent='center'>
          <Text fontSize='$lg' fontWeight='$bold' color='$white'>
            Открыть вклад
          </Text>
          <ModalCloseButton />
        </ModalHeader>

        <ModalBody>
          {!product ? (
            <Spinner size='large' />
          ) : (
            <VStack space='lg'>
              <Box>
                <Text color='$textLight300'>{product.bank}</Text>
                <Text fontSize='$lg' fontWeight='$semibold' color='$textLight0' mb='$1'>
                  {product.productName}
                </Text>
                <Text color={isPriority ? '$green400' : '$textLight300'} fontWeight='$bold'>
                  Ставка {formattedRate}
                </Text>
                <Text color='$textLight300'>Срок: {product.termMonths} мес.</Text>
                <Text color='$textLight300'>
                  Лимиты: {!product.maxAmount && 'От '}
                  {product.minAmount} {productSymbol}{' '}
                  {product.maxAmount && `— ${product.maxAmount} ${productSymbol}`}
                </Text>
              </Box>

              <FormControl isInvalid={Boolean(accountError)}>
                <FormControlLabel>
                  <FormControlLabelText color='$textLight300'>
                    С какого счёта пополнить
                  </FormControlLabelText>
                </FormControlLabel>
                <Pressable
                  onPress={toggleAccountList}
                  style={{
                    borderRadius: 12,
                    borderWidth: 1,
                    borderColor: '#525252',
                    paddingVertical: 12,
                    paddingHorizontal: 14,
                  }}
                >
                  <Text color='$textLight0'>
                    {selectedAccount
                      ? `${selectedAccount.bank} •• ${selectedAccount.number.slice(-4)} — ${selectedAccount.balance} ${CURRENCY_SYMBOL[selectedAccount.currency] ?? selectedAccount.currency}`
                      : eligibleAccounts.length
                        ? 'Выберите счёт'
                        : 'Нет подходящих счетов'}
                  </Text>
                </Pressable>

                {isAccountListOpen && (
                  <Box
                    borderRadius={12}
                    borderWidth={1}
                    borderColor='$backgroundDark700'
                    mt='$2'
                    maxHeight={200}
                    overflow='hidden'
                  >
                    {isAccountsLoadingState ? (
                      <Spinner />
                    ) : (
                      <ScrollView>{eligibleAccounts.map(renderAccountItem)}</ScrollView>
                    )}
                  </Box>
                )}
                {accountError ? (
                  <FormControlError>
                    <FormControlErrorText>{accountError}</FormControlErrorText>
                  </FormControlError>
                ) : null}
              </FormControl>

              <FormControl isInvalid={Boolean(amountError)}>
                <FormControlLabel>
                  <FormControlLabelText color='$textLight300'>
                    Сумма пополнения
                  </FormControlLabelText>
                </FormControlLabel>
                <Input borderColor='$backgroundDark700' borderRadius={12} borderWidth={1}>
                  <InputField
                    placeholder={`От ${product.minAmount}`}
                    keyboardType='decimal-pad'
                    value={amountValue}
                    placeholderTextColor='$textLight400'
                    color='$textLight0'
                    onChangeText={setAmountValue}
                  />
                </Input>
                {amountError ? (
                  <FormControlError>
                    <FormControlErrorText>{amountError}</FormControlErrorText>
                  </FormControlError>
                ) : null}
              </FormControl>

              {generalError ? (
                <Text color='$red500' fontSize='$sm'>
                  {generalError}
                </Text>
              ) : null}
            </VStack>
          )}
        </ModalBody>

        <ModalFooter>
          <Button
            flex={1}
            bg='$green600'
            borderRadius={12}
            isDisabled={disabledSubmit || isSubmitting}
            onPress={handleSubmit}
          >
            {isSubmitting ? (
              <Spinner color='$white' />
            ) : (
              <Text color='white' fontWeight='$bold'>
                Открыть
              </Text>
            )}
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default memo(DepositCreateModal)
