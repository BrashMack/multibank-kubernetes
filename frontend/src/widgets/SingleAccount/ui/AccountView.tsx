import { useRoute } from '@react-navigation/native'
import { useQuery } from '@tanstack/react-query'
import React, { useEffect } from 'react'

import useService from '~/entities/useService'
import ExtendedAccountView from '~/features/ExtendedAccount/ExtendedAccountView/ui/ExtendedAccountView'
import { ExtendedAccountData } from '~/shared/typings/Accounts/AccountsTypes'
import { AccountPageRouteProp } from '~/shared/typings/App/AppTypes'
import { useAccountId } from '~/widgets/TransactionsFeed/context/useAccountContext'

const AccountView: React.FC = () => {
  const { getAccountById } = useService()
  //   const navigation = useNavigation<StackNavigationProp<RootStackParamList>>()
  const route = useRoute<AccountPageRouteProp>()
  const { accountId } = route.params
  const { setAccountId } = useAccountId()

  //   useEffect(() => {
  //     if (!accountId) navigation.navigate('main')
  //   }, [accountId, navigation])

  useEffect(() => {
    if (accountId) {
      setAccountId(accountId)
    }
  }, [accountId, setAccountId])

  const { data } = useQuery({
    queryKey: ['account', accountId],
    queryFn: ({ signal }) => getAccountById(accountId, signal),
    enabled: !!accountId,
    retry: 2,
  })

  return <ExtendedAccountView {...(data as ExtendedAccountData)} />
}

export default AccountView
