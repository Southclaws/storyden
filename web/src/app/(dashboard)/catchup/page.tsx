import type { Metadata } from "next";

import { categoryList } from "@/api/openapi-server/categories";
import { threadList } from "@/api/openapi-server/threads";
import { getServerSession } from "@/auth/server-session";
import {
  UnauthenticatedBanner,
  UnreadyBanner,
} from "@/components/site/Unready";
import { getSettings } from "@/lib/settings/settings-server";
import * as z from "@/lib/validation";
import { CatchupScreen } from "@/screens/catchup/CatchupScreen";
import { catchupViews, getCatchupThreadParams } from "@/screens/catchup/query";

export const metadata: Metadata = {
  title: "Catch Up",
  robots: { index: false, follow: false },
};

const QuerySchema = z.object({
  page: z.optional(
    z
      .pipe(z.string().check(z.regex(/^[1-9]\d*$/)), z.transform(Number))
      .check(z.refine(Number.isSafeInteger)),
  ),
  view: z.optional(z.enum(catchupViews)),
  category: z.optional(z.string()),
});

type Props = {
  searchParams: Promise<z.input<typeof QuerySchema>>;
};

export default async function Page({ searchParams }: Props) {
  try {
    const {
      page = 1,
      view = "unread",
      category = "",
    } = QuerySchema.parse(await searchParams);
    const [session, settings] = await Promise.all([
      getServerSession(),
      getSettings(),
    ]);

    if (!session) {
      return <UnauthenticatedBanner initialSettings={settings} />;
    }

    const [threads, categories] = await Promise.all([
      threadList(getCatchupThreadParams(page, view, category), {
        cache: "no-store",
      }),
      categoryList(),
    ]);

    return (
      <CatchupScreen
        initialSession={session}
        initialPage={page}
        initialView={view}
        initialCategory={category}
        initialCategoryList={categories.data}
        initialThreadList={threads.data}
      />
    );
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
