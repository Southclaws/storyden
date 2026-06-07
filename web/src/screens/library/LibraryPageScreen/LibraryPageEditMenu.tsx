"use client";

import { MenuSelectionDetails, Portal } from "@ark-ui/react";
import { useRouter } from "next/navigation";
import { parseAsString, useQueryState } from "nuqs";
import { useState } from "react";

import {
  useNodeVersionDraftGet,
  useNodeVersionList,
} from "@/api/openapi-client/nodes";
import {
  NodeVersion,
  NodeVersionStatus,
  NodeWithChildren,
} from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { Timestamp } from "@/components/site/Timestamp";
import { Button, ButtonGroup } from "@/components/ui/button";
import { CancelIcon } from "@/components/ui/icons/Cancel";
import { ChevronDownIcon } from "@/components/ui/icons/Chevron";
import { DraftIcon } from "@/components/ui/icons/Draft";
import { EditIcon } from "@/components/ui/icons/Edit";
import { InfoIcon } from "@/components/ui/icons/Info";
import { VersionsIcon } from "@/components/ui/icons/Versions";
import * as Menu from "@/components/ui/menu";
import { Box, HStack, LStack, WStack, styled } from "@/styled-system/jsx";

import { useLibraryPath } from "../useLibraryPath";

import { PageVersionStatusBadge } from "./PageVersionStatusBadge";
import { useLibraryPagePermissions } from "./permissions";
import { useEditState } from "./useEditState";

const VERSION_MENU_LIMIT = 5;

type Props = {
  node: NodeWithChildren;
};

export function LibraryPageEditMenu({ node }: Props) {
  const router = useRouter();
  const libraryPath = useLibraryPath();
  const [open, setOpen] = useState(false);
  const [, setReviewVersionID] = useQueryState("version", {
    ...parseAsString,
    clearOnDefault: true,
  });
  const { editing, saving, startDirectEdit, startProposalEdit, stopEditing } =
    useEditState();
  const { isAllowedToDirectEdit, isAllowedToEdit, isAllowedToProposeEdit } =
    useLibraryPagePermissions();
  const { data, mutate } = useNodeVersionList(node.id, undefined, {
    swr: {
      enabled: open && !editing,
    },
  });
  const { data: visibleDraftProbe, mutate: mutateVisibleDraftProbe } =
    useNodeVersionDraftGet(node.slug, {
      swr: {
        enabled: isAllowedToEdit && !editing,
        shouldRetryOnError: false,
      },
    });

  const versions = data?.versions ?? [];
  const versionsLoaded = data !== undefined;
  const menuVersions = versions.slice(0, VERSION_MENU_LIMIT);
  const hasMoreVersions = versions.length > VERSION_MENU_LIMIT;
  const versionHistoryHref = `/l/versions/${libraryPath.join("/")}`;
  const visibleDraft =
    versions.find((version) => version.status === NodeVersionStatus.draft) ??
    visibleDraftProbe;
  const hasVisibleDraft = Boolean(visibleDraft);
  const directEditDisabled = !versionsLoaded || Boolean(visibleDraft) || saving;
  const showEditActions =
    isAllowedToEdit && (isAllowedToDirectEdit || isAllowedToProposeEdit);

  function handleDirectEdit() {
    if (visibleDraft) {
      return;
    }

    startDirectEdit();
    setOpen(false);
  }

  async function handleDraftEdit() {
    await startProposalEdit(visibleDraft);
    await mutate();
    await mutateVisibleDraftProbe();
    setOpen(false);
  }

  function handleViewPage() {
    stopEditing();
  }

  async function handleSelect({ value }: MenuSelectionDetails) {
    if (value === "see-all-versions") {
      setOpen(false);
      router.push(versionHistoryHref);
      return;
    }

    if (!value.startsWith("version:")) {
      return;
    }

    setReviewVersionID(value.replace("version:", ""));
  }

  function handleOpenVersion(versionID: string) {
    setReviewVersionID(versionID);
  }

  if (editing) {
    return (
      <Button
        type="button"
        variant="subtle"
        loading={saving}
        disabled={saving}
        onClick={handleViewPage}
      >
        <CancelIcon width="4" height="4" />
        View
      </Button>
    );
  }

  return (
    <Menu.Root
      size="lg"
      open={open}
      onOpenChange={(details) => setOpen(details.open)}
      onSelect={handleSelect}
      positioning={{
        // NOTE: The trigger for this is in the top-right of the Library Page
        // screen. So we want the menu to expand downwards and to the left.
        placement: "bottom-end",
      }}
    >
      <Menu.Trigger asChild>
        <Button type="button" variant="subtle" position="relative">
          <EditIcon width="4" height="4" />
          {isAllowedToEdit ? (
            <>
              Edit
              {hasVisibleDraft && (
                <Box
                  bgColor="visibility.draft.fg"
                  borderRadius="full"
                  w="2"
                  h="2"
                  pointerEvents="none"
                />
              )}
            </>
          ) : (
            "Versions"
          )}

          <ChevronDownIcon width="4" height="4" />
        </Button>
      </Menu.Trigger>

      <Portal>
        <Menu.Positioner>
          <Menu.Content
            w={{ base: "72", sm: "96" }}
            maxW="full"
            overflowY="auto"
            p="1"
            gap="2"
            userSelect="none"
          >
            {showEditActions && (
              <>
                <ButtonGroup
                  role="group"
                  aria-label="Edit actions"
                  px="1"
                  py="1"
                  attached
                >
                  {isAllowedToDirectEdit && (
                    <Button
                      type="button"
                      role="menuitem"
                      variant="subtle"
                      size="sm"
                      flex="1"
                      disabled={directEditDisabled}
                      aria-disabled={directEditDisabled}
                      onClick={handleDirectEdit}
                    >
                      <EditIcon width="4" height="4" />
                      Quick edit
                    </Button>
                  )}

                  {isAllowedToProposeEdit && (
                    <Button
                      type="button"
                      role="menuitem"
                      variant="subtle"
                      size="sm"
                      flex="1"
                      disabled={saving}
                      aria-disabled={saving}
                      onClick={handleDraftEdit}
                    >
                      <DraftIcon width="4" height="4" />
                      {visibleDraft ? "Edit draft" : "Create draft"}
                    </Button>
                  )}
                </ButtonGroup>

                <Menu.Separator />
              </>
            )}

            {versions.length === 0 ? (
              <styled.p px="1" py="1" color="fg.muted" fontSize="xs">
                No versions or drafts yet.
              </styled.p>
            ) : (
              <Menu.ItemGroup id="versions" gap="1">
                {menuVersions.map((version) => (
                  <Menu.Item
                    key={version.id}
                    value={`version:${version.id}`}
                    h="auto"
                    p="0"
                    mx="0"
                    alignItems="stretch"
                  >
                    <VersionMenuItem
                      version={version}
                      onOpenVersion={handleOpenVersion}
                    />
                  </Menu.Item>
                ))}

                {hasMoreVersions && (
                  <Menu.Item value="see-all-versions">
                    <HStack gap="1">
                      <VersionsIcon width="4" height="4" />
                      See all version history
                    </HStack>
                  </Menu.Item>
                )}
              </Menu.ItemGroup>
            )}
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}

function VersionMenuItem({
  version,
  onOpenVersion,
}: {
  version: NodeVersion;
  onOpenVersion: (versionID: string) => void;
}) {
  const isApplied = version.status === NodeVersionStatus.applied;
  const buttonLabel = isApplied ? "View changes" : "Review";

  return (
    <LStack
      w="full"
      gap="2"
      borderRadius="sm"
      px="2"
      py="2"
      _hover={{ bgColor: "bg.muted" }}
    >
      <WStack alignItems="start">
        <PageVersionStatusBadge status={version.status} />
        <Button
          type="button"
          variant="subtle"
          size="xs"
          onClick={() => onOpenVersion(version.id)}
        >
          {buttonLabel}
        </Button>
      </WStack>

      <WStack alignItems="start">
        <MemberBadge
          profile={version.author}
          size="xs"
          name="handle"
          avatar="visible"
        />
        <styled.span color="fg.muted" fontSize="xs">
          <Timestamp created={version.updated_at} /> ago
        </styled.span>
      </WStack>
    </LStack>
  );
}
