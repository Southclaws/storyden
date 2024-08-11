import { getPatternStyles, patternFns } from '../helpers.mjs';
import { css } from '../css/index.mjs';

const CardBoxConfig = {
transform(props28) {
  const { kind, display } = props28;
  const padding = kind === "edge" ? "0" : "2";
  return {
    display,
    flexDirection: "column",
    gap: "1",
    width: "full",
    boxShadow: "sm",
    borderRadius: "lg",
    backgroundColor: "bg.default",
    padding
  };
}}

export const getCardBoxStyle = (styles = {}) => {
  const _styles = getPatternStyles(CardBoxConfig, styles)
  return CardBoxConfig.transform(_styles, patternFns)
}

export const CardBox = (styles) => css(getCardBoxStyle(styles))
CardBox.raw = getCardBoxStyle