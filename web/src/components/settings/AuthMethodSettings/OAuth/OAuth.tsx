import { AccountAuthMethod, AuthProvider } from "@/api/openapi-schema";
import { Timestamp } from "@/components/site/Timestamp";
import { CardBox } from "@/components/ui/card-box";
import { LinkButton } from "@/components/ui/link-button";
import { SectionHeading } from "@/components/ui/section-heading";
import { Text } from "@/components/ui/text";
import { OAuthProvider } from "@/lib/auth/oauth";
import { LStack, VStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

type Props = {
  active: AccountAuthMethod[];
  available: OAuthProvider[];
};

export function OAuth({ active, available }: Props) {
  return (
    <LStack>
      <SectionHeading>Linked accounts</SectionHeading>

      <Text variant="supporting">
        You can link as many accounts as you want to. Linked accounts allow you
        to log in easily and may also provide additional features.
      </Text>

      <SectionHeading>Active</SectionHeading>

      {active.length ? (
        <styled.ul className={lstack()}>
          {active.map((v) => (
            <styled.li key={v.id} w="full">
              <CardBox>
                <Text
                  variant="supporting"
                  color="text.default"
                  fontWeight="semibold"
                >
                  {v.name}
                </Text>

                <WStack color="text.subtle" alignItems="end">
                  <span>
                    Added&nbsp;
                    <Timestamp created={v.created_at} large />
                  </span>

                  <styled.pre fontSize="sm">id:{v.identifier}</styled.pre>
                </WStack>
              </CardBox>
            </styled.li>
          ))}
        </styled.ul>
      ) : (
        <Text variant="supporting">You currently have no linked accounts.</Text>
      )}

      <SectionHeading>Available</SectionHeading>

      {available.length ? (
        <styled.ul className={lstack()}>
          {available.map((v) => (
            <styled.li key={v.provider} w="full">
              <CardBox>
                <WStack alignItems="center">
                  <Text
                    variant="supporting"
                    color="text.default"
                    fontWeight="semibold"
                  >
                    {v.name}
                  </Text>

                  <LinkButton href={v.link} variant="subtle" size="sm">
                    Link with {v.name}
                  </LinkButton>
                </WStack>
              </CardBox>
            </styled.li>
          ))}
        </styled.ul>
      ) : (
        <Text variant="supporting">
          {active.length > 0
            ? "There are no more authentication providers available."
            : "There are currently no authentication providers available."}
        </Text>
      )}
    </LStack>
  );
}
