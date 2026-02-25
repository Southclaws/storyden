import { z } from "zod";

import { UnreadyBanner } from "src/components/site/Unready";

import { FeedScreen } from "@/screens/feed/FeedScreen";

type Props = {
  searchParams: Promise<Query>;
};

const QuerySchema = z.object({
  page: z
    .string()
    .transform((v) => parseInt(v, 10))
    .optional(),
});
type Query = z.infer<typeof QuerySchema>;

export default async function Page({ searchParams }: Props) {
  try {
    const { page } = QuerySchema.parse(await searchParams);

    return <FeedScreen page={page ?? 1} />;
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
