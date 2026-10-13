import { memo } from 'react'
import Svg, { Path } from 'react-native-svg'

const NavBurger = memo(() => {
  return (
    <Svg width='24' height='24' viewBox='0 0 24 24' fill='none'>
      <Path
        d='M4 6H20M4 12H20M4 18H12'
        stroke='#1B1D20'
        stroke-width='2'
        stroke-linecap='round'
        stroke-linejoin='round'
      />
    </Svg>
  )
})

export default NavBurger
