import { categoryList } from "@/api/openapi-server/categories";
import { nodeList } from "@/api/openapi-server/nodes";
import { Box, styled } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";

import { Unready } from "../../Unready";
import { ContentNavigationList } from "../ContentNavigationList/ContentNavigationList";

import { AdminZone } from "./AdminZone/AdminZone";

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
        alignItems="end"
        flexDirection="column"
        borderRadius="md"
        className={Floating()}
      >
        <AdminZone />
        <Box id="desktop-nav-box" w="full" height="full" p="2">
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
