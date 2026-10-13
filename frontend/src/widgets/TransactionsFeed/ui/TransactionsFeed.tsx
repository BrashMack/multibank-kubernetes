import { Box, Spinner } from '@gluestack-ui/themed'
import { useInfiniteQuery } from '@tanstack/react-query'
import { useCallback } from 'react'
import { FlatList } from 'react-native'

import useService from '~/entities/useService'
import TransactionItem from '~/features/ExtendedAccount/TransactionItem/ui/TransactionItem'
import { Transaction } from '~/shared/typings/Transactions/TransactionsTypes'

import { useAccountId } from '../context/useAccountContext'

const LIMIT = 10

const TransactionsFeed = () => {
  const { getTransactionsByAccount } = useService()
  const { accountId } = useAccountId()
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage } = useInfiniteQuery({
    queryKey: ['transactions', accountId],
    queryFn: ({ pageParam = 1, signal }) =>
      getTransactionsByAccount(
        {
          accountId,
          page: Number(pageParam),
          limit: LIMIT,
        },
        signal,
      ),
    initialPageParam: 1,
    enabled: !!accountId,
    getNextPageParam: (lastPage, allPages) =>
      lastPage.transactions.length < LIMIT ? undefined : allPages.length + 1,
  })

  const items = data?.pages.flatMap(p => p.transactions) ?? []

  const renderItem = useCallback(
    ({ item }: { item: Transaction }) => <TransactionItem {...item} />,
    [],
  )

  return (
    <FlatList
      data={items}
      renderItem={renderItem}
      keyExtractor={item => `${item.transactionId}__${item.createdAt}`}
      onEndReached={() => {
        if (hasNextPage && !isFetchingNextPage) {
          fetchNextPage()
        }
      }}
      onEndReachedThreshold={0.5}
      ListFooterComponent={
        isFetchingNextPage ? (
          <Box py='$4' alignItems='center'>
            <Spinner size='small' />
          </Box>
        ) : null
      }
    />
  )
}

export default TransactionsFeed
