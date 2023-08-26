"use client";

import { Box, Divider, VStack } from "@chakra-ui/react";

import { SIDEBAR_WIDTH, useNavigation } from "../useNavigation";

import { css } from "@/styled-system/css";
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
      className={css({
        display: "flex",
        position: "fixed",
        width: SIDEBAR_WIDTH,
        height: "full",
        justifyContent: "end",
        bgColor: "gray.100",
      })}
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
