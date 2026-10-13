import axios, { AxiosInstance } from 'axios'

import { clearProfileIfUnauthorized } from '~/features/Auth/handler/unauthorizedHandler'
import { API_BASE_URL } from '~/shared/config/config'

export const getBaseUrl = () => {
  return `${API_BASE_URL}/api/v1`
}

export const axiosInstance: AxiosInstance = axios.create({
  baseURL: getBaseUrl(),
  withCredentials: true,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const axiosInstanceMultipart: AxiosInstance = axios.create({
  baseURL: getBaseUrl(),
  withCredentials: true,
  timeout: 10000,
  headers: {
    'Content-Type': 'multipart/form-data',
  },
})

axiosInstance.interceptors.response.use(
  response => response,
  error => {
    if (error?.response?.status) {
      clearProfileIfUnauthorized(error.response.status)
    }
    return Promise.reject(error)
  },
)
