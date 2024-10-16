"use client";

import { CategoryList } from "@/components/category/CategoryList/CategoryList";
import { Divider, LStack, styled } from "@/styled-system/jsx";

import { MembersAnchor } from "../Anchors/Members";
import { LibraryNavigationTree } from "../LibraryNavigationTree/LibraryNavigationTree";
import { useNavigation } from "../useNavigation";

export function ContentNavigationList() {
  const { nodeSlug } = useNavigation();

  return (
    <styled.nav
      display="flex"
      flexDir="column"
      gap="4"
      height="full"
      width="full"
      alignItems="start"
      justifyContent="space-between"
      overflowY="scroll"
    >
      <LStack gap="1">
        <CategoryList />
        <LibraryNavigationTree
          label="Library"
          href="/l"
          currentNode={nodeSlug}
          visibility={["draft", "review", "unlisted", "published"]}
        />
      </LStack>

      <LStack gap="1">
        <Divider />
        <MembersAnchor />
      </LStack>
    </styled.nav>
  );
}
