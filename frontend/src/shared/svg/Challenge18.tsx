import { memo } from 'react'
import Svg, { Path } from 'react-native-svg'

import { ColorfulSvg } from '../typings/Common/common'

const Challenge = memo((props: ColorfulSvg) => {
  const { color } = props

  return (
    <Svg width='16' height='16' viewBox='0 0 16 16' fill='none'>
      <Path
        d='M12.0007 2.66797H4.00065C2.52789 2.66797 1.33398 3.86188 1.33398 5.33464V12.0013C1.33398 13.4741 2.52789 14.668 4.00065 14.668H12.0007C13.4734 14.668 14.6673 13.4741 14.6673 12.0013V5.33464C14.6673 3.86188 13.4734 2.66797 12.0007 2.66797Z'
        stroke={`${color}`}
        stroke-linecap='round'
        stroke-linejoin='round'
      />
      <Path
        d='M5.33398 1.33398V4.00065M10.6673 1.33398V4.00065M1.33398 6.66732H14.6673'
        stroke={`${color}`}
        stroke-linecap='round'
        stroke-linejoin='round'
      />
    </Svg>
  )
})

export default Challenge
