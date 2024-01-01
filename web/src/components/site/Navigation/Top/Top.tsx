import { Sidebar } from "src/components/graphics/Sidebar/Sidebar";
import { Button } from "src/theme/components/Button";

import styles from "../navigation.module.css";

import { ComposeAction } from "../Anchors/Compose";
import { Title } from "../components/Title";
import { Toolbar } from "../components/Toolbar";
import { useNavigation } from "../useNavigation";

import { cx } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";
import { FrostedGlass } from "@/styled-system/patterns";

type Props = {
  onToggleSidebar: (t: boolean) => void;
  sidebarState: boolean;
};

export function Top({ onToggleSidebar, sidebarState }: Props) {
  const { title } = useNavigation();

  return (
    <HStack
      className={cx(FrostedGlass(), styles["topbar"])}
      justify="space-between"
      alignItems="center"
      px="4"
    >
      <HStack className={styles["topbar-left"]}>
        {/* TODO: Action? */}
        <Button
          size="sm"
          kind="ghost"
          onClick={() => onToggleSidebar(!sidebarState)}
        >
          <Sidebar open={sidebarState} />
        </Button>
        <ComposeAction>new</ComposeAction>
      </HStack>

      <HStack className={styles["topbar-middle"]} justify="center">
        {/* TODO: Semantic search */}
        {/* <Input placeholder="Search content, knowledgebase and links" /> */}
        <Title>{title}</Title>
      </HStack>

      <HStack className={styles["topbar-right"]}>
        <Toolbar />
      </HStack>
    </HStack>
  );
}
