import { formatDate } from "date-fns";
import { createRelativeLink } from "fumadocs-ui/mdx";
import Link from "next/link";
import { notFound } from "next/navigation";

import { blog } from "@/lib/source";
import { getMDXComponents } from "@/mdx-components";

export default async function Page(props: {
  params: Promise<{ slug?: string[] }>;
}) {
  const params = await props.params;
  const page = blog.getPage(params.slug);
  if (!page) notFound();

  const { body: MDXContent } = await page.data.load();

  return (
    <article className="blog-article">
      <header>
        <Link className="story-text-link" href="/blog">
          ← From the workshop
        </Link>
        <h1>{page.data.title}</h1>
        <p>{page.data.description}</p>
        <time dateTime={formatDate(page.data.date, "yyyy-MM-dd")}>
          {formatDate(page.data.date, "d MMMM yyyy")}
        </time>
      </header>
      <div className="prose">
        <MDXContent
          components={{
            a: createRelativeLink(blog, page),
            ...getMDXComponents(),
          }}
        />
      </div>
      <footer>
        <Link className="story-text-link" href="/blog">
          Back to the workshop →
        </Link>
      </footer>
    </article>
  );
}

export async function generateStaticParams() {
  return blog.generateParams();
}

export async function generateMetadata(props: {
  params: Promise<{ slug?: string[] }>;
}) {
  const params = await props.params;
  const page = blog.getPage(params.slug);
  if (!page) notFound();

  const { slug = [] } = params;

  if (!page) notFound();

  const image = ["/blog-og", ...slug, "image.png"].join("/");

  return {
    title: page.data.title,
    description: page.data.description,
    openGraph: {
      images: image,
    },
    twitter: {
      card: "summary_large_image",
      images: image,
    },
  };
}
