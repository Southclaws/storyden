import { threadGet } from "@/api/openapi-server/threads";
import { UnreadyBanner } from "@/components/site/Unready";
import { getSettings } from "@/lib/settings/settings-server";
import { ThreadScreen } from "@/screens/thread/ThreadScreen/ThreadScreen";

export type Props = {
  params: Promise<{
    slug: string;
  }>;
};

export default async function Page(props: Props) {
  const { slug } = await props.params;

  try {
    const { data } = await threadGet(slug);

    return <ThreadScreen slug={slug} thread={data} />;
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
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
