import { AccountAuthMethod, AuthProvider } from "src/api/openapi/schemas";
import { Heading3 } from "src/theme/components/Heading/Index";
import { Link } from "src/theme/components/Link";

import { VStack, styled } from "@/styled-system/jsx";

type Props = {
  active: AccountAuthMethod[];
  available: AuthProvider[];
};

export function OAuth({ active, available }: Props) {
  return (
    <VStack alignItems="start">
      <Heading3 size="sm">Linked accounts</Heading3>

      <styled.p>
        You can link as many accounts as you want to. Linked accounts allow you
        to log in easily and may also provide additional features.
      </styled.p>

      <Heading3 size="xs">Active</Heading3>

      {active.length ? (
        <styled.ul display="flex" flexDir="column" gap="2" w="full">
          {active.map((v) => (
            <styled.li key={v.id}>
              <styled.span>{v.name}</styled.span>
            </styled.li>
          ))}
        </styled.ul>
      ) : (
        <styled.p>You currently have no linked accounts.</styled.p>
      )}

      <Heading3 size="xs">Available</Heading3>

      <styled.ul display="flex" flexDir="column" gap="2" w="full">
        {available.map((v) => (
          <styled.li key={v.provider}>
            <Link href={v.link} variant="subtle">
              {v.name}
            </Link>
          </styled.li>
        ))}
      </styled.ul>
    </VStack>
  );
}
