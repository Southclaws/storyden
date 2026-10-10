import { stableStringify, toHash } from '../helpers';

export const keyframes = (keyframe) => {
  const prefix = null
  const wrap = (base) => (prefix ? prefix + '-' + base : base)
  const block = keyframe && typeof keyframe === 'object' ? keyframe : {}
  return wrap('kf_' + toHash(stableStringify(block)))
}