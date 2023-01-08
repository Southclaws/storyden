import { useToast } from "@chakra-ui/react";
import { postsCreate } from "src/api/openapi/posts";
import { APIError, Thread } from "src/api/openapi/schemas";
import { useSession } from "src/auth";

export function useThread(thread: Thread) {
  const account = useSession();
  const toast = useToast();

  function onReply(md: string) {
    postsCreate(thread.id, {
      body: md,
    }).catch((e: APIError) =>
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
