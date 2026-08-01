import React, { PropsWithChildren } from "react";

import { getServerSession } from "@/auth/server-session";
import { allowsPublicRegistration } from "@/lib/settings/registration";
import { getSettings } from "@/lib/settings/settings-server";

import { CommandPalette } from "../CommandPalette/CommandPalette";
import { Onboarding } from "../Onboarding/Onboarding";
import { VerificationBanner } from "../VerificationBanner/VerificationBanner";

import { DesktopCommandBar } from "./DesktopCommandBar";
import { NavigationChrome } from "./NavigationChrome";
import { NavigationPane } from "./NavigationPane/NavigationPane";

export async function Navigation({ children }: PropsWithChildren) {
  const globalSettings = await getSettings();
  const canRegister = allowsPublicRegistration(
    globalSettings.registration_mode,
  );
  const sessionAccount = await getServerSession();

  return (
    <div id="navigation__container" className="navigation__container">
      <div id="navigation__scroll" className="navigation__grid">
        <div className="navigation__main">
          {/*  */}
          <Onboarding />
          <VerificationBanner
            session={sessionAccount}
            settings={globalSettings}
          />
          {children}
          {/*  */}
        </div>
      </div>

      <div id="navigation__fixed" className="navigation__fixed">
        <DesktopCommandBar />

        <NavigationChrome canRegister={canRegister}>
          <NavigationPane
            initialSession={sessionAccount}
            initialSettings={globalSettings}
          />
        </NavigationChrome>
      </div>

      <CommandPalette />
    </div>
  );
}
