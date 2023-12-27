"use client";

import { LinkIcon, UsersIcon } from "@heroicons/react/24/outline";

import { CategoryCreateTrigger } from "src/components/category/CategoryCreate/CategoryCreateTrigger";
import { Link } from "src/theme/components/Link";

import { CategoryList } from "../../../category/CategoryList/CategoryList";
import { navbarStyles } from "../common";
import { Authbar } from "../components/Authbar";
import { useNavigation } from "../useNavigation";

import { Box, Divider, VStack, styled } from "@/styled-system/jsx";

export function Left2() {
  return (
    <VStack className={navbarStyles} justify="space-between" px="4"></VStack>
  );
}

export function Left() {
  const { isAdmin } = useNavigation();

  return (
    <styled.header
      display="flex"
      height="full"
      justifyContent="end"
      bgColor="accent.200"
      className={navbarStyles}
    >
      <Box id="desktop-nav-box" w="full" height="full" p="4">
        <styled.nav
          display="flex"
          flexDir="column"
          height="full"
          gap="2"
          justifyContent="space-between"
          alignItems="start"
        >
          <VStack width="full" alignItems="start">
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
