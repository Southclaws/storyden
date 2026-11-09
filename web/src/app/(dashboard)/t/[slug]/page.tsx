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
import { getThread } from "@/lib/metadata/thread";
import { getSettings } from "@/lib/settings/settings-server";
import { getConclusionReply } from "@/lib/thread/conclusion-server";
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

  const [initialConclusionReply, session, settings] = await Promise.all([
    getConclusionReply(data),
    getServerSession(),
    getSettings(),
  ]);

  const url = canonicalPage(`/t/${encodeURIComponent(data.slug)}`);
  const isPublished = data.visibility === "published";
  const isValidPage = (page ?? 1) <= Math.max(1, data.replies.total_pages);
  const structuredData = isPublished &&
    isValidPage && {
      "@context": "https://schema.org",
      "@type": "DiscussionForumPosting",
      headline: data.title,
      description: data.description,
      text: data.body,
      url,
      mainEntityOfPage: url,
      datePublished: data.createdAt,
      dateModified: data.updatedAt,
      author: {
        "@type": "Person",
        name: data.author.name,
        url: canonicalPage(`/m/${encodeURIComponent(data.author.handle)}`),
      },
      commentCount: data.replies.results,
      comment: data.replies.replies
        .filter((reply) => reply.visibility === "published")
        .map((reply) => ({
          "@type": "Comment",
          text: reply.body,
          url: getPermalinkForPost(
            encodeURIComponent(data.slug),
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
        values={[
          data.body,
          initialConclusionReply?.body,
          ...data.replies.replies.map((reply) => reply.body),
        ]}
      >
        <ThreadScreen
          initialSession={session}
          initialPage={page}
          slug={slug}
          thread={data}
          initialSettings={settings}
          initialSignatureConfig={settings.metadata.signatures}
          initialConclusionReply={initialConclusionReply}
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
  const totalPages = Math.max(1, thread.replies.total_pages);
  if (thread.visibility !== "published" || page > totalPages) {
    return { title, robots: { index: false, follow: false } };
  }

  const path = `/t/${encodeURIComponent(thread.slug)}`;
  const url = canonicalPage(path, page);

  return {
    title,
    description: thread.description,
    alternates: { canonical: url },
    pagination: pageLinks(path, page, totalPages),
    openGraph: {
      images: [socialImage(undefined, settings.title)],
      type: "article",
      title: thread.title,
      description: thread.description,
      url,
      publishedTime: thread.createdAt,
      modifiedTime: thread.updatedAt,
    },
  };
}
