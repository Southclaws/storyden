import { formatDate } from "date-fns";
import { match } from "ts-pattern";

import { RequestError } from "@/api/common";
import { useAccountView } from "@/api/openapi-client/accounts";
import { MemberIdent } from "@/components/member/MemberBadge/MemberIdent";
import { Badge } from "@/components/ui/badge";
import { AdminIcon } from "@/components/ui/icons/Admin";
import { WarningIcon } from "@/components/ui/icons/Warning";
import {
  Box,
  CardBox,
  Flex,
  HStack,
  LStack,
  styled,
} from "@/styled-system/jsx";
import { deriveError } from "@/utils/error";

import { AccountPurgeTrigger } from "./AccountPurgeModal/AccountPurgeTrigger";

type Props = {
  accountId: string;
};

function useProfileAccountManagement({ accountId }: Props) {
  const { data: account, error } = useAccountView(accountId);
  if (!account) {
    if (
      error !== undefined &&
      error instanceof RequestError &&
      error.status === 403
    ) {
      return {
        ready: false as const,
        error:
          "You do not have permission to view additional details for an administrator account.",
      };
    }

    return {
      ready: false as const,
      error: deriveError(error),
    };
  }

  const emailVerifiedStatus =
    account.email_addresses.length === 0
      ? "not_applicable"
      : account.email_addresses.some((e) => e.verified)
        ? "verified"
        : "not_verified";

  return {
    ready: true as const,
    account,
    emailVerifiedStatus,
  };
}

export function ProfileAccountManagement({ accountId }: Props) {
  const { ready, error, account, emailVerifiedStatus } =
    useProfileAccountManagement({ accountId });
  if (!ready || !account) {
    return (
      <CardBox
        borderColor="border.warning"
        borderWidth="thin"
        borderStyle="dashed"
        borderRadius="sm"
        p="2"
      >
        <HStack
          alignItems="center"
          color="fg.subtle"
          role="alert"
          aria-atomic="true"
        >
          <Box w="5" flexShrink="0">
            <WarningIcon aria-hidden="true" />
          </Box>
          <p id="error__message">{error}</p>
        </HStack>
      </CardBox>
    );
  }

  const emailVerifiedStatusBadge = match(emailVerifiedStatus)
    .with("not_applicable", () => (
      <Badge colorPalette="gray">No emails to verify</Badge>
    ))
    .with("verified", () => <Badge colorPalette="green">Verified</Badge>)
    .with("not_verified", () => <Badge colorPalette="gray">Unverified</Badge>);

  return (
    <CardBox
      p="0"
      borderColor="border.warning"
      borderWidth="thin"
      borderStyle="dashed"
      borderRadius="sm"
    >
      <Box bgColor="bg.warning" borderTopRadius="sm" pl="3" pr="2" py="1">
        <HStack
          gap="1"
          color="fg.warning"
          fontSize="xs"
          justifyContent="space-between"
        >
          <HStack gap="1">
            <AdminIcon w="4" />
            <p>Account information</p>
          </HStack>
          <AccountPurgeTrigger accountId={account.id} handle={account.handle} />
        </HStack>
      </Box>
      <Flex
        p="3"
        gap="4"
        direction={{ base: "column", md: "row" }}
        alignItems="start"
      >
        <LStack flex="1" gap="3" flexShrink="1" flexGrow="1" minW="0">
          <LStack gap="1">
            <styled.p fontSize="xs" fontWeight="semibold" color="fg.muted">
              Account Status
            </styled.p>
            <Box fontSize="sm">{emailVerifiedStatusBadge.run()}</Box>
          </LStack>

          <LStack gap="1">
            <styled.p fontSize="xs" fontWeight="semibold" color="fg.muted">
              Joined at
            </styled.p>
            <styled.p fontSize="sm">
              {formatDate(new Date(account.joined), "PPPppp")}
            </styled.p>
          </LStack>

          {account.suspended && (
            <LStack gap="1">
              <styled.p fontSize="xs" fontWeight="semibold" color="fg.muted">
                Suspended
              </styled.p>
              <styled.p fontSize="sm" color="fg.destructive">
                {formatDate(new Date(account.suspended), "PPPppp")}
              </styled.p>
            </LStack>
          )}

          {account.invited_by && (
            <LStack gap="1">
              <styled.p fontSize="xs" fontWeight="semibold" color="fg.muted">
                Invited By
              </styled.p>
              <MemberIdent
                size="sm"
                name="full-vertical"
                profile={account.invited_by}
              />
            </LStack>
          )}
        </LStack>

        <LStack flex="1" gap="3" flexShrink="1" flexGrow="1" minW="0">
          <LStack gap="1" minW="0">
            <styled.p fontSize="xs" fontWeight="semibold" color="fg.muted">
              Email Addresses
            </styled.p>
            {account.email_addresses.length > 0 ? (
              <LStack gap="1" minW="0">
                {account.email_addresses.map((email) => (
                  <HStack
                    key={email.id}
                    gap="2"
                    fontSize="sm"
                    flexWrap="wrap"
                    minW="0"
                    width="full"
                  >
                    <styled.code
                      fontFamily="mono"
                      w="full"
                      minW="0"
                      textOverflow="ellipsis"
                      overflow="hidden"
                    >
                      {email.email_address}
                    </styled.code>
                    {email.verified ? (
                      <Badge colorPalette="green" size="sm">
                        Verified
                      </Badge>
                    ) : (
                      <Badge colorPalette="gray" size="sm">
                        Unverified
                      </Badge>
                    )}
                  </HStack>
                ))}
              </LStack>
            ) : (
              <styled.p fontSize="sm" color="fg.subtle">
                No email addresses
              </styled.p>
            )}
          </LStack>
        </LStack>
      </Flex>
    </CardBox>
  );
}
