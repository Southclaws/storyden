import {
  getGetDocsUrl,
  getGetSpecUrl,
  getInfo,
} from "@/api/openapi-server/misc";
import { getAICatalog, mcpCardURL } from "@/lib/discovery/public";
import { resolveDiscoverySettings } from "@/lib/settings/discovery";

export const dynamic = "force-dynamic";

export async function GET() {
  try {
    const { data: info } = await getInfo({
      cache: "no-store",
    });

    const settings = resolveDiscoverySettings(info.discovery);

    if (!settings.llms_txt) {
      return new Response("Not found", {
        status: 404,
        headers: { "Cache-Control": "no-store" },
      });
    }

    const catalog = settings.ai_catalog ? await getAICatalog() : null;

    const web = new URL(info.web_address);
    const api = new URL(info.api_address);

    const authGuide = new URL("/auth.md", web);
    const apiCatalog = new URL("/.well-known/api-catalog", web);
    const openAPI = new URL(`/api${getGetSpecUrl()}`, api);
    const apiDocs = new URL(`/api${getGetDocsUrl()}`, api);

    const hasMCP = catalog?.entries.some(
      (entry) => entry.type === "application/mcp-server-card+json"
    );

    let mcpSection = "";

    if (settings.ai_catalog && catalog === null) {
      mcpSection = `## MCP

MCP discovery: unavailable from this instance.
`;
    } else if (hasMCP && settings.mcp_server_card) {
      const aiCatalog = new URL("/.well-known/ai-catalog.json", web);
      const card = mcpCardURL(catalog);

      mcpSection = `## MCP

- [AI catalog](${aiCatalog}):
  Discover agent services offered by this instance.
`;

      if (card) {
        mcpSection += `
- [MCP Server Card](${card}):
  Connect an MCP client, then list tools at runtime.
`;
      }
    }

    const links = [
      `- [Community](${web.href}):
  Read public discussions and pages in the web interface.`,
    ];

    if (settings.auth_md) {
      links.push(`- [Authentication guide](${authGuide}):
  Choose member access keys or available OAuth flows.`);
    }

    if (settings.api_catalog) {
      links.push(`- [API catalog](${apiCatalog}):
  Find this instance's API services.`);
    }

    links.push(
      `- [OpenAPI specification](${openAPI}):
  Generate a client or inspect request and response schemas.`,
      `- [API documentation](${apiDocs}):
  Explore endpoints interactively.`,
      `- [Storyden CLI](https://www.storyden.org/docs/introduction/cli):
  Install and use \`sd\` for command line based access.`,
      `- [OAuth and client registration](https://www.storyden.org/docs/introduction/oauth/client-registration):
  Learn operator-controlled authorization options.`
    );

    const body = `# ${info.title}

> ${info.description}

## Explore and integrate

${links.join("\n\n")}

${mcpSection}`;

    return new Response(body, {
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
