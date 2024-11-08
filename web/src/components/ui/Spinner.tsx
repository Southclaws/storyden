import { Styles, css, cx } from "@/styled-system/css";

import styles from "./spinner.module.css";

export function Spinner(props: Styles) {
  const styleProps = css(props);
  return <div className={cx(styleProps, styles["spinner"])} />;
}
