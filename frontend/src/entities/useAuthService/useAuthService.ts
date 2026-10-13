import { useCallback, useMemo } from 'react'

import { axiosInstance } from '~/shared'
import { AuthResponse } from '~/shared/typings/Auth/AuthTypes'

const useAuthService = () => {
  const login = useCallback(
    async (login: string, password: string, signal?: AbortSignal): Promise<AuthResponse> => {
      const response = await axiosInstance.post(
        `/login`,
        {
          login,
          password,
        },
        {
          signal,
        },
      )

      return response.data
    },
    [],
  )

  return useMemo(() => ({ login }), [login])
}

export default useAuthService
