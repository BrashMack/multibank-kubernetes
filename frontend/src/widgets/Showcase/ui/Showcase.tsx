import { Box, Button, ScrollView, Spinner, Text } from '@gluestack-ui/themed'
import { useNavigation } from '@react-navigation/native'
import { StackNavigationProp } from '@react-navigation/stack'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import React, { memo, useCallback, useMemo, useState } from 'react'
import { Pressable } from 'react-native'

import useService from '~/entities/useService'
import DepositCreateModal from '~/features/Deposit/ui/DepositCreateModal'
import PriorityDepositCard from '~/features/Deposit/ui/PriorityDepositCard'
import { ErrorState } from '~/features/ServiceStates/ErrorState'
import { CURRENCY_SYMBOL } from '~/shared/constants/commonConstants'
import { RootStackParamList } from '~/shared/typings/App/AppTypes'
import {
  DepositAgreement,
  DepositProduct,
  GetDepositAgreementsResponse,
} from '~/shared/typings/Deposits/DepositsTypes'

interface OfferCardProps extends DepositProduct {
  onPress: (product: DepositProduct) => void
}

const OfferCard: React.FC<OfferCardProps> = ({ onPress, ...product }) => {
  const { productName, interestRate, minAmount, maxAmount, currency, termMonths, bank } = product
  const symbol = useMemo(() => currency && (CURRENCY_SYMBOL[currency] ?? ''), [currency])
  const min = useMemo(() => Number(minAmount), [minAmount])
  const max = useMemo(() => Number(maxAmount), [maxAmount])

  return (
    <Box maxWidth={550}>
      <Pressable onPress={() => onPress(product)}>
        <Box
          w='100%'
          h={120}
          borderRadius={16}
          borderWidth={1}
          borderColor='$backgroundDark700'
          bg='$backgroundDark900'
          p='$4'
          mb='$4'
        >
          <Box flex={1} flexDirection='column' justifyContent='flex-start' alignItems='flex-end'>
            <Text color='$textLight0' fontWeight='$bold' fontSize='$lg'>
              {productName}
            </Text>
            <Text color='$textLight300' fontSize='$sm' mb='$2'>
              {bank}
            </Text>
          </Box>
          <Box flex={1} flexDirection='column' justifyContent='flex-end' alignItems='flex-start'>
            <Text color='$textLight300' fontWeight='$semibold'>
              Ставка {interestRate}%
            </Text>
            <Text color='$textLight300'>
              От {min.toFixed(0)} {symbol} {maxAmount && `до ${max.toFixed(0)} ${symbol} `}на{' '}
              {termMonths} мес.
            </Text>
          </Box>
        </Box>
      </Pressable>
    </Box>
  )
}

const Showcase: React.FC = () => {
  const navigation = useNavigation<StackNavigationProp<RootStackParamList>>()
  const queryClient = useQueryClient()
  const { getDepositProducts, getPriorityDeposit } = useService()

  const [selectedProduct, setSelectedProduct] = useState<DepositProduct | null>(null)
  const [isModalOpen, setModalOpen] = useState(false)
  const [isPriority, setIsPriority] = useState(false)

  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: ['deposit-products'],
    queryFn: ({ signal }) => getDepositProducts(signal),
  })

  const { data: priorityDeposit } = useQuery({
    queryKey: ['deposit-priority'],
    queryFn: ({ signal }) => getPriorityDeposit(signal),
  })

  const offers = useMemo(() => data?.deposits ?? [], [data?.deposits])

  const openModal = useCallback((product: DepositProduct) => {
    setSelectedProduct(product)
    setModalOpen(true)
  }, [])

  const openModalPriority = useCallback((product: DepositProduct) => {
    setSelectedProduct(product)
    setIsPriority(true)
    setModalOpen(true)
  }, [])

  const closeModal = useCallback(() => {
    setModalOpen(false)
    setIsPriority(false)
    setSelectedProduct(null)
  }, [])

  const handleDepositCreated = useCallback(
    (deposit: DepositAgreement) => {
      navigation.navigate('main')
      queryClient.setQueryData(['deposits-agreements'], (prev?: GetDepositAgreementsResponse) => {
        if (!prev) return { deposits: [deposit] }
        const exists = prev.deposits.some(item => item.agreementId === deposit.agreementId)
        if (exists) return prev
        return { deposits: [deposit, ...prev.deposits] }
      })
    },
    [queryClient, navigation],
  )

  return (
    <>
      <Box flex={1} bg='$backgroundDark950' p='$4'>
        <Button
          mb='$4'
          px='$4'
          maxWidth='$32'
          borderRadius={16}
          bg='$blue600'
          onPress={() => navigation.navigate('main' as never)}
          shadowColor='#000'
          shadowOpacity={0.3}
          shadowRadius={6}
          elevation={5}
          justifyContent='center'
          alignItems='center'
        >
          <Text color='$white' fontSize='$md' fontWeight='$semibold'>
            На главную
          </Text>
        </Button>
        <Text color='$white' fontSize='$2xl' fontWeight='$bold' mb='$4'>
          Доступные вклады
        </Text>
        {isError && <ErrorState message='Не удалось загрузить предложения' onRetry={refetch} />}
        {!isError && (
          <ScrollView showsVerticalScrollIndicator={false}>
            {priorityDeposit && (
              <PriorityDepositCard {...priorityDeposit} onPress={openModalPriority} />
            )}
            {isLoading && (
              <Box py='$8' alignItems='center'>
                <Spinner size='large' />
              </Box>
            )}
            {!isLoading && offers.length === 0 && (
              <Text color='$textLight300'>Нет доступных предложений.</Text>
            )}
            {!isLoading &&
              offers
                .filter(
                  product =>
                    product.productId !== priorityDeposit?.productId ||
                    product.bank !== priorityDeposit?.bank,
                )
                .map(product => (
                  <OfferCard key={`${product.productId}`} {...product} onPress={openModal} />
                ))}
          </ScrollView>
        )}
      </Box>
      <DepositCreateModal
        isOpen={isModalOpen}
        product={selectedProduct}
        isPriority={isPriority}
        onClose={closeModal}
        onSuccess={handleDepositCreated}
      />
    </>
  )
}

memo(OfferCard)
export default memo(Showcase)
