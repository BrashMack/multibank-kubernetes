import { CameraView } from 'expo-camera'
import { useEffect, useRef } from 'react'
import { AppState, Platform, Pressable, StatusBar, StyleSheet, Text, View } from 'react-native'

export interface QRScannerProps {
  permitted?: boolean
  onClose?: () => void
  onCodeScan?: (data: string) => any
}

const QRScanner: React.FC<QRScannerProps> = ({
  permitted,
  onClose = () => {},
  onCodeScan = (data: string) => data,
}) => {
  const qrLock = useRef(false)
  const appState = useRef(AppState.currentState)

  useEffect(() => {
    const subscription = AppState.addEventListener('change', nextAppState => {
      if (appState.current.match(/inactive|background/) && nextAppState === 'active') {
        qrLock.current = false
      }
      appState.current = nextAppState
    })

    return () => {
      subscription.remove()
      qrLock.current = true
    }
  }, [])

  return (
    <View style={styles.container}>
      {(Platform.OS === 'android' || Platform.OS === 'ios' || Platform.OS === 'web') && (
        <StatusBar hidden />
      )}

      <Pressable style={styles.closeButton} onPress={onClose}>
        <Text style={styles.closeText}>X</Text>
      </Pressable>

      {!permitted ? (
        <View style={styles.permissionContainer}>
          <Text style={styles.permissionText}>
            Доступ к камере не предоставлен. Разрешите использование камеры в настройках.
          </Text>
        </View>
      ) : (
        <>
          <CameraView
            style={StyleSheet.absoluteFill}
            facing='back'
            onBarcodeScanned={({ data }) => {
              if (data && !qrLock.current) {
                qrLock.current = true
                setTimeout(() => {
                  try {
                    JSON.parse(data)
                    onCodeScan(data)
                  } catch {
                    return
                  }
                }, 300)
              }
            }}
          />

          {/* <QROverlay /> */}
        </>
      )}
    </View>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    zIndex: 9999,
    ...StyleSheet.absoluteFillObject,
    backgroundColor: 'black',
  },

  closeButton: {
    position: 'absolute',
    top: 40,
    right: 20,
    zIndex: 10000,
    width: 42,
    height: 42,
    borderRadius: 21,
    backgroundColor: 'rgba(0,0,0,0.6)',
    alignItems: 'center',
    justifyContent: 'center',
  },
  closeText: {
    color: 'white',
    fontSize: 32,
    marginTop: -3,
  },

  permissionContainer: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    paddingHorizontal: 20,
  },
  permissionText: {
    color: 'white',
    fontSize: 18,
    textAlign: 'center',
  },
})

export default QRScanner
