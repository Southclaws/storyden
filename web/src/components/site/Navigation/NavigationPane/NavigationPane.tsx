"use client";

import { Box, styled } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";

import { ContentNavigationList } from "../ContentNavigationList/ContentNavigationList";

export function NavigationPane() {
  return (
    <styled.header
      display="flex"
      height="full"
      justifyContent="end"
      borderRadius="md"
      className={Floating()}
    >
      <Box id="desktop-nav-box" w="full" height="full" p="2" pr="0">
        <ContentNavigationList />
      </Box>
    </styled.header>
  );
}
