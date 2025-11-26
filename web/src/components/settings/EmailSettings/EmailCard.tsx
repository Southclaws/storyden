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
import { withUndo } from "@/lib/thread/undo";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { CardBox, HStack, WStack } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

type Props = {
  email: AccountEmailAddress;
};
export function EmailCard({ email }: Props) {
  const { mutate } = useSWRConfig();

  async function handleRemove() {
    await handle(
      async () => {
        await withUndo({
          message: "Email address deleted",
          duration: 5000,
          toastId: `email-${email.id}`,
          action: async () => {
            await accountEmailRemove(email.id);
          },
          onUndo: () => {},
        });
      },
      {
        cleanup: async () => {
          await mutate(getAccountGetKey());
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
          <Heading size="sm">{email.email_address}</Heading>
          {email.verified ? (
            <Badge
              borderColor="border.success"
              backgroundColor="bg.success"
              color="fg.success"
            >
              Verified
            </Badge>
          ) : (
            <Link href="/auth/verify/email?returnURL=/settings">
              <Badge borderColor="border.error" backgroundColor="bg.error" color="fg.error">
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
            size="xs"
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
