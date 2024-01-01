"use client";

import { usePathname } from "next/navigation";
import { PropsWithChildren, useState } from "react";

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
  const [showLeftBar, setShowLeftBar] = useState(false);

  return (
    <>
      {/* MOBILE */}
      <Box
        id="mobile-nav-container"
        width="full"
        display={{
          base: isNavpillShown(pathname) ? "none" : "unset",
          md: "none",
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
        className={styles["desktop-nav-container"]}
        display={{
          base: "none",
          md: "block",
        }}
        w="full"
        data-leftbar-shown={showLeftBar}
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
          <Top sidebarState={showLeftBar} onToggleSidebar={setShowLeftBar} />

          <Box className={styles["leftbar"]}>
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
