import React, { PropsWithChildren } from "react";

import { getServerSession } from "@/auth/server-session";
import { allowsPublicRegistration } from "@/lib/settings/registration";
import { getSettings } from "@/lib/settings/settings-server";
import { Box } from "@/styled-system/jsx";

import { CommandPalette } from "../CommandPalette/CommandPalette";
import { Onboarding } from "../Onboarding/Onboarding";
import { VerificationBanner } from "../VerificationBanner/VerificationBanner";

import { DesktopCommandBar } from "./DesktopCommandBar";
import { MobileCommandBar } from "./MobileCommandBar/MobileCommandBar";
import { NavigationPane } from "./NavigationPane/NavigationPane";

export async function Navigation({ children }: PropsWithChildren) {
  const globalSettings = await getSettings();
  const canRegister = allowsPublicRegistration(
    globalSettings.registration_mode,
  );
  const sessionAccount = await getServerSession();

  return (
    <Box id="navigation__container" className="navigation__container" w="full">
      <Box id="navigation__scroll" className="navigation__grid">
        <Box className="navigation__main">
          {/*  */}
          <Onboarding />
          <VerificationBanner
            session={sessionAccount}
            settings={globalSettings}
          />
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
        right="0"
        height="dvh"
        pointerEvents="none"
      >
        <DesktopCommandBar />

        <Box className="navigation__grid navigation__fixed-grid">
          <Box id="navigation__leftbar" className="navigation__leftbar">
            <NavigationPane
              initialSession={sessionAccount}
              initialSettings={globalSettings}
            />
          </Box>

          <Box className="navigation__navpill">
            <MobileCommandBar canRegister={canRegister} />
          </Box>
        </Box>
      </Box>

      <CommandPalette />
    </Box>
  );
}
