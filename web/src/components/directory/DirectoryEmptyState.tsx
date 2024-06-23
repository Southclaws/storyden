import { Empty } from "src/components/site/Empty";

import { useSession } from "@/auth";
import { Center } from "@/styled-system/jsx";

export function DirectoryEmptyState() {
  const session = useSession();
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
