import { joeiBold, workSans } from "@/fonts/og";
import { blog } from "@/lib/source";
import { notFound } from "next/navigation";
import { ImageResponse } from "next/og";

export async function GET(
  _req: Request,
  { params }: { params: Promise<{ slug: string[] }> }
) {
  const { slug } = await params;
  const page = blog.getPage(slug.slice(0, -1));
  if (!page) notFound();

  return new ImageResponse(
    (
      <div
        style={{
          display: "flex",
          height: "640px",
          width: "1280px",
          alignItems: "center",
          justifyContent: "center",
          backgroundImage: 'url("https://www.storyden.org/docs_og.png")',
        }}
      >
        <div
          style={{
            display: "flex",
            flexDirection: "column",
            width: "100%",
            height: "100%",
            padding: "6rem",
          }}
        >
          <div
            style={{
              color: "white",
              fontFamily: "Joie Grotesk",
              fontSize: "4rem",
              fontWeight: "700",
            }}
          >
            {page.data.title}
          </div>
          <div
            style={{
              color: "#d8dbcd",
              fontFamily: "Work Sans",
              fontSize: "3rem",
              fontWeight: "400",
            }}
          >
            {page.data.description}
          </div>
        </div>
      </div>
    ),
    {
      height: 640,
      width: 1280,
      fonts: [
        {
          name: "Joie Grotesk",
          data: await joeiBold(),
          style: "normal",
          weight: 700,
        },
        {
          name: "Work Sans",
          data: await workSans(),
          style: "normal",
          weight: 400,
        },
      ],
    }
  );
}

export function generateStaticParams() {
  return blog.generateParams().map((page) => ({
    ...page,
    slug: [...page.slug, "image.png"],
  }));
}
