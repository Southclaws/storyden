import { AuthProvider } from "src/api/openapi-schema";

import { LinkButton } from "@/components/ui/link-button";
import { styled } from "@/styled-system/jsx";

import { getProviders } from "../../../../lib/auth/providers";

export async function AuthSelection() {
  const { oauth: providers } = await getProviders();

  // sort by alphabetical because lazy
  // TODO: allow the order to be configured by the admin.
  providers?.sort((a, b) => a.provider.localeCompare(b.provider));

  return (
    <styled.ul
      w="full"
      display="flex"
      flexDir="column"
      gap="2"
      alignItems="center"
    >
      <p>AUTH SELECTION</p>

      {providers?.map((v: AuthProvider) => (
        <styled.li w="full" key={v.provider}>
          <LinkButton
            size="sm"
            variant="ghost"
            w="full"
            href={v.link ?? `/auth/${v.provider}`}
          >
            {v.name}
          </LinkButton>
        </styled.li>
      ))}
    </styled.ul>
  );
}
