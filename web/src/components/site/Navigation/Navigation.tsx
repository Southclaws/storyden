"use client";

import { usePathname } from "next/navigation";
import { PropsWithChildren } from "react";

import { Navpill } from "src/components/site/Navigation/Navpill/Navpill";

import styles from "./navigation.module.css";

import { Box } from "@/styled-system/jsx";

import { Left } from "./Left/Left";
import { Top } from "./Top/Top";

const ROUTES_WITHOUT_NAVPILL = ["/new"];

const isNavpillShown = (path: string | null) =>
  ROUTES_WITHOUT_NAVPILL.includes(path ?? "");

export function Navigation({ children }: PropsWithChildren) {
  const pathname = usePathname();

  return (
    <>
      {/* MOBILE */}
      <Box
        id="mobile-nav-container"
        display={{
          base: isNavpillShown(pathname) ? "none" : "unset",
          lg: "none",
        }}
      >
        <Box p="3">
          {/*  */}
          {children}
          {/*  */}
        </Box>

        <Navpill />
      </Box>

      {/* DESKTOP */}
      <Box
        id="desktop-nav-container"
        display={{
          base: "none",
          lg: "block",
        }}
        w="full"
      >
        <Box id="navigation__scroll" className={styles["navgrid"]}>
          <Box className={styles["main"]}>
            {/*  */}
            {children}
            {/*  */}
          </Box>
        </Box>

        <Box
          id="navigation__fixed"
          position="fixed"
          zIndex="overlay"
          top="0"
          left="0"
          height="dvh"
          className={styles["navgrid"]}
          pointerEvents="none"
        >
          <Box className={styles["topbar"]}>
            <Top />
          </Box>

          <Box className={styles["leftbar"]} pb="2">
            <Left />
          </Box>

          <Box className={styles["rightbar"]}>
            {/* RIGHT BAR NOT DONE YET */}
            {/* <Right /> */}
          </Box>
        </Box>
      </Box>
    </>
  );
}
