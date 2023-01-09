import { useToast } from "@chakra-ui/react";
import { postsCreate } from "src/api/openapi/posts";
import { APIError, Thread } from "src/api/openapi/schemas";
import { getThreadsGetKey } from "src/api/openapi/threads";
import { useSession } from "src/auth";
import { mutate } from "swr";

export function useThread(thread: Thread) {
  const account = useSession();
  const toast = useToast();

  function onReply(md: string) {
    postsCreate(thread.id, {
      body: md,
    })
      .then(() => {
        mutate(getThreadsGetKey(`${thread.id}-${thread.slug}`));
      })
      .catch((e: APIError) =>
        toast({
          title: "Error",
          status: "error",
          description: e.message,
        })
      );
  }

  return {
    loggedIn: !!account,
    onReply,
  };
}
