import { Empty } from "src/components/site/Empty";

import { Center } from "@/styled-system/jsx";

export function DirectoryEmptyState() {
  return (
    <Center h="full">
      <Empty>
        This community knowledgebase is empty.
        <br />
        {session ? (
          <>Be the first to contribute!</>
        ) : (
          <>Please log in to contribute.</>
        )}
      </Empty>
    </Center>
  );
}
