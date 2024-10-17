import { z } from "zod";

import { ThreadIndexScreen } from "src/screens/thread/ThreadIndexScreen/ThreadIndexScreen";

import { threadList } from "@/api/openapi-server/threads";

type Props = {
  searchParams: Promise<Query>;
};

const QuerySchema = z.object({
  q: z.string().optional(),
  page: z
    .string()
    .transform((v) => parseInt(v, 10))
    .optional(),
});

type Query = z.infer<typeof QuerySchema>;

export default async function Page(props: Props) {
  const searchParams = await props.searchParams;
  const { data } = await threadList({
    ...(searchParams.q ? { q: searchParams.q } : {}),
    ...(searchParams.page ? { page: searchParams.page?.toString() } : {}),
  });

  return (
    <ThreadIndexScreen
      threads={data}
      page={searchParams.page}
      query={searchParams.q}
    />
  );
}
