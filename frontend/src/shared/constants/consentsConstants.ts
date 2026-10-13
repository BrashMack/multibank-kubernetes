import AwesomeBankLogo from '~/shared/assets/images/banks/awesome-bank.png'
import SmartBankLogo from '~/shared/assets/images/banks/smart-bank.png'
import VirtualBankLogo from '~/shared/assets/images/banks/virtual-bank.png'

export const BANKS = {
  VIRTUAL_BANK: 'VirtualBank',
  SMART_BANK: 'SmartBank',
  AWESOME_BANK: 'AwesomeBank',
}

export const CONSENTS_TYPES = {
  READ: 'read',
}

export const AWESOME_BANK = {
  label: 'AwesomeBank',
  value: BANKS.AWESOME_BANK,
  logo: AwesomeBankLogo,
}

export const SMART_BANK = {
  label: 'SmartBank',
  value: BANKS.SMART_BANK,
  logo: SmartBankLogo,
}

export const VIRTUAL_BANK = {
  label: 'VirtualBank',
  value: BANKS.VIRTUAL_BANK,
  logo: VirtualBankLogo,
}

export const BANK_CHOOSE = [AWESOME_BANK, SMART_BANK, VIRTUAL_BANK]
