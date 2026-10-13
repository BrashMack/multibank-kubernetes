export const isEmptyObject = (object: object) => {
  const arrayKeys = Object.keys(object)
  if (arrayKeys.length) {
    return false
  }

  return true
}
