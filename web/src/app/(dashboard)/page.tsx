import { Suspense } from "react";
import { z } from "zod";

import { UnreadyBanner } from "@/components/site/Unready";
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

export default function Page({ searchParams }: Props) {
  return (
    <Suspense fallback={<UnreadyBanner />}>
      <FeedScreenWrapper searchParams={searchParams} />
    </Suspense>
  );
}

async function FeedScreenWrapper({ searchParams }: Props) {
  const { page } = QuerySchema.parse(await searchParams);
  return <FeedScreen page={page ?? 1} />;
}
