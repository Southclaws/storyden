import React, { ReactNode, useState } from "react";
import { toast } from "sonner";
import { useSWRConfig } from "swr";

import { ModalDrawer } from "src/components/site/Modaldrawer/Modaldrawer";
import { useDisclosure } from "src/utils/useDisclosure";

import { handle } from "@/api/client";
import {
  getAccountWarningListKey,
  useAccountWarningCreate,
} from "@/api/openapi-client/accounts";
import { ProfileReference } from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";
import { useI18n } from "@/i18n/provider";
import { HStack, VStack, styled } from "@/styled-system/jsx";

type MemberWarningTriggerProps = {
  children?: ReactNode;
  profile: ProfileReference;
};

export function MemberWarningTrigger({
  children,
  profile,
}: MemberWarningTriggerProps) {
  const { mutate } = useSWRConfig();
  const { onOpen, onClose, isOpen } = useDisclosure();
  const { t } = useI18n();
  const { trigger: createWarning, isMutating: loading } =
    useAccountWarningCreate(profile.id);
  const [reason, setReason] = useState("");

  async function issueWarning() {
    if (!reason.trim()) {
      toast.error(t("Please provide a warning reason."));
      return;
    }

    await handle(async () => {
      await createWarning({
        reason: reason.trim(),
      });
      await mutate(getAccountWarningListKey(profile.id));
      toast.success(t("Warning issued to {{name}}.", { name: profile.name }));
      setReason("");
      onClose();
    });
  }

  const triggerNode = React.isValidElement<{
    onClick?: React.MouseEventHandler;
  }>(children) ? (
    React.cloneElement(children, { onClick: onOpen })
  ) : (
    <Button colorPalette="orange" onClick={onOpen}>
      {t("Warn")}
    </Button>
  );

  return (
    <>
      {triggerNode}
      <ModalDrawer
        isOpen={isOpen}
        onClose={onClose}
        title={t("Issue warning to {{name}}", { name: profile.name })}
      >
        <VStack alignItems="start" gap="3">
          <styled.p fontSize="sm" color="fg.subtle">
            {t("Warnings are recorded for internal moderation history.")}
          </styled.p>
          <styled.textarea
            rows={5}
            value={reason}
            onChange={(e) => setReason(e.currentTarget.value)}
            placeholder={t("Clear, specific reason for this warning")}
            width="full"
            borderWidth="thin"
            borderRadius="sm"
            borderColor="border.default"
            padding="2"
          />

          <HStack w="full">
            <Button
              type="button"
              flexGrow="1"
              onClick={onClose}
              disabled={loading}
            >
              {t("Cancel")}
            </Button>
            <Button
              type="button"
              flexGrow="1"
              colorPalette="orange"
              onClick={issueWarning}
              loading={loading}
            >
              {t("Issue warning")}
            </Button>
          </HStack>
        </VStack>
      </ModalDrawer>
    </>
  );
}
