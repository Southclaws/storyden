import type { Metadata } from "next";

import { LinkButton } from "@/components/ui/link-button";
import { PageHeader } from "@/components/ui/page-header";
import { SectionHeading } from "@/components/ui/section-heading";
import { Text } from "@/components/ui/text";
import {
  getAICatalog,
  getPublicInfo,
  mcpCardURL,
} from "@/lib/discovery/public";
import { LStack, WStack } from "@/styled-system/jsx";

export const metadata: Metadata = {
  title: "Developers",
  description: "Build an integration with this Storyden community.",
};

export default async function Page() {
  const [info, catalog] = await Promise.all([getPublicInfo(), getAICatalog()]);
  const api = new URL(info.apiAddress);
  const card = mcpCardURL(catalog, info.apiAddress);

  return (
    <LStack gap="8">
      <PageHeader
        title="Developers"
        description={`Connect an app or agent to ${info.title}.`}
      />
      <LStack gap="3">
        <SectionHeading>REST API</SectionHeading>
        <Text>
          Explore the API contract and use an access key
          {info.oauthEnabled ? " or OAuth token" : ""} for authenticated
          requests.
        </Text>
        <WStack gap="2" flexWrap="wrap">
          <LinkButton href={new URL("/api/docs", api).href} variant="outline">
            API reference
          </LinkButton>
          <LinkButton
            href={new URL("/api/openapi.json", api).href}
            variant="subtle"
          >
            OpenAPI schema
          </LinkButton>
          <LinkButton href="/auth.md" variant="subtle">
            Authentication guide
          </LinkButton>
        </WStack>
      </LStack>
      {card && (
        <LStack gap="3">
          <SectionHeading>MCP</SectionHeading>
          <Text>
            Connect an MCP client to this community. The Server Card describes
            the endpoint; available tools appear after you connect.
          </Text>
          <WStack gap="2" flexWrap="wrap">
            <LinkButton href={card} variant="outline">
              Server Card
            </LinkButton>
            <LinkButton href="/.well-known/ai-catalog.json" variant="subtle">
              AI catalog
            </LinkButton>
          </WStack>
        </LStack>
      )}
      {catalog === null && (
        <Text variant="supporting">
          MCP discovery is unavailable from this instance.
        </Text>
      )}
      <LStack gap="3">
        <SectionHeading>Further reading</SectionHeading>
        <Text>
          Storyden’s shared documentation covers OAuth clients, access keys, the
          CLI, and the MCP transport.
        </Text>
        <WStack gap="2" flexWrap="wrap">
          <LinkButton
            href="https://www.storyden.org/docs/introduction/oauth"
            variant="subtle"
          >
            OAuth
          </LinkButton>
          <LinkButton
            href="https://www.storyden.org/docs/reference/access-keys"
            variant="subtle"
          >
            Access keys
          </LinkButton>
          <LinkButton
            href="https://www.storyden.org/docs/introduction/cli"
            variant="subtle"
          >
            Storyden CLI
          </LinkButton>
          <LinkButton
            href="https://www.storyden.org/docs/introduction/mcp"
            variant="subtle"
          >
            MCP guide
          </LinkButton>
        </WStack>
      </LStack>
    </LStack>
  );
}
