"use client";

import { BookOpenIcon, UsersIcon } from "@heroicons/react/24/outline";

import { CategoryCreateTrigger } from "@/components/category/CategoryCreate/CategoryCreateTrigger";
import { LinkButton } from "@/components/ui/link-button";
import { Box, LStack, styled } from "@/styled-system/jsx";

import { CategoryList } from "../../../category/CategoryList/CategoryList";
import { useNavigation } from "../useNavigation";

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

      <LStack gap="1">
        <LinkButton w="full" size="xs" variant="ghost" href="/directory">
          <BookOpenIcon />
          Knowledgebase
        </LinkButton>

        <LinkButton w="full" size="xs" variant="ghost" href="/p">
          <UsersIcon />
          Member directory
        </LinkButton>
      </LStack>
    </styled.nav>
  );
}
