import { getRSS } from "@/lib/rss";

export const revalidate = false;

export async function GET() {
  try {
    const rssContent = await getRSS();
    return new Response(rssContent, {
      headers: {
        "Content-Type": "application/xml",
      },
    });
  } catch (error) {
    console.error("RSS generation failed:", error);
    return new Response("RSS generation failed", { status: 500 });
  }
}
