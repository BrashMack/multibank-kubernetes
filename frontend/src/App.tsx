/* eslint-disable @typescript-eslint/no-require-imports */

import { config } from '@gluestack-ui/config'
import { GluestackUIProvider } from '@gluestack-ui/themed'
import { useFonts } from 'expo-font'
import * as SplashScreen from 'expo-splash-screen'
import { useCallback } from 'react'
import { LogBox, View } from 'react-native'

import Layout from './app/Layout'

// убрать при разработке
LogBox.ignoreAllLogs()

SplashScreen.preventAutoHideAsync()

export default function App() {
  const [fontsLoaded] = useFonts({
    MontserratRegular: require('~/../assets/fonts/Montserrat-Regular.ttf'),
    MontserratSemibold: require('~/../assets/fonts/Montserrat-SemiBold.ttf'),
    MontserratMedium: require('~/../assets/fonts/Montserrat-Medium.ttf'),
  })

  const onLayoutRootView = useCallback(async () => {
    if (fontsLoaded) {
      await SplashScreen.hideAsync()
    }
  }, [fontsLoaded])

  if (!fontsLoaded) {
    return null
  }

  return (
    <GluestackUIProvider config={config}>
      <View style={{ flex: 1 }} onLayout={onLayoutRootView}>
        <Layout />
      </View>
    </GluestackUIProvider>
  )
}
