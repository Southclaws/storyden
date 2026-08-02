"use client";

import { CategoryListOKResponse, NodeListResult } from "@/api/openapi-schema";
import { CategoryList } from "@/components/category/CategoryList/CategoryList";
import { LStack, styled } from "@/styled-system/jsx";
import { useOverflowGradient } from "@/utils/useOverflowGradient";

import { CollectionsAnchor } from "../Anchors/Collections";
import { LinksAnchor } from "../Anchors/Link";
import { MembersAnchor } from "../Anchors/Members";
import { RolesAnchor } from "../Anchors/Roles";
import { LibraryNavigationTree } from "../LibraryNavigationTree/LibraryNavigationTree";
import { useNavigation } from "../useNavigation";

type Props = {
  initialNodeList?: NodeListResult;
  initialCategoryList?: CategoryListOKResponse;
};

export function ContentNavigationList(props: Props) {
  const { nodeSlug } = useNavigation();
  const scrollViewportRef = useOverflowGradient<HTMLDivElement>();

  return (
    <styled.nav
      aria-label="Site navigation"
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
        ref={scrollViewportRef}
        className="navigation__scroll-viewport"
        gap="1"
        overflowY="scroll"
        style={{
          scrollbarWidth: "none",
        }}
      >
        <CategoryList initialCategoryList={props.initialCategoryList} />
        <LibraryNavigationTree
          initialNodeList={props.initialNodeList}
          currentNode={nodeSlug}
          visibility={["draft", "review", "unlisted", "published"]}
        />
      </LStack>

      <LStack gap="1">
        <CollectionsAnchor />
        <LinksAnchor />
        <MembersAnchor />
        <RolesAnchor />
      </LStack>
    </styled.nav>
  );
}
