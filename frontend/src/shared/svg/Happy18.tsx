import { memo } from 'react'
import Svg, { ClipPath, Defs, G, Path, Rect } from 'react-native-svg'

import { ColorfulSvg } from '../typings/Common/common'

const Happy = memo((props: ColorfulSvg) => {
  const { color } = props

  return (
    <Svg width='18' height='18' viewBox='0 0 18 18' fill='none'>
      <G clip-path='url(#clip0_952_13547)'>
        <Path
          d='M9.00021 17.3569C13.6158 17.3569 17.3574 13.6153 17.3574 8.99972C17.3574 4.3842 13.6158 0.642578 9.00021 0.642578C4.38469 0.642578 0.643066 4.3842 0.643066 8.99972C0.643066 13.6153 4.38469 17.3569 9.00021 17.3569Z'
          stroke={color}
          stroke-linecap='round'
          stroke-linejoin='round'
        />
        <Path
          d='M4.75732 10.2852C5.40018 12.5994 7.97161 14.0137 10.2859 13.3709C11.7002 12.8566 12.8573 11.6994 13.2431 10.2852'
          stroke={color}
          stroke-linecap='round'
          stroke-linejoin='round'
        />
        <Path
          d='M6.10756 7.07157C5.93004 7.07157 5.78613 6.92766 5.78613 6.75014C5.78613 6.57262 5.93004 6.42871 6.10756 6.42871'
          stroke={color}
          stroke-linecap='round'
          stroke-linejoin='round'
        />
        <Path
          d='M6.10742 7.07157C6.28494 7.07157 6.42885 6.92766 6.42885 6.75014C6.42885 6.57262 6.28494 6.42871 6.10742 6.42871'
          stroke={color}
          stroke-linecap='round'
          stroke-linejoin='round'
        />
        <Path
          d='M11.8932 7.07157C11.7157 7.07157 11.5718 6.92766 11.5718 6.75014C11.5718 6.57262 11.7157 6.42871 11.8932 6.42871'
          stroke={color}
          stroke-linecap='round'
          stroke-linejoin='round'
        />
        <Path
          d='M11.8931 7.07157C12.0706 7.07157 12.2145 6.92766 12.2145 6.75014C12.2145 6.57262 12.0706 6.42871 11.8931 6.42871'
          stroke={color}
          stroke-linecap='round'
          stroke-linejoin='round'
        />
      </G>
      <Defs>
        <ClipPath id='clip0_952_13547'>
          <Rect width='18' height='18' fill={color} />
        </ClipPath>
      </Defs>
    </Svg>
  )
})

export default Happy
