import {
  MarkdownCopyButton as LLMCopyButton,
  ViewOptionsPopover as ViewOptions,
} from "@/components/ai/page-actions";
import { source } from "@/lib/source";
import { getMDXComponents } from "@/mdx-components";
import { HStack } from "@/styled-system/jsx";
import { createRelativeLink } from "fumadocs-ui/mdx";
import {
  DocsBody,
  DocsDescription,
  DocsPage,
  DocsTitle,
} from "fumadocs-ui/page";
import { Metadata } from "next";
import { notFound } from "next/navigation";

export default async function Page(props: {
  params: Promise<{ slug?: string[] }>;
}) {
  const params = await props.params;
  const page = source.getPage(params.slug);
  if (!page) notFound();

  const markdownUrl = `${page.url}.md`;
  const githubUrl = `https://github.com/Southclaws/storyden/blob/main/home/content/docs/${page.path}`;

  const MDXContent = page.data.body;

  return (
    <DocsPage
      toc={page.data.toc}
      full={page.data.full}
      editOnGithub={{
        owner: "Southclaws",
        repo: "storyden",
        sha: "main",
        path: `home/content/docs/${page.path}`,
      }}
    >
      <DocsTitle>{page.data.title}</DocsTitle>
      <DocsDescription>{page.data.description}</DocsDescription>

      <HStack>
        <LLMCopyButton markdownUrl={markdownUrl} />
        <ViewOptions markdownUrl={markdownUrl} githubUrl={githubUrl} />
      </HStack>

      <DocsBody>
        <MDXContent
          components={{
            // this allows you to link to other pages with relative file paths
            a: createRelativeLink(source, page),
            ...getMDXComponents(),
          }}
        />
      </DocsBody>
    </DocsPage>
  );
}

export async function generateStaticParams() {
  return source.generateParams();
}

export async function generateMetadata(props: {
  params: Promise<{ slug?: string[] }>;
}) {
  const params = await props.params;
  const { slug = [] } = params;
  const page = source.getPage(params.slug);
  if (!page) notFound();

  const image = ["/docs-og", ...slug, "image.png"].join("/");

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
  } satisfies Metadata;
}
