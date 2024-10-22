import { redirect } from "next/navigation";

import { oAuthProviderCallback } from "src/api/openapi-server/auth";
import { UnreadyBanner } from "src/components/site/Unready";

import { OAuthCallback } from "@/api/openapi-schema";
import { LinkButton } from "@/components/ui/link-button";
import { HStack, VStack } from "@/styled-system/jsx";

export type Props = {
  params: Promise<{
    provider: string;
  }>;
  searchParams: Promise<OAuthCallback>;
};

export default async function Page(props: Props) {
  try {
    const { provider } = await props.params;
    const query = await props.searchParams;
    const { data } = await oAuthProviderCallback(provider, query);

    const { id } = data;

    redirect(`/?id=${id}`);
  } catch (e) {
    return (
      <VStack w="full" height="dvh" justify="center" p="10">
        <UnreadyBanner error={e} />
        <HStack>
          <LinkButton href="/register" variant="outline">
            Back to register
          </LinkButton>
          <LinkButton href="/login" variant="outline">
            Back to login
          </LinkButton>
        </HStack>
      </VStack>
    );
  }
}
