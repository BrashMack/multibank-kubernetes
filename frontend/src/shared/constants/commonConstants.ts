import { ImageSourcePropType } from 'react-native'

import Logout from '~/shared/assets/images/common/Logout.png'
import QRLogo from '~/shared/assets/images/common/QRLogo.png'
import JCBLogo from '~/shared/assets/images/payment_systems/JCB.png'
import MasterCardLogo from '~/shared/assets/images/payment_systems/MasterCard.png'
import MIRLogo from '~/shared/assets/images/payment_systems/MIR.png'
import UnionPayLogo from '~/shared/assets/images/payment_systems/UnionPay.png'
import VisaLogo from '~/shared/assets/images/payment_systems/Visa.png'

export const DEFAULT_SKELETON_COUNT = 2
export const CARDS_COUNT_KEY = 'cards-count'
export const ACCOUNTS_COUNT_KEY = 'accounts-count'
export const DEPOSITS_COUNT_KEY = 'deposits-count'

export const CURRENCY_SYMBOL: Record<string, string> = {
  RUB: '₽',
  USD: '$',
  EUR: '€',
  GBP: '£',
}

export const PAYMENT_SYSTEM: Record<string, ImageSourcePropType> = {
  MasterCard: MasterCardLogo,
  Visa: VisaLogo,
  МИР: MIRLogo,
  JCB: JCBLogo,
  UnionPay: UnionPayLogo,
}

export const LOGO_SIZES: Record<string, LogoSize> = {
  МИР: { width: 35, height: 35, top: '$2.5' },
  Visa: { width: 33, height: 33, top: '$2.5' },
  MasterCard: { width: 28, height: 28, top: '$1.5' },
  UnionPay: { width: 29, height: 29, top: '$1.5' },
  JCB: { width: 28, height: 28, top: '$1.5' },
  defaultSize: { width: 30, height: 30, top: '$1.5' },
}

export interface LogoSize {
  width: number
  height: number
  top: '$1.5' | '$2.5'
}

export const INTERFACE_ICONS: Record<string, ImageSourcePropType> = {
  QRButton: QRLogo,
  LogoutButton: Logout,
}
