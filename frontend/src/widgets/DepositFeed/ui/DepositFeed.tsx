/* eslint-disable react-hooks/refs */
import { Box, ScrollView, Text } from '@gluestack-ui/themed'
import { useNavigation } from '@react-navigation/native'
import { StackNavigationProp } from '@react-navigation/stack'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import React, { memo, useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { Animated, LayoutAnimation, Platform, Pressable, StyleSheet, UIManager } from 'react-native'
import Svg, { Path } from 'react-native-svg'

import useService from '~/entities/useService'
import DepositCard from '~/features/Deposit/ui/DepositCard'
import DepositCreateModal from '~/features/Deposit/ui/DepositCreateModal'
import PriorityDepositCard from '~/features/Deposit/ui/PriorityDepositCard'
import RequestConsentDeposit from '~/features/Deposit/ui/RequestConsentDeposit'
import { ErrorState } from '~/features/ServiceStates/ErrorState'
import SkeletonCard from '~/features/ServiceStates/SkeletonCard'
import { DEFAULT_SKELETON_COUNT, DEPOSITS_COUNT_KEY } from '~/shared/constants/commonConstants'
import { RootStackParamList } from '~/shared/typings/App/AppTypes'
import {
  DepositAgreement,
  DepositProduct,
  GetDepositAgreementsResponse,
} from '~/shared/typings/Deposits/DepositsTypes'
import { getDataFromStore, setDataInStore } from '~/shared/utils/tokenSaving/tokenSaving'
import { isVoid } from '~/shared/utils/validators/isVoid'

if (Platform.OS === 'android' && UIManager.setLayoutAnimationEnabledExperimental) {
  UIManager.setLayoutAnimationEnabledExperimental(true)
}

const DepositFeed: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false)
  const [selectedProduct, setSelectedProduct] = useState<DepositProduct | null>(null)
  const [isModalOpen, setModalOpen] = useState(false)
  const [isPriority, setIsPriority] = useState(false)
  const arrowAnimation = useRef(new Animated.Value(0)).current

  const navigation = useNavigation<StackNavigationProp<RootStackParamList>>()
  const queryClient = useQueryClient()
  const { getDepositAgreements, getPriorityDeposit } = useService()

  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: ['deposits-agreements'],
    queryFn: ({ signal }) => getDepositAgreements(signal),
  })

  const { data: priorityDeposit } = useQuery({
    queryKey: ['deposit-priority'],
    queryFn: ({ signal }) => getPriorityDeposit(signal),
  })

  useEffect(() => {
    Animated.timing(arrowAnimation, {
      toValue: isOpen ? 1 : 0,
      duration: 200,
      useNativeDriver: true,
    }).start()
  }, [arrowAnimation, isOpen])

  const deposits = useMemo(() => data?.deposits ?? [], [data])
  const depositsLength = deposits.length
  const skeletonCount = useMemo(() => {
    const storedCount = Number(getDataFromStore(DEPOSITS_COUNT_KEY))
    return !Number.isNaN(storedCount) && storedCount > 0 ? storedCount : DEFAULT_SKELETON_COUNT
  }, [])

  useEffect(() => {
    if (!isLoading && !isVoid(depositsLength)) {
      const nextCount = depositsLength > 0 ? depositsLength : DEFAULT_SKELETON_COUNT
      setDataInStore(DEPOSITS_COUNT_KEY, String(nextCount))
    }
  }, [depositsLength, isLoading])

  const rotate = arrowAnimation.interpolate({
    inputRange: [0, 1],
    outputRange: ['0deg', '90deg'],
  })

  const handleToggle = useCallback(() => {
    LayoutAnimation.configureNext(LayoutAnimation.Presets.easeInEaseOut)
    setIsOpen(prev => !prev)
  }, [])

  const onPriorityDepositPress = useCallback((product: DepositProduct) => {
    setSelectedProduct(product)
    setIsPriority(true)
    setModalOpen(true)
  }, [])

  const handleModalClose = useCallback(() => {
    setModalOpen(false)
    setIsPriority(false)
    setSelectedProduct(null)
  }, [])

  const handleDepositCreated = useCallback(
    (deposit: DepositAgreement) => {
      queryClient.setQueryData(['deposits-agreements'], (prev?: GetDepositAgreementsResponse) => {
        if (!prev) {
          return { deposits: [deposit] }
        }
        const exists = prev.deposits.some(item => item.agreementId === deposit.agreementId)
        if (exists) {
          return prev
        }
        return { deposits: [deposit, ...prev.deposits] }
      })
    },
    [queryClient],
  )

  const onConsentDeposit = useCallback(() => {
    navigation.navigate('showcase')
  }, [navigation])

  return (
    <>
      <Box flex={1} p='$4' pt='$1'>
        <Pressable style={styles.toggle} onPress={handleToggle}>
          <Animated.View style={{ transform: [{ rotate }] }}>
            <Svg width={28} height={28} viewBox='0 0 24 24'>
              <Path d='M8 5V19L18 12L8 5Z' fill='#ffffffff' />
            </Svg>
          </Animated.View>
          <Text color='$white' fontSize='$2xl' fontWeight='$bold' ml='$2'>
            Вклады
          </Text>
        </Pressable>

        {isOpen && (
          <Box mt='$4'>
            {isError && <ErrorState message='Не удалось загрузить вклады' onRetry={refetch} />}
            {!isError && (
              <>
                {priorityDeposit && (
                  <PriorityDepositCard
                    key={`${priorityDeposit.productId}__${priorityDeposit.bank}`}
                    {...priorityDeposit}
                    onPress={onPriorityDepositPress}
                  />
                )}
                <ScrollView horizontal showsHorizontalScrollIndicator={false}>
                  {isLoading &&
                    Array.from({ length: skeletonCount }).map((_, index) => (
                      <SkeletonCard key={`deposit-skeleton__${index}`} source='deposit' />
                    ))}
                  {!isLoading &&
                    deposits.length > 0 &&
                    deposits.map(deposit => <DepositCard key={deposit.agreementId} {...deposit} />)}
                  <RequestConsentDeposit onPress={onConsentDeposit} />
                </ScrollView>
              </>
            )}
          </Box>
        )}
      </Box>
      <DepositCreateModal
        isOpen={isModalOpen}
        product={selectedProduct}
        isPriority={isPriority}
        onClose={handleModalClose}
        onSuccess={handleDepositCreated}
      />
    </>
  )
}

const styles = StyleSheet.create({
  toggle: {
    flexDirection: 'row',
    alignItems: 'center',
  },
})

export default memo(DepositFeed)
