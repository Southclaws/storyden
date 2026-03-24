"use client";

import Link from "next/link";

import { Unready } from "src/components/site/Unready";

import { InvitedByFilter } from "@/components/library/members/MemberFilters/InvitedByFilter";
import { JoinedDateFilter } from "@/components/library/members/MemberFilters/JoinedDateFilter";
import { RoleFilter } from "@/components/library/members/MemberFilters/RoleFilter";
import { SortMenu } from "@/components/library/members/MemberFilters/SortMenu";
import { MemberList } from "@/components/library/members/MemberList";
import { PaginatedSearch } from "@/components/site/PaginatedSearch/PaginatedSearch";
import { Button } from "@/components/ui/button";
import { Flex, VStack, WStack } from "@/styled-system/jsx";

import { Props, useMemberIndexScreen } from "./useMemberIndexScreen";

export function MemberIndexScreen(props: Props) {
  const { ready, data, error } = useMemberIndexScreen(props);
  if (!ready) {
    return <Unready error={error} />;
  }

  return (
    <VStack>
      <PaginatedSearch
        index="/m"
        initialQuery={props.query}
        initialPage={props.page}
        totalPages={data.total_pages}
        pageSize={data.page_size}
      />

      <Flex
        w="full"
        gap="2"
        flexDir={{
          base: "column",
          md: "row",
        }}
      >
        <RoleFilter />
        <InvitedByFilter />
        <JoinedDateFilter />
        <SortMenu />
      </Flex>

      {props.adminModeAvailable && (
        <WStack justifyContent="end">
          <Link href="/m?mode=admin">
            <Button size="sm" variant="subtle" bg="bg.warning">
              Admin mode
            </Button>
          </Link>
        </WStack>
      )}

      <MemberList profiles={data.profiles} />
    </VStack>
  );
}
