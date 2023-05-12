import { Divider, List } from "@chakra-ui/react";
import { ThreadReference } from "src/api/openapi/schemas";
import { ThreadListItem } from "./ThreadListItem";

type Props = { threads: ThreadReference[] };

export function ThreadList(props: Props) {
  const children = props.threads.map((t) => (
    <>
      <Divider />
      <ThreadListItem key={t.id} thread={t} />
    </>
  ));

  return (
    <List width="full" display="flex" flexDirection="column" gap={2}>
      {children}
    </List>
  );
}
