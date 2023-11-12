"use client";

import { CategoryCreateTrigger } from "src/components/category/CategoryCreate/CategoryCreateTrigger";
import { Divider } from "src/theme/components";

import { CategoryList } from "../../../category/CategoryList/CategoryList";
import { useNavigation } from "../useNavigation";

import { Box, VStack, styled } from "@/styled-system/jsx";

import { Authbar } from "./components/Authbar";
import { Title } from "./components/Title";
import { Toolbar } from "./components/Toolbar";

export function Sidebar() {
  const { isAdmin, title } = useNavigation();

  return (
    <styled.header
      display="flex"
      position="fixed"
      width={{
        md: "1/4",
        lg: "1/3",
      }}
      height="full"
      justifyContent="end"
      bgColor="accent.200"
    >
      <Box
        id="desktop-nav-box"
        maxWidth="xs"
        minWidth={{
          base: "full",
          lg: "xs",
        }}
        height="full"
        p="4"
      >
        <styled.nav
          display="flex"
          flexDir="column"
          height="full"
          gap="2"
          justifyContent="space-between"
          alignItems="start"
        >
          <VStack width="full" alignItems="start" overflow="hidden">
            <Title>{title}</Title>

            <Toolbar />

            <Divider
              // TODO: make this clever based on accent colour.
              borderColor="oklch(0 0 0deg / 5%)"
            />

            <Box
              overflowY="scroll"
              width="full"
              css={{
                touchAction: "none",
                scrollbarWidth: "none",
                "&::-webkit-scrollbar": {
                  display: "none",
                },
              }}
            >
              <CategoryList />
            </Box>

            {isAdmin && <CategoryCreateTrigger />}
          </VStack>

          <VStack alignItems="start">
            <Authbar />
          </VStack>
        </styled.nav>
      </Box>
    </styled.header>
  );
}
