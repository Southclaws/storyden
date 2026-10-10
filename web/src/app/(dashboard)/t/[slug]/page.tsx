import type { Metadata } from "next";
import { z } from "zod";

import { threadGet } from "@/api/openapi-server/threads";
import { getServerSession } from "@/auth/server-session";
import {
  canonicalPage,
  getPublicThread,
  hiddenMetadata,
  jsonLD,
  pageLinks,
  pageNumber,
  threadSchema,
} from "@/lib/metadata/public";
import { getSettings } from "@/lib/settings/settings-server";
import { ThreadScreen } from "@/screens/thread/ThreadScreen/ThreadScreen";

export type Props = {
  params: Promise<{
    slug: string;
  }>;
  searchParams: Promise<{ page?: string }>;
};

const QuerySchema = z.object({
  page: z
    .string()
    .regex(/^[1-9]\d*$/)
    .transform(Number)
    .refine(Number.isSafeInteger)
    .optional(),
});

export default async function Page(props: Props) {
  const { slug } = await props.params;
  const searchParams = await props.searchParams;

  const { page } = QuerySchema.parse(searchParams);

  const { data } = await threadGet(slug, {
    page: page?.toString(),
  });

  const session = await getServerSession();
  const settings = await getSettings();
  const publicThread =
    data.visibility === "published" && (page === undefined || page === 1)
      ? await getPublicThread(slug)
      : undefined;

  return (
    <>
      {publicThread && (
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{
            __html: jsonLD(threadSchema(publicThread)),
          }}
        />
      )}
      <ThreadScreen
        initialSession={session}
        initialPage={page}
        slug={slug}
        thread={data}
        initialSettings={settings}
        initialSignatureConfig={settings.metadata.signatures}
      />
    </>
  );
}

export async function generateMetadata(props: Props): Promise<Metadata> {
  const params = await props.params;
  const { page: rawPage } = await props.searchParams;
  const page = rawPage === undefined ? 1 : pageNumber(rawPage);
  const thread = page ? await getPublicThread(params.slug) : undefined;

  const totalPages = Math.max(1, thread?.replies.total_pages ?? 1);
  if (!thread || !page || page > totalPages) {
    return {
      title: "Thread Not Found",
      description: "The thread you are looking for does not exist.",
      ...hiddenMetadata,
    };
  }

  const settings = await getSettings();
  const path = `/t/${encodeURIComponent(thread.slug)}`;
  const url = canonicalPage(path, String(page));

  return {
    title: `${thread.title} | ${settings.title}`,
    description: thread.description,
    alternates: {
      canonical: url,
      types: { "text/markdown": canonicalPage(`${path}.md`, String(page)) },
    },
    pagination: pageLinks(path, page, totalPages),
    openGraph: {
      type: "article",
      title: thread.title,
      description: thread.description,
      url,
      publishedTime: thread.createdAt,
      modifiedTime: thread.updatedAt,
    },
  };
}
