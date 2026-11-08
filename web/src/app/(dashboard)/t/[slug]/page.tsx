import type { Metadata } from "next";
import { notFound } from "next/navigation";

import { getServerSession } from "@/auth/server-session";
import { ServerRenderedContent } from "@/components/content/ContentRenderer/ServerRenderedContent";
import { getPermalinkForPost } from "@/components/thread/utils";
import {
  canonicalPage,
  jsonLD,
  pageLinks,
  socialImage,
} from "@/lib/metadata/public";
import { getPublicThread, getThread } from "@/lib/metadata/thread";
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

export default async function Page(props: Props) {
  const { slug } = await props.params;
  const { page } = QuerySchema.parse(await props.searchParams);
  const data = await getThread(slug, page ?? 1);
  if (!data) {
    notFound();
  }

  const session = await getServerSession();
  const settings = await getSettings();
  const result = await getPublicThread(slug, page ?? 1);
  const publicThread = result.status === "public" ? result.data : undefined;

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
      <ServerRenderedContent
        values={[data.body, ...data.replies.replies.map((reply) => reply.body)]}
      >
        <ThreadScreen
          initialSession={session}
          initialPage={page}
          slug={slug}
          thread={data}
          initialSettings={settings}
          initialSignatureConfig={settings.metadata.signatures}
        />
      </ServerRenderedContent>
    </>
  );
}

export async function generateMetadata(props: Props): Promise<Metadata> {
  const { slug } = await props.params;
  const query = QuerySchema.safeParse(await props.searchParams);
  if (!query.success) {
    return {
      title: "Thread Not Found",
      robots: { index: false, follow: false },
    };
  }

  const page = query.data.page ?? 1;
  const thread = await getThread(slug, page);
  if (!thread) {
    return {
      title: "Thread Not Found",
      robots: { index: false, follow: false },
    };
  }

  const settings = await getSettings();
  const title = `${thread.title} | ${settings.title}`;
  const result = await getPublicThread(slug, page);
  if (result.status === "not-public") {
    return { title, robots: { index: false, follow: false } };
  }

  if (result.status === "unavailable") {
    return { title };
  }

  const publicThread = result.data;
  const path = `/t/${encodeURIComponent(publicThread.slug)}`;
  const url = canonicalPage(path, page);
  const totalPages = Math.max(1, publicThread.replies.total_pages);

  return {
    title,
    description: publicThread.description,
    alternates: { canonical: url },
    pagination: pageLinks(path, page, totalPages),
    openGraph: {
      images: [socialImage(undefined, settings.title)],
      type: "article",
      title: publicThread.title,
      description: publicThread.description,
      url,
      publishedTime: publicThread.createdAt,
      modifiedTime: publicThread.updatedAt,
    },
  };
}
