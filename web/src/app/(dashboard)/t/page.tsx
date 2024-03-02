import { z } from "zod";

import {
  ThreadListOKResponse,
  ThreadListParams,
} from "src/api/openapi/schemas";
import { server } from "src/api/server";
import { ThreadIndexScreen } from "src/screens/thread/ThreadIndexScreen/ThreadIndexScreen";

type Props = {
  searchParams: Query;
};

const QuerySchema = z.object({
  q: z.string().optional(),
  page: z
    .string()
    .transform((v) => parseInt(v, 10))
    .optional(),
});

type Query = z.infer<typeof QuerySchema>;

export default async function Page({ searchParams }: Props) {
  const response = await server<ThreadListOKResponse>({
    url: "/v1/threads",
    params: {
      ...(searchParams.q ? { q: searchParams.q } : {}),
      ...(searchParams.page ? { page: searchParams.page?.toString() } : {}),
    } satisfies ThreadListParams,
  });

  return (
    <ThreadIndexScreen
      threads={response}
      page={searchParams.page}
      query={searchParams.q}
    />
  );
}
