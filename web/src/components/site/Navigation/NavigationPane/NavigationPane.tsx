import { categoryList } from "@/api/openapi-server/categories";
import { nodeList } from "@/api/openapi-server/nodes";
import { Box, styled } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";

import { Unready } from "../../Unready";
import { ContentNavigationList } from "../ContentNavigationList/ContentNavigationList";

export async function NavigationPane() {
  try {
    const { data: initialNodeList } = await nodeList({
      // NOTE: This doesn't work due to a bug in Orval.
      // visibility: ["draft", "review", "unlisted", "published"],
    });
    const { data: initialCategoryList } = await categoryList();

    return (
      <styled.header
        display="flex"
        height="full"
        justifyContent="end"
        borderRadius="md"
        className={Floating()}
      >
        <Box id="desktop-nav-box" w="full" height="full" p="2" pr="0">
          <ContentNavigationList
            initialNodeList={initialNodeList}
            initialCategoryList={initialCategoryList}
          />
        </Box>
      </styled.header>
    );
  } catch (e) {
    return <Unready error={e} />;
  }
}
