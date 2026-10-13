import { memo } from 'react'
import Svg, { Ellipse } from 'react-native-svg'

const MyPathIcon = memo(() => {
  return (
    <Svg width='30' height='34' viewBox='0 0 30 34' fill='none'>
      <Ellipse cx='15.0007' cy='3.44523' rx='3.48647' ry='3.44523' fill='#E4B237' />
      <Ellipse cx='3.60951' cy='9.98161' rx='3.46666' ry='3.44523' fill='#749978' />
      <Ellipse cx='3.60951' cy='23.6044' rx='3.46666' ry='3.44523' fill='#9D89FF' />
      <Ellipse
        cx='3.48647'
        cy='3.44523'
        rx='3.48647'
        ry='3.44523'
        transform='matrix(1 0 0 -1 11.5142 33.4285)'
        fill='#9A5371'
      />
      <Ellipse cx='26.3903' cy='23.6044' rx='3.46666' ry='3.44523' fill='#D15A46' />
      <Ellipse cx='26.3903' cy='9.98161' rx='3.46666' ry='3.44523' fill='#FFBE81' />
    </Svg>
  )
})

export default MyPathIcon
