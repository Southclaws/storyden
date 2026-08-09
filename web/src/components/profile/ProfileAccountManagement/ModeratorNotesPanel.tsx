import { useEffect, useState } from "react";
import { useSWRConfig } from "swr";

import { handle } from "@/api/client";
import { mutateTransaction } from "@/api/mutate";
import {
  getAccountModerationNoteListKey,
  useAccountModerationNoteCreate,
  useAccountModerationNoteDelete,
  useAccountModerationNoteList,
} from "@/api/openapi-client/accounts";
import {
  AccountModerationNoteListOKResponse,
  ModerationNote,
} from "@/api/openapi-schema";
import { DeletedMemberIdent } from "@/components/member/MemberBadge/DeletedMemberIdent";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { Timestamp } from "@/components/site/Timestamp";
import { useConfirmation } from "@/components/site/useConfirmation";
import { Button } from "@/components/ui/button";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { Text } from "@/components/ui/text";
import { Box, Flex, HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { deriveError } from "@/utils/error";

type ModeratorNotesPanelProps = {
  accountId: string;
  canManageModerationNotes: boolean;
  canViewModerationNotes: boolean;
};

export function ModeratorNotesPanel({
  accountId,
  canManageModerationNotes,
  canViewModerationNotes,
}: ModeratorNotesPanelProps) {
  const [noteDraft, setNoteDraft] = useState("");

  useEffect(() => {
    setNoteDraft("");
  }, [accountId]);

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

  async function submitNote() {
    if (!canManageModerationNotes) {
      return;
    }

    const content = noteDraft.trim();
    if (!content) {
      return;
    }

    await handle(
      async () => {
        await createNote({ content });
        setNoteDraft("");

        if (canViewModerationNotes) {
          await mutateNotes();
        }
      },
      {
        promiseToast: {
          loading: "Saving note...",
          success: "Moderator note added.",
        },
      },
    );
  }

  return (
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
        <Text variant="metadata" color="status.danger.content">
          {deriveError(createNoteError)}
        </Text>
      ) : null}

      {canViewModerationNotes ? (
        <LStack gap="3" w="full" minW="0">
          {notesError ? (
            <Text variant="metadata" color="status.danger.content">
              {deriveError(notesError)}
            </Text>
          ) : notesLoading ? (
            <Text variant="supporting">Loading moderation notes...</Text>
          ) : (notesData?.notes?.length ?? 0) === 0 ? (
            <Text variant="supporting">No moderation notes yet.</Text>
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
        <Text variant="supporting">
          You can add notes, but you do not have permission to view note
          history.
        </Text>
      )}
    </LStack>
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
                notes: current.notes.filter(
                  (moderationNote) => moderationNote.id !== note.id,
                ),
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
      bgColor="background.inset"
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
          <Timestamp
            created={note.created_at}
            color="text.subtle"
            fontSize="xs"
          />

          {canManageModerationNotes &&
            (isConfirming ? (
              <HStack gap="1">
                <Button
                  variant="subtle"
                  bgColor="status.danger.surface"
                  onClick={handleConfirmAction}
                  loading={isMutating}
                >
                  Confirm delete
                </Button>
                <Button
                  variant="subtle"
                  onClick={handleCancelAction}
                  disabled={isMutating}
                >
                  Cancel
                </Button>
              </HStack>
            ) : (
              <Button
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

      <Text variant="supporting" color="text.default" whiteSpace="pre-wrap">
        {note.content}
      </Text>

      {deleteNoteError ? (
        <Text variant="metadata" mt="2" color="status.danger.content">
          {deriveError(deleteNoteError)}
        </Text>
      ) : null}
    </Box>
  );
}
