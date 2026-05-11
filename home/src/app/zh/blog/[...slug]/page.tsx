import { blogZh } from "@/lib/source";
import { getMDXComponents } from "@/mdx-components";
import { styled, VStack } from "@/styled-system/jsx";
import { linkButton } from "@/styled-system/patterns";
import { createRelativeLink } from "fumadocs-ui/mdx";
import { notFound } from "next/navigation";

export default async function Page(props: {
  params: Promise<{ slug?: string[] }>;
}) {
  const params = await props.params;
  const page = blogZh.getPage(params.slug);
  if (!page) notFound();

  const { body: MDXContent } = await page.data.load();

  return (
    <VStack>
      <styled.article
        className="prose"
        w="full"
        px="2"
        maxW="prose"
        pt="20"
        alignItems="start"
      >
        <h1 className="mb-2 text-3xl font-bold text-white">
          {page.data.title}
        </h1>

        <MDXContent
          components={{
            a: createRelativeLink(blogZh, page),
            ...getMDXComponents(),
          }}
        />

        <styled.hr w="full" />

        <styled.a
          href="/zh/blog"
          className={linkButton({
            backgroundColor: "Shades.newspaper",
            color: "Shades.iron",
            p: "2",
            height: "auto",
            lineHeight: "tight",
            fontWeight: "medium",
          })}
        >
          返回索引
        </styled.a>
      </styled.article>
    </VStack>
  );
}

export async function generateStaticParams() {
  return blogZh.generateParams();
}

export async function generateMetadata(props: {
  params: Promise<{ slug?: string[] }>;
}) {
  const params = await props.params;
  const page = blogZh.getPage(params.slug);
  if (!page) notFound();

  const { slug = [] } = params;
  const image = ["/zh/blog-og", ...slug, "image.png"].join("/");

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
