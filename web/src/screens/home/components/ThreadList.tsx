import { Divider, List } from "@chakra-ui/react";
import { Fragment } from "react";

import { ThreadReference } from "src/api/openapi/schemas";

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
    <List width="full" display="flex" flexDirection="column">
      {props.threads.map((t) => (
        <Fragment key={t.id}>
          <Divider />
          <ThreadListItem key={t.id} thread={t} />
        </Fragment>
      ))}
    </List>
  );
}
