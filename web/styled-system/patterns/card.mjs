import { mapObject } from '../helpers.mjs';
import { css } from '../css/index.mjs';

const CardConfig = {
transform(props) {
  const { kind } = props;
  const padding = kind === "edge" ? "0" : "2";
  return {
    display: "flex",
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

export const getCardStyle = (styles = {}) => CardConfig.transform(styles, { map: mapObject })

export const Card = (styles) => css(getCardStyle(styles))
Card.raw = getCardStyle