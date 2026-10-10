/**
 * Conditionally join classNames into a single string
 */
export function cx(...args) {
  let str = '',
    i = 0,
    arg

  for (; i < arguments.length; ) {
    if ((arg = arguments[i++]) && typeof arg === 'string') {
      str && (str += ' ')
      str += arg
    }
  }
  return str
}