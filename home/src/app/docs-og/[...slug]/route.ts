import { joeiBold } from "@/fonts/og";
import { source } from "@/lib/source";
import { generateOGImage } from "fumadocs-ui/og";
import { notFound } from "next/navigation";

export async function GET(
  _req: Request,
  { params }: { params: Promise<{ slug: string[] }> }
) {
  const { slug } = await params;
  const page = source.getPage(slug.slice(0, -1));
  if (!page) notFound();

  return generateOGImage({
    title: page.data.title,
    description: page.data.description,
    site: "Storyden",
    primaryColor: "#D68E4D",
    fonts: [
      {
        data: await joeiBold(),
        name: "Joie Grotesk",
        weight: 700,
        style: "normal",
      },
    ],
  });
}

export function generateStaticParams() {
  return source.generateParams().map((page) => ({
    ...page,
    slug: [...page.slug, "image.png"],
  }));
}
