import { PublicProfileList } from "src/api/openapi-schema";
import { Empty } from "src/components/site/Empty";

import { styled } from "@/styled-system/jsx";

import { MemberCard } from "./MemberCard";

type Props = {
  profiles: PublicProfileList;
  onChange: () => void;
};

export function MemberList(props: Props) {
  if (props.profiles.length === 0) {
    return <Empty>no members were found</Empty>;
  }

  return (
    <styled.table
      display="grid"
      w="full"
      style={{
        gridTemplateColumns:
          "minmax(150px, 2fr) minmax(150px, 1fr) min-content",
      }}
      gap="2"
    >
      <styled.tbody display="contents">
        {props.profiles.map((v) => (
          <MemberCard key={v.id} {...v} onChange={props.onChange} />
        ))}
      </styled.tbody>
    </styled.table>
  );
}
