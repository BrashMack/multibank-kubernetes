import { Box, ScrollView, Text } from '@gluestack-ui/themed'
import { useNavigation } from '@react-navigation/native'
import { StackNavigationProp } from '@react-navigation/stack'
import { useQuery } from '@tanstack/react-query'
import React, { useCallback, useEffect, useMemo } from 'react'

import useService from '~/entities/useService'
import AccountCard from '~/features/Accounts/AccountCard/ui/AccountCard'
import { ErrorState } from '~/features/ServiceStates/ErrorState'
import SkeletonCard from '~/features/ServiceStates/SkeletonCard'
import { ACCOUNTS_COUNT_KEY, DEFAULT_SKELETON_COUNT } from '~/shared/constants/commonConstants'
import { RootStackParamList } from '~/shared/typings/App/AppTypes'
import { getDataFromStore, setDataInStore } from '~/shared/utils/tokenSaving/tokenSaving'
import { isVoid } from '~/shared/utils/validators/isVoid'
import { useAccountId } from '~/widgets/TransactionsFeed/context/useAccountContext'

const AccountsFeed: React.FC = () => {
  const { getAccounts } = useService()
  const { setAccountId } = useAccountId()
  const navigation = useNavigation<StackNavigationProp<RootStackParamList>>()
  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: ['accounts'],
    queryFn: ({ signal }) => getAccounts(signal),
  })

  const accountsLength = data?.accounts?.length
  const skeletonCount = useMemo(() => {
    const storedCount = Number(getDataFromStore(ACCOUNTS_COUNT_KEY))
    return !Number.isNaN(storedCount) && storedCount > 0 ? storedCount : DEFAULT_SKELETON_COUNT
  }, [])

  useEffect(() => {
    if (!isLoading && !isVoid(accountsLength)) {
      const nextCount = accountsLength > 0 ? accountsLength : DEFAULT_SKELETON_COUNT
      setDataInStore(ACCOUNTS_COUNT_KEY, String(nextCount))
    }
  }, [accountsLength, isLoading])

  const onAccountPress = useCallback(
    (accountId: string) => {
      setAccountId(accountId)
      navigation.navigate('account', { accountId })
    },
    [navigation, setAccountId],
  )

  return (
    <Box p='$4' pt='$1'>
      <Text color='$white' fontSize='$2xl' fontWeight='$bold' mb='$4'>
        Счета
      </Text>
      {isError && <ErrorState message='Не удалось загрузить счета' onRetry={refetch} />}
      {!isError && (
        <ScrollView horizontal showsHorizontalScrollIndicator={false}>
          {isLoading &&
            Array.from({ length: skeletonCount }).map((_, index) => (
              <SkeletonCard key={`account-skeleton__${index}`} />
            ))}
          {!isLoading &&
            data?.accounts?.map(account => (
              <AccountCard
                key={`${account.accountId}__${account.number}`}
                {...account}
                onPress={onAccountPress}
              />
            ))}
        </ScrollView>
      )}
    </Box>
  )
}

export default AccountsFeed
