export const revalidate = false;

export function GET() {
  return new Response(
    `# Storyden

> Open-source software for forums, knowledge libraries, and communities. Each deployment has its own web and API origins, settings, permissions, and optional integrations.

## Start here

- [Overview](/docs/introduction): Product concepts and setup.
- [REST API](/docs/api): API reference and authentication.
- [MCP](/docs/introduction/mcp): Discover and connect to an instance's optional MCP endpoint.
- [Storyden CLI](/docs/introduction/cli): Build and use the sd command-line client.
- [OAuth](/docs/introduction/oauth): Authorization flows and client registration.
- [Access keys](/docs/reference/access-keys): Bearer tokens for existing members.
- [Configuration](/docs/reference/configuration): Deployment settings.
- [Complete documentation text](/llms-full.txt): All docs pages in one file.

For a particular community, start with its own /llms.txt, /developers, and /.well-known catalogs. These describe that instance's addresses and enabled integrations.
`,
    {
      headers: { "Content-Type": "text/plain; charset=utf-8" },
    },
  );
}
