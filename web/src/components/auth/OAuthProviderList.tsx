import { filter } from "lodash/fp";

import { AuthProvider } from "src/api/openapi-schema";

import { LinkButton } from "@/components/ui/link-button";
import { Text } from "@/components/ui/text";
import { getProviders } from "@/lib/auth/providers";
import { Divider, VStack, WStack, styled } from "@/styled-system/jsx";

const hasLink = (provider: AuthProvider): provider is WithLink => {
  return provider.link !== undefined;
};

const filterWithLink = (list: AuthProvider[]): WithLink[] =>
  filter(hasLink)(list);

type WithLink = AuthProvider & { link: string };

export async function OAuthProviderList() {
  const { oauth: all } = await getProviders();

  const oauth = filterWithLink(all);

  if (oauth.length === 0) {
    return null;
  }

  // sort by alphabetical because lazy
  // TODO: allow the order to be configured by the admin.
  oauth?.sort((a, b) => a.provider.localeCompare(b.provider));

  return (
    <VStack w="full">
      <WStack alignItems="center" textWrap="nowrap" color="fg.subtle">
        <Divider />
        <Text>or via third party</Text>
        <Divider />
      </WStack>

      <styled.ul
        w="full"
        display="flex"
        flexDir="column"
        gap="2"
        alignItems="center"
      >
        {oauth.map((v) => (
          <styled.li w="full" key={v.provider}>
            <OAuthProvider provider={v} />
          </styled.li>
        ))}
      </styled.ul>
    </VStack>
  );
}

function OAuthProvider({ provider }: { provider: WithLink }) {
  return (
    <LinkButton size="sm" variant="outline" w="full" href={provider.link}>
      {provider.name}
    </LinkButton>
  );
}
