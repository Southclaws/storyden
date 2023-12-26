import { LinkListResult } from "src/api/openapi/schemas";
import { LinkRef } from "src/components/directory/links/LinkRef";
import { Empty } from "src/components/feed/common/PostRef/Empty";

import { styled } from "@/styled-system/jsx";

export function LinkRefList(props: { result: LinkListResult }) {
  if (props.result.links.length === 0) {
    return <Empty>no links were found</Empty>;
  }

  return (
    <styled.ol display="flex" flexDir="column" gap="4">
      {props.result.links.map((v) => (
        <styled.li key={v.url}>
          <LinkRef {...v} />
        </styled.li>
      ))}
    </styled.ol>
  );
}
