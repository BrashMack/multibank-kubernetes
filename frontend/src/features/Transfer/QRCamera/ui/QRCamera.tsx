import { Image } from '@gluestack-ui/themed'
import { useNavigation } from '@react-navigation/native'
import { StackNavigationProp } from '@react-navigation/stack'
import { useCameraPermissions } from 'expo-camera'
import { memo, useCallback, useState } from 'react'
import { Modal, StyleSheet, Text, TouchableOpacity, View } from 'react-native'

import QRScanner from '~/shared/components/QRScanner/QRScanner'
import { INTERFACE_ICONS } from '~/shared/constants/commonConstants'
import { RootStackParamList } from '~/shared/typings/App/AppTypes'

const QRCamera: React.FC = () => {
  const [permission, requestPermission] = useCameraPermissions()
  const [showScanner, setShowScanner] = useState<boolean>(false)

  const navigation = useNavigation<StackNavigationProp<RootStackParamList>>()

  const isPermissionGranted = Boolean(permission?.granted)

  const handleOpenScanner = async () => {
    if (!isPermissionGranted) {
      const res = await requestPermission()
      if (!res.granted) return
    }
    setShowScanner(true)
  }

  const onScannerClose = useCallback(() => setShowScanner(false), [])

  const onCodeScan = useCallback(
    (data: string) => {
      onScannerClose()
      navigation.navigate('qrPayment', { ...JSON.parse(data) })
    },
    [navigation, onScannerClose],
  )

  return (
    <>
      <View style={styles.container}>
        {!isPermissionGranted && <Text style={styles.warning}>Нет доступа к камере</Text>}

        <TouchableOpacity onPress={handleOpenScanner}>
          <Image
            source={INTERFACE_ICONS.QRButton}
            alt='QR scaner'
            width={50}
            height={50}
            resizeMode='contain'
          />
        </TouchableOpacity>
      </View>

      {showScanner && isPermissionGranted && (
        <Modal
          visible
          transparent={false}
          animationType='fade'
          onRequestClose={onScannerClose}
          statusBarTranslucent
        >
          <QRScanner onClose={onScannerClose} onCodeScan={onCodeScan} permitted />
        </Modal>
      )}
    </>
  )
}

const styles = StyleSheet.create({
  container: {
    alignItems: 'flex-start',
    justifyContent: 'center',
  },
  warning: {
    color: 'white',
    marginBottom: 12,
    fontSize: 16,
  },
})

export default memo(QRCamera)
