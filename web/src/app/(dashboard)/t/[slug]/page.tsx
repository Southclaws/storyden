import type { Metadata } from "next";
import { cache } from "react";

import { threadGet } from "@/api/openapi-server/threads";
import { getServerSession } from "@/auth/server-session";
import { getPermalinkForPost } from "@/components/thread/utils";
import { canonicalPage, jsonLD, pageLinks } from "@/lib/metadata/public";
import { getSettings } from "@/lib/settings/settings-server";
import * as z from "@/lib/validation";
import { ThreadScreen } from "@/screens/thread/ThreadScreen/ThreadScreen";

export type Props = {
  params: Promise<{ slug: string }>;
  searchParams: Promise<Query>;
};

const QuerySchema = z.object({
  page: z.optional(
    z
      .pipe(z.string().check(z.regex(/^[1-9]\d*$/)), z.transform(Number))
      .check(z.refine(Number.isSafeInteger)),
  ),
});
type Query = z.input<typeof QuerySchema>;

// Metadata and JSON-LD share one anonymous read for this render.
const getPublicThread = cache(async (slug: string, page: number) => {
  try {
    const { data } = await threadGet(
      slug,
      { page: String(page) },
      {
        anonymous: true,
        cache: "no-store",
      },
    );

    if (
      data.visibility === "published" &&
      page <= Math.max(1, data.replies.total_pages)
    ) {
      return data;
    }
  } catch {
    // Public metadata is unavailable if the anonymous read fails.
  }

  return undefined;
});

export default async function Page(props: Props) {
  const { slug } = await props.params;
  const { page } = QuerySchema.parse(await props.searchParams);
  const { data } = await threadGet(slug, { page: page?.toString() });
  const session = await getServerSession();
  const settings = await getSettings();
  const publicThread =
    data.visibility === "published"
      ? await getPublicThread(slug, page ?? 1)
      : undefined;
  const url = canonicalPage(`/t/${encodeURIComponent(data.slug)}`);
  const structuredData = publicThread && {
    "@context": "https://schema.org",
    "@type": "DiscussionForumPosting",
    headline: publicThread.title,
    description: publicThread.description,
    text: publicThread.body,
    url,
    mainEntityOfPage: url,
    datePublished: publicThread.createdAt,
    dateModified: publicThread.updatedAt,
    author: {
      "@type": "Person",
      name: publicThread.author.name,
      url: canonicalPage(
        `/m/${encodeURIComponent(publicThread.author.handle)}`,
      ),
    },
    commentCount: publicThread.replies.results,
    comment: publicThread.replies.replies
      .filter(
        (reply) =>
          reply.visibility === "published" &&
          data.replies.replies.some(
            (visibleReply) => visibleReply.id === reply.id,
          ),
      )
      .map((reply) => ({
        "@type": "Comment",
        text: reply.body,
        url: getPermalinkForPost(
          encodeURIComponent(publicThread.slug),
          reply.id,
          page,
        ),
        datePublished: reply.createdAt,
        dateModified: reply.updatedAt,
        author: {
          "@type": "Person",
          name: reply.author.name,
          url: canonicalPage(`/m/${encodeURIComponent(reply.author.handle)}`),
        },
      })),
  };

  return (
    <>
      {structuredData && (
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: jsonLD(structuredData) }}
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
  const { slug } = await props.params;
  const query = QuerySchema.safeParse(await props.searchParams);
  const page = query.success ? (query.data.page ?? 1) : 1;
  const thread = query.success ? await getPublicThread(slug, page) : undefined;
  const totalPages = Math.max(1, thread?.replies.total_pages ?? 0);

  if (!thread || page > totalPages) {
    return {
      title: "Thread Not Found",
      description: "The thread you are looking for does not exist.",
      robots: { index: false, follow: false },
    };
  }

  const settings = await getSettings();
  const path = `/t/${encodeURIComponent(thread.slug)}`;
  const url = canonicalPage(path, page);

  return {
    title: `${thread.title} | ${settings.title}`,
    description: thread.description,
    alternates: { canonical: url },
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
