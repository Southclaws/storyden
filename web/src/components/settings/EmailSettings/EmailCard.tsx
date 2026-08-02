"use client";

import Link from "next/link";
import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import {
  accountEmailRemove,
  getAccountGetKey,
} from "@/api/openapi-client/accounts";
import { AccountEmailAddress } from "@/api/openapi-schema";
import { CancelAction } from "@/components/site/Action/Cancel";
import { useConfirmation } from "@/components/site/useConfirmation";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { CardBox } from "@/components/ui/card-box";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { Text } from "@/components/ui/text";
import { HStack, WStack } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

type Props = {
  email: AccountEmailAddress;
};
export function EmailCard({ email }: Props) {
  const { mutate } = useSWRConfig();

  async function handleRemove() {
    await handle(
      async () => {
        await accountEmailRemove(email.id);
      },
      {
        async cleanup() {
          await mutate(getAccountGetKey());
        },
        promiseToast: {
          loading: "Deleting email address...",
          success: "Email deleted",
        },
      },
    );
  }

  const { isConfirming, handleCancelAction, handleConfirmAction } =
    useConfirmation(handleRemove);

  return (
    <CardBox key={email.email_address} className={lstack()} gap="4">
      <WStack alignItems="center">
        <HStack>
          <Text variant="supporting" color="text.default" fontWeight="semibold">
            {email.email_address}
          </Text>
          {email.verified ? (
            <Badge
              borderColor="status.success.border"
              backgroundColor="status.success.surface"
              color="status.success.content"
            >
              Verified
            </Badge>
          ) : (
            <Link href="/auth/verify/email?returnURL=/settings">
              <Badge
                borderColor="status.danger.border"
                backgroundColor="status.danger.surface"
                color="status.danger.content"
              >
                Verify this email
              </Badge>
            </Link>
          )}
        </HStack>

        <HStack gap="0">
          <Button
            style={{
              borderBottomRightRadius: isConfirming ? "0" : undefined,
              borderTopRightRadius: isConfirming ? "0" : undefined,
            }}
            variant="subtle"
            onClick={handleConfirmAction}
          >
            {isConfirming ? (
              "Are you sure?"
            ) : (
              <>
                <DeleteIcon /> delete email
              </>
            )}
          </Button>

          {isConfirming && (
            <CancelAction
              variant="subtle"
              borderLeftRadius="none"
              onClick={handleCancelAction}
            />
          )}
        </HStack>
      </WStack>
    </CardBox>
  );
}
