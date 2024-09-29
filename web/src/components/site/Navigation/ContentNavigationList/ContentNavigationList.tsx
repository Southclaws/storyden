"use client";

import { BookOpenIcon, UsersIcon } from "@heroicons/react/24/outline";

import { CategoryCreateTrigger } from "@/components/category/CategoryCreate/CategoryCreateTrigger";
import { CategoryList } from "@/components/category/CategoryList/CategoryList";
import { LinkButton } from "@/components/ui/link-button";
import { LStack, styled } from "@/styled-system/jsx";

import { DatagraphNavTree } from "../DatagraphNavTree/DatagraphNavTree";
import { useNavigation } from "../useNavigation";

export function ContentNavigationList() {
  const { isAdmin, nodeSlug } = useNavigation();

  return (
    <styled.nav
      display="flex"
      flexDir="column"
      height="full"
      width="full"
      alignItems="start"
      justifyContent="space-between"
      overflowY="scroll"
    >
      <LStack gap="1">
        <CategoryList />

        <DatagraphNavTree currentNode={nodeSlug} />
      </LStack>

      <LStack gap="1">
        {isAdmin && <CategoryCreateTrigger />}

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
