"use client";

import { formatDistanceToNowStrict } from "date-fns";
import { useState } from "react";
import { useSWRConfig } from "swr";

import { useRobotsList } from "@/api/openapi-client/robots";
import {
  getTrailListKey,
  trailActionRunCancel,
  trailRunNow,
  trailUpdate,
  useTrailGet,
  useTrailRunList,
} from "@/api/openapi-client/trails";
import {
  Trail,
  TrailActionRun,
  TrailMutableProps,
  TrailRun,
} from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { BackAction } from "@/components/site/Action/Back";
import { EmptyState } from "@/components/site/EmptyState";
import { Unready } from "@/components/site/Unready";
import { Button, ButtonGroup } from "@/components/ui/button";
import { CardBox } from "@/components/ui/card-box";
import { CalendarIcon } from "@/components/ui/icons/Calendar";
import { RobotIcon } from "@/components/ui/icons/Robot";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { LinkButton } from "@/components/ui/link-button";
import { PageHeader } from "@/components/ui/page-header";
import { Text } from "@/components/ui/text";
import { describeTrailSchedule, formatOccurrence } from "@/lib/trails";
import { Box, HStack, LStack, WStack, styled } from "@/styled-system/jsx";

import {
  TrailActionRunStatusBadge,
  TrailRunStatusBadge,
  TrailStatusBadge,
} from "./TrailStatusBadge";

type RobotReference = {
  id: string;
  name: string;
};

type TrailDetailViewProps = {
  trail: Trail;
  runs?: TrailRun[];
  runsError?: unknown;
  robots?: RobotReference[];
  busy?: boolean;
  error?: string;
  onRunNow: () => void | Promise<void>;
  onChangeState: (status: TrailMutableProps["status"]) => void | Promise<void>;
  onCancelAction: (
    run: TrailRun,
    action: TrailActionRun,
  ) => void | Promise<void>;
};

export function TrailDetailScreen({ trailId }: { trailId: string }) {
  const { mutate } = useSWRConfig();
  const trailQuery = useTrailGet(trailId);
  const runsQuery = useTrailRunList(trailId, {
    swr: { refreshInterval: 3000 },
  });
  const robotsQuery = useRobotsList();
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string>();

  if (!trailQuery.data) return <Unready error={trailQuery.error} />;

  const trail = trailQuery.data;

  async function changeState(status: TrailMutableProps["status"]) {
    setBusy(true);
    setError(undefined);

    try {
      await trailUpdate(trail.id, writePayload(trail, status));
      await Promise.all([trailQuery.mutate(), mutate(getTrailListKey())]);
    } catch (cause) {
      setError(
        cause instanceof Error
          ? cause.message
          : "Could not change this Trail's state.",
      );
    }

    setBusy(false);
  }

  async function runNow() {
    setBusy(true);
    setError(undefined);

    try {
      await trailRunNow(trail.id);
      await runsQuery.mutate();
    } catch (cause) {
      setError(
        cause instanceof Error ? cause.message : "Could not start this Trail.",
      );
    }

    setBusy(false);
  }

  async function cancelAction(run: TrailRun, action: TrailActionRun) {
    setError(undefined);

    try {
      await trailActionRunCancel(trail.id, run.id, action.id);
      await runsQuery.mutate();
    } catch (cause) {
      setError(
        cause instanceof Error
          ? cause.message
          : "Could not cancel this Trail action.",
      );
    }
  }

  return (
    <TrailDetailView
      trail={trail}
      runs={runsQuery.data?.runs}
      runsError={runsQuery.error}
      robots={robotsQuery.data?.robots}
      busy={busy}
      error={error}
      onRunNow={runNow}
      onChangeState={changeState}
      onCancelAction={cancelAction}
    />
  );
}

export function TrailDetailView({
  trail,
  runs,
  runsError,
  robots = [],
  busy = false,
  error,
  onRunNow,
  onChangeState,
  onCancelAction,
}: TrailDetailViewProps) {
  const robotNames = new Map(robots.map((robot) => [robot.id, robot.name]));
  const editable = trail.status !== "archived" && trail.status !== "finished";

  return (
    <LStack gap="6" alignItems="stretch" minWidth="0">
      <PageHeader
        title={trail.name}
        description={trail.description || undefined}
        badge={<TrailStatusBadge status={trail.status} />}
        back={<BackAction fallbackHref="/robots/trails" />}
        actions={
          <HStack
            role="group"
            aria-label="Trail actions"
            flexWrap="wrap"
            justifyContent="end"
          >
            {trail.status !== "archived" && (
              <Button
                variant="ghost"
                intent="destructive"
                disabled={busy}
                onClick={() => onChangeState("archived")}
              >
                Archive
              </Button>
            )}

            <ButtonGroup attached>
              {trail.status === "active" ? (
                <Button
                  variant="outline"
                  disabled={busy}
                  onClick={() => onChangeState("paused")}
                >
                  Pause
                </Button>
              ) : trail.status === "paused" ? (
                <Button
                  variant="outline"
                  disabled={busy}
                  onClick={() => onChangeState("active")}
                >
                  Resume
                </Button>
              ) : null}

              {editable ? (
                <LinkButton
                  variant="outline"
                  href={`/robots/trails/${trail.id}/edit`}
                >
                  Edit
                </LinkButton>
              ) : (
                <Button variant="outline" disabled>
                  Edit
                </Button>
              )}

              <Button
                variant="solid"
                disabled={busy || trail.status === "archived"}
                onClick={onRunNow}
              >
                Run now
              </Button>
            </ButtonGroup>
          </HStack>
        }
      />

      {error && (
        <styled.p role="alert" color="status.danger.content">
          {error}
        </styled.p>
      )}

      <TrailOverview trail={trail} />

      <styled.section aria-labelledby="trail-run-history">
        <WStack alignItems="end" gap="3" mb="3">
          <styled.h2 id="trail-run-history" fontSize="lg" fontWeight="semibold">
            Run history
          </styled.h2>

          {runs && runs.length > 0 && (
            <Text as="span" variant="metadata">
              {runs.length} {pluralise(runs.length, "run")}
            </Text>
          )}
        </WStack>

        {!runs ? (
          <Unready error={runsError} />
        ) : runs.length === 0 ? (
          <EmptyState hideContributionLabel>
            This Trail has not run yet.
          </EmptyState>
        ) : (
          <LStack gap="3" alignItems="stretch">
            {runs.map((run) => (
              <TrailRunCard
                key={run.id}
                run={run}
                robotNames={robotNames}
                onCancelAction={onCancelAction}
              />
            ))}
          </LStack>
        )}
      </styled.section>
    </LStack>
  );
}

function TrailOverview({ trail }: { trail: Trail }) {
  return (
    <styled.section
      aria-label="Trail details"
      borderBottomWidth="thin"
      borderColor="border.default"
      paddingY="4"
      paddingX={{ base: "0", sm: "1" }}
    >
      <styled.div
        display="grid"
        gridTemplateColumns={{
          base: "minmax(0, 1fr) minmax(0, 1fr)",
          md: "minmax(0, 1fr) auto auto",
        }}
        columnGap={{ base: "4", md: "8" }}
        rowGap="3"
        alignItems="start"
      >
        <HStack
          alignItems="start"
          gap="3"
          minWidth="0"
          gridColumn={{ base: "1 / -1", md: "auto" }}
        >
          <Box
            display="flex"
            alignItems="center"
            justifyContent="center"
            width="9"
            height="9"
            flexShrink="0"
            borderRadius="md"
            background="background.inset"
            color="text.muted"
          >
            <CalendarIcon width="4" height="4" />
          </Box>

          <Box minWidth="0">
            <Text
              as="span"
              variant="metadata"
              fontWeight="bold"
              textTransform="uppercase"
              display="block"
            >
              Schedule
            </Text>
            <Text
              as="span"
              fontWeight="semibold"
              overflowWrap="anywhere"
              display="block"
            >
              {describeTrailSchedule(trail.trigger.schedule)}
            </Text>
            <Text as="span" variant="metadata" display="block">
              {trail.trigger.schedule.timezone}
            </Text>
          </Box>
        </HStack>

        <Occurrence
          label="Last run"
          value={trail.last_occurrence_at}
          empty="Not run yet"
        />

        <Occurrence
          label="Next run"
          value={trail.next_occurrence_at}
          empty="No future run"
          emphasis
        />
      </styled.div>

      <WStack alignItems="center" gap="4" flexWrap="wrap" marginTop="4">
        <HStack gap="2" minWidth="0">
          <Text as="span" variant="metadata" flexShrink="0">
            Created by
          </Text>
          <MemberBadge
            profile={trail.created_by}
            size="xs"
            name="handle"
            as="link"
          />
        </HStack>

        <Text as="span" variant="metadata" flexShrink="0">
          {trail.actions.length} {pluralise(trail.actions.length, "action")}
        </Text>
      </WStack>
    </styled.section>
  );
}

function Occurrence({
  label,
  value,
  empty,
  emphasis = false,
}: {
  label: string;
  value?: string;
  empty: string;
  emphasis?: boolean;
}) {
  return (
    <Box minWidth="0">
      <Text
        as="span"
        variant="metadata"
        fontWeight="bold"
        textTransform="uppercase"
        display="block"
      >
        {label}
      </Text>
      {value ? (
        <styled.time
          dateTime={value}
          fontSize="sm"
          fontWeight={emphasis ? "semibold" : "normal"}
          display="block"
        >
          {formatOccurrence(value)}
        </styled.time>
      ) : (
        <Text as="span" variant="supporting" display="block">
          {empty}
        </Text>
      )}
    </Box>
  );
}

function TrailRunCard({
  run,
  robotNames,
  onCancelAction,
}: {
  run: TrailRun;
  robotNames: ReadonlyMap<string, string>;
  onCancelAction: TrailDetailViewProps["onCancelAction"];
}) {
  const occurrence = run.scheduled_for ?? run.createdAt;
  const relativeOccurrence = formatDistanceToNowStrict(new Date(occurrence), {
    addSuffix: true,
  });
  const runType = run.trigger.kind === "manual" ? "Manual" : "Scheduled";

  return (
    <CardBox
      as="article"
      id={`run-${run.id}`}
      padding="0"
      gap="0"
      overflow="hidden"
    >
      <WStack alignItems="start" gap="4" padding="3">
        <LStack gap="0.5" alignItems="start" minWidth="0">
          <styled.h3 fontSize="sm" fontWeight="semibold">
            <styled.time dateTime={occurrence}>
              {relativeOccurrence}
            </styled.time>
          </styled.h3>

          <Text as="span" variant="metadata">
            {runType} run at{" "}
            <styled.time dateTime={occurrence}>
              {formatOccurrence(occurrence)}
            </styled.time>
          </Text>
        </LStack>

        <TrailRunStatusBadge status={run.status} />
      </WStack>

      {run.status === "skipped" ? (
        <Text
          padding="3"
          borderTopWidth="thin"
          borderColor="border.default"
          variant="supporting"
        >
          Skipped while the scheduler was offline.
        </Text>
      ) : (
        <styled.ul listStyle="none">
          {run.actions.map((action) => (
            <TrailActionRunRow
              key={action.id}
              run={run}
              action={action}
              showStatus={run.actions.length > 1}
              robotNames={robotNames}
              onCancelAction={onCancelAction}
            />
          ))}
        </styled.ul>
      )}
    </CardBox>
  );
}

function TrailActionRunRow({
  run,
  action,
  showStatus,
  robotNames,
  onCancelAction,
}: {
  run: TrailRun;
  action: TrailActionRun;
  showStatus: boolean;
  robotNames: ReadonlyMap<string, string>;
  onCancelAction: TrailDetailViewProps["onCancelAction"];
}) {
  const robotAction =
    action.action.action.type === "robot_run"
      ? action.action.action
      : undefined;
  const invocation =
    action.target?.type === "robot_run" ? action.target : undefined;
  const output = invocation?.output;
  const robotName = robotAction
    ? (robotNames.get(robotAction.robot_ref) ?? "Robot invocation")
    : "Trail action";

  return (
    <styled.li
      borderTopWidth="thin"
      borderColor="border.default"
      padding="3"
      minWidth="0"
    >
      <WStack alignItems="start" gap="3" flexWrap="wrap">
        <HStack alignItems="start" gap="2" minWidth="0" flex="1">
          <Box
            display="flex"
            alignItems="center"
            justifyContent="center"
            width="8"
            height="8"
            flexShrink="0"
            borderRadius="md"
            background="background.inset"
            color="text.muted"
          >
            <RobotIcon width="4" height="4" />
          </Box>

          <LStack gap="0.5" alignItems="start" minWidth="0">
            <Text as="span" fontWeight="semibold" overflowWrap="anywhere">
              {robotName}
            </Text>

            {robotAction && (
              <Text
                as="span"
                variant="metadata"
                overflowWrap="anywhere"
                lineClamp="2"
              >
                {robotAction.instruction}
              </Text>
            )}
          </LStack>
        </HStack>

        {(showStatus || action.status === "running") && (
          <HStack gap="2" flexShrink="0">
            {showStatus && <TrailActionRunStatusBadge status={action.status} />}

            {action.status === "running" && (
              <Button
                size="sm"
                variant="ghost"
                intent="destructive"
                onClick={() => onCancelAction(run, action)}
              >
                Cancel
              </Button>
            )}
          </HStack>
        )}
      </WStack>

      {(output?.summary || output?.attention || action.error || invocation) && (
        <LStack
          gap="2"
          alignItems="stretch"
          marginTop="3"
          paddingLeft={{ base: "0", sm: "10" }}
          minWidth="0"
        >
          {output?.summary && (
            <Text overflowWrap="anywhere">{output.summary}</Text>
          )}

          {output?.attention && (
            <HStack
              alignItems="start"
              gap="2"
              borderWidth="thin"
              borderColor="status.warning.border"
              borderRadius="md"
              background="status.warning.surface"
              color="status.warning.content"
              padding="2"
            >
              <WarningIcon width="4" height="4" flexShrink="0" />
              <LStack gap="0" alignItems="start" minWidth="0">
                <Text as="span" fontWeight="semibold" color="current">
                  {humanise(output.attention.reason)}
                </Text>
                <Text as="span" variant="supporting" color="current">
                  {output.attention.message}
                </Text>
              </LStack>
            </HStack>
          )}

          {action.error && (
            <Text color="status.danger.content" overflowWrap="anywhere">
              {action.error}
            </Text>
          )}

          {invocation && (
            <LinkButton
              href={`/robots/chats/${invocation.robot_session_id}`}
              variant="subtle"
              width="fit"
            >
              <RobotIcon />
              View Robot session
            </LinkButton>
          )}
        </LStack>
      )}
    </styled.li>
  );
}

function pluralise(count: number, singular: string): string {
  return count === 1 ? singular : `${singular}s`;
}

function humanise(value: string): string {
  const label = value.replaceAll("_", " ");
  return label.charAt(0).toUpperCase() + label.slice(1);
}

function writePayload(
  trail: Trail,
  status: TrailMutableProps["status"],
): TrailMutableProps {
  return {
    name: trail.name,
    description: trail.description,
    status,
    trigger: trail.trigger,
    actions: trail.actions.map(({ action }) => action),
  };
}
