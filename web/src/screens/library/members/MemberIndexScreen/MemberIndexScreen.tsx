"use client";

import { Unready } from "src/components/site/Unready";

import { InvitedByFilter } from "@/components/library/members/MemberFilters/InvitedByFilter";
import { JoinedDateFilter } from "@/components/library/members/MemberFilters/JoinedDateFilter";
import { RoleFilter } from "@/components/library/members/MemberFilters/RoleFilter";
import { SortMenu } from "@/components/library/members/MemberFilters/SortMenu";
import { MemberList } from "@/components/library/members/MemberList";
import { PaginatedSearch } from "@/components/site/PaginatedSearch/PaginatedSearch";
import { Flex, VStack } from "@/styled-system/jsx";

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

      <MemberList profiles={data.profiles} />
    </VStack>
  );
}
