import { ThreadReference } from "src/api/openapi-schema";

import { styled } from "@/styled-system/jsx";

import { ThreadReferenceCard } from "./ThreadCard";

type Props = {
  threads: ThreadReference[];
};

export function ThreadReferenceList(props: Props) {
  return (
    <styled.ol width="full" display="flex" flexDirection="column" gap="3">
      {props.threads.map((t) => (
        <ThreadReferenceCard key={t.id} thread={t} />
      ))}
    </styled.ol>
  );
}
