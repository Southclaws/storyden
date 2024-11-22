import { AccountAuthMethod, AuthProvider } from "src/api/openapi-schema";

import { Timestamp } from "@/components/site/Timestamp";
import { Heading } from "@/components/ui/heading";
import { LinkButton } from "@/components/ui/link-button";
import { OAuthProvider } from "@/lib/auth/oauth";
import { CardBox, LStack, VStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

type Props = {
  active: AccountAuthMethod[];
  available: OAuthProvider[];
};

export function OAuth({ active, available }: Props) {
  return (
    <LStack>
      <Heading size="sm">Linked accounts</Heading>

      <styled.p>
        You can link as many accounts as you want to. Linked accounts allow you
        to log in easily and may also provide additional features.
      </styled.p>

      <Heading size="sm" color="fg.subtle">
        Active
      </Heading>

      {active.length ? (
        <styled.ul className={lstack()}>
          {active.map((v) => (
            <styled.li key={v.id} w="full">
              <CardBox>
                <Heading size="sm">{v.name}</Heading>

                <WStack color="fg.muted" alignItems="end">
                  <styled.span>
                    Added&nbsp;
                    <Timestamp created={v.created_at} large />
                  </styled.span>

                  <styled.pre fontSize="sm">id:{v.identifier}</styled.pre>
                </WStack>
              </CardBox>
            </styled.li>
          ))}
        </styled.ul>
      ) : (
        <styled.p color="fg.muted">
          You currently have no linked accounts.
        </styled.p>
      )}

      <Heading size="sm" color="fg.subtle">
        Available
      </Heading>

      {available.length ? (
        <styled.ul className={lstack()}>
          {available.map((v) => (
            <styled.li key={v.provider} w="full">
              <CardBox>
                <WStack alignItems="center">
                  <Heading size="sm">{v.name}</Heading>

                  <LinkButton href={v.link} variant="subtle" size="sm">
                    Link with {v.name}
                  </LinkButton>
                </WStack>
              </CardBox>
            </styled.li>
          ))}
        </styled.ul>
      ) : (
        <styled.p color="fg.muted">
          {active.length > 0
            ? "There are no more authentication providers available."
            : "There are currently no authentication providers available."}
        </styled.p>
      )}
    </LStack>
  );
}
