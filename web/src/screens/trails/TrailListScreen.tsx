"use client";

import { formatDistanceToNow } from "date-fns";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect } from "react";

import { useTrailList, useTrailRunList } from "@/api/openapi-client/trails";
import { Trail, TrailRun } from "@/api/openapi-schema";
import { EmptyState } from "@/components/site/EmptyState";
import { Unready } from "@/components/site/Unready";
import { CardBox } from "@/components/ui/card-box";
import { RunningAnimatedIcon } from "@/components/ui/icons/RunningAnimatedIcon";
import { TrailIcon } from "@/components/ui/icons/Trail";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { linkOverlay } from "@/styled-system/patterns";

export function TrailListScreen() {
  const { data, error } = useTrailList();
  const targetRunID = useSearchParams().get("run");
  if (!data) return <Unready error={error} />;
  if (data.trails.length === 0) {
    return (
      <EmptyState icon={<TrailIcon />} hideContributionLabel>
        No Trails yet. Create one to run Robots automatically on a schedule.
      </EmptyState>
    );
  }

  return (
    <LStack gap="3" alignItems="stretch">
      {data.trails.map((trail) => (
        <TrailRow key={trail.id} trail={trail} targetRunID={targetRunID} />
      ))}
    </LStack>
  );
}

function TrailRow({
  trail,
  targetRunID,
}: {
  trail: Trail;
  targetRunID: string | null;
}) {
  const router = useRouter();
  const runs = useTrailRunList(trail.id, {
    swr: { refreshInterval: trail.status === "active" ? 5000 : 0 },
  });
  useEffect(() => {
    if (targetRunID && runs.data?.runs.some((run) => run.id === targetRunID)) {
      router.replace(`/robots/trails/${trail.id}#run-${targetRunID}`);
    }
  }, [router, runs.data, targetRunID, trail.id]);

  const latest = runs.data?.runs[0];

  return (
    <TrailCard
      trail={trail}
      latest={latest}
      activityState={runs.data ? "ready" : runs.error ? "failed" : "loading"}
    />
  );
}

export function TrailCard({
  trail,
  latest,
  activityState = "ready",
}: {
  trail: Trail;
  latest?: TrailRun;
  activityState?: "ready" | "loading" | "failed";
}) {
  const loading = activityState === "loading";
  const failed = activityState === "failed";

  return (
    <CardBox as="article" id={trail.id} position="relative" overflow="hidden">
      <WStack alignItems="center" gap="4" minWidth="0">
        <HStack minWidth="0" flex="1" overflow="hidden">
          <styled.h2
            minWidth="0"
            overflow="hidden"
            overflowWrap="anywhere"
            lineClamp="1"
            fontSize="sm"
            fontWeight="semibold"
            lineHeight="normal"
          >
            <Link className={linkOverlay()} href={`/robots/trails/${trail.id}`}>
              {trail.name}
            </Link>
          </styled.h2>
        </HStack>

        <TrailOperationalStatus
          trail={trail}
          latest={latest}
          loading={loading}
          failed={failed}
        />
      </WStack>

      {trail.description && (
        <styled.p
          minWidth="0"
          overflow="hidden"
          overflowWrap="anywhere"
          lineClamp="1"
        >
          {trail.description}
        </styled.p>
      )}

      <WStack alignItems="center" gap="4" minWidth="0">
        <ActionCount count={trail.actions.length} />
        <LatestActivity
          next={trail.next_occurrence_at}
          latest={latest}
          loading={loading}
          failed={failed}
        />
      </WStack>
    </CardBox>
  );
}

function TrailOperationalStatus({
  trail,
  latest,
  loading,
  failed,
}: {
  trail: Trail;
  latest?: TrailRun;
  loading: boolean;
  failed: boolean;
}) {
  if (failed) {
    return (
      <StatusLabel color="status.danger.content" icon={<WarningIcon />}>
        Status unavailable
      </StatusLabel>
    );
  }

  if (loading) {
    return <StatusLabel>Checking status</StatusLabel>;
  }

  if (latest?.status === "attention_required") {
    return (
      <StatusLabel color="status.danger.content" icon={<WarningIcon />} strong>
        Needs attention
      </StatusLabel>
    );
  }

  if (latest?.status === "running") {
    return (
      <StatusLabel color="accent.text" icon={<RunningAnimatedIcon />} strong>
        Running now
      </StatusLabel>
    );
  }

  if (latest?.status === "queued") {
    return <StatusLabel color="accent.text">Queued</StatusLabel>;
  }

  return <StatusLabel>{capitalise(trail.status)}</StatusLabel>;
}

function StatusLabel({
  children,
  color = "text.muted",
  icon,
  strong = false,
}: {
  children: string;
  color?: "accent.text" | "status.danger.content" | "text.muted";
  icon?: React.ReactNode;
  strong?: boolean;
}) {
  return (
    <HStack
      gap="1.5"
      color={color}
      flexShrink="0"
      fontSize="xs"
      fontWeight={strong ? "semibold" : "medium"}
      textWrap="nowrap"
    >
      {children}
      {icon && (
        <styled.span
          display="inline-flex"
          alignItems="center"
          justifyContent="center"
          width="4"
          height="4"
          flexShrink="0"
          css={{
            "& svg": {
              width: "full",
              height: "full",
            },
          }}
          aria-hidden
        >
          {icon}
        </styled.span>
      )}
    </HStack>
  );
}

function ActionCount({ count }: { count: number }) {
  return (
    <styled.span
      color="text.muted"
      flexShrink="0"
      fontSize="xs"
      textWrap="nowrap"
    >
      {count} {count === 1 ? "action" : "actions"}
    </styled.span>
  );
}

function LatestActivity({
  latest,
  loading,
  failed,
  next,
}: {
  latest?: TrailRun;
  loading: boolean;
  failed: boolean;
  next?: string;
}) {
  if (failed) {
    return <ActivityLabel>Activity unavailable</ActivityLabel>;
  }

  if (loading) {
    return <ActivityLabel>Checking activity</ActivityLabel>;
  }

  if (!latest) {
    return <ActivityLabel>Not run yet</ActivityLabel>;
  }

  const occurredAt =
    latest.finished_at ?? latest.scheduled_for ?? latest.createdAt;
  const lastRanLabel = formatDistanceToNow(new Date(occurredAt), {
    addSuffix: true,
  });
  const nextRunLabel = next
    ? formatDistanceToNow(new Date(next), {
        addSuffix: true,
      })
    : undefined;
  const nextRunText = nextRunLabel ? `, next ${nextRunLabel}` : "";

  const label = (() => {
    switch (latest.status) {
      case "running":
        return `Started ${lastRanLabel}`;
      case "queued":
        return `Queued ${lastRanLabel}`;
      case "cancelled":
        return `Last run cancelled ${lastRanLabel}`;
      case "skipped":
        return `Last run skipped ${lastRanLabel}`;
      case "completed":
      case "attention_required":
        return `Last ran ${lastRanLabel}${nextRunText}`;
    }
  })();

  return <ActivityLabel dateTime={occurredAt}>{label}</ActivityLabel>;
}

function ActivityLabel({
  dateTime,
  children,
}: {
  dateTime?: string;
  children: React.ReactNode;
}) {
  if (dateTime) {
    return (
      <styled.time
        dateTime={dateTime}
        color="text.muted"
        fontSize="xs"
        minWidth="0"
        overflow="hidden"
        textAlign="right"
        textOverflow="ellipsis"
        textWrap="nowrap"
      >
        {children}
      </styled.time>
    );
  }

  return (
    <styled.span
      color="text.muted"
      fontSize="xs"
      minWidth="0"
      overflow="hidden"
      textAlign="right"
      textOverflow="ellipsis"
      textWrap="nowrap"
    >
      {children}
    </styled.span>
  );
}

function capitalise(value: string): string {
  return value.charAt(0).toUpperCase() + value.slice(1);
}
