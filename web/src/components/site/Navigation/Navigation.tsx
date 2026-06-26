import React, { PropsWithChildren, Suspense } from "react";

import { getServerSession } from "@/auth/server-session";
import { parseMemberSettings } from "@/lib/settings/member-settings";
import { allowsPublicRegistration } from "@/lib/settings/registration";
import { getSettings } from "@/lib/settings/settings-server";
import { Box } from "@/styled-system/jsx";

import { CommandPalette } from "../CommandPalette/CommandPalette";
import { Onboarding } from "../Onboarding/Onboarding";
import { VerificationBanner } from "../VerificationBanner/VerificationBanner";

import styles from "./navigation.module.css";

import { ContextPane } from "./ContextPane";
import { DesktopCommandBar } from "./DesktopCommandBar";
import { MobileCommandBar } from "./MobileCommandBar/MobileCommandBar";
import { NavigationPane } from "./NavigationPane/NavigationPane";
import { getServerSidebarState } from "./NavigationPane/server";
import { type Settings } from "@/lib/settings/settings";

type Props = {
  contextpane: React.ReactNode;
};

export async function Navigation({
  contextpane,
  children,
}: PropsWithChildren<Props>) {
  const globalSettings = await getSettings();
  const canRegister = allowsPublicRegistration(
    globalSettings.registration_mode,
  );

  return (
    <Box
      id="navigation__container"
      className={styles["navigation__container"]}
      w="full"
    >
      <Box id="navigation__scroll" className={styles["navgrid"]}>
        <Box className={styles["main"]}>
          <Suspense>
            <NavigationBanners settings={globalSettings} />
          </Suspense>
          {children}
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
        <Suspense>
          <DesktopCommandBar />
        </Suspense>

        <Suspense>
          <NavigationLeftBar settings={globalSettings} />
        </Suspense>

        <Box
          id="navigation__rightbar"
          className={styles["rightbar"]}
        >
          <ContextPane>{contextpane}</ContextPane>
        </Box>

        <Box className={styles["navpill"]}>
          <MobileCommandBar canRegister={canRegister} />
        </Box>
      </Box>

      <CommandPalette />
    </Box>
  );
}

async function NavigationBanners({ settings }: { settings: Settings }) {
  const session = await getServerSession();
  return (
    <>
      <Onboarding />
      <VerificationBanner session={session} settings={settings} />
    </>
  );
}

async function NavigationLeftBar({ settings }: { settings: Settings }) {
  const session = await getServerSession();
  const memberSettings = session
    ? parseMemberSettings(session, settings.metadata)
    : undefined;
  const sidebarDefaultState =
    memberSettings?.meta.sidebar.defaultState ??
    settings.metadata.sidebar.defaultState;
  const showLeftBar = await getServerSidebarState(sidebarDefaultState);

  return (
    <Box
      id="navigation__leftbar"
      className={styles["leftbar"]}
      aria-hidden={!showLeftBar}
      inert={!showLeftBar}
    >
      <NavigationPane
        initialSession={session}
        initialSettings={settings}
      />
    </Box>
  );
}
