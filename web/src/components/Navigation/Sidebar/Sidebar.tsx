"use client";

import { Box, Divider, VStack } from "@chakra-ui/react";

import { useNavigation } from "../useNavigation";

import { styled } from "@/styled-system/jsx";

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
      bgColor="gray.100"
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
        <VStack
          as="nav"
          height="full"
          gap={2}
          justifyContent="space-between"
          alignItems="start"
        >
          <VStack width="full" alignItems="start" overflow="hidden">
            <Title>{title}</Title>

            <Toolbar />

            <Divider />

            <Box overflowY="scroll" width="full">
              <CategoryList />
            </Box>

            {isAdmin && <CategoryCreate />}
          </VStack>

          <VStack alignItems="start">
            <Authbar />
          </VStack>
        </VStack>
      </Box>
    </styled.header>
  );
}
