"use client";

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
          <Heading size="sm">{email.email_address}</Heading>
          {email.verified && (
            <Badge
              borderColor="green.6"
              backgroundColor="green.5"
              color="green.11"
            >
              Verified
            </Badge>
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
