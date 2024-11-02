"use client";

import { Unready } from "src/components/site/Unready";

import { MemberList } from "@/components/library/members/MemberList";
import { PaginatedSearch } from "@/components/site/PaginatedSearch/PaginatedSearch";
import { VStack } from "@/styled-system/jsx";

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

      <MemberList profiles={data.profiles} />
    </VStack>
  );
}
