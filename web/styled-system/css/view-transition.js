import { stableStringify, toHash } from '../helpers';

export const viewTransition = (options) => {
  const prefix = null
  if (typeof options === 'string') {
    const base = 'vt_' + options
    return prefix ? prefix + '-' + base : base
  }
  const slots = ['group', 'imagePair', 'old', 'new']
  const filtered = {}
  if (options && typeof options === 'object') {
    for (const key of slots) {
      if (key in options) filtered[key] = options[key]
    }
  }
  const base = 'vt_' + toHash(stableStringify(filtered))
  return prefix ? prefix + '-' + base : base
}