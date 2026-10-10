import { getPatternStyles, patternFns } from './runtime';
import { memo } from '../helpers';
import { css } from '../css/index';

const linkButtonConfig = {transform:(props) => ({
	backgroundColor: "white",
	alignItems: "center",
	appearance: "none",
	borderRadius: "lg",
	boxShadow: "xs",
	cursor: "pointer",
	display: "inline-flex",
	fontWeight: "semibold",
	minWidth: "0",
	justifyContent: "center",
	outline: "none",
	position: "relative",
	transitionDuration: "normal",
	transitionProperty: "background, border-color, color, box-shadow",
	transitionTimingFunction: "default",
	userSelect: "none",
	verticalAlign: "middle",
	whiteSpace: "nowrap",
	_hover: {
		background: "gray.100",
		boxShadow: "md"
	},
	_focusVisible: {
		outlineOffset: "2px",
		outline: "2px solid",
		outlineColor: "border.outline"
	},
	_active: { backgroundColor: "gray.200" },
	h: "11",
	minW: "11",
	textStyle: "md",
	px: "5",
	gap: "2",
	"& svg": {
		width: "4",
		height: "4"
	},
	...props
})}

export function linkButtonRaw(styles) {
  const s = getPatternStyles(linkButtonConfig, styles || {})
  return linkButtonConfig.transform(s, patternFns)
}

export const linkButton = /* @__PURE__ */ Object.assign(/* @__PURE__ */ memo(function linkButton(styles = {}) {
  return css(linkButtonRaw(styles))
}), { raw: linkButtonRaw, propKeys: [] })