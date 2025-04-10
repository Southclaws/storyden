import { getPatternStyles, patternFns } from '../helpers.mjs';
import { css } from '../css/index.mjs';

const linkButtonConfig = {
transform:(props) => ({
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
  _active: {
    backgroundColor: "gray.200"
  },
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

export const getLinkButtonStyle = (styles = {}) => {
  const _styles = getPatternStyles(linkButtonConfig, styles)
  return linkButtonConfig.transform(_styles, patternFns)
}

export const linkButton = (styles) => css(getLinkButtonStyle(styles))
linkButton.raw = getLinkButtonStyle