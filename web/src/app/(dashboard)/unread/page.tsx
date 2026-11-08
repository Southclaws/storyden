import type { Metadata } from "next";

import { threadList } from "@/api/openapi-server/threads";
import { getServerSession } from "@/auth/server-session";
import {
  UnauthenticatedBanner,
  UnreadyBanner,
} from "@/components/site/Unready";
import { getSettings } from "@/lib/settings/settings-server";
import * as z from "@/lib/validation";
import { UnreadScreen } from "@/screens/unread/UnreadScreen";

export const metadata: Metadata = {
  title: "Unread",
  robots: { index: false, follow: false },
};

const QuerySchema = z.object({
  page: z.optional(
    z
      .pipe(z.string().check(z.regex(/^[1-9]\d*$/)), z.transform(Number))
      .check(z.refine(Number.isSafeInteger)),
  ),
});

type Props = {
  searchParams: Promise<z.input<typeof QuerySchema>>;
};

export default async function Page({ searchParams }: Props) {
  try {
    const { page = 1 } = QuerySchema.parse(await searchParams);
    const [session, settings] = await Promise.all([
      getServerSession(),
      getSettings(),
    ]);

    if (!session) {
      return <UnauthenticatedBanner initialSettings={settings} />;
    }

    const { data } = await threadList(
      { page: page.toString(), read_status: "unread", ignore_pinned: true },
      { cache: "no-store" },
    );

    return (
      <UnreadScreen
        initialSession={session}
        initialPage={page}
        initialThreadList={data}
      />
    );
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
