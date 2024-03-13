"use client";

import { ContentNavigationList } from "../ContentNavigationList/ContentNavigationList";

import { Box, styled } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";

export function Left() {
  return (
    <styled.header
      display="flex"
      height="full"
      justifyContent="end"
      borderRadius="md"
      className={Floating()}
    >
      <Box id="desktop-nav-box" w="full" height="full" p="4">
        <ContentNavigationList />
      </Box>
    </styled.header>
  );
}
