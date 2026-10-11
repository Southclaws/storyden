import { InstanceCapability } from "@/api/openapi-schema";
import { getInfo } from "@/api/openapi-server/misc";

export const dynamic = "force-dynamic";

export async function GET() {
  try {
    const { data: info } = await getInfo({ cache: "no-store" });

    const oauthEnabled = info.capabilities.includes(InstanceCapability.oauth);
    let oauthSection = "OAuth is not advertised by this instance.";

    if (oauthEnabled) {
      const api = new URL(info.api_address);
      const oauthDiscovery = new URL(
        "/.well-known/oauth-authorization-server",
        api
      );
      const openIDDiscovery = new URL("/.well-known/openid-configuration", api);

      oauthSection = `This instance offers OAuth 2.0 access tokens.

- OAuth discovery: ${oauthDiscovery}
- OpenID discovery: ${openIDDiscovery}

OAuth clients use authorization code with PKCE or device authorization.
An administrator may also create a confidential client.

Dynamic client registration is controlled by the instance operator and may be
disabled. Creating a separate autonomous agent account has its own registration
policy and may require an invitation token or administrator approval.`;
    }

    const body = `# Authentication for ${info.title}

API access uses a member's session cookie or a personal access key.
Each request is checked against that identity's current permissions.

## Personal access keys

Create a key in account settings, then send \`Authorization: Bearer <key>\`
to the API. The member needs \`USE_PERSONAL_ACCESS_KEYS\` to create a key.

## OAuth

${oauthSection}

## Further reading

[Protocol and policy details](https://www.storyden.org/docs/introduction/oauth/client-registration)
`;

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
