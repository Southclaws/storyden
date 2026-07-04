"use client";

import { useSearchParams } from "next/navigation";
import { parseAsInteger, useQueryState } from "nuqs";

import { InvitedByFilter } from "@/components/library/members/MemberFilters/InvitedByFilter";
import { JoinedDateFilter } from "@/components/library/members/MemberFilters/JoinedDateFilter";
import { RoleFilter } from "@/components/library/members/MemberFilters/RoleFilter";
import { SortMenu } from "@/components/library/members/MemberFilters/SortMenu";
import { MemberList } from "@/components/library/members/MemberList";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { Unready } from "@/components/site/Unready";
import { Group } from "@/components/ui/group";
import { AdminIcon } from "@/components/ui/icons/Admin";
import { LinkButton } from "@/components/ui/link-button";
import { Flex, VStack } from "@/styled-system/jsx";

import { SearchInput } from "../SearchInput";

import { Props, useMemberIndexScreen } from "./useMemberIndexScreen";

export function MemberIndexScreen(props: Props) {
  const { ready, data, error } = useMemberIndexScreen(props);
  const searchParams = useSearchParams();
  const [page] = useQueryState(
    "page",
    parseAsInteger.withDefault(props.page ?? 1),
  );
  const currentParams = Object.fromEntries(searchParams.entries());

  if (!ready) {
    return <Unready error={error} />;
  }

  return (
    <VStack>
      <Group w="full">
        <SearchInput index="/m" initialQuery={props.query} />
        {props.adminModeAvailable && (
          <LinkButton
            size="md"
            variant="subtle"
            bg="bg.warning"
            href="/m?mode=admin"
          >
            <AdminIcon />
          </LinkButton>
        )}
      </Group>

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

      <PaginationControls
        path="/m"
        params={currentParams}
        currentPage={page}
        totalPages={data.total_pages}
        pageSize={data.page_size}
      />
    </VStack>
  );
}
