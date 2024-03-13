"use client";

import { BookOpenIcon, UsersIcon } from "@heroicons/react/24/outline";

import { CategoryCreateTrigger } from "src/components/category/CategoryCreate/CategoryCreateTrigger";
import { Link } from "src/theme/components/Link";

import { CategoryList } from "../../../category/CategoryList/CategoryList";
import { useNavigation } from "../useNavigation";

import { Box, LStack, styled } from "@/styled-system/jsx";

export function ContentNavigationList() {
  const { isAdmin } = useNavigation();

  return (
    <styled.nav
      display="flex"
      flexDir="column"
      height="full"
      width="full"
      gap="2"
      alignItems="start"
      justifyContent="space-between"
      overflowY="scroll"
    >
      <LStack>
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
      </LStack>

      <LStack>
        <Link w="full" size="xs" href="/directory">
          <BookOpenIcon />
          Knowledgebase
        </Link>

        <Link w="full" size="xs" href="/p">
          <UsersIcon />
          Member directory
        </Link>
      </LStack>
    </styled.nav>
  );
}
