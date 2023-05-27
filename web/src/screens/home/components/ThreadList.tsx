import { Divider, List } from "@chakra-ui/react";
import { ThreadReference } from "src/api/openapi/schemas";
import { ThreadListItem } from "./ThreadListItem";
import { Fragment } from "react";

type Props = { threads: ThreadReference[] };

export function ThreadList(props: Props) {
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
