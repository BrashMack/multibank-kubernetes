import {
  Box,
  Button,
  FormControl,
  FormControlError,
  FormControlErrorText,
  FormControlLabel,
  Heading,
  Input,
  InputField,
  Spinner,
  Text,
  VStack,
} from '@gluestack-ui/themed'
import { useNavigation } from '@react-navigation/native'
import type { StackNavigationProp } from '@react-navigation/stack'
import { Formik } from 'formik'
import React, { memo, useCallback } from 'react'
import * as Yup from 'yup'

import { isEnabledSignUp } from '~/shared/config/featureToggles'
import { RootStackParamList } from '~/shared/typings/App/AppTypes'

interface AuthFormProps {
  onFormSubmit: (username: string, password: string) => void
  isLoginForm?: boolean
  errorMessage?: string | null
}

const AuthForm: React.FC<AuthFormProps> = ({ onFormSubmit, isLoginForm = false, errorMessage }) => {
  const navigation = useNavigation<StackNavigationProp<RootStackParamList>>()

  // почему при isLoginForm - true кидает на signup?
  const handleSwitchScreen = useCallback(() => {
    if (isLoginForm) {
      navigation.navigate('signup')
    } else {
      navigation.navigate('login')
    }
  }, [navigation, isLoginForm])

  const validationSchema = Yup.object().shape({
    username: Yup.string().trim().required('Введите имя пользователя').min(3, 'Минимум 3 символа'),
    password: Yup.string().required('Введите пароль').min(6, 'Минимум 6 символов'),
    ...(isLoginForm
      ? {}
      : {
          confirmPassword: Yup.string()
            .required('Повторите пароль')
            .oneOf([Yup.ref('password')], 'Пароли не совпадают'),
        }),
  })

  return (
    <Box
      flex={1}
      alignItems='center'
      justifyContent='center'
      px='$6'
      py='$8'
      bg='$backgroundDark950'
    >
      <VStack
        w='$full'
        maxWidth='$80'
        space='lg'
        p='$6'
        rounded='$2xl'
        bg='$backgroundDark800'
        shadowColor='$backgroundDark900'
        shadowOpacity={0.4}
        shadowRadius={10}
      >
        <Heading size='xl' textAlign='center' color='$textLight0' mb='$2'>
          {isLoginForm ? 'Вход' : 'Регистрация'}
        </Heading>

        <Formik
          initialValues={{
            username: '',
            password: '',
            confirmPassword: '',
          }}
          validationSchema={validationSchema}
          onSubmit={(values, { setSubmitting }) => {
            onFormSubmit(values.username.trim(), values.password.trim())
            setSubmitting(false)
          }}
        >
          {({ handleChange, handleSubmit, values, errors, touched, isSubmitting }) => (
            <VStack space='md'>
              <FormControl isInvalid={!!(touched.username && errors.username)}>
                <FormControlLabel>
                  <Text color='$textLight500' mb='$1'>
                    Имя пользователя
                  </Text>
                </FormControlLabel>
                <Input bg='$backgroundDark700' borderColor='$backgroundDark600'>
                  <InputField
                    placeholder='Введите имя пользователя'
                    placeholderTextColor='$textLight400'
                    value={values.username}
                    onChangeText={handleChange('username')}
                    autoCapitalize='none'
                    color='$textLight0'
                  />
                </Input>
                {touched.username && errors.username && (
                  <FormControlError>
                    <FormControlErrorText>{errors.username}</FormControlErrorText>
                  </FormControlError>
                )}
              </FormControl>

              <FormControl isInvalid={!!(touched.password && errors.password)} mb='$4'>
                <FormControlLabel>
                  <Text color='$textLight500' mb='$1'>
                    Пароль
                  </Text>
                </FormControlLabel>
                <Input bg='$backgroundDark700' borderColor='$backgroundDark600'>
                  <InputField
                    placeholder='Введите пароль'
                    placeholderTextColor='$textLight400'
                    value={values.password}
                    onChangeText={handleChange('password')}
                    secureTextEntry
                    color='$textLight0'
                  />
                </Input>
                {touched.password && errors.password && (
                  <FormControlError>
                    <FormControlErrorText>{errors.password}</FormControlErrorText>
                  </FormControlError>
                )}
              </FormControl>

              {!isLoginForm && (
                <FormControl isInvalid={!!(touched.confirmPassword && errors.confirmPassword)}>
                  <FormControlLabel>
                    <Text color='$textLight500' mb='$1'>
                      Повторите пароль
                    </Text>
                  </FormControlLabel>
                  <Input bg='$backgroundDark700' borderColor='$backgroundDark600'>
                    <InputField
                      placeholder='Введите пароль повторно'
                      placeholderTextColor='$textLight400'
                      value={values.confirmPassword}
                      onChangeText={handleChange('confirmPassword')}
                      secureTextEntry
                      color='$textLight0'
                    />
                  </Input>
                  {touched.confirmPassword && errors.confirmPassword && (
                    <FormControlError>
                      <FormControlErrorText>{errors.confirmPassword}</FormControlErrorText>
                    </FormControlError>
                  )}
                </FormControl>
              )}

              <Button
                onPress={handleSubmit as any}
                bg='$primary500'
                rounded='$xl'
                size='lg'
                isDisabled={isSubmitting}
              >
                {isSubmitting ? (
                  <Spinner color='$textLight0' />
                ) : (
                  <Text color='$textLight0' fontWeight='$bold'>
                    {isLoginForm ? 'Войти' : 'Зарегистрироваться'}
                  </Text>
                )}
              </Button>

              {isLoginForm && errorMessage && (
                <Text color='$error500' textAlign='center' position='relative'>
                  {errorMessage}
                </Text>
              )}

              {isEnabledSignUp && (
                <Button variant='link' onPress={handleSwitchScreen} mt='$2' alignSelf='center'>
                  <Text color='$primary400' fontWeight='$bold'>
                    {isLoginForm ? 'Нет аккаунта? Зарегистрироваться' : 'Уже есть аккаунт? Войти'}
                  </Text>
                </Button>
              )}
            </VStack>
          )}
        </Formik>
      </VStack>
    </Box>
  )
}

export default memo(AuthForm)
