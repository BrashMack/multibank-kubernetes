import { Box, Image, Text } from '@gluestack-ui/themed'
import React, { memo, useCallback, useEffect, useMemo } from 'react'
import { Animated, Easing, ImageSourcePropType, Linking, Pressable } from 'react-native'

import { AWESOME_BANK, SMART_BANK, VIRTUAL_BANK } from '~/shared/constants/consentsConstants'
import { ConsentCardData } from '~/shared/typings/Consents/WaitConsentCardTypes'

const AnimatedText = Animated.createAnimatedComponent(Text)
const AnimatedImage = Animated.createAnimatedComponent(Image)

interface WaitConsentCardProps {
  data: ConsentCardData
}

const BANK_LINKS: Record<string, string> = {
  [AWESOME_BANK.label]: 'https://abank.open.bankingapi.ru/client',
  [SMART_BANK.label]: 'https://sbank.open.bankingapi.ru/client',
  [VIRTUAL_BANK.label]: 'https://vbank.open.bankingapi.ru/client',
}

const WaitConsentCard: React.FC<WaitConsentCardProps> = ({ data }) => {
  const fadeAnim = useMemo(() => new Animated.Value(0), [])
  const bank = String(data)

  const url = useMemo(() => BANK_LINKS[bank], [bank])

  const bankLogo = useMemo<ImageSourcePropType>(() => {
    switch (bank) {
      case AWESOME_BANK.label:
        return AWESOME_BANK.logo
      case SMART_BANK.label:
        return SMART_BANK.logo
      case VIRTUAL_BANK.label:
        return VIRTUAL_BANK.logo
      default:
        return VIRTUAL_BANK.logo
    }
  }, [bank])

  useEffect(() => {
    const loop = Animated.loop(
      Animated.sequence([
        Animated.timing(fadeAnim, {
          toValue: 1,
          duration: 2000,
          easing: Easing.inOut(Easing.sin),
          useNativeDriver: false,
        }),
        Animated.timing(fadeAnim, {
          toValue: 0,
          duration: 2000,
          easing: Easing.inOut(Easing.sin),
          useNativeDriver: false,
        }),
      ]),
    )
    loop.start()
    return () => loop.stop()
  }, [fadeAnim])

  const backgroundColor = fadeAnim.interpolate({
    inputRange: [0, 1],
    outputRange: ['#1f1f1f', '#4b4b4b'],
  })

  const textColor = fadeAnim.interpolate({
    inputRange: [0, 1],
    outputRange: ['#555', '#999999ff'],
  })

  const logoOpacity = fadeAnim.interpolate({
    inputRange: [0, 1],
    outputRange: [0.6, 1],
  })

  const handlePress = useCallback(() => {
    if (url) Linking.openURL(url)
  }, [url])

  return (
    <Pressable onPress={handlePress}>
      <Animated.View
        style={{
          width: 180,
          height: 110,
          borderRadius: 12,
          margin: 8,
          backgroundColor,
          shadowColor: '#000',
          shadowOpacity: 0.3,
          shadowRadius: 6,
          overflow: 'hidden',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <Box position='absolute' top='$2' left='$3'>
          <AnimatedText fontWeight='$medium' fontSize='$sm' style={{ color: textColor }}>
            {bank}
          </AnimatedText>
        </Box>
        <AnimatedImage
          source={bankLogo}
          alt='bank logo'
          top='$1'
          width={64}
          height={64}
          resizeMode='contain'
          tintColor='$lightgrey'
          style={{ opacity: logoOpacity }}
        />
      </Animated.View>
    </Pressable>
  )
}

export default memo(WaitConsentCard)
