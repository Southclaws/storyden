"use client";

import { LinkIcon, UsersIcon } from "@heroicons/react/24/outline";

import { CategoryCreateTrigger } from "src/components/category/CategoryCreate/CategoryCreateTrigger";
import { Link } from "src/theme/components/Link";

import { CategoryList } from "../../../category/CategoryList/CategoryList";
import { useNavigation } from "../useNavigation";

import { Box, Divider, VStack, styled } from "@/styled-system/jsx";

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
          <VStack width="full" alignItems="start">
            <Title>{title}</Title>

            <Toolbar />

            <Divider />

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

            <Divider />

            <Link w="full" size="xs" href="/l">
              <LinkIcon />
              Link directory
            </Link>

            <Link w="full" size="xs" href="/p">
              <UsersIcon />
              Member directory
            </Link>
          </VStack>

          <VStack alignItems="start">
            <Authbar />
          </VStack>
        </styled.nav>
      </Box>
    </styled.header>
  );
}
