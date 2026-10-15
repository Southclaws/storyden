import { toast } from "sonner";
import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import {
  getAccountWarningListKey,
  useAccountWarningCreate,
} from "@/api/openapi-client/accounts";
import { ProfileReference } from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";
import { Text } from "@/components/ui/text";
import { Textarea } from "@/components/ui/textarea";
import { HStack, VStack } from "@/styled-system/jsx";

type Props = {
  profile: ProfileReference;
  reason: string;
  onReasonChange: (reason: string) => void;
  onClose: () => void;
};

export function MemberWarningScreen({
  profile,
  reason,
  onReasonChange,
  onClose,
}: Props) {
  const { mutate } = useSWRConfig();
  const { trigger: createWarning, isMutating: loading } =
    useAccountWarningCreate(profile.id);

  async function issueWarning() {
    if (!reason.trim()) {
      toast.error("Please provide a warning reason.");
      return;
    }

    await handle(async () => {
      await createWarning({
        reason: reason.trim(),
      });
      await mutate(getAccountWarningListKey(profile.id));
      toast.success(`Warning issued to ${profile.name}.`);
      onReasonChange("");
      onClose();
    });
  }

  return (
    <VStack alignItems="start" gap="3">
      <Text variant="supporting">
        Warnings are recorded for internal moderation history.
      </Text>
      <Textarea
        rows={5}
        value={reason}
        onChange={(e) => onReasonChange(e.currentTarget.value)}
        placeholder="Clear, specific reason for this warning"
      />

      <HStack w="full">
        <Button type="button" flexGrow="1" onClick={onClose} disabled={loading}>
          Cancel
        </Button>
        <Button
          type="button"
          flexGrow="1"
          intent="warning"
          variant="solid"
          onClick={issueWarning}
          loading={loading}
        >
          Issue warning
        </Button>
      </HStack>
    </VStack>
  );
}
