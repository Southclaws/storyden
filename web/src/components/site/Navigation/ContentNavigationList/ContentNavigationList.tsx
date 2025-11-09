import { Suspense } from "react";

import { LStack, styled } from "@/styled-system/jsx";

import { CollectionsAnchor } from "../Anchors/Collections";
import { LinksAnchor } from "../Anchors/Link";
import { MembersAnchor } from "../Anchors/Members";
import { CategoryListServer } from "@/components/category/CategoryList/CategoryListServer";
import { LibraryNavigationTreeServer } from "../LibraryNavigationTree/LibraryNavigationTreeServer";

export function ContentNavigationList() {
  return (
    <styled.nav
      display="flex"
      flexDir="column"
      gap="4"
      height="full"
      width="full"
      minH="0"
      alignItems="start"
      justifyContent="space-between"
    >
      <LStack
        gap="1"
        overflowY="scroll"
        style={{
          scrollbarWidth: "none",
        }}
      >
        <Suspense>
          <CategoryListServer />
        </Suspense>
        <Suspense>
          <LibraryNavigationTreeServer />
        </Suspense>
      </LStack>

      <LStack gap="1">
        <CollectionsAnchor />
        <LinksAnchor />
        <MembersAnchor />
      </LStack>
    </styled.nav>
  );
}
