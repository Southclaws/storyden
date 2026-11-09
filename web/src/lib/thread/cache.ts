import { useSWRConfig } from "swr";

import {
  getThreadGetKey,
  getThreadListKey,
} from "@/api/openapi-client/threads";
import type {
  ThreadGetResponse,
  ThreadListOKResponse,
  ThreadReference,
} from "@/api/openapi-schema";

type ThreadState = Partial<
  Pick<
    ThreadReference,
    "read_status" | "lifecycle_state" | "conclusion_reply_id" | "updatedAt"
  >
>;

export function useThreadCache(thread: ThreadReference | ThreadGetResponse) {
  const { mutate } = useSWRConfig();

  async function updateCachedThread(state: ThreadState) {
    const threadPaths = [
      getThreadGetKey(thread.slug)[0],
      getThreadGetKey(thread.id)[0],
    ];

    // Patch root state across pages without replacing their replies. SWR's
    // fallback data may still be absent from the shared cache on first render.
    await mutate(
      (key) => Array.isArray(key) && threadPaths.includes(key[0]),
      (cached?: ThreadGetResponse) => {
        const current = cached ?? ("replies" in thread ? thread : undefined);
        if (!current) {
          return;
        }

        return { ...current, ...state };
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

            return { ...item, ...state };
          }),
        };
      },
      { revalidate: false },
    );
  }

  return { updateCachedThread };
}
