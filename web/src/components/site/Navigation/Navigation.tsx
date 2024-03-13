"use client";

import { PropsWithChildren } from "react";
import { useLocalStorage } from "usehooks-ts";

import { NAVIGATION_SIDEBAR_STATE_KEY } from "src/local/state-keys";

import styles from "./navigation.module.css";

import { Onboarding } from "../Onboarding/Onboarding";

import { Box } from "@/styled-system/jsx";

import { Left } from "./Left/Left";
import { Navpill } from "./Navpill/Navpill";
import { Top } from "./Top/Top";

export function Navigation({ children }: PropsWithChildren) {
  const [showLeftBar, setShowLeftBar] = useLocalStorage(
    NAVIGATION_SIDEBAR_STATE_KEY,
    false,
    {
      initializeWithValue: false,
    },
  );

  return (
    <Box
      className={styles["navigation__container"]}
      w="full"
      data-leftbar-shown={showLeftBar}
    >
      <Box id="navigation__scroll" className={styles["navgrid"]}>
        <Box className={styles["main"]}>
          {/*  */}
          <Onboarding />
          {children}
          {/*  */}
        </Box>
      </Box>

      <Box
        id="navigation__fixed"
        position="fixed"
        zIndex="docked"
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

        <Box className={styles["navpill"]}>
          <Navpill />
        </Box>
      </Box>
    </Box>
  );
}
