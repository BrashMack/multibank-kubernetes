import { memo } from 'react'
import Svg, { ClipPath, Defs, G, Path, Rect } from 'react-native-svg'

const RingActive = memo(() => {
  return (
    <Svg width='24' height='24' viewBox='0 0 24 24' fill='none'>
      <G clip-path='url(#clip0_308_3335)'>
        <Path
          d='M17.0003 11C19.2095 11 21.0003 9.20914 21.0003 7C21.0003 4.79086 19.2095 3 17.0003 3C14.7912 3 13.0003 4.79086 13.0003 7C13.0003 9.20914 14.7912 11 17.0003 11Z'
          fill='#A05035'
        />
        <Path
          d='M16.1894 13.429V15.6198C16.1894 16.1785 16.4113 16.7144 16.8065 17.1095C17.2015 17.5045 17.8698 17.7265 18.4285 17.7265H2.7142C3.27295 17.7265 3.94117 17.5045 4.33627 17.1095C4.73136 16.7144 4.95333 16.1785 4.95333 15.6198L4.95332 10.4756C4.95332 8.98559 5.54522 7.55662 6.5988 6.50303C7.65239 5.44944 9.08136 4.85754 10.5713 4.85754'
          stroke='#1B1D20'
          stroke-width='1.5'
          stroke-linecap='round'
          stroke-linejoin='round'
        />
        <Path
          d='M9.14276 21.2861H11.9999'
          stroke='#1B1D20'
          stroke-width='1.5'
          stroke-linecap='round'
          stroke-linejoin='round'
        />
      </G>
      <Defs>
        <ClipPath id='clip0_308_3335'>
          <Rect width='20' height='20' fill='white' transform='translate(2 2)' />
        </ClipPath>
      </Defs>
    </Svg>
  )
})

export default RingActive
