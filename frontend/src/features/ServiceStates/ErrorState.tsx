import { Box, Button, Text } from '@gluestack-ui/themed'
import React from 'react'

export const ErrorState = ({ onRetry, message }: { onRetry: () => void; message: string }) => (
  <Box flex={1} alignItems='center' justifyContent='center' p='$4' w='$full'>
    <Text color='$white' fontWeight='$semibold' mb='$3'>
      {message}
    </Text>
    <Button backgroundColor='$black' borderRadius='$12' onPress={onRetry}>
      <Text color='$white'>Повторить попытку</Text>
    </Button>
  </Box>
)
