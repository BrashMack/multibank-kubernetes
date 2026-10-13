import { Box, Text } from '@gluestack-ui/themed'
import React from 'react'

import AccountView from '~/widgets/SingleAccount/ui/AccountView'
import TransactionsFeed from '~/widgets/TransactionsFeed/ui/TransactionsFeed'

const ExtendedAccountPage: React.FC = () => {
  return (
    <>
      <AccountView />
      <Box bg='$backgroundDark950' px='$4' flex={1}>
        <Text color='$white' fontSize='$2xl' px='$4' fontWeight='$bold' mb='$4'>
          История транзакций
        </Text>
        <TransactionsFeed />
      </Box>
    </>
  )
}

export default ExtendedAccountPage
