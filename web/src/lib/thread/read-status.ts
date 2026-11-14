import { threadReadStatusUpdate } from "@/api/openapi-client/threads";
import type {
  ThreadGetResponse,
  ThreadReadStatusMutableProps,
  ThreadReference,
} from "@/api/openapi-schema";

import { useThreadCache } from "./cache";

export function useThreadReadStatus(
  thread: ThreadReference | ThreadGetResponse,
) {
  const { updateCachedThread } = useThreadCache(thread);

  async function updateReadStatus(
    status: ThreadReadStatusMutableProps["status"],
  ) {
    const readStatus = await threadReadStatusUpdate(thread.slug, { status });
    await updateCachedThread({ read_status: readStatus });

    return readStatus;
  }

  async function recordRead() {
    await updateCachedThread({
      read_status: {
        last_read_at: new Date().toISOString(),
        replies_since: 0,
      },
    });
  }

  return { updateReadStatus, recordRead };
}
