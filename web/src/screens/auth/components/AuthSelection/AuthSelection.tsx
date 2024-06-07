import { AuthProvider } from "src/api/openapi/schemas";

import { LinkButton } from "@/components/ui/link-button";
import { styled } from "@/styled-system/jsx";

import { getProviders } from "../../providers";

export async function AuthSelection() {
  const { providers } = await getProviders();

  // sort by alphabetical because lazy
  // TODO: allow the order to be configured by the admin.
  providers?.sort((a, b) => a.provider.localeCompare(b.provider));

  return (
    <styled.ul
      display="flex"
      flexDir="column"
      gap="2"
      alignItems="center"
      w="1/2"
    >
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
