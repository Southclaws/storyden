import { LinkListResult } from "src/api/openapi/schemas";
import { LinkCard } from "src/components/directory/links/LinkCard";
import { Empty } from "src/components/feed/common/PostRef/Empty";

import { styled } from "@/styled-system/jsx";

type Props = {
  links: LinkListResult;
  show?: number;
};

export function LinkCardList({ links, show }: Props) {
  if (links.links.length === 0) {
    return <Empty>no links were found</Empty>;
  }

  const shown = show ? links.links.slice(0, show) : links.links;

  return (
    <styled.ol w="full" display="flex" flexDir="column" gap="4">
      {shown.map((v) => (
        <styled.li key={v.url}>
          <LinkCard {...v} />
        </styled.li>
      ))}
    </styled.ol>
  );
}
