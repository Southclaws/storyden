import { useSWRConfig } from "swr";

import { getReplyGetKey } from "@/api/openapi-client/replies";
import { threadUpdate } from "@/api/openapi-client/threads";
import type {
  Reply,
  ThreadGetResponse,
  ThreadLifecycleState,
  ThreadMutableProps,
  ThreadReference,
} from "@/api/openapi-schema";

import { useThreadCache } from "./cache";

export function useThreadLifecycle(
  thread: ThreadReference | ThreadGetResponse,
) {
  const { updateCachedThread } = useThreadCache(thread);
  const { mutate } = useSWRConfig();

  async function updateLifecycleState(lifecycleState: ThreadLifecycleState) {
    await publishState({ lifecycle_state: lifecycleState });
  }

  async function updateConclusionReply(reply: Reply | null) {
    await publishState({ conclusion_reply_id: reply?.id ?? null }, reply);
  }

  async function publishState(
    partial: Pick<
      ThreadMutableProps,
      "lifecycle_state" | "conclusion_reply_id"
    >,
    reply?: Reply | null,
  ) {
    const updated = await threadUpdate(thread.id, partial);
    if (reply && updated.conclusion_reply_id === reply.id) {
      await mutate(getReplyGetKey(thread.slug, reply.id), reply, {
        revalidate: false,
      });
    }

    await updateCachedThread({
      lifecycle_state: updated.lifecycle_state,
      conclusion_reply_id: updated.conclusion_reply_id,
      updatedAt: updated.updatedAt,
    });
  }

  return { updateLifecycleState, updateConclusionReply };
}
