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
        <DatagraphNavTree currentNode={nodeSlug} />
      </LStack>

      <LStack gap="1">
        <MembersAnchor />
      </LStack>
    </styled.nav>
  );
}
