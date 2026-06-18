import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const cqConfig = {transform(props) {
	const { name, type, ...rest } = props;
	return {
		containerType: type,
		containerName: name,
		...rest
	};
},defaultValues:{type:'inline-size'}}

export function cqRaw(styles) {
  const s = getPatternStyles(cqConfig, styles || {})
  return cqConfig.transform(s, patternFns)
}

export const cq = /* @__PURE__ */ Object.assign(function cq(styles = {}) {
  return css(cqRaw(styles))
}, { raw: cqRaw })