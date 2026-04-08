import { formatDate } from "date-fns";
import { useEffect, useState } from "react";
import { useSWRConfig } from "swr";
import { match } from "ts-pattern";

import { RequestError } from "@/api/common";
import { mutateTransaction } from "@/api/mutate";
import {
  getAccountModerationNoteListKey,
  useAccountModerationNoteCreate,
  useAccountModerationNoteDelete,
  useAccountModerationNoteList,
  useAccountView,
} from "@/api/openapi-client/accounts";
import {
  AccountModerationNoteListOKResponse,
  ModerationNote,
  Permission,
} from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { DeletedMemberIdent } from "@/components/member/MemberBadge/DeletedMemberIdent";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { MemberIdent } from "@/components/member/MemberBadge/MemberIdent";
import { Timestamp } from "@/components/site/Timestamp";
import { useConfirmation } from "@/components/site/useConfirmation";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { AdminIcon } from "@/components/ui/icons/Admin";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { WarningIcon } from "@/components/ui/icons/Warning";
import * as Tabs from "@/components/ui/tabs";
import {
  Box,
  CardBox,
  Flex,
  HStack,
  LStack,
  WStack,
  styled,
} from "@/styled-system/jsx";
import { deriveError } from "@/utils/error";
import { hasPermission } from "@/utils/permissions";

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
  const session = useSession();
  const canViewModerationNotes = hasPermission(
    session,
    Permission.VIEW_MODERATION_NOTES,
  );
  const canManageModerationNotes = hasPermission(
    session,
    Permission.MANAGE_MODERATION_NOTES,
  );
  const hasModerationNotesAccess =
    canViewModerationNotes || canManageModerationNotes;

  const [noteDraft, setNoteDraft] = useState("");

  useEffect(() => {
    setNoteDraft("");
  }, [accountId]);

  const { ready, error, account, emailVerifiedStatus } =
    useProfileAccountManagement({ accountId });
  const {
    data: notesData,
    error: notesError,
    isLoading: notesLoading,
    mutate: mutateNotes,
  } = useAccountModerationNoteList(accountId, {
    swr: {
      enabled: canViewModerationNotes,
    },
  });
  const {
    trigger: createNote,
    isMutating,
    error: createNoteError,
  } = useAccountModerationNoteCreate(accountId);

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

  async function submitNote() {
    if (!canManageModerationNotes) {
      return;
    }

    const content = noteDraft.trim();
    if (!content) {
      return;
    }

    try {
      await createNote({ content });
      setNoteDraft("");

      if (canViewModerationNotes) {
        await mutateNotes();
      }
    } catch {
      // useSWRMutation stores errors on the hook state, no extra action needed
    }
  }

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
      <Tabs.Root
        variant="line"
        size="sm"
        pt="2"
        colorPalette="amber"
        defaultValue="overview"
        w="full"
        lazyMount={true}
      >
        <Tabs.List>
          <Tabs.Trigger value="overview">Overview</Tabs.Trigger>
          {hasModerationNotesAccess && (
            <Tabs.Trigger value="moderator_notes">Moderator Notes</Tabs.Trigger>
          )}
          <Tabs.Indicator />
        </Tabs.List>

        <Tabs.Content value="overview" px="3" pb="3">
          <Flex
            gap={{ base: "4", md: "6" }}
            direction={{ base: "column", md: "row" }}
            alignItems="start"
          >
            <LStack flex="1" gap="4" flexShrink="1" flexGrow="1" minW="0">
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
                  <styled.p
                    fontSize="xs"
                    fontWeight="semibold"
                    color="fg.muted"
                  >
                    Suspended
                  </styled.p>
                  <styled.p fontSize="sm" color="fg.destructive">
                    {formatDate(new Date(account.suspended), "PPPppp")}
                  </styled.p>
                </LStack>
              )}

              {account.invited_by && (
                <LStack gap="1">
                  <styled.p
                    fontSize="xs"
                    fontWeight="semibold"
                    color="fg.muted"
                  >
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

            <LStack flex="1" gap="2" flexShrink="1" flexGrow="1" minW="0">
              <styled.p fontSize="xs" fontWeight="semibold" color="fg.muted">
                Email Addresses
              </styled.p>
              {account.email_addresses.length > 0 ? (
                <LStack gap="2" minW="0">
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
          </Flex>
        </Tabs.Content>

        {hasModerationNotesAccess && (
          <Tabs.Content value="moderator_notes" px="3" pb="3">
            <LStack gap="2" minW="0">
              {canManageModerationNotes && (
                <LStack gap="1">
                  <styled.textarea
                    value={noteDraft}
                    onChange={(e) => setNoteDraft(e.target.value)}
                    placeholder="Add an internal moderator note…"
                    aria-label="Internal moderator note"
                    rows={3}
                    maxLength={2000}
                    borderWidth="thin"
                    borderRadius="sm"
                    p="2"
                    fontSize="sm"
                    w="full"
                  />
                  <WStack justifyContent="end">
                    <Button
                      size="xs"
                      variant="subtle"
                      onClick={submitNote}
                      loading={isMutating}
                      disabled={!noteDraft.trim()}
                      alignSelf="start"
                    >
                      Add note
                    </Button>
                  </WStack>
                </LStack>
              )}

              {createNoteError ? (
                <styled.p fontSize="sm" color="fg.destructive">
                  {deriveError(createNoteError)}
                </styled.p>
              ) : null}

              {canViewModerationNotes ? (
                <LStack gap="3" w="full" minW="0">
                  {notesError ? (
                    <styled.p fontSize="sm" color="fg.destructive">
                      {deriveError(notesError)}
                    </styled.p>
                  ) : notesLoading ? (
                    <styled.p fontSize="sm" color="fg.subtle">
                      Loading moderation notes...
                    </styled.p>
                  ) : (notesData?.notes?.length ?? 0) === 0 ? (
                    <styled.p fontSize="sm" color="fg.subtle">
                      No moderation notes yet.
                    </styled.p>
                  ) : (
                    (notesData?.notes ?? []).map((note) => (
                      <ModerationNoteCard
                        key={note.id}
                        accountId={accountId}
                        note={note}
                        canManageModerationNotes={canManageModerationNotes}
                      />
                    ))
                  )}
                </LStack>
              ) : (
                <styled.p fontSize="sm" color="fg.subtle">
                  You can add notes, but you do not have permission to view note
                  history.
                </styled.p>
              )}
            </LStack>
          </Tabs.Content>
        )}
      </Tabs.Root>
    </CardBox>
  );
}

type ModerationNoteCardProps = {
  accountId: string;
  note: ModerationNote;
  canManageModerationNotes: boolean;
};

function ModerationNoteCard({
  accountId,
  note,
  canManageModerationNotes,
}: ModerationNoteCardProps) {
  const { mutate } = useSWRConfig();
  const {
    trigger: deleteNote,
    isMutating,
    error: deleteNoteError,
  } = useAccountModerationNoteDelete(accountId, note.id);

  const { isConfirming, handleConfirmAction, handleCancelAction } =
    useConfirmation(async () => {
      await mutateTransaction(
        mutate,
        [
          {
            key: getAccountModerationNoteListKey(accountId),
            optimistic: (
              current: AccountModerationNoteListOKResponse | undefined,
            ) => {
              if (!current) {
                return current;
              }

              return {
                ...current,
                notes: current.notes.filter((moderationNote) => {
                  return moderationNote.id !== note.id;
                }),
              };
            },
          },
        ],
        () => deleteNote({}),
        { revalidate: true },
      );
    });

  return (
    <Box
      borderWidth="thin"
      borderRadius="sm"
      p="2"
      bgColor="bg.subtle"
      w="full"
    >
      <Flex
        justifyContent="space-between"
        alignItems={{ base: "start", md: "center" }}
        direction={{ base: "column", md: "row" }}
        gap="2"
        mb="2"
      >
        {note.author ? (
          <MemberBadge profile={note.author} size="sm" name="full-horizontal" />
        ) : (
          <DeletedMemberIdent size="sm" />
        )}

        <HStack gap="2" alignItems="center" flexShrink="0">
          <Timestamp created={note.created_at} color="fg.muted" fontSize="xs" />

          {canManageModerationNotes &&
            (isConfirming ? (
              <HStack gap="1">
                <Button
                  size="xs"
                  variant="subtle"
                  bgColor="bg.destructive"
                  onClick={handleConfirmAction}
                  loading={isMutating}
                >
                  Confirm delete
                </Button>
                <Button
                  size="xs"
                  variant="subtle"
                  onClick={handleCancelAction}
                  disabled={isMutating}
                >
                  Cancel
                </Button>
              </HStack>
            ) : (
              <Button
                size="xs"
                variant="ghost"
                aria-label="Delete moderation note"
                title="Delete note"
                onClick={handleConfirmAction}
              >
                <DeleteIcon />
              </Button>
            ))}
        </HStack>
      </Flex>

      <styled.p fontSize="sm" whiteSpace="pre-wrap">
        {note.content}
      </styled.p>

      {deleteNoteError ? (
        <styled.p mt="2" fontSize="sm" color="fg.destructive">
          {deriveError(deleteNoteError)}
        </styled.p>
      ) : null}
    </Box>
  );
}
