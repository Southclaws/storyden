import React, { PropsWithChildren, Suspense } from "react";

import { Box } from "@/styled-system/jsx";

import { Onboarding } from "../Onboarding/Onboarding";

import styles from "./navigation.module.css";

import { ContextPane } from "./ContextPane";
import { DesktopCommandBar } from "./DesktopCommandBar";
import { MobileCommandBar } from "./MobileCommandBar/MobileCommandBar";
import { NavigationContainer } from "./NavigationContainer";
import { NavigationPane } from "./NavigationPane/NavigationPane";

type Props = {
  contextpane: React.ReactNode;
};

export async function Navigation({
  contextpane,
  children,
}: PropsWithChildren<Props>) {
  return (
    <Suspense
      fallback={
        <Box
          id="navigation__container"
          className={styles["navigation__container"]}
          w="full"
          data-leftbar-shown="false"
        >
          <NavigationContent contextpane={contextpane}>
            {children}
          </NavigationContent>
        </Box>
      }
    >
      <NavigationContainer>
        <NavigationContent contextpane={contextpane}>
          {children}
        </NavigationContent>
      </NavigationContainer>
    </Suspense>
  );
}

type NavigationContentProps = {
  contextpane: React.ReactNode;
  sidebarShown?: string;
};

async function NavigationContent({
  contextpane,
  children,
  sidebarShown,
}: PropsWithChildren<NavigationContentProps>) {
  return (
    <>
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
          <Suspense>
            <MobileCommandBar />
          </Suspense>
        </Box>
      </Box>
    </>
  );
}
