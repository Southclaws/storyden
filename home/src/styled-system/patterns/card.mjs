import { getPatternStyles, patternFns } from '../helpers.mjs';
import { css } from '../css/index.mjs';

const CardConfig = {
transform(props) {
  const { kind, display } = props;
  const padding = kind === "edge" ? "0" : "2";
  return {
    display,
    flexDirection: "column",
    gap: "1",
    width: "full",
    overflow: "hidden",
    boxShadow: "sm",
    borderRadius: "lg",
    backgroundColor: "bg.default",
    padding
  };
}}

export const getCardStyle = (styles = {}) => {
  const _styles = getPatternStyles(CardConfig, styles)
  return CardConfig.transform(_styles, patternFns)
}

export const Card = (styles) => css(getCardStyle(styles))
Card.raw = getCardStyle