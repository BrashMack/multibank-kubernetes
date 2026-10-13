let onUnauthorizedHandler: null | (() => void) = null

export const registerOnUnauthorized = (handler: () => void) => {
  onUnauthorizedHandler = handler
}

export const clearProfileIfUnauthorized = (status: number) => {
  if (status === 401 && onUnauthorizedHandler) {
    onUnauthorizedHandler()
  }
}
