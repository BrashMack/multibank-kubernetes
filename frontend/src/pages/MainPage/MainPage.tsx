import { Box } from '@gluestack-ui/themed'
import React from 'react'

import AccountsFeed from '~/widgets/AccountFeed/ui/AccountsFeed'
import CardFeed from '~/widgets/CardFeed/ui/CardFeed'
import DepositFeed from '~/widgets/DepositFeed/ui/DepositFeed'
import NavBar from '~/widgets/NavBar/NavBar'

const MainPage: React.FC = () => {
  return (
    <Box flex={1} bg='$backgroundDark950'>
      <NavBar />
      <CardFeed />
      <AccountsFeed />
      <DepositFeed />
    </Box>
  )
}

export default MainPage
