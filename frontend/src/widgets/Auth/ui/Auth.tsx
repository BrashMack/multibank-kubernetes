import { useNavigation } from '@react-navigation/native'
import { StackNavigationProp } from '@react-navigation/stack'
import { useMutation } from '@tanstack/react-query'
import { AxiosError, isAxiosError } from 'axios'
import React, { useState } from 'react'

import useService from '~/entities/useService'
import AuthForm from '~/features/Auth/AuthForm/ui/AuthForm'
import { ACCESS_TOKEN_KEY } from '~/shared/constants/authConstants'
import { RootStackParamList } from '~/shared/typings/App/AppTypes'
import { setDataInStore } from '~/shared/utils/tokenSaving/tokenSaving'

interface AuthProps {
  isLoginForm?: boolean
}

const Auth: React.FC<AuthProps> = ({ isLoginForm = false }) => {
  const { login: loginRequest } = useService()
  const navigation = useNavigation<StackNavigationProp<RootStackParamList>>()
  const [authError, setAuthError] = useState<string | null>(null)
  const mutation = useMutation({
    mutationFn: ({ login, password }: { login: string; password: string }) =>
      loginRequest(login, password),
    onSuccess: ({ token }) => {
      setAuthError(null)
      navigation.navigate('main')
      setDataInStore(ACCESS_TOKEN_KEY, token)
    },
    onError: (error: AxiosError | unknown) => {
      if (isAxiosError(error) && error.response?.status === 401) {
        setAuthError('Неверный логин или пароль')
        return
      }

      setAuthError('Не удалось выполнить вход. Попробуйте ещё раз.')
    },
  })

  const onFormSubmit = (username: string, password: string) => {
    setAuthError(null)
    mutation.mutateAsync({ login: username, password })
  }

  return <AuthForm onFormSubmit={onFormSubmit} isLoginForm={isLoginForm} errorMessage={authError} />
}

export default Auth
