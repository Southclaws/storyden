import { getPatternStyles, patternFns } from './runtime';
import { css } from '../css/index';

const CardBoxConfig = {transform(props) {
	const { kind, display, ...rest } = props;
	return {
		display,
		flexDirection: "column",
		gap: "1",
		width: "full",
		boxShadow: "sm",
		borderRadius: "lg",
		backgroundColor: "bg.default",
		padding: kind === "edge" ? "0" : "2",
		...rest
	};
}}

export function CardBoxRaw(styles) {
  const s = getPatternStyles(CardBoxConfig, styles || {})
  return CardBoxConfig.transform(s, patternFns)
}

export const CardBox = /* @__PURE__ */ Object.assign(function CardBox(styles = {}) {
  return css(CardBoxRaw(styles))
}, { raw: CardBoxRaw })