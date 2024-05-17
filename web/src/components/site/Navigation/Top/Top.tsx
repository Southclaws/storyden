import { Sidebar } from "src/components/graphics/Sidebar/Sidebar";

import styles from "../navigation.module.css";

import { Search } from "../Search/Search";
import { Title } from "../components/Title";
import { Toolbar } from "../components/Toolbar";
import { useNavigation } from "../useNavigation";

import { Button } from "@/components/ui/button";
import { cx } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";

type Props = {
  onToggleSidebar: (t: boolean) => void;
  sidebarState: boolean;
};

export function Top({ onToggleSidebar, sidebarState }: Props) {
  const { title } = useNavigation();

  function handleToggle() {
    onToggleSidebar(!sidebarState);
  }

  return (
    <HStack
      className={cx(Floating(), styles["topbar"])}
      justify="space-between"
      alignItems="center"
      px="4"
    >
      <HStack className={styles["topbar-left"]}>
        {/* TODO: Action? */}
        <Button size="md" p="0" variant="ghost" onClick={handleToggle}>
          <Sidebar open={sidebarState} />
        </Button>
        <Search />
      </HStack>

      <HStack className={styles["topbar-middle"]} justify="space-around">
        <Title>{title}</Title>
      </HStack>

      <HStack className={styles["topbar-right"]}>
        <Toolbar />
      </HStack>
    </HStack>
  );
}
