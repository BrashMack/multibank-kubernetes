import React, { createContext, useContext, useState } from 'react'

import { AccountIdContext } from '~/shared/typings/Accounts/AccountsTypes'

const AcoountIdContext = createContext<AccountIdContext>({
  setAccountId: (accountId: string) => accountId,
  accountId: '',
})

export const useAccountId = () => useContext(AcoountIdContext)

export const AccountIdProvider: React.FC<React.PropsWithChildren<unknown>> = ({ children }) => {
  const [accountId, setId] = useState<string>('')

  const setAccountId = (accountId: string) => {
    setId(accountId)
  }

  return (
    <AcoountIdContext.Provider value={{ accountId, setAccountId }}>
      {children}
    </AcoountIdContext.Provider>
  )
}
