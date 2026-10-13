import { memo } from 'react'
import Svg, { Path } from 'react-native-svg'

import { ColorfulSvg } from '../typings/Common/common'

const Article = memo((props: ColorfulSvg) => {
  const { color } = props

  return (
    <Svg width='16' height='16' viewBox='0 0 16 16' fill='none'>
      <Path
        d='M4.66667 5.33366H11.3333M4.66667 8.00033H11.3333M4.66667 10.667H11.3333M2 4.00033C2 3.6467 2.14048 3.30756 2.39052 3.05752C2.64057 2.80747 2.97971 2.66699 3.33333 2.66699H12.6667C13.0203 2.66699 13.3594 2.80747 13.6095 3.05752C13.8595 3.30756 14 3.6467 14 4.00033V12.0003C14 12.3539 13.8595 12.6931 13.6095 12.9431C13.3594 13.1932 13.0203 13.3337 12.6667 13.3337H3.33333C2.97971 13.3337 2.64057 13.1932 2.39052 12.9431C2.14048 12.6931 2 12.3539 2 12.0003V4.00033Z'
        stroke={`${color}`}
        stroke-linecap='round'
        stroke-linejoin='round'
      />
    </Svg>
  )
})

export default Article
