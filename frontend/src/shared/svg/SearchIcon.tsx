import { memo } from 'react'
import Svg, { Path } from 'react-native-svg'

const SearchIcon = memo(() => {
  return (
    <Svg width='24' height='24' viewBox='0 0 24 24' fill='none'>
      <Path
        d='M10.5709 18.4289C14.9103 18.4289 18.428 14.9111 18.428 10.5717C18.428 6.23236 14.9103 2.7146 10.5709 2.7146C6.23149 2.7146 2.71373 6.23236 2.71373 10.5717C2.71373 14.9111 6.23149 18.4289 10.5709 18.4289Z'
        stroke='#393856'
        stroke-width='2'
        stroke-linecap='round'
        stroke-linejoin='round'
      />
      <Path
        d='M21.2852 21.2861L16.2852 16.2861'
        stroke='#393856'
        stroke-width='2'
        stroke-linecap='round'
        stroke-linejoin='round'
      />
    </Svg>
  )
})

export default SearchIcon
