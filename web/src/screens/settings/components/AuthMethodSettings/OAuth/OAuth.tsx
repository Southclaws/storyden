import { AccountAuthMethod, AuthProvider } from "src/api/openapi/schemas";

import { Heading } from "@/components/ui/heading";
import { LinkButton } from "@/components/ui/link-button";
import { VStack, styled } from "@/styled-system/jsx";

type Props = {
  active: AccountAuthMethod[];
  available: AuthProvider[];
};

export function OAuth({ active, available }: Props) {
  return (
    <VStack alignItems="start">
      <Heading size="sm">Linked accounts</Heading>

      <styled.p>
        You can link as many accounts as you want to. Linked accounts allow you
        to log in easily and may also provide additional features.
      </styled.p>

      <Heading size="xs">Active</Heading>

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

      <Heading size="xs">Available</Heading>

      <styled.ul display="flex" flexDir="column" gap="2" w="full">
        {available.map((v) => (
          <styled.li key={v.provider}>
            <LinkButton href={v.link} variant="subtle">
              {v.name}
            </LinkButton>
          </styled.li>
        ))}
      </styled.ul>
    </VStack>
  );
}
