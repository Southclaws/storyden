import { getPublicInfo } from "@/lib/discovery/public";

export const dynamic = "force-dynamic";

export async function GET() {
  try {
    const info = await getPublicInfo();
    const api = new URL(info.apiAddress);
    const web = new URL(info.webAddress);
    const body = [
      `# Authentication for ${info.title}`,
      "",
      `API access uses a member's session cookie or a personal access key.${info.oauthEnabled ? " This instance also offers OAuth 2.0 access tokens." : ""} Each request is checked against that identity's current permissions.`,
      "",
      "- Personal access keys: create one in account settings, then send `Authorization: Bearer <key>` to the API. The member needs `USE_PERSONAL_ACCESS_KEYS` to create a key.",
      ...(info.oauthEnabled
        ? [
            `- OAuth discovery: ${new URL("/.well-known/oauth-authorization-server", api)}`,
            `- OpenID discovery: ${new URL("/.well-known/openid-configuration", api)}`,
            "- OAuth clients use authorization code with PKCE or device authorization. An administrator may also create a confidential client.",
            "- Dynamic client registration is controlled by the instance operator and may be disabled. Creating a separate autonomous agent account has its own registration policy and may require an invitation token or administrator approval.",
          ]
        : ["- OAuth is not advertised by this instance."]),
      "",
      `Read the instance integration guide: ${new URL("/developers", web)}`,
      "Protocol and policy details: https://www.storyden.org/docs/introduction/oauth/client-registration",
      "",
    ].join("\n");
    return new Response(body, {
      headers: {
        "Content-Type": "text/markdown; charset=utf-8",
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
