import { AuthProvider } from "src/api/openapi/schemas";
import { Link } from "src/theme/components/Link";

import { getProviders } from "../../providers";

import { styled } from "@/styled-system/jsx";

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
          <Link
            size="sm"
            kind="neutral"
            w="full"
            href={v.link ?? `/auth/${v.provider}`}
          >
            {v.name}
          </Link>
        </styled.li>
      ))}
    </styled.ul>
  );
}
