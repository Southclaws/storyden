import { PropsWithChildren } from "react";

import { Navigation } from "src/components/Navigation/Navigation";

import { css } from "@/styled-system/css";

export function Default(props: PropsWithChildren) {
  return (
    <div
      className={css({
        display: "flex",
        minHeight: "100vh",
        width: "full",
        flexDirection: "row",
        bgColor: "white",
      })}
      vaul-drawer-wrapper=""
    >
      <Navigation />

      <main
        className={css({
          width: "full",
          maxW: {
            base: "full",
            lg: "container.md",
          },
          px: 4,
          py: 2,
          backgroundColor: "white",
        })}
      >
        {props.children}
        <div
          className={css({
            height: "6rem",
          })}
        ></div>
      </main>
    </div>
  );
}
