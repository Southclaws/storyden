import { Suspense } from "react";

import { Box, styled } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";

import { ContentNavigationList } from "../ContentNavigationList/ContentNavigationList";

import { AdminZone } from "./AdminZone/AdminZone";

export function NavigationPane() {
  return (
    <styled.header
      display="flex"
      height="full"
      alignItems="end"
      flexDirection="column"
      borderRadius="md"
      className={Floating()}
    >
      <Suspense>
        <AdminZone />
      </Suspense>
      <Box id="desktop-nav-box" w="full" height="full" minH="0" p="2">
        <ContentNavigationList />
      </Box>
    </styled.header>
  );
}
