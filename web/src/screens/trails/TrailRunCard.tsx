import { TrailActionRun, TrailRun } from "@/api/openapi-schema";
import * as Alert from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { CardBox } from "@/components/ui/card-box";
import { RobotIcon } from "@/components/ui/icons/Robot";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { LinkButton } from "@/components/ui/link-button";
import { RelativeTime } from "@/components/ui/relative-time";
import { Text } from "@/components/ui/text";
import { formatOccurrence } from "@/lib/trails";
import { Box, HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { humanise } from "@/utils/text";

import {
  TrailActionRunStatusBadge,
  TrailRunStatusBadge,
} from "./TrailStatusBadge";

type TrailRunCardProps = {
  run: TrailRun;
  robotNames: ReadonlyMap<string, string>;
  onCancelAction: (
    run: TrailRun,
    action: TrailActionRun,
  ) => void | Promise<void>;
};

export function TrailRunCard({
  run,
  robotNames,
  onCancelAction,
}: TrailRunCardProps) {
  const occurrence = run.scheduled_for ?? run.createdAt;
  let runType = "Run";

  switch (run.trigger?.kind) {
    case "manual":
      runType = "Manual run";
      break;
    case "event":
      runType = "Event run";
      break;
    case "scheduled":
      runType = "Scheduled run";
      break;
  }

  return (
    <CardBox
      as="article"
      id={`run-${run.id}`}
      kind="edge"
      gap="0"
      overflow="hidden"
    >
      <WStack alignItems="start" gap="4" padding="3">
        <LStack gap="0.5" alignItems="start" minWidth="0">
          <styled.h3 fontSize="sm" fontWeight="semibold">
            <RelativeTime value={occurrence} precision="strict" />
          </styled.h3>

          <Text as="span" variant="metadata">
            {runType} at{" "}
            <styled.time dateTime={occurrence}>
              {formatOccurrence(occurrence)}
            </styled.time>
          </Text>
        </LStack>

        <TrailRunStatusBadge status={run.status} />
      </WStack>

      {!run.trigger && (
        <Text
          padding="3"
          borderTopWidth="thin"
          borderColor="border.default"
          variant="supporting"
        >
          Trigger details are unavailable for this run.
        </Text>
      )}

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
  onCancelAction: TrailRunCardProps["onCancelAction"];
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
            <Alert.Root tone="warning">
              <Alert.Icon asChild>
                <WarningIcon />
              </Alert.Icon>
              <Alert.Content>
                <Alert.Title>{humanise(output.attention.reason)}</Alert.Title>
                <Alert.Description>
                  {output.attention.message}
                </Alert.Description>
              </Alert.Content>
            </Alert.Root>
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
