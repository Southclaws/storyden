import { Fragment } from "react";

import { ThreadReference } from "src/api/openapi/schemas";

import { Divider, styled } from "@/styled-system/jsx";

import { EmptyState } from "./EmptyState";
import { ThreadListItem } from "./ThreadListItem";

type Props = {
  threads: ThreadReference[];
  showEmptyState: boolean;
};

export function ThreadList(props: Props) {
  if (props.showEmptyState && props.threads.length === 0) {
    return <EmptyState />;
  }

  return (
    <styled.ul width="full" display="flex" flexDirection="column">
      {props.threads.map((t) => (
        <Fragment key={t.id}>
          <Divider />
          <ThreadListItem key={t.id} thread={t} />
        </Fragment>
      ))}
    </styled.ul>
  );
}
