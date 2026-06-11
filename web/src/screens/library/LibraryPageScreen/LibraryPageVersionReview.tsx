"use client";

import { parseAsString, useQueryState } from "nuqs";
import { ReactNode, useMemo, useState } from "react";

import { handle } from "@/api/client";
import {
  nodeVersionDelete,
  nodeVersionUpdateStatus,
  useNodeVersionGet,
} from "@/api/openapi-client/nodes";
import {
  NodeVersion,
  NodeVersionStatus,
  NodeWithChildren,
  Permission,
  Property,
  PropertyMutation,
} from "@/api/openapi-schema";
import { useSession } from "@/auth";
import { ContentDiffView } from "@/components/content/DiffViewer/ContentDiffView";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { MemberIdent } from "@/components/member/MemberBadge/MemberIdent";
import { Timestamp } from "@/components/site/Timestamp";
import { Unready } from "@/components/site/Unready";
import * as Alert from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { InfoIcon } from "@/components/ui/icons/Info";
import { SaveIcon } from "@/components/ui/icons/Save";
import * as Table from "@/components/ui/table";
import { css } from "@/styled-system/css";
import { Box, HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { hasPermission } from "@/utils/permissions";
import { useIsTextWrapping } from "@/utils/useIsTextWrapping";

export function LibraryPageVersionReview({
  node,
  versionID,
  onApplied,
}: {
  node: NodeWithChildren;
  versionID: string;
  onApplied?: () => Promise<void> | void;
}) {
  const [, setReviewVersionID] = useQueryState("version", {
    ...parseAsString,
    clearOnDefault: true,
  });
  const { data, error, mutate } = useNodeVersionGet(node.id, versionID, {
    swr: {
      enabled: Boolean(versionID),
    },
  });

  async function handleClose() {
    await setReviewVersionID(null);
  }

  if (error) {
    return (
      <ReviewShell>
        <ReviewHeader
          title="Version unavailable"
          subtitle="This edit may have been removed or you may not have access."
          onClose={handleClose}
        />
      </ReviewShell>
    );
  }

  if (!data) {
    return (
      <ReviewShell>
        <ReviewHeader
          title="Loading version"
          subtitle="Fetching the checkpoint snapshot."
          onClose={handleClose}
        />
      </ReviewShell>
    );
  }

  return (
    <VersionReviewPanel
      node={node}
      version={data}
      onClose={handleClose}
      onMutate={mutate}
      onApplied={onApplied}
    />
  );
}

function VersionReviewPanel({
  node,
  version,
  onClose,
  onMutate,
  onApplied,
}: {
  node: NodeWithChildren;
  version: NodeVersion;
  onClose: () => Promise<void> | void;
  onMutate: () => Promise<NodeVersion | undefined>;
  onApplied?: () => Promise<void> | void;
}) {
  const session = useSession();
  const isLibraryManager = hasPermission(session, Permission.MANAGE_LIBRARY);
  const [applying, setApplying] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const isDraft = version.status === NodeVersionStatus.draft;
  const canApply = isLibraryManager && isDraft;
  const canDelete =
    isDraft && (isLibraryManager || session?.id === version.author.id);
  const previousVersionID = isDraft ? undefined : version.previous?.id;
  const {
    data: previousVersion,
    error: previousVersionError,
    isLoading: previousVersionLoading,
  } = useNodeVersionGet(node.id, previousVersionID ?? "", {
    swr: {
      enabled: Boolean(previousVersionID),
    },
  });
  const comparisonBase = isDraft ? node : previousVersion;

  async function handleApply() {
    setApplying(true);
    try {
      await handle(
        () =>
          nodeVersionUpdateStatus(node.id, version.id, {
            status: NodeVersionStatus.applied,
          }),
        {
          promiseToast: {
            loading: "Applying draft...",
            success: "Draft applied",
          },
        },
      );
      await onMutate();
      await onApplied?.();
      await onClose();
    } finally {
      setApplying(false);
    }
  }

  async function handleDelete() {
    setDeleting(true);
    try {
      await handle(() => nodeVersionDelete(node.id, version.id), {
        promiseToast: {
          loading: "Deleting draft...",
          success: "Draft deleted",
        },
      });
      await onClose();
    } finally {
      setDeleting(false);
    }
  }

  return (
    <ReviewShell>
      <ReviewHeader
        title={
          isDraft ? (
            <styled.span color="fg.subtle">
              Reviewing draft for{" "}
              <styled.span color="fg.default" fontWeight="semibold">
                {node.name}
              </styled.span>
            </styled.span>
          ) : (
            <styled.span color="fg.subtle">
              Viewing checkpoint for{" "}
              <styled.span color="fg.default" fontWeight="semibold">
                {node.name}
              </styled.span>
            </styled.span>
          )
        }
        onClose={onClose}
        controls={
          <HStack gap="1">
            {canDelete && (
              <Button
                type="button"
                size="xs"
                variant="ghost"
                color="fg.destructive"
                loading={deleting}
                onClick={handleDelete}
              >
                <DeleteIcon width="4" height="4" />
                Delete
              </Button>
            )}

            {canApply && (
              <Button
                type="button"
                size="xs"
                variant="subtle"
                loading={applying}
                onClick={handleApply}
              >
                <SaveIcon width="4" height="4" />
                Apply
              </Button>
            )}
          </HStack>
        }
      >
        <WStack>
          <MemberBadge
            profile={version.author}
            size="xs"
            name="handle"
            avatar="visible"
          />

          <Timestamp created={version.updated_at} />
        </WStack>

        {!isDraft && !previousVersionID && <NoPreviousVersionAlert />}

        {!isDraft && previousVersionID && (
          <VersionComparisonContext
            previous={previousVersion}
            error={previousVersionError}
            loading={previousVersionLoading}
          />
        )}

        {comparisonBase && (
          <FieldDiff original={comparisonBase} snapshot={version} />
        )}

        <LStack gap="4">
          {comparisonBase &&
            version.content !== undefined &&
            version.content !== null && (
              <ContentDiffView
                originalHTML={comparisonBase.content ?? ""}
                modifiedHTML={version.content}
              />
            )}
        </LStack>
      </ReviewHeader>
    </ReviewShell>
  );
}

function VersionComparisonContext({
  previous,
  error,
  loading,
}: {
  previous?: NodeVersion;
  error?: unknown;
  loading?: boolean;
}) {
  if (loading || error || !previous) {
    return <Unready error={error} />;
  }

  return (
    <LStack
      borderWidth="thin"
      borderStyle="solid"
      borderColor="border.muted"
      borderRadius="sm"
      bgColor="bg.default"
      color="fg.muted"
      fontSize="sm"
      lineHeight="tight"
      p="2"
    >
      <HStack w="full" flexWrap="wrap" gap="1">
        <InfoIcon width="4" height="4" />
        <span>Comparing with previous version by</span>
        <Box display="inline-block">
          <MemberIdent
            profile={previous.author}
            size="xs"
            name="handle"
            avatar="visible"
          />
        </Box>
        <span>
          from <Timestamp created={previous.updated_at} /> ago
        </span>
      </HStack>
    </LStack>
  );
}

function NoPreviousVersionAlert() {
  return (
    <Alert.Root>
      <Alert.Icon asChild>
        <InfoIcon />
      </Alert.Icon>
      <Alert.Content>
        <Alert.Title>Earlier page state is unavailable</Alert.Title>
        <Alert.Description>
          This is the first tracked version of the page.
        </Alert.Description>
      </Alert.Content>
    </Alert.Root>
  );
}

function ReviewShell({ children }: { children: ReactNode }) {
  return (
    <LStack w="full" gap="4">
      {children}
    </LStack>
  );
}

function ReviewHeader({
  children,
  title,
  subtitle,
  onClose,
  controls,
}: {
  children?: ReactNode;
  title: ReactNode;
  subtitle?: string;
  onClose: () => Promise<void> | void;
  controls?: ReactNode;
}) {
  return (
    <LStack
      borderWidth="thin"
      borderStyle="dashed"
      borderColor="visibility.draft.border"
      borderRadius="sm"
      bgColor="bg.subtle"
      p="2"
      gap="2"
    >
      <WStack alignItems="start">
        <LStack gap="0" minW="0">
          <Heading size="sm">{title}</Heading>
          {subtitle && (
            <styled.span color="fg.muted" fontSize="sm">
              {subtitle}
            </styled.span>
          )}
        </LStack>

        <HStack gap="1">
          {controls}
          <Button type="button" size="xs" variant="ghost" onClick={onClose}>
            Close
          </Button>
        </HStack>
      </WStack>

      {children}
    </LStack>
  );
}

function FieldDiff({
  original,
  snapshot,
}: {
  original: FieldDiffSource;
  snapshot: NodeVersion;
}) {
  const rows = useMemo(() => {
    const titleDiff = computeTitleDiff(original.name, snapshot.name);

    const slugDiff = computeSlugDiff(original.slug, snapshot.slug);

    const propDiff = computePropertyDiff(
      original.properties,
      snapshot.properties,
    );

    return [...titleDiff, ...slugDiff, ...propDiff];
  }, [original, snapshot]);

  if (rows.length === 0) {
    return null;
  }

  return (
    <LStack gap="2" w="full" minW="0">
      <Table.Root
        size="sm"
        tableLayout={{
          base: "auto",
          md: "fixed",
        }}
        w="full"
        overflow="hidden"
      >
        <Table.Body>
          {rows.map((row) => (
            <FieldDiffRow key={row.key} row={row} />
          ))}
        </Table.Body>
      </Table.Root>
    </LStack>
  );
}

type VersionProperty = Property | PropertyMutation;

type FieldDiffSource = {
  name: string;
  slug: string;
  content?: string | null;
  properties: VersionProperty[];
};

const diffBeforeLabelStyles = css({
  backgroundColor: "red.4",
  color: "red.12/70",
});

const diffAfterLabelStyles = css({
  backgroundColor: "green.3",
  color: "green.12/70",
});

const diffBeforeStyles = css({
  textDecoration: "line-through",
  backgroundColor: "red.4",
  color: "red.12",
});

const diffAfterStyles = css({
  backgroundColor: "green.3",
  color: "green.12",
});

function FieldDiffRow({ row }: { row: FieldDiffRow }) {
  switch (row.kind) {
    case "added":
      return (
        <Table.Row key={row.key}>
          <Table.Cell className={diffAfterLabelStyles}>
            <styled.span fontWeight="medium">{row.kind}</styled.span>
            <br />
            <span className={diffAfterStyles}>{row.after.name}</span>
          </Table.Cell>
          <Table.Cell className={diffAfterStyles}>{row.after.value}</Table.Cell>
        </Table.Row>
      );

    case "removed":
      return (
        <Table.Row key={row.key}>
          <Table.Cell className={diffBeforeLabelStyles}>
            <styled.span fontWeight="medium">{row.kind}</styled.span>
            <br />
            <span className={diffBeforeStyles}>{row.before.name}</span>
          </Table.Cell>
          <Table.Cell className={diffBeforeStyles}>
            {row.before.value}
          </Table.Cell>
        </Table.Row>
      );

    case "changed":
      const nameChanged =
        isFieldProperty(row) && row.before.name !== row.after.name;
      const valueChanged = isFieldProperty(row)
        ? row.before.value !== row.after.value
        : row.before !== row.after;
      // nOTE: we don;t do types right now, but in future we will...
      // const afterType = row.after.type ?? row.before.type;
      // const typeChanged = row.before.type !== afterType;

      const beforeName = isFieldProperty(row) ? row.before.name : row.key;
      const afterName = isFieldProperty(row) ? row.after.name : row.key;
      const beforeValue = isFieldProperty(row) ? row.before.value : row.before;
      const afterValue = isFieldProperty(row) ? row.after.value : row.after;

      return (
        <Table.Row key={row.key}>
          {/* <Table.Cell fontWeight="medium" color="fg.muted">
            {row.kind}
          </Table.Cell> */}
          <Table.Cell>
            <styled.span fontWeight="medium" color="fg.muted">
              {row.kind}
            </styled.span>
            <br />
            {nameChanged ? (
              <WrappedDiff t1={beforeName} t2={afterName} />
            ) : (
              beforeName || "(unnamed)"
            )}
          </Table.Cell>
          <Table.Cell>
            {valueChanged ? (
              <WrappedDiff t1={beforeValue} t2={afterValue} />
            ) : (
              beforeValue || "(empty)"
            )}
          </Table.Cell>
        </Table.Row>
      );
  }
}

const diffArrowVerticalStyles = css({
  display: "inline-block",
  fontWeight: "medium",
  color: "fg.muted",
  textAlign: "center",
  w: "full",
  transform: "rotate(90deg)",
});

const diffArrowHorizontalStyles = css({
  display: "inline-block",
  fontWeight: "medium",
  color: "fg.muted",
  transform: "rotate(0deg)",
  paddingX: "1",
});

function WrappedDiff({ t1, t2 }: { t1: string; t2: string }) {
  const [ref, wrapped] = useIsTextWrapping<HTMLSpanElement>();

  return (
    <styled.div position="relative" lineHeight="tight">
      {/* NOTE: creates a hidden reference instance of the full text that's
          immutable, we use this to measure the true length and wrapping state
          of the content then use that to mutate the actual rendered content */}

      <styled.span
        ref={ref}
        position="absolute"
        top="0"
        left="0"
        width="full"
        visibility="hidden"
        bgColor="visibility.draft.bg"
        pointerEvents="none"
        aria-hidden
      >
        {`${t1} → ${t2}`}
      </styled.span>

      <span className={diffBeforeStyles}>{t1}</span>
      {wrapped && <br />}
      <styled.span
        className={
          wrapped ? diffArrowVerticalStyles : diffArrowHorizontalStyles
        }
        style={{
          // silly lil anim nobody will notice but it's cool when you resize!
          transition: "350ms cubic-bezier(0.34, 1.56, 0.64, 1)",
        }}
      >
        →
      </styled.span>
      {wrapped && <br />}
      <span className={diffAfterStyles}>{t2}</span>
    </styled.div>
  );
}

type FieldDiffRow = TitleDiffRow | SlugDiffRow | PropertyDiffRow;

function isFieldProperty(row: FieldDiffRow): row is PropertyDiffRow {
  return row.key !== "Title" && row.key !== "Slug";
}

type TitleDiffRow = {
  key: "Title";
  kind: "changed";
  before: string;
  after: string;
};

type SlugDiffRow = {
  key: "Slug";
  kind: "changed";
  before: string;
  after: string;
};

type PropertyDiffRow =
  | {
      key: string;
      kind: "added";
      after: PropertyMutation;
    }
  | {
      key: string;
      kind: "removed";
      before: VersionProperty;
    }
  | {
      key: string;
      kind: "changed";
      before: VersionProperty;
      after: PropertyMutation;
    };

function computeTitleDiff(before: string, after: string): FieldDiffRow[] {
  if (before === after) {
    return [];
  }

  return [
    {
      key: "Title",
      kind: "changed",
      before,
      after,
    },
  ];
}

function computeSlugDiff(before: string, after: string): FieldDiffRow[] {
  if (before === after) {
    return [];
  }

  return [
    {
      key: "Slug",
      kind: "changed",
      before,
      after,
    },
  ];
}

function propertyEqual(before: VersionProperty, after: PropertyMutation) {
  return (
    before.name === after.name &&
    before.type === (after.type ?? before.type) &&
    before.value === after.value
  );
}

function propertyKey(property: VersionProperty) {
  return property.fid ?? property.name;
}

function computePropertyDiff(
  original: VersionProperty[],
  snapshot: PropertyMutation[],
): FieldDiffRow[] {
  const originalByFID = new Map(
    original.map((property) => [property.fid, property]),
  );
  const originalByName = new Map(
    original.map((property) => [property.name, property]),
  );
  const consumedOriginals = new Set<VersionProperty>();
  const changes: PropertyDiffRow[] = [];

  for (const after of snapshot) {
    const before =
      (after.fid ? originalByFID.get(after.fid) : undefined) ??
      originalByName.get(after.name);

    if (before) {
      consumedOriginals.add(before);
    }

    if (!before) {
      changes.push({
        key: propertyKey(after),
        kind: "added",
        after,
      });
      continue;
    }

    if (propertyEqual(before, after)) {
      continue;
    }

    changes.push({
      key: propertyKey(after),
      kind: "changed",
      before,
      after,
    });
  }

  for (const before of original) {
    if (consumedOriginals.has(before)) {
      continue;
    }

    changes.push({
      key: propertyKey(before),
      kind: "removed",
      before,
    });
  }

  return changes;
}
