import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const aspectRatioConfig = {transform(props, { map }) {
	const { ratio = 4 / 3, ...rest } = props;
	return {
		position: "relative",
		_before: {
			content: `""`,
			display: "block",
			height: "0",
			paddingBottom: map(ratio, (r) => `${1 / r * 100}%`)
		},
		"&>*": {
			display: "flex",
			justifyContent: "center",
			alignItems: "center",
			overflow: "hidden",
			position: "absolute",
			inset: "0",
			width: "100%",
			height: "100%"
		},
		"&>img, &>video": { objectFit: "cover" },
		...rest
	};
}}

export function aspectRatioRaw(styles) {
  const s = getPatternStyles(aspectRatioConfig, styles || {})
  return aspectRatioConfig.transform(s, patternFns)
}

export const aspectRatio = /* @__PURE__ */ Object.assign(function aspectRatio(styles = {}) {
  return css(aspectRatioRaw(styles))
}, { raw: aspectRatioRaw })