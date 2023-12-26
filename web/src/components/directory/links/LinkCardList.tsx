import { LinkListResult } from "src/api/openapi/schemas";
import { LinkCard } from "src/components/directory/links/LinkCard";
import { Empty } from "src/components/feed/common/PostRef/Empty";

import { styled } from "@/styled-system/jsx";

export function LinkCardList(props: { result: LinkListResult }) {
  if (props.result.links.length === 0) {
    return <Empty>no links were found</Empty>;
  }

  return (
    <styled.ol w="full" display="flex" flexDir="column" gap="4">
      {props.result.links.map((v) => (
        <styled.li key={v.url}>
          <LinkCard {...v} />
        </styled.li>
      ))}
    </styled.ol>
  );
}
