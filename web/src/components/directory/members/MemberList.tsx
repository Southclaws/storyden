import { PublicProfileList } from "src/api/openapi/schemas";
import { Empty } from "src/components/feed/common/PostRef/Empty";

import { styled } from "@/styled-system/jsx";

import { MemberCard } from "./MemberCard";

export function MemberList(props: { profiles: PublicProfileList }) {
  if (props.profiles.length === 0) {
    return <Empty>no members were found</Empty>;
  }

  return (
    <styled.table
      w="full"
      tableLayout="fixed"
      borderCollapse="separate"
      borderSpacingY="2"
    >
      <styled.tbody>
        {props.profiles.map((v) => (
          <MemberCard key={v.id} {...v} />
        ))}
      </styled.tbody>
    </styled.table>
  );
}
