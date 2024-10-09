"use client";

import { CategoryList } from "@/components/category/CategoryList/CategoryList";
import { LStack, styled } from "@/styled-system/jsx";

import { MembersAnchor } from "../Anchors/Members";
import { DatagraphNavTree } from "../DatagraphNavTree/DatagraphNavTree";
import { useNavigation } from "../useNavigation";

export function ContentNavigationList() {
  const { nodeSlug } = useNavigation();

  return (
    <styled.nav
      display="flex"
      flexDir="column"
      height="full"
      width="full"
      alignItems="start"
      justifyContent="space-between"
      overflowY="scroll"
    >
      <LStack gap="1">
        <CategoryList />
        <DatagraphNavTree
          label="Library"
          href="/l"
          currentNode={nodeSlug}
          visibility={["published"]}
        />
        <DatagraphNavTree
          label="Private"
          href="/drafts"
          currentNode={nodeSlug}
          visibility={["draft", "review", "unlisted"]}
        />
      </LStack>

      <LStack gap="1">
        <MembersAnchor />
      </LStack>
    </styled.nav>
  );
}
