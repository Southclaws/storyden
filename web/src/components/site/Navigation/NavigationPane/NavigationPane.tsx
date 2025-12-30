import { CategoryListOKResponse, NodeListResult } from "@/api/openapi-schema";
import { categoryListCached } from "@/lib/category/server-category-list";
import { nodeListCached } from "@/lib/library/server-node-list";
import { Box, styled } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";

import { ContentNavigationList } from "../ContentNavigationList/ContentNavigationList";

import { AdminZone } from "./AdminZone/AdminZone";

export async function NavigationPane() {
  try {
    const { data: initialNodeList } = await nodeListCached({
      // NOTE: This doesn't work due to a bug in Orval.
      // visibility: ["draft", "review", "unlisted", "published"],
    });
    const { data: initialCategoryList } = await categoryListCached();

    return (
      <NavigationPaneContent
        initialNodeList={initialNodeList}
        initialCategoryList={initialCategoryList}
      />
    );
  } catch (e) {
    return <NavigationPaneContent />;
  }
}

type Props = {
  initialNodeList?: NodeListResult;
  initialCategoryList?: CategoryListOKResponse;
};

function NavigationPaneContent({
  initialNodeList,
  initialCategoryList,
}: Props) {
  return (
    <styled.header
      display="flex"
      height="full"
      alignItems="end"
      flexDirection="column"
      borderRadius="md"
      className={Floating()}
    >
      <AdminZone />
      <Box id="desktop-nav-box" w="full" height="full" minH="0" p="2">
        <ContentNavigationList
          initialNodeList={initialNodeList}
          initialCategoryList={initialCategoryList}
        />
      </Box>
    </styled.header>
  );
}
