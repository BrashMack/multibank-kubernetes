import { memo } from 'react'
import Svg, { Path } from 'react-native-svg'

const HomeFilled = memo(() => {
  return (
    <Svg width='24' height='24' viewBox='0 0 24 24' fill='none'>
      <Path
        fill-rule='evenodd'
        clip-rule='evenodd'
        d='M2.45455 10.6356C2.16463 10.9058 2 11.2843 2 11.6806V19.8571C2 21.0406 2.95939 22 4.14286 22H10.5714V17.7143C10.5714 16.9253 11.211 16.2857 12 16.2857C12.789 16.2857 13.4286 16.9253 13.4286 17.7143V22H19.8571C21.0406 22 22 21.0406 22 19.8571V11.6806C22 11.2843 21.8354 10.9058 21.5454 10.6356L12.4649 2.17196C12.1974 1.94268 11.8026 1.94268 11.5351 2.17196L2.45455 10.6356Z'
        fill='#393856'
      />
    </Svg>
  )
})

export default HomeFilled
