import { Box, Image, Text } from '@gluestack-ui/themed'
import React, { memo, useCallback, useMemo } from 'react'
import { Pressable } from 'react-native'

import {
  CURRENCY_SYMBOL,
  LOGO_SIZES,
  LogoSize,
  PAYMENT_SYSTEM,
} from '~/shared/constants/commonConstants'

interface CardProps {
  cardId: string
  cardNumber: string
  cardName?: string
  balance?: string
  currency?: 'RUB' | 'USD' | 'EUR' | 'GBP'
  bank?: string
  backgroundImage: string
  type: 'debit' | 'credit'
  onPress?: (cardId: string) => void
}

const Card: React.FC<CardProps> = ({
  cardId,
  cardNumber,
  cardName,
  balance,
  currency,
  bank,
  backgroundImage,
  type,
  onPress,
}) => {
  const maskedNumber = useMemo(() => (cardNumber ? `•• ${cardNumber.slice(-4)}` : ''), [cardNumber])
  const symbol = useMemo(() => currency && (CURRENCY_SYMBOL[currency] ?? ''), [currency])
  const paymentSystemKey = useMemo(() => cardName?.split(' ')[0].split('-')[0] ?? '', [cardName])
  const paymentLogo = useMemo(
    () => paymentSystemKey && (PAYMENT_SYSTEM[paymentSystemKey] ?? ''),
    [paymentSystemKey],
  )
  const logoSize = useMemo<LogoSize>(
    () => LOGO_SIZES[paymentSystemKey] ?? LOGO_SIZES.defaultSize,
    [paymentSystemKey],
  )
  const handlePress = useCallback(() => onPress?.(cardId), [onPress, cardId])

  return (
    <Pressable accessibilityRole='button' onPress={handlePress}>
      <Box
        w={180}
        h={110}
        borderRadius={12}
        overflow='hidden'
        m='$2'
        shadowColor='#000'
        shadowOpacity={0.2}
        shadowRadius={4}
        bg='$primary500'
        display='flex'
        flexDirection='column'
      >
        <Image
          source={{ uri: backgroundImage }}
          alt='Card background'
          w='100%'
          h='100%'
          resizeMode='cover'
          position='absolute'
        />
        <Box
          flex={2}
          px='$3'
          py='$2'
          flexDirection='row'
          justifyContent='space-between'
          alignItems='flex-start'
          bg='rgba(0,0,0,0.25)'
        >
          <Text color='$white' fontWeight='$medium' fontSize='$sm'>
            {bank}
          </Text>
          <Text color='$white' fontWeight='$medium' fontSize='$sm'>
            {maskedNumber}
          </Text>
        </Box>
        <Box flex={1} p='$2' justifyContent='center' alignItems='center' bg='rgba(0,0,0,0.25)'>
          <Box flexDirection='row' justifyContent='center' alignItems='baseline'>
            <Text color='$white' fontWeight='$bold' fontSize='$md'>
              {balance}
            </Text>
            {symbol ? (
              <Text color='$white' fontWeight='$bold' fontSize='$md' ml='$1'>
                {symbol}
              </Text>
            ) : null}
          </Box>
        </Box>
        <Box
          flex={2}
          px='$3'
          py='$2'
          flexDirection='row'
          justifyContent='space-between'
          alignItems='flex-end'
          bg='rgba(0,0,0,0.25)'
        >
          <Image
            source={paymentLogo}
            alt='payment system'
            width={logoSize.width}
            height={logoSize.height}
            resizeMode='contain'
            position='relative'
            top={logoSize.top}
            tintColor='$lightgrey'
          />
          <Text color='$lightgrey' fontWeight='$medium' fontSize='$xs'>
            {type}
          </Text>
        </Box>
      </Box>
    </Pressable>
  )
}

export default memo(Card)
