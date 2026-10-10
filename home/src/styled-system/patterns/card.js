import { getPatternStyles, patternFns } from './runtime';
import { memo } from '../helpers';
import { css } from '../css/index';

const cardConfig = {transform(props) {
	const { kind, display } = props;
	return {
		display,
		flexDirection: "column",
		gap: "1",
		width: "full",
		overflow: "hidden",
		boxShadow: "sm",
		borderRadius: "lg",
		backgroundColor: "bg.default",
		padding: kind === "edge" ? "0" : "2"
	};
}}

export function cardRaw(styles) {
  const s = getPatternStyles(cardConfig, styles || {})
  return cardConfig.transform(s, patternFns)
}

export const card = /* @__PURE__ */ Object.assign(/* @__PURE__ */ memo(function card(styles = {}) {
  return css(cardRaw(styles))
}), { raw: cardRaw, propKeys: ["display","kind"] })