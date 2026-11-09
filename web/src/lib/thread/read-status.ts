import { useSWRConfig } from "swr";

import {
  getThreadListKey,
  threadReadStatusUpdate,
} from "@/api/openapi-client/threads";
import type {
  ThreadGetResponse,
  ThreadReadStatusMutableProps,
  ThreadReference,
} from "@/api/openapi-schema";

import { useThreadCache } from "./cache";

export function useThreadReadStatus(
  thread: ThreadReference | ThreadGetResponse,
) {
  const { mutate } = useSWRConfig();
  const { updateCachedThread } = useThreadCache(thread);

  async function updateReadStatus(
    status: ThreadReadStatusMutableProps["status"],
  ) {
    const readStatus = await threadReadStatusUpdate(thread.slug, { status });
    await updateCachedThread({ read_status: readStatus });
    await mutate(
      (key) =>
        Array.isArray(key) &&
        key[0] === getThreadListKey()[0] &&
        key[1]?.read_status !== undefined,
    );

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
