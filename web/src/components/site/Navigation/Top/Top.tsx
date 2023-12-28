import styles from "../navigation.module.css";

import { ComposeAction } from "../Anchors/Compose";
import { Title } from "../components/Title";
import { Toolbar } from "../components/Toolbar";
import { useNavigation } from "../useNavigation";

import { cx } from "@/styled-system/css";
import { HStack } from "@/styled-system/jsx";
import { FrostedGlass } from "@/styled-system/patterns";

export function Top() {
  const { title } = useNavigation();

  return (
    <HStack
      className={cx(FrostedGlass(), styles["topbar"])}
      justify="space-between"
      alignItems="center"
      px="4"
    >
      <HStack className={styles["topbar-left"]} justify="space-between">
        <Title>{title}</Title>
        <ComposeAction>new post</ComposeAction>
      </HStack>

      {/* TODO: Semantic search */}
      {/* <HStack className={styles["topbar-middle"]}>
        <Input placeholder="Search content, knowledgebase and links" />
      </HStack> */}

      <HStack className={styles["topbar-right"]}>
        <Toolbar />
      </HStack>
    </HStack>
  );
}
