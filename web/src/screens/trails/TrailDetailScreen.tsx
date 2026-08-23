"use client";

import { Portal } from "@ark-ui/react";
import Link from "next/link";
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
import { BackAction } from "@/components/site/Action/Back";
import { MoreAction } from "@/components/site/Action/More";
import { EmptyState } from "@/components/site/EmptyState";
import { Unready } from "@/components/site/Unready";
import { Button, ButtonGroup } from "@/components/ui/button";
import { LinkButton } from "@/components/ui/link-button";
import * as Menu from "@/components/ui/menu";
import { PageHeader } from "@/components/ui/page-header";
import { SectionHeading } from "@/components/ui/section-heading";
import { Text } from "@/components/ui/text";
import { HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { pluralise } from "@/utils/text";

import { TrailOverview } from "./TrailOverview";
import { TrailRunCard } from "./TrailRunCard";

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
    <LStack gap="2" alignItems="stretch" minWidth="0">
      <PageHeader
        title={trail.name}
        back={<BackAction fallbackHref="/robots/trails" />}
        actions={
          <TrailActions
            trail={trail}
            editable={editable}
            busy={busy}
            onRunNow={onRunNow}
            onChangeState={onChangeState}
          />
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
          <SectionHeading id="trail-run-history">Run history</SectionHeading>

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

function TrailActions({
  trail,
  editable,
  busy,
  onRunNow,
  onChangeState,
}: Pick<TrailDetailViewProps, "busy" | "onRunNow" | "onChangeState"> & {
  trail: Trail;
  editable: boolean;
}) {
  const canChangeState = trail.status === "active" || trail.status === "paused";
  const canArchive = trail.status !== "archived";
  const hasMenuActions = editable || canChangeState || canArchive;

  return (
    <>
      <ButtonGroup
        attached
        role="group"
        aria-label="Trail actions"
        display={{ base: "flex", md: "none" }}
      >
        {hasMenuActions && (
          <Menu.Root lazyMount positioning={{ placement: "bottom-end" }}>
            <Menu.Trigger asChild>
              <MoreAction
                variant="outline"
                aria-label="More Trail actions"
                disabled={busy}
              />
            </Menu.Trigger>

            <Portal>
              <Menu.Positioner>
                <Menu.Content minW="36">
                  {editable && (
                    <Menu.Item value="edit" asChild>
                      <Link href={`/robots/trails/${trail.id}/edit`}>Edit</Link>
                    </Menu.Item>
                  )}

                  {trail.status === "active" && (
                    <Menu.Item
                      value="pause"
                      onClick={() => onChangeState("paused")}
                    >
                      Pause
                    </Menu.Item>
                  )}

                  {trail.status === "paused" && (
                    <Menu.Item
                      value="resume"
                      onClick={() => onChangeState("active")}
                    >
                      Resume
                    </Menu.Item>
                  )}

                  {canArchive && canChangeState && <Menu.Separator />}

                  {canArchive && (
                    <Menu.Item
                      value="archive"
                      onClick={() => onChangeState("archived")}
                    >
                      Archive
                    </Menu.Item>
                  )}
                </Menu.Content>
              </Menu.Positioner>
            </Portal>
          </Menu.Root>
        )}

        <Button
          variant="solid"
          disabled={busy || trail.status === "archived"}
          onClick={onRunNow}
        >
          Run now
        </Button>
      </ButtonGroup>

      <HStack
        role="group"
        aria-label="Trail actions"
        display={{ base: "none", md: "flex" }}
        flexWrap="wrap"
        justifyContent="end"
      >
        {canArchive && (
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
    </>
  );
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
