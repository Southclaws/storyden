import { Box } from "@chakra-ui/react";
import { ThreadReference } from "src/api/openapi/schemas";
import { ThreadListItem } from "./ThreadListItem";

type Props = { threads: ThreadReference[] };

export function ThreadList(props: Props) {
  const children = props.threads.map((t) => (
    <ThreadListItem key={t.id} thread={t} />
  ));

  return <Box>{children}</Box>;
}
