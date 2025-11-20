import { sva } from "@/styled-system/css";
import { styled } from "@/styled-system/jsx";

import styles from "./spinner.module.css";

const spinner = sva({
  slots: ["root"],
  base: {
    root: {},
  },
  variants: {
    size: {
      sm: {
        root: {
          w: "4",
          h: "4",
          borderWidth: "[2px]",
        },
      },
      md: {
        root: {
          w: "5",
          h: "5",
          borderWidth: "[3px]",
        },
      },
      lg: {
        root: {
          w: "6",
          h: "6",
          borderWidth: "[5px]",
        },
      },
    },
  },
  defaultVariants: {
    size: "md",
  },
});

type SpinnerProps = {
  size?: "sm" | "md" | "lg";
};

export function Spinner({ size }: SpinnerProps) {
  const classes = spinner({ size });

  return <styled.div className={`${styles["spinner"]} ${classes.root}`} />;
}
