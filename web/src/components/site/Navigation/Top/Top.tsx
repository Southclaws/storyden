import { getServerSession } from "@/auth/server-session";
import { cx } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";
import { getInfo } from "@/utils/info";

import { Search } from "../Search/Search";
import { SidebarToggle } from "../Sidebar/SidebarToggle";
import { getServerSidebarState } from "../Sidebar/server";
import { Title } from "../components/Title";
import { Toolbar } from "../components/Toolbar";
import styles from "../navigation.module.css";

export async function Top() {
  const { title } = await getInfo();
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
        <Toolbar session={session} />
      </HStack>
    </HStack>
  );
}
