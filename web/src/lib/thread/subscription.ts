import { z } from "zod";

import { Thread } from "@/api/openapi-schema";

import { useSubscribe } from "../subscribe/useSubscribe";

import { useThreadMutations } from "./mutation";

export const ThreadUpdateSchema = z.object({
  type: z.literal("thread_update"),
});
export type ThreadUpdate = z.infer<typeof ThreadUpdateSchema>;

export function useThreadSubscription(thread: Thread) {
  const { revalidate } = useThreadMutations(thread);

  useSubscribe<ThreadUpdate>("thread", ThreadUpdateSchema, async (event) => {
    await revalidate();
  });
}
