"use client";

import { LinkIcon, UsersIcon } from "@heroicons/react/24/outline";

import { CategoryCreateTrigger } from "src/components/category/CategoryCreate/CategoryCreateTrigger";
import { Link } from "src/theme/components/Link";

import { CategoryList } from "../../../category/CategoryList/CategoryList";
import { useNavigation } from "../useNavigation";

import { Box, Divider, styled } from "@/styled-system/jsx";
import { Floating } from "@/styled-system/patterns";

export function Left() {
  const { isAdmin } = useNavigation();

  return (
    <styled.header
      display="flex"
      height="full"
      justifyContent="end"
      bgColor="accent.200"
      borderRadius="md"
      className={Floating()}
    >
      <Box id="desktop-nav-box" w="full" height="full" p="4">
        <styled.nav
          display="flex"
          flexDir="column"
          height="full"
          gap="2"
          alignItems="start"
          overflowY="scroll"
        >
          <Box
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
        </styled.nav>
      </Box>
    </styled.header>
  );
}
