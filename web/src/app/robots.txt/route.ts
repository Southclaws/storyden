import { getDiscoverySettings } from "@/lib/settings/discovery-server";
import { generateRobotsTxt } from "@/lib/settings/robots-txt";

export const dynamic = "force-dynamic";

export async function GET(): Promise<Response> {
  try {
    const settings = await getDiscoverySettings();

    if (settings.robots_txt !== "") {
      return new Response(settings.robots_txt, {
        headers: {
          "Content-Type": "text/plain; charset=utf-8",
          "Cache-Control": "no-store",
        },
      });
    }

    const body = await generateRobotsTxt(
      settings.search_indexing,
      settings.training_crawling,
    );

    return new Response(body, {
      headers: {
        "Content-Type": "text/plain; charset=utf-8",
        "Cache-Control": "no-store",
      },
    });
  } catch (error) {
    console.error("Failed to determine public crawler access", error);

    return new Response("Crawler policy temporarily unavailable", {
      status: 503,
      headers: { "Cache-Control": "no-store" },
    });
  }
}
