import { useNavigation } from '@react-navigation/native'
import { StackNavigationProp } from '@react-navigation/stack'
import { memo, useEffect, useState } from 'react'
import { ActivityIndicator, View } from 'react-native'

import { enableAuthGuard } from '~/shared/config/featureToggles'
import { ACCESS_TOKEN_KEY } from '~/shared/constants/authConstants'
import { RootStackParamList } from '~/shared/typings/App/AppTypes'
import { getDataFromStoreAsync } from '~/shared/utils/tokenSaving/tokenSaving'

const AuthGuard: React.FC<React.PropsWithChildren<unknown>> = ({ children }) => {
  const [loading, setLoading] = useState(true)
  const [authenticated, setAuthenticated] = useState(false)
  const navigation = useNavigation<StackNavigationProp<RootStackParamList>>()

  useEffect(() => {
    const checkAuth = async () => {
      try {
        setLoading(true)
        const token = await getDataFromStoreAsync(ACCESS_TOKEN_KEY)
        setAuthenticated(!!token)

        if (!token && enableAuthGuard) {
          navigation.navigate('login')
        }
      } catch {
        if (enableAuthGuard) navigation.navigate('login')
      } finally {
        setLoading(false)
      }
    }

    checkAuth()
  }, [navigation])

  if (loading) {
    return (
      <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center' }}>
        <ActivityIndicator size='large' />
      </View>
    )
  }

  return authenticated || !enableAuthGuard ? <>{children}</> : null
}

export default memo(AuthGuard)
