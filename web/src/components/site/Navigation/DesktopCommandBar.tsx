import { getServerSession } from "@/auth/server-session";
import { hasCapability } from "@/lib/settings/capabilities";
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
  const { title, capabilities } = await getSettings();
  const initialSidebarState = await getServerSidebarState();

  const session = await getServerSession();

  const isSemdexEnabled = hasCapability("semdex", capabilities);

  return (
    <HStack
      className={cx(Floating(), styles["topbar"])}
      borderRadius="md"
      justify="space-between"
      alignItems="center"
      px="1"
    >
      <HStack className={styles["topbar-left"]}>
        <SidebarToggle initialValue={initialSidebarState} />
        <SearchAnchor />
        {isSemdexEnabled && <AskAnchor />}
      </HStack>

      <HStack className={styles["topbar-middle"]} justify="space-around">
        <Title>{title}</Title>
      </HStack>

      <HStack className={styles["topbar-right"]}>
        <MemberActions session={session} />
      </HStack>
    </HStack>
  );
}
