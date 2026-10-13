import React, { memo, useEffect, useMemo } from 'react'
import { Animated, Easing } from 'react-native'

interface SkeletonCardProps {
  source?: string
}

const SkeletonCard: React.FC<SkeletonCardProps> = ({ source }) => {
  const fadeAnim = useMemo(() => new Animated.Value(0), [])

  const params = useMemo(() => {
    const width = source === 'deposit' ? 135 : 180
    const height = source === 'deposit' ? 135 : 110
    const radius = source === 'deposit' ? 16 : 12

    return {
      width: width,
      height: height,
      radius: radius,
    }
  }, [source])

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
    <Animated.View
      style={{
        width: params.width,
        height: params.height,
        borderRadius: params.radius,
        margin: 8,
        backgroundColor,
        shadowColor: '#000',
        shadowOpacity: 0.3,
        shadowRadius: 6,
      }}
    />
  )
}

export default memo(SkeletonCard)
