import {
  Box,
  Button,
  FormControl,
  FormControlError,
  FormControlErrorText,
  FormControlLabel,
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
import { useMutation, useQueryClient } from '@tanstack/react-query'
import axios from 'axios'
import React, { memo, useCallback, useMemo, useState } from 'react'
import { Alert } from 'react-native'

import useService from '~/entities/useService'
import { MakeTransactionRequest } from '~/shared/typings/Transactions/TransactionsTypes'

interface TransferActionProps {
  sourceAccountNumber?: string
  onSuccess?: () => void | Promise<void>
}

const sanitizeAccountNumber = (value: string) => value.replace(/\s+/g, '')

const TransferAction: React.FC<TransferActionProps> = ({ sourceAccountNumber, onSuccess }) => {
  const { makeTransaction } = useService()

  const [isTransferModalOpen, setTransferModalOpen] = useState(false)
  const [destinationAccount, setdestinationAccount] = useState('')
  const [amountValue, setAmountValue] = useState('')
  const [destinationError, setDestinationError] = useState('')
  const [amountError, setAmountError] = useState('')
  const [generalError, setGeneralError] = useState('')
  const queryClient = useQueryClient()

  const hasSourceAccount = Boolean(sourceAccountNumber)

  const sanitizedSource = sourceAccountNumber ? sanitizeAccountNumber(sourceAccountNumber) : ''
  const maskedSourceNumber = useMemo(
    () => (sanitizedSource ? `•••• ${sanitizedSource.slice(-4)}` : '—'),
    [sanitizedSource],
  )

  const closeTransferModal = useCallback(() => {
    setTransferModalOpen(false)
    setdestinationAccount('')
    setAmountValue('')
    setDestinationError('')
    setAmountError('')
    setGeneralError('')
  }, [])

  const { mutate: submitTransaction, isPending: isTransactionPending } = useMutation({
    mutationFn: (payload: MakeTransactionRequest) => makeTransaction(payload),
    onSuccess: async () => {
      Alert.alert('Перевод создан', 'Запрос на перевод отправлен.')
      await queryClient.invalidateQueries({ queryKey: ['transactions'] })
      await queryClient.invalidateQueries({ queryKey: ['cards'] })
      await queryClient.invalidateQueries({ queryKey: ['accounts'] })
      await queryClient.invalidateQueries({ queryKey: ['account'] })
      await queryClient.invalidateQueries({ queryKey: ['card'] })

      closeTransferModal()
      if (onSuccess) {
        await onSuccess()
      }
    },
    onError: error => {
      if (axios.isAxiosError(error)) {
        const backendMessage =
          (typeof error.response?.data === 'object' &&
            (error.response?.data?.error || error.response?.data?.message)) ||
          ''
        const normalizedBackendMessage =
          typeof backendMessage === 'string' ? backendMessage.toLowerCase() : ''

        if (
          error.response?.status === 404 &&
          normalizedBackendMessage &&
          normalizedBackendMessage.includes('to')
        ) {
          setGeneralError('Счёт получателя не существует')
          return
        }

        if (
          error.response?.status === 400 &&
          normalizedBackendMessage &&
          normalizedBackendMessage.includes('from is empty')
        ) {
          setGeneralError('Отсутствует номер счёта отправителя')
          return
        }

        if (
          error.response?.status === 404 &&
          normalizedBackendMessage &&
          normalizedBackendMessage.includes('from')
        ) {
          setGeneralError('Счёт отправителя не существует')
          return
        }

        if (error.response?.status === 422) {
          setGeneralError('Недостаточно средств на счёте')
          return
        }

        if (backendMessage) {
          setGeneralError(backendMessage)
          return
        }
      }

      setGeneralError('Не удалось выполнить перевод. Попробуйте позже.')
    },
  })

  const openTransferModal = useCallback(() => {
    if (!hasSourceAccount) return
    setGeneralError('')
    setDestinationError('')
    setAmountError('')
    setTransferModalOpen(true)
  }, [hasSourceAccount])

  const handleTransferSubmit = useCallback(() => {
    setDestinationError('')
    setAmountError('')
    setGeneralError('')

    if (!sanitizedSource) {
      setGeneralError('Нет счёта для списания')
      return
    }
    const sanitizedDestination = sanitizeAccountNumber(destinationAccount)
    const isDestinationLengthValid =
      sanitizedDestination.length >= 18 && sanitizedDestination.length <= 20
    if (!sanitizedDestination || !isDestinationLengthValid) {
      setDestinationError('Укажите корректный номер счёта получателя')
    }

    const normalizedAmountString = amountValue.replace(',', '.').trim()
    const amountFormatRegex = /^\d+(\.\d{0,2})?$/
    const isAmountFormatValid = amountFormatRegex.test(normalizedAmountString)
    const normalizedAmountNumber = Number(normalizedAmountString)

    if (
      !normalizedAmountString ||
      !isAmountFormatValid ||
      Number.isNaN(normalizedAmountNumber) ||
      normalizedAmountNumber <= 0
    ) {
      setAmountError('Введите корректную сумму перевода')
    }

    if (
      !sanitizedDestination ||
      !isDestinationLengthValid ||
      !normalizedAmountString ||
      !isAmountFormatValid ||
      Number.isNaN(normalizedAmountNumber) ||
      normalizedAmountNumber <= 0
    ) {
      return
    }

    const [integerPart, decimals = ''] = normalizedAmountString.split('.')
    const formattedAmount = `${integerPart}.${decimals.padEnd(2, '0')}`

    submitTransaction({
      cardIdFrom: sanitizedSource,
      cardIdTo: sanitizedDestination,
      amount: formattedAmount,
    })
  }, [amountValue, destinationAccount, sanitizedSource, submitTransaction])

  return (
    <Box w='100%'>
      <Button
        mt='$2'
        px='$4'
        borderRadius={16}
        bg='$green600'
        isDisabled={!hasSourceAccount || isTransactionPending}
        onPress={openTransferModal}
        shadowColor='#000'
        shadowOpacity={0.3}
        shadowRadius={6}
        elevation={5}
        justifyContent='center'
        alignItems='center'
      >
        {isTransactionPending ? (
          <Spinner color='$white' />
        ) : (
          <Text color='$white' fontSize='$md' fontWeight='$semibold'>
            Перевести
          </Text>
        )}
      </Button>

      {!!generalError && hasSourceAccount && (
        <Text color='$red500' mt='$2'>
          {generalError}
        </Text>
      )}

      <Modal isOpen={isTransferModalOpen} onClose={closeTransferModal}>
        <ModalBackdrop />
        <ModalContent bg='$backgroundDark900'>
          <ModalHeader>
            <Text color='$white' fontSize='$lg' fontWeight='$bold'>
              Перевод средств
            </Text>
            <ModalCloseButton onPress={closeTransferModal} />
          </ModalHeader>
          <ModalBody>
            <Text color='$textLight400' mb='$4'>
              Со счёта {maskedSourceNumber}
            </Text>
            <VStack gap='$4'>
              <FormControl isInvalid={!!destinationError}>
                <FormControlLabel>
                  <Text color='$textLight500'>Счёт получателя</Text>
                </FormControlLabel>
                <Input bg='$backgroundDark800' borderColor='$backgroundDark700'>
                  <InputField
                    placeholder='Введите номер счёта'
                    placeholderTextColor='$textLight400'
                    color='$textLight0'
                    keyboardType='number-pad'
                    value={destinationAccount}
                    onChangeText={text => setdestinationAccount(text)}
                  />
                </Input>
                {!!destinationError && (
                  <FormControlError>
                    <FormControlErrorText>{destinationError}</FormControlErrorText>
                  </FormControlError>
                )}
              </FormControl>

              <FormControl isInvalid={!!amountError}>
                <FormControlLabel>
                  <Text color='$textLight500'>Сумма</Text>
                </FormControlLabel>
                <Input bg='$backgroundDark800' borderColor='$backgroundDark700'>
                  <InputField
                    placeholder='0.00'
                    placeholderTextColor='$textLight400'
                    color='$textLight0'
                    keyboardType='decimal-pad'
                    value={amountValue}
                    onChangeText={text => setAmountValue(text)}
                  />
                </Input>
                {!!amountError && (
                  <FormControlError>
                    <FormControlErrorText>{amountError}</FormControlErrorText>
                  </FormControlError>
                )}
              </FormControl>

              {!!generalError && (
                <Text color='$red500' fontSize='$sm'>
                  {generalError}
                </Text>
              )}
            </VStack>
          </ModalBody>
          <ModalFooter gap='$3'>
            <Button variant='outline' borderColor='$backgroundDark600' onPress={closeTransferModal}>
              <Text color='$textLight400'>Отмена</Text>
            </Button>
            <Button bg='$green600' onPress={handleTransferSubmit} isDisabled={isTransactionPending}>
              {isTransactionPending ? (
                <Spinner color='$white' />
              ) : (
                <Text color='$white' fontWeight='$semibold'>
                  Отправить
                </Text>
              )}
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </Box>
  )
}

export default memo(TransferAction)
