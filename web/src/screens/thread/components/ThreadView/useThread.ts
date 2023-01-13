import { useToast } from "@chakra-ui/react";
import { useState } from "react";
import { postsCreate } from "src/api/openapi/posts";
import { Thread } from "src/api/openapi/schemas";
import { useThreadsGet } from "src/api/openapi/threads";
import { useSession } from "src/auth";
import { errorToast } from "src/components/ErrorBanner";

export function useThread(thread: Thread) {
  const account = useSession();
  const toast = useToast();
  const { mutate } = useThreadsGet(thread.slug);
  const [isLoading, setLoading] = useState(false);

  const loggedIn = !!account;

  async function onReply(md: string) {
    if (!loggedIn) return;

    setLoading(true);
    await postsCreate(thread.id, { body: md }).catch(errorToast(toast));
    mutate();
    setLoading(false);
  }

  return {
    loggedIn,
    onReply,
    isLoading,
  };
}
