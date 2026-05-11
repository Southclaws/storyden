import { useEffect, useState } from "react";
import { toast } from "sonner";
import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import { mutateTransaction } from "@/api/mutate";
import {
  getAccountWarningListKey,
  useAccountWarningDelete,
  useAccountWarningList,
  useAccountWarningUpdate,
} from "@/api/openapi-client/accounts";
import {
  AccountWarningListOKResponse,
  ProfileReference,
  Warning,
} from "@/api/openapi-schema";
import { DeletedMemberIdent } from "@/components/member/MemberBadge/DeletedMemberIdent";
import { MemberIdent } from "@/components/member/MemberBadge/MemberIdent";
import { MemberWarningTrigger } from "@/components/member/MemberWarning/MemberWarningTrigger";
import { Timestamp } from "@/components/site/Timestamp";
import { Unready } from "@/components/site/Unready";
import { useConfirmation } from "@/components/site/useConfirmation";
import { Button } from "@/components/ui/button";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { useI18n } from "@/i18n/provider";
import { Box, Flex, HStack, LStack, WStack, styled } from "@/styled-system/jsx";

type WarningsPanelProps = {
  accountId: string;
  profile: ProfileReference;
  canManageWarnings?: boolean;
};

export function WarningsPanel({
  accountId,
  profile,
  canManageWarnings = false,
}: WarningsPanelProps) {
  const { t } = useI18n();
  const { data, error } = useAccountWarningList(accountId);
  if (!data) {
    return <Unready error={error} />;
  }

  return (
    <LStack gap="3">
      {data.total === 0 || data.warnings.length === 0 ? (
        <Box borderWidth="thin" borderRadius="sm" p="2">
          <styled.p fontSize="sm" color="fg.subtle">
            {t("No warning history for this account yet.")}
          </styled.p>
        </Box>
      ) : (
        <LStack gap="2">
          {data.warnings.map((warning) => (
            <WarningRecordCard
              key={warning.id}
              accountId={accountId}
              warning={warning}
              canManageWarnings={canManageWarnings}
            />
          ))}
        </LStack>
      )}

      {canManageWarnings && (
        <WStack justifyContent="flex-end">
          <MemberWarningTrigger profile={profile}>
            <Button size="xs" colorPalette="orange" variant="subtle">
              {t("Issue warning")}
            </Button>
          </MemberWarningTrigger>
        </WStack>
      )}
    </LStack>
  );
}

function WarningRecordCard({
  accountId,
  warning,
  canManageWarnings,
}: {
  accountId: string;
  warning: Warning;
  canManageWarnings: boolean;
}) {
  const { t } = useI18n();
  const { mutate } = useSWRConfig();
  const [isEditing, setIsEditing] = useState(false);
  const [reasonDraft, setReasonDraft] = useState(warning.reason ?? "");

  useEffect(() => {
    setReasonDraft(warning.reason ?? "");
  }, [warning.id, warning.reason]);

  const { trigger: updateWarning, isMutating: isUpdating } =
    useAccountWarningUpdate(accountId, warning.id);
  const { trigger: deleteWarning, isMutating: isDeleting } =
    useAccountWarningDelete(accountId, warning.id);

  const { isConfirming, handleConfirmAction, handleCancelAction } =
    useConfirmation(async () => {
      await handle(
        async () => {
          await mutateTransaction(
            mutate,
            [
              {
                key: getAccountWarningListKey(accountId),
                optimistic: (
                  current: AccountWarningListOKResponse | undefined,
                ) => {
                  if (!current) {
                    return current;
                  }

                  return {
                    ...current,
                    total: Math.max(0, current.total - 1),
                    warnings: current.warnings.filter(
                      (item) => item.id !== warning.id,
                    ),
                  };
                },
              },
            ],
            () => deleteWarning({}),
            { revalidate: true },
          );
        },
        {
          promiseToast: {
            loading: t("Deleting warning..."),
            success: t("Warning deleted."),
          },
        },
      );
    });

  async function handleSave() {
    const trimmed = reasonDraft.trim();
    if (!trimmed) {
      toast.error(t("Reason cannot be empty."));
      return;
    }

    await handle(
      async () => {
        await mutateTransaction(
          mutate,
          [
            {
              key: getAccountWarningListKey(accountId),
              optimistic: (
                current: AccountWarningListOKResponse | undefined,
              ) => {
                if (!current) {
                  return current;
                }

                return {
                  ...current,
                  warnings: current.warnings.map((item) =>
                    item.id === warning.id
                      ? { ...item, reason: trimmed }
                      : item,
                  ),
                };
              },
            },
          ],
          () => updateWarning({ reason: trimmed }),
          { revalidate: true },
        );

        setIsEditing(false);
      },
      {
        promiseToast: {
          loading: t("Saving warning..."),
          success: t("Warning updated."),
        },
      },
    );
  }

  function cancelEditing() {
    setReasonDraft(warning.reason ?? "");
    setIsEditing(false);
  }

  return (
    <Box
      borderWidth="thin"
      borderRadius="sm"
      bgColor="bg.subtle"
      p="2"
      w="full"
    >
      <Flex
        alignItems={{ base: "start", md: "center" }}
        justifyContent="space-between"
        direction={{ base: "column", md: "row" }}
        gap="2"
        mb="2"
      >
        <HStack gap="2" alignItems="center" flexWrap="wrap">
          {warning.issued_by ? (
            <MemberIdent
              profile={warning.issued_by}
              size="sm"
              name="full-horizontal"
            />
          ) : (
            <DeletedMemberIdent size="sm" label="hidden" />
          )}
        </HStack>

        <HStack gap="1" flexWrap="wrap">
          <Timestamp
            created={warning.issued_at}
            color="fg.muted"
            fontSize="xs"
          />

          {canManageWarnings && (
            <>
              {isEditing ? (
                <>
                  <Button
                    size="xs"
                    variant="subtle"
                    onClick={handleSave}
                    loading={isUpdating}
                  >
                    {t("Save")}
                  </Button>
                  <Button
                    size="xs"
                    variant="ghost"
                    onClick={cancelEditing}
                    disabled={isUpdating}
                  >
                    {t("Cancel")}
                  </Button>
                </>
              ) : (
                <Button
                  size="xs"
                  variant="ghost"
                  onClick={() => setIsEditing(true)}
                >
                  {t("Edit")}
                </Button>
              )}

              {isConfirming ? (
                <HStack gap="1">
                  <Button
                    size="xs"
                    variant="subtle"
                    bgColor="bg.destructive"
                    onClick={handleConfirmAction}
                    loading={isDeleting}
                  >
                    {t("Confirm delete")}
                  </Button>
                  <Button
                    size="xs"
                    variant="ghost"
                    onClick={handleCancelAction}
                    disabled={isDeleting}
                  >
                    {t("Cancel")}
                  </Button>
                </HStack>
              ) : (
                <Button
                  size="xs"
                  variant="ghost"
                  aria-label={t("Delete warning")}
                  title={t("Delete warning")}
                  onClick={handleConfirmAction}
                >
                  <DeleteIcon />
                </Button>
              )}
            </>
          )}
        </HStack>
      </Flex>

      {isEditing && canManageWarnings ? (
        <styled.textarea
          value={reasonDraft}
          onChange={(event) => setReasonDraft(event.target.value)}
          rows={4}
          maxLength={2000}
          borderWidth="thin"
          borderRadius="sm"
          p="2"
          fontSize="sm"
          w="full"
        />
      ) : (
        <styled.p fontSize="sm" whiteSpace="pre-wrap">
          {warning.reason?.trim() || t("No reason recorded.")}
        </styled.p>
      )}
    </Box>
  );
}
