import React, { PropsWithChildren } from "react";

import { getServerSession } from "@/auth/server-session";
import { parseMemberSettings } from "@/lib/settings/member-settings";
import { getSettings } from "@/lib/settings/settings-server";
import { Box } from "@/styled-system/jsx";

import { Onboarding } from "../Onboarding/Onboarding";

import styles from "./navigation.module.css";

import { ContextPane } from "./ContextPane";
import { DesktopCommandBar } from "./DesktopCommandBar";
import { MobileCommandBar } from "./MobileCommandBar/MobileCommandBar";
import { NavigationPane } from "./NavigationPane/NavigationPane";
import { getServerSidebarState } from "./NavigationPane/server";

type Props = {
  contextpane: React.ReactNode;
};

export async function Navigation({
  contextpane,
  children,
}: PropsWithChildren<Props>) {
  const globalSettings = await getSettings();
  const sessionAccount = await getServerSession();
  const session = sessionAccount
    ? parseMemberSettings(sessionAccount, globalSettings.metadata)
    : undefined;

  const sidebarDefaultState =
    session?.meta.sidebar.defaultState ?? globalSettings.metadata.sidebar.defaultState;
  const showLeftBar = await getServerSidebarState(sidebarDefaultState);

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
        <DesktopCommandBar />

        <Box className={styles["leftbar"]}>
          <NavigationPane />
        </Box>

        <Box className={styles["rightbar"]}>
          <ContextPane>{contextpane}</ContextPane>
        </Box>

        <Box className={styles["navpill"]}>
          <MobileCommandBar />
        </Box>
      </Box>
    </Box>
  );
}
