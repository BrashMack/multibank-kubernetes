import { Box, Image } from '@gluestack-ui/themed'
import { useNavigation } from '@react-navigation/native'
import { StackNavigationProp } from '@react-navigation/stack'
import { useQueryClient } from '@tanstack/react-query'
import React, { memo, useCallback } from 'react'
import { TouchableOpacity } from 'react-native'

import QRCamera from '~/features/Transfer/QRCamera/ui/QRCamera'
import { INTERFACE_ICONS } from '~/shared/constants/commonConstants'
import { RootStackParamList } from '~/shared/typings/App/AppTypes'
import { resetStore } from '~/shared/utils/tokenSaving/tokenSaving'

const NavBar: React.FC = () => {
  const queryClient = useQueryClient()
  const navigation = useNavigation<StackNavigationProp<RootStackParamList>>()

  const handleLogout = useCallback(() => {
    resetStore()
    queryClient.clear()
    navigation.navigate('login' as never)
  }, [navigation, queryClient])

  return (
    <Box
      flexDirection='row'
      justifyContent='space-between'
      alignItems='flex-start'
      bg='$backgroundDark950'
      px='$4'
      pb='$2'
      overflow='visible'
    >
      <QRCamera />
      <TouchableOpacity onPress={handleLogout}>
        <Image
          source={INTERFACE_ICONS.LogoutButton}
          alt='Logout'
          width={50}
          height={50}
          resizeMode='contain'
        />
      </TouchableOpacity>
    </Box>
  )
}

export default memo(NavBar)
