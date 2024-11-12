import { getServerSession } from "@/auth/server-session";
import { getSettings } from "@/lib/settings/settings-server";
import { cx } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";

import { Search } from "../../search/Search/Search";

import styles from "./navigation.module.css";

import { MemberActions } from "./MemberActions";
import { SidebarToggle } from "./NavigationPane/SidebarToggle";
import { getServerSidebarState } from "./NavigationPane/server";
import { Title } from "./Title";

export async function DesktopCommandBar() {
  const { title } = await getSettings();
  const initialSidebarState = await getServerSidebarState();

  const session = await getServerSession();

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
        <Search />
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
