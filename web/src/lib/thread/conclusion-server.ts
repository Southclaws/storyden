import "server-only";

import type { Reply, Thread } from "@/api/openapi-schema";
import { replyGet } from "@/api/openapi-server/replies";
import { isNotPublic } from "@/lib/metadata/public";

export async function getConclusionReply(
  thread: Thread,
): Promise<Reply | undefined> {
  const id = thread.conclusion_reply_id;
  if (!id) {
    return undefined;
  }

  const current = thread.replies.replies.find((reply) => reply.id === id);
  if (current) {
    return current;
  }

  try {
    const { data } = await replyGet(thread.slug, id);
    return data;
  } catch (error) {
    if (!isNotPublic(error)) {
      throw error;
    }
  }

  return undefined;
}
