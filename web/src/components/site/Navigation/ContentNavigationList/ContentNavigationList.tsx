"use client";

import { CategoryListOKResponse, NodeListResult } from "@/api/openapi-schema";
import { CategoryList } from "@/components/category/CategoryList/CategoryList";
import { Divider, LStack, styled } from "@/styled-system/jsx";

import { MembersAnchor } from "../Anchors/Members";
import { LibraryNavigationTree } from "../LibraryNavigationTree/LibraryNavigationTree";
import { useNavigation } from "../useNavigation";

type Props = {
  initialNodeList?: NodeListResult;
  initialCategoryList?: CategoryListOKResponse;
};

export function ContentNavigationList(props: Props) {
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
        <CategoryList initialCategoryList={props.initialCategoryList} />
        <LibraryNavigationTree
          initialNodeList={props.initialNodeList}
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
