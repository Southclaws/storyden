import { z } from "zod";

import { UnreadyBanner } from "@/components/site/Unready";

import { FeedScreen } from "@/screens/feed/FeedScreen";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

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
