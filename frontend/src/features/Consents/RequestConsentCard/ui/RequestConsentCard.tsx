/* eslint-disable react-hooks/refs */
import React, { memo, useEffect, useRef } from 'react'
import { Animated, Easing, Pressable } from 'react-native'
import Svg, { Line } from 'react-native-svg'

interface RequestConsentCardProps {
  onPress?: () => void
}

const RequestConsentCard: React.FC<RequestConsentCardProps> = ({ onPress }) => {
  const fadeAnim = useRef(new Animated.Value(0)).current

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

  return (
    <Pressable accessibilityRole='button' onPress={onPress}>
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
          justifyContent: 'center',
          alignItems: 'center',
        }}
      >
        <Svg width={48} height={48} viewBox='0 0 48 48'>
          <Line
            x1='24'
            y1='8'
            x2='24'
            y2='40'
            stroke='#D3D3D3'
            strokeWidth='4'
            strokeLinecap='round'
          />
          <Line
            x1='8'
            y1='24'
            x2='40'
            y2='24'
            stroke='#D3D3D3'
            strokeWidth='4'
            strokeLinecap='round'
          />
        </Svg>
      </Animated.View>
    </Pressable>
  )
}

export default memo(RequestConsentCard)
