import { z } from "zod";

import { threadGet } from "@/api/openapi-server/threads";
import { UnreadyBanner } from "@/components/site/Unready";
import { getSettings } from "@/lib/settings/settings-server";
import { ThreadScreen } from "@/screens/thread/ThreadScreen/ThreadScreen";

export type Props = {
  params: Promise<{
    slug: string;
  }>;
  searchParams: Promise<Query>;
};

const QuerySchema = z.object({
  page: z
    .string()
    .transform((v) => parseInt(v, 10))
    .optional(),
});

type Query = z.infer<typeof QuerySchema>;

export default async function Page(props: Props) {
  const { slug } = await props.params;
  const searchParams = await props.searchParams;

  const { page } = QuerySchema.parse(searchParams);

  const { data } = await threadGet(slug);

  return <ThreadScreen initialPage={page} slug={slug} thread={data} />;
}

export async function generateMetadata(props: Props) {
  const params = await props.params;
  try {
    const settings = await getSettings();
    const { data } = await threadGet(params.slug);

    return {
      title: `${data.title} | ${settings.title}`,
      description: data.description,
    };
  } catch (e) {
    return {
      title: "Thread Not Found",
      description: "The thread you are looking for does not exist.",
    };
  }
}
