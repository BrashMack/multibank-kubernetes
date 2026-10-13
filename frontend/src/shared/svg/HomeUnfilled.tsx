import { memo } from 'react'
import Svg, { Path } from 'react-native-svg'

const HomeUnfilled = memo(() => {
  return (
    <Svg width='24' height='24' viewBox='0 0 24 24' fill='none'>
      <Path
        d='M21.5714 12.2C21.5729 12.0017 21.5329 11.8053 21.4541 11.6233C21.3754 11.4413 21.2597 11.2777 21.1143 11.1429L12.2857 3L3.45718 11.1429C3.31178 11.2777 3.19605 11.4413 3.11735 11.6233C3.03864 11.8053 2.99869 12.0017 3.00003 12.2V20.1429C3.00003 20.5217 3.15054 20.8851 3.41845 21.153C3.68636 21.4209 4.04973 21.5714 4.4286 21.5714H20.1429C20.5217 21.5714 20.8851 21.4209 21.153 21.153C21.421 20.8851 21.5714 20.5217 21.5714 20.1429V12.2Z'
        stroke='#908C93'
        stroke-width='1.5'
        stroke-linecap='round'
        stroke-linejoin='round'
      />
      <Path
        d='M12.2857 21.5715V15.8572'
        stroke='#908C93'
        stroke-width='1.5'
        stroke-linecap='round'
        stroke-linejoin='round'
      />
    </Svg>
  )
})

export default HomeUnfilled
