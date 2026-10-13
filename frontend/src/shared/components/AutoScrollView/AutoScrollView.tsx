import React, { forwardRef, useCallback, useImperativeHandle, useRef, type Ref } from 'react'
import type { LayoutChangeEvent, ScrollViewProps } from 'react-native'
import { ScrollView } from 'react-native'

const AutoScrollView = forwardRef<ScrollView, ScrollViewProps>(
  ({ children, onContentSizeChange, onLayout, ...rest }, externalRef: Ref<ScrollView>) => {
    const scrollViewRef = useRef<ScrollView>(null)
    const contentHeightRef = useRef(0)
    const containerHeightRef = useRef(0)

    useImperativeHandle(externalRef, () => scrollViewRef.current as ScrollView)

    const handleLayout = useCallback(
      (event: LayoutChangeEvent) => {
        containerHeightRef.current = event.nativeEvent.layout.height
        onLayout?.(event)
      },
      [onLayout],
    )

    const scrollToBottom = useCallback(() => {
      scrollViewRef.current?.scrollToEnd({ animated: true })
    }, [])

    const handleContentSizeChange = useCallback<Required<ScrollViewProps>['onContentSizeChange']>(
      (contentWidth, contentHeight) => {
        const shouldAutoScroll =
          contentHeight > containerHeightRef.current && contentHeight > contentHeightRef.current

        contentHeightRef.current = contentHeight

        if (shouldAutoScroll) {
          setTimeout(scrollToBottom, 0)
        }

        onContentSizeChange?.(contentWidth, contentHeight)
      },
      [onContentSizeChange, scrollToBottom],
    )

    return (
      <ScrollView
        ref={scrollViewRef}
        onLayout={handleLayout}
        onContentSizeChange={handleContentSizeChange}
        {...rest}
      >
        {children}
      </ScrollView>
    )
  },
)

AutoScrollView.displayName = 'AutoScrollView'

export default AutoScrollView
