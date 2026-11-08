import { useSWRConfig } from "swr";

import {
  getThreadGetKey,
  getThreadListKey,
  threadReadStatusUpdate,
} from "@/api/openapi-client/threads";
import type {
  ReadStatus,
  ThreadGetResponse,
  ThreadListOKResponse,
  ThreadReadStatusMutableProps,
  ThreadReference,
} from "@/api/openapi-schema";

export function useThreadReadStatus(
  thread: ThreadReference | ThreadGetResponse,
) {
  const { mutate } = useSWRConfig();

  async function updateReadStatus(
    status: ThreadReadStatusMutableProps["status"],
  ) {
    const readStatus = await threadReadStatusUpdate(thread.slug, { status });
    await updateCachedReadStatus(readStatus);

    return readStatus;
  }

  async function recordRead() {
    await updateCachedReadStatus({
      last_read_at: new Date().toISOString(),
      replies_since: 0,
    });
  }

  async function updateCachedReadStatus(readStatus: ReadStatus) {
    const threadPaths = [
      getThreadGetKey(thread.slug)[0],
      getThreadGetKey(thread.id)[0],
    ];

    await mutate(
      (key) => Array.isArray(key) && threadPaths.includes(key[0]),
      (cached?: ThreadGetResponse) => {
        const current = cached ?? ("replies" in thread ? thread : undefined);
        if (!current) {
          return;
        }

        return { ...current, read_status: readStatus };
      },
      { revalidate: false },
    );

    await mutate(
      (key) => Array.isArray(key) && key[0] === getThreadListKey()[0],
      (cached?: ThreadListOKResponse) => {
        if (!cached) {
          return;
        }

        return {
          ...cached,
          threads: cached.threads.map((item) => {
            if (item.id !== thread.id) {
              return item;
            }

            return { ...item, read_status: readStatus };
          }),
        };
      },
      { revalidate: false },
    );
  }

  return { updateReadStatus, recordRead };
}
