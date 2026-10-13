import { NavigationContainer } from '@react-navigation/native'
import { createStackNavigator } from '@react-navigation/stack'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import React, { useEffect } from 'react'
import { ScrollView, Text } from 'react-native'

import { registerOnUnauthorized } from '~/features/Auth/handler/unauthorizedHandler'
import BankConsentionPage from '~/pages/BankConsentionPage/BankConsentionPage'
import ExtendedAccountPage from '~/pages/ExtendedAccountPage/ExtendedAccountPage'
import ExtendedCardPage from '~/pages/ExtendedCardPage/ExtendedCardPage'
import LoginPage from '~/pages/LoginPage/LoginPage'
import MainPage from '~/pages/MainPage/MainPage'
import QRPaymentPage from '~/pages/QRPaymentPage/QRPaymentPage'
import ShowcasePage from '~/pages/ShowcasePage/ShowcasePage'
import SignupPage from '~/pages/SignupPage/SignupPage'
import AutoScrollView from '~/shared/components/AutoScrollView/AutoScrollView'
import AuthGuard from '~/shared/components/Guard/AuthGuard'
import LayoutWrapper from '~/shared/components/LayoutWrapper/LayoutWrapper'
import { isEnabledSignUp } from '~/shared/config/featureToggles'
import { ACCESS_TOKEN_KEY } from '~/shared/constants/authConstants'
import { RootStackParamList } from '~/shared/typings/App/AppTypes'
import { setDataInStore } from '~/shared/utils/tokenSaving/tokenSaving'
import { AccountIdProvider } from '~/widgets/TransactionsFeed/context/useAccountContext'

const client = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 0,
      retryDelay: attemptIndex => Math.min(1000 * 2 ** attemptIndex, 30000),
      refetchOnWindowFocus: true,
      refetchOnReconnect: true,
      refetchOnMount: true,
    },
    mutations: {
      retry: 0,
      retryDelay: attemptIndex => Math.min(1000 * 2 ** attemptIndex, 30000),
    },
  },
})

const Stack = createStackNavigator<RootStackParamList>()

const linking: any = {
  prefixes: ['phoenix://', 'https://phoenix.com'],
  config: {
    screens: {
      login: 'login',
      main: 'main',
      consention: 'consention',
      card: 'card/:cardId',
      account: 'account/:accountId',
      qrPayment: 'qr-payment',
      signup: 'signup',
      showcase: 'showcase',
    },
  },
}

function Layout() {
  useEffect(() => {
    registerOnUnauthorized(() => {
      setDataInStore(ACCESS_TOKEN_KEY, '')
    })
  }, [])

  return (
    <QueryClientProvider client={client}>
      <LayoutWrapper>
        <AccountIdProvider>
          <NavigationContainer linking={linking} fallback={<Text>Loading...</Text>}>
            <Stack.Navigator
              initialRouteName='main'
              screenOptions={{
                headerShown: false,
                gestureEnabled: true,
              }}
            >
              <Stack.Screen name='login' component={LoginPage} options={{ title: 'Войти' }} />

              <Stack.Screen name='main' options={{ title: 'Главная' }}>
                {() => (
                  <AuthGuard>
                    <AutoScrollView
                      contentContainerStyle={{ flexGrow: 1 }}
                      showsVerticalScrollIndicator={false}
                    >
                      <MainPage />
                    </AutoScrollView>
                  </AuthGuard>
                )}
              </Stack.Screen>

              <Stack.Screen name='consention' options={{ title: 'Запрос на доступ' }}>
                {() => (
                  <AuthGuard>
                    <ScrollView
                      contentContainerStyle={{ flexGrow: 1 }}
                      showsVerticalScrollIndicator={false}
                    >
                      <BankConsentionPage />
                    </ScrollView>
                  </AuthGuard>
                )}
              </Stack.Screen>

              <Stack.Screen name='card' options={{ title: 'Информация о карте' }}>
                {() => (
                  <AuthGuard>
                    <ScrollView
                      contentContainerStyle={{ flexGrow: 1 }}
                      showsVerticalScrollIndicator={false}
                    >
                      <ExtendedCardPage />
                    </ScrollView>
                  </AuthGuard>
                )}
              </Stack.Screen>

              <Stack.Screen name='account' options={{ title: 'Информация о счёте' }}>
                {() => (
                  <AuthGuard>
                    <ScrollView
                      contentContainerStyle={{ flexGrow: 1 }}
                      showsVerticalScrollIndicator={false}
                    >
                      <ExtendedAccountPage />
                    </ScrollView>
                  </AuthGuard>
                )}
              </Stack.Screen>

              <Stack.Screen name='qrPayment' options={{ title: 'Перевод по QR' }}>
                {() => (
                  <AuthGuard>
                    <ScrollView
                      contentContainerStyle={{ flexGrow: 1 }}
                      showsVerticalScrollIndicator={false}
                    >
                      <QRPaymentPage />
                    </ScrollView>
                  </AuthGuard>
                )}
              </Stack.Screen>

              <Stack.Screen name='showcase' options={{ title: 'Витрина продуктов' }}>
                {() => (
                  <AuthGuard>
                    <ScrollView
                      contentContainerStyle={{ flexGrow: 1 }}
                      showsVerticalScrollIndicator={false}
                    >
                      <ShowcasePage />
                    </ScrollView>
                  </AuthGuard>
                )}
              </Stack.Screen>

              {isEnabledSignUp && (
                <Stack.Screen
                  name='signup'
                  component={SignupPage}
                  options={{ title: 'Регистрация' }}
                />
              )}
            </Stack.Navigator>
          </NavigationContainer>
        </AccountIdProvider>
      </LayoutWrapper>
    </QueryClientProvider>
  )
}

export default Layout
