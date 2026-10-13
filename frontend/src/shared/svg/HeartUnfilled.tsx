import { memo } from 'react'
import Svg, { ClipPath, Defs, G, Path, Rect } from 'react-native-svg'

const HeartUnfilled = memo(() => {
  return (
    <Svg width='25' height='24' viewBox='0 0 25 24' fill='none'>
      <G clip-path='url(#clip0_308_3926)'>
        <Path
          d='M12.5066 19.6902L4.6861 12.6064C0.435802 8.35614 6.68373 0.195582 12.5066 6.7977C18.3295 0.195582 24.5492 8.38448 20.3272 12.6064L12.5066 19.6902Z'
          stroke='black'
          stroke-opacity='0.5'
          stroke-width='1.42857'
          stroke-linecap='round'
          stroke-linejoin='round'
        />
      </G>
      <Defs>
        <ClipPath id='clip0_308_3926'>
          <Rect width='20' height='20' fill='white' transform='translate(2.5 2.00098)' />
        </ClipPath>
      </Defs>
    </Svg>
  )
})

export default HeartUnfilled
