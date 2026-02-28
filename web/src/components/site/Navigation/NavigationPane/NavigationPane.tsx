import {
  type Account,
  type CategoryListOKResponse,
  type NodeListResult,
} from "@/api/openapi-schema";
import { categoryListCached } from "@/lib/category/server-category-list";
import { nodeListCached } from "@/lib/library/server-node-list";
import { type Settings } from "@/lib/settings/settings";
import { Box, styled } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";

import { ContentNavigationList } from "../ContentNavigationList/ContentNavigationList";

import { AdminZone } from "./AdminZone/AdminZone";

type ServerProps = {
  initialSession?: Account;
  initialSettings?: Settings;
};

export async function NavigationPane({
  initialSession,
  initialSettings,
}: ServerProps) {
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
        initialSession={initialSession}
        initialSettings={initialSettings}
      />
    );
  } catch (e) {
    return (
      <NavigationPaneContent
        initialSession={initialSession}
        initialSettings={initialSettings}
      />
    );
  }
}

type Props = {
  initialNodeList?: NodeListResult;
  initialCategoryList?: CategoryListOKResponse;
  initialSession?: Account;
  initialSettings?: Settings;
};

function NavigationPaneContent({
  initialNodeList,
  initialCategoryList,
  initialSession,
  initialSettings,
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
      <AdminZone
        initialSession={initialSession}
        initialSettings={initialSettings}
      />
      <Box id="desktop-nav-box" w="full" height="full" minH="0" p="2">
        <ContentNavigationList
          initialNodeList={initialNodeList}
          initialCategoryList={initialCategoryList}
        />
      </Box>
    </styled.header>
  );
}
