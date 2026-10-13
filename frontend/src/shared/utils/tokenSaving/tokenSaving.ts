import * as SecureStore from 'expo-secure-store'
import { Platform } from 'react-native'

import { ACCESS_TOKEN_KEY, REFRESH_TOKEN_KEY } from '~/shared/constants/authConstants'
import {
  ACCOUNTS_COUNT_KEY,
  CARDS_COUNT_KEY,
  DEPOSITS_COUNT_KEY,
} from '~/shared/constants/commonConstants'

const CURRENT_PLATFORM =
  Platform.OS === 'android' || Platform.OS === 'ios'
    ? 'mobile'
    : Platform.OS === 'web'
      ? 'web'
      : 'unknown'

export const setDataInStore = (key: string, value: string) => {
  if (CURRENT_PLATFORM === 'mobile') {
    SecureStore.setItem(key, value)
  } else if (CURRENT_PLATFORM === 'web') {
    localStorage.setItem(key, value)
  }
  return
}

export const getDataFromStore = (key: string): string | null => {
  if (CURRENT_PLATFORM === 'mobile') {
    return SecureStore.getItem(key)
  } else if (CURRENT_PLATFORM === 'web') {
    return localStorage.getItem(key)
  }
  return null
}

export const getDataFromStoreAsync = async (key: string): Promise<string | null> => {
  if (CURRENT_PLATFORM === 'mobile') {
    return await SecureStore.getItemAsync(key)
  } else if (CURRENT_PLATFORM === 'web') {
    return localStorage.getItem(key)
  }
  return null
}

export const setDataInStoreAsync = async (key: string, value: string): Promise<void> => {
  if (CURRENT_PLATFORM === 'mobile') {
    SecureStore.setItemAsync(key, value)
  } else if (CURRENT_PLATFORM === 'web') {
    localStorage.setItem(key, value)
  }
  return
}

export const resetStore = () => {
  if (CURRENT_PLATFORM === 'mobile') {
    SecureStore.deleteItemAsync(ACCESS_TOKEN_KEY)
    SecureStore.deleteItemAsync(REFRESH_TOKEN_KEY)
    SecureStore.deleteItemAsync(CARDS_COUNT_KEY)
    SecureStore.deleteItemAsync(ACCOUNTS_COUNT_KEY)
    SecureStore.deleteItemAsync(DEPOSITS_COUNT_KEY)
  } else if (CURRENT_PLATFORM === 'web') {
    localStorage.clear()
  }
  return
}
