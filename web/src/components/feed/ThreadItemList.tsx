import { ThreadReference } from "src/api/openapi-schema";

import { styled } from "@/styled-system/jsx";

import { ThreadItem } from "./ThreadItem";

type Props = {
  threads: ThreadReference[];
};

export function ThreadItemList(props: Props) {
  return (
    <styled.ol width="full" display="flex" flexDirection="column" gap="3">
      {props.threads.map((t) => (
        <ThreadItem key={t.id} thread={t} />
      ))}
    </styled.ol>
  );
}
