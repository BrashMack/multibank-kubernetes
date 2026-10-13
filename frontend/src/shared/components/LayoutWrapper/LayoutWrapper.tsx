import { Box } from '@gluestack-ui/themed'
import React, { memo } from 'react'

const LayoutWrapper: React.FC<React.PropsWithChildren<unknown>> = ({ children }) => {
  return (
    <Box pt='$12' pb='$8' flex={1} bgColor='$backgroundDark950'>
      {children}
    </Box>
  )
}

export default memo(LayoutWrapper)
