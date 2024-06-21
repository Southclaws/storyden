import { PropsWithChildren } from "react";

import { Box } from "@/styled-system/jsx";

import { Onboarding } from "../Onboarding/Onboarding";

import styles from "./navigation.module.css";

import { Left } from "./Left/Left";
import { Navpill } from "./Navpill/Navpill";
import { getServerSidebarState } from "./Sidebar/server";
import { Top } from "./Top/Top";

export async function Navigation({ children }: PropsWithChildren) {
  const showLeftBar = await getServerSidebarState();

  return (
    <Box
      id="navigation__container"
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
        <Top />

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
