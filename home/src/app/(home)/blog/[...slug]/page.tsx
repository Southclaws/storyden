import { blog } from "@/lib/source";
import { getMDXComponents } from "@/mdx-components";
import { styled, VStack } from "@/styled-system/jsx";
import { linkButton } from "@/styled-system/patterns";
import { createRelativeLink } from "fumadocs-ui/mdx";
import { notFound } from "next/navigation";

export default async function Page(props: {
  params: Promise<{ slug?: string[] }>;
}) {
  const params = await props.params;
  const page = blog.getPage(params.slug);
  if (!page) notFound();

  const { body: MDXContent, toc } = await page.data.load();

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

        {/* <InlineTOC items={toc} /> */}

        <MDXContent
          components={{
            a: createRelativeLink(blog, page),
            ...getMDXComponents(),
          }}
        />

        <styled.hr w="full" />

        <styled.a
          href="/blog"
          className={linkButton({
            backgroundColor: "Shades.newspaper",
            color: "Shades.iron",
            p: "2",
            height: "auto",
            lineHeight: "tight",
            fontWeight: "medium",
          })}
        >
          Back to index
        </styled.a>
      </styled.article>
    </VStack>
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

  return {
    title: page.data.title,
    description: page.data.description,
  };
}
