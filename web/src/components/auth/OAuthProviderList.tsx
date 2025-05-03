import { AuthProvider } from "@/api/openapi-schema";
import { LinkButton } from "@/components/ui/link-button";
import { Text } from "@/components/ui/text";
import { OAuthProvider, filterWithLink } from "@/lib/auth/oauth";
import { getProviders } from "@/lib/auth/providers";
import { Divider, VStack, WStack, styled } from "@/styled-system/jsx";

type Props = {
  providers: AuthProvider[];
};

export async function OAuthProviderList({ providers }: Props) {
  const oauth = filterWithLink(providers);

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
            <OAuthProviderLink provider={v} />
          </styled.li>
        ))}
      </styled.ul>
    </VStack>
  );
}

function OAuthProviderLink({ provider }: { provider: OAuthProvider }) {
  return (
    <LinkButton size="sm" variant="outline" w="full" href={provider.link}>
      {provider.name}
    </LinkButton>
  );
}
