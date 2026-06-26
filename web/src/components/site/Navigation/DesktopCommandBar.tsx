import { Suspense } from "react";

import { getServerSession } from "@/auth/server-session";
import { hasCapability } from "@/lib/settings/capabilities";
import { allowsPublicRegistration } from "@/lib/settings/registration";
import { getSettings } from "@/lib/settings/settings-server";
import { cx } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";

import styles from "./navigation.module.css";

import { AskAnchor } from "./Anchors/Ask";
import { SearchAnchor } from "./Anchors/Search";
import { MemberActions } from "./MemberActions";
import { SidebarToggle } from "./NavigationPane/SidebarToggle";
import { getServerSidebarState } from "./NavigationPane/server";
import { Title } from "./Title";

export async function DesktopCommandBar() {
  const { title, capabilities, registration_mode } = await getSettings();

  const isSemdexEnabled = hasCapability("semdex", capabilities);
  const canRegister = allowsPublicRegistration(registration_mode);

  return (
    <HStack
      className={cx(Floating(), styles["topbar"])}
      borderRadius="md"
      justify="space-between"
      alignItems="center"
      px="1"
    >
      <HStack className={styles["topbar-left"]}>
        <Suspense>
          <SidebarToggleWithState />
        </Suspense>
        <SearchAnchor />
        {isSemdexEnabled && <AskAnchor />}
      </HStack>

      <HStack className={styles["topbar-middle"]} justify="space-around">
        <Title>{title}</Title>
      </HStack>

      <HStack className={styles["topbar-right"]}>
        <Suspense>
          <MemberActionsWithSession canRegister={canRegister} />
        </Suspense>
      </HStack>
    </HStack>
  );
}

async function SidebarToggleWithState() {
  const initialSidebarState = await getServerSidebarState();
  return <SidebarToggle initialValue={initialSidebarState} />;
}

async function MemberActionsWithSession({
  canRegister,
}: {
  canRegister: boolean;
}) {
  const session = await getServerSession();
  return <MemberActions session={session} canRegister={canRegister} />;
}
