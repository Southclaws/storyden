"use client";

import { Divider } from "@chakra-ui/react";

import { useNavigation } from "../useNavigation";

import { Box, VStack, styled } from "@/styled-system/jsx";

import { Authbar } from "./components/Authbar";
import { CategoryCreate } from "./components/CategoryCreate/CategoryCreate";
import { CategoryList } from "./components/CategoryList/CategoryList";
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
      bgColor="var(--accent-colour-muted)"
    >
      <Box
        id="desktop-nav-box"
        maxWidth="2xs"
        minWidth={{
          base: "full",
          lg: "3xs",
        }}
        height="full"
        p={4}
      >
        <styled.nav
          display="flex"
          flexDir="column"
          height="full"
          gap={2}
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

            <Box overflowY="scroll" width="full">
              <CategoryList />
            </Box>

            {isAdmin && <CategoryCreate />}
          </VStack>

          <VStack alignItems="start">
            <Authbar />
          </VStack>
        </styled.nav>
      </Box>
    </styled.header>
  );
}
