import { stableStringify, toHash } from '../helpers';

export const positionTry = (options) => {
  const prefix = null
  const wrap = (base) => '--' + (prefix ? prefix + '-' + base : base)
  if (typeof options === 'string') {
    return wrap('pt_' + options)
  }
  const block = options && typeof options === 'object' ? options : {}
  return wrap('pt_' + toHash(stableStringify(block)))
}