import { Canvas, DiffRect, rect, rrect } from '@shopify/react-native-skia'
import { memo, useEffect, useState } from 'react'
import { Dimensions, Platform, StyleSheet, View } from 'react-native'

const { width, height } = Dimensions.get('window')
const innerDimension = 300
const innerLeft = width / 2 - innerDimension / 2
const innerTop = height / 2 - innerDimension / 2

const QROverlay = () => {
  const isWeb = Platform.OS === 'web'
  const [shapes, setShapes] = useState<{ outer: any; inner: any } | null>(null)

  useEffect(() => {
    if (!isWeb) {
      const outer = rrect(rect(0, 0, width, height), 0, 0)
      const inner = rrect(rect(innerLeft, innerTop, innerDimension, innerDimension), 50, 50)
      ;(() => setShapes({ outer, inner }))()
    }
  }, [isWeb])

  if (isWeb) {
    return (
      <View style={[StyleSheet.absoluteFillObject, { backgroundColor: 'rgba(0,0,0,0.5)' }]}>
        <View
          style={{
            position: 'absolute',
            top: innerTop,
            left: innerLeft,
            width: innerDimension,
            height: innerDimension,
            backgroundColor: 'transparent',
            borderRadius: 50,
            borderWidth: 3,
            borderColor: '#00FFAA',
          }}
        />
      </View>
    )
  }

  if (!shapes) return null // ждём mount

  return (
    <Canvas style={StyleSheet.absoluteFillObject}>
      <DiffRect inner={shapes.inner} outer={shapes.outer} color='black' opacity={0.5} />
    </Canvas>
  )
}

export default memo(QROverlay)
