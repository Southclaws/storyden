import { Suspense } from "react";

import { cx } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";

import styles from "./navigation.module.css";

import { AskServer } from "./Anchors/AskServer";
import { SearchAnchor } from "./Anchors/Search";
import { MemberActionsServer } from "./MemberActionsServer";
import { SidebarToggleServer } from "./NavigationPane/SidebarToggleServer";
import { Title } from "./Title";

export function DesktopCommandBar() {
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
          <SidebarToggleServer />
        </Suspense>
        <SearchAnchor />
        <Suspense>
          <AskServer />
        </Suspense>
      </HStack>

      <HStack className={styles["topbar-middle"]} justify="space-around">
        <Suspense>
          <Title />
        </Suspense>
      </HStack>

      <HStack className={styles["topbar-right"]}>
        <Suspense>
          <MemberActionsServer />
        </Suspense>
      </HStack>
    </HStack>
  );
}
