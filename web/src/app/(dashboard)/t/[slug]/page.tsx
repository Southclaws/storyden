import type { Metadata } from "next";
import { notFound } from "next/navigation";

import { getServerSession } from "@/auth/server-session";
import { ServerRenderedContent } from "@/components/content/ContentRenderer/ServerRenderedContent";
import { getPermalinkForPost } from "@/components/thread/utils";
import {
  type PageQuery,
  PageQuerySchema,
  canonicalPage,
  jsonLD,
} from "@/lib/metadata/public";
import {
  getPublicThread,
  getThread,
  getThreadMetadata,
} from "@/lib/metadata/thread";
import { getSettings } from "@/lib/settings/settings-server";
import { ThreadScreen } from "@/screens/thread/ThreadScreen/ThreadScreen";

export type Props = {
  params: Promise<{ slug: string }>;
  searchParams: Promise<PageQuery>;
};

export default async function Page(props: Props) {
  const { slug } = await props.params;
  const { page } = PageQuerySchema.parse(await props.searchParams);
  const data = await getThread(slug, page ?? 1);
  if (!data) {
    notFound();
  }

  const session = await getServerSession();
  const settings = await getSettings();
  let publicThread;

  if (data.visibility === "published") {
    const result = await getPublicThread(slug, page ?? 1);
    if (result.status === "public") {
      publicThread = result.data;
    }
  }

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
  const query = PageQuerySchema.safeParse(await props.searchParams);
  if (!query.success) {
    return {
      title: "Thread Not Found",
      robots: { index: false, follow: false },
    };
  }

  return getThreadMetadata(slug, query.data.page ?? 1);
}
