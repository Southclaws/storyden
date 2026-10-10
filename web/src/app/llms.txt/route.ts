import {
  getAICatalog,
  getPublicInfo,
  mcpCardURL,
} from "@/lib/discovery/public";

export const dynamic = "force-dynamic";

export async function GET() {
  try {
    const [info, catalog] = await Promise.all([
      getPublicInfo(),
      getAICatalog(),
    ]);
    const web = new URL(info.webAddress);
    const api = new URL(info.apiAddress);
    const lines = [
      `# ${info.title}`,
      "",
      `> ${info.description}`,
      "",
      "## Explore and integrate",
      "",
      `- [Community](${web.href}): Read public discussions and pages in the web interface.`,
      `- [Developer guide](${new URL("/developers", web)}): Choose an API or MCP integration for this instance.`,
      `- [Authentication guide](${new URL("/auth.md", web)}): Choose member access keys or available OAuth flows.`,
      `- [API catalog](${new URL("/.well-known/api-catalog", web)}): Find this instance's API services.`,
      `- [OpenAPI specification](${new URL("/api/openapi.json", api)}): Generate a client or inspect request and response schemas.`,
      `- [API documentation](${new URL("/api/docs", api)}): Explore endpoints interactively.`,
      "- [Storyden CLI](https://www.storyden.org/docs/introduction/cli): Build and use `sd` for terminal workflows.",
      "- [OAuth and client registration](https://www.storyden.org/docs/introduction/oauth/client-registration): Learn operator-controlled authorization options.",
    ];
    const card = mcpCardURL(catalog, info.apiAddress);
    if (card) {
      lines.push(
        `- [AI catalog](${new URL("/.well-known/ai-catalog.json", web)}): Discover agent services offered by this instance.`,
        `- [MCP Server Card](${card}): Connect an MCP client, then list tools at runtime.`,
      );
    } else if (catalog === null) {
      lines.push("- MCP discovery: unavailable from this instance");
    }
    return new Response(`${lines.join("\n")}\n`, {
      headers: {
        "Content-Type": "text/plain; charset=utf-8",
        "Cache-Control": "no-store",
      },
    });
  } catch {
    return new Response("Instance discovery unavailable\n", {
      status: 503,
      headers: {
        "Content-Type": "text/plain; charset=utf-8",
        "Cache-Control": "no-store",
      },
    });
  }
}
