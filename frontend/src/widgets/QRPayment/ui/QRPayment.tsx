import { useNavigation, useRoute } from '@react-navigation/native'
import { StackNavigationProp } from '@react-navigation/stack'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import axios from 'axios'
import { useCallback, useEffect, useState } from 'react'

import useService from '~/entities/useService'
import QRPaymentCardAutocomplete, {
  QRTransaction,
} from '~/features/QRPayment/QRPaymentCardAutocomplete/ui/QRPaymentCardAutocomplete'
import { QRPaymentPageRouteProp, RootStackParamList } from '~/shared/typings/App/AppTypes'
import { MakeTransactionRequest } from '~/shared/typings/Transactions/TransactionsTypes'

const QRPayment: React.FC = () => {
  const { makeTransaction, getAccounts } = useService()
  const queryClient = useQueryClient()
  const navigation = useNavigation<StackNavigationProp<RootStackParamList>>()
  const [errorMsg, setErrorMsg] = useState<string>('')

  const handleError = (error: any) => {
    if (axios.isAxiosError(error)) {
      const backendMessage =
        (typeof error.response?.data === 'object' &&
          (error.response?.data?.error || error.response?.data?.message)) ||
        ''

      const normalized = typeof backendMessage === 'string' ? backendMessage.toLowerCase() : ''

      if (error.response?.status === 400 && normalized.includes('from is empty')) {
        setErrorMsg('Отсутствует номер счёта отправителя')
        return
      }
      if (error.response?.status === 404 && normalized.includes('to')) {
        setErrorMsg('Счёт получателя не существует')
        return
      }
      if (error.response?.status === 404 && normalized.includes('from')) {
        setErrorMsg('Счёт отправителя не существует')
        return
      }
      if (error.response?.status === 422) {
        setErrorMsg('Недостаточно средств на счёте')
        return
      }
      if (backendMessage) {
        setErrorMsg(backendMessage)
        return
      }
    }
    setErrorMsg('Не удалось выполнить перевод. Попробуйте позже.')
  }

  const { mutateAsync: submitIntermediateTransaction, isPending: isLoadingInter } = useMutation({
    mutationFn: (payload: MakeTransactionRequest) => makeTransaction(payload),
    onError: handleError,
  })

  const { mutateAsync: submitTransaction, isPending: isLoading } = useMutation({
    mutationFn: (payload: MakeTransactionRequest) => makeTransaction(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['accounts'] })
      await queryClient.invalidateQueries({ queryKey: ['cards'] })
      await queryClient.invalidateQueries({ queryKey: ['transactions'] })
      navigation.navigate('main')
    },
    onError: handleError,
  })

  const route = useRoute<QRPaymentPageRouteProp>()

  const { accountIdTo, amount, currency } = route.params ?? {
    accountIdTo: '',
    amount: '',
    currency: '',
  }

  useEffect(() => {
    if (!accountIdTo || !amount || !currency) {
      navigation.navigate('main')
    }
  }, [accountIdTo, amount, currency, navigation])

  const { data } = useQuery({
    queryKey: ['accounts'],
    queryFn: ({ signal }) => getAccounts(signal),
  })

  const onContinueClick = useCallback(
    async (transaction: QRTransaction) => {
      const { deltaNum, delta, secondaryId, primaryId, accountIdTo, amount } = transaction

      try {
        if (deltaNum < 0) {
          await submitIntermediateTransaction({
            cardIdFrom: secondaryId,
            cardIdTo: primaryId,
            amount: delta,
          })
        }

        if (Number(amount)) {
          await submitTransaction({
            cardIdFrom: primaryId,
            cardIdTo: accountIdTo,
            amount,
          })
        }
      } catch {
        return
      }
    },
    [submitIntermediateTransaction, submitTransaction],
  )

  return (
    <QRPaymentCardAutocomplete
      accounts={data?.accounts?.length ? [...data.accounts] : []}
      accountIdTo={accountIdTo}
      amount={amount}
      currency={currency}
      onContinueClick={onContinueClick}
      errorMsg={errorMsg}
      isLoading={isLoading || isLoadingInter}
    />
  )
}

export default QRPayment
