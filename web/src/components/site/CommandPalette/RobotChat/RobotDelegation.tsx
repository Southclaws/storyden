import type { ReactNode } from "react";

import { RobotDelegationData, StorydenUIMessage } from "@/api/robots-types";
import { ChevronDownIcon } from "@/components/ui/icons/Chevron";
import { RobotIcon } from "@/components/ui/icons/Robot";
import { Text } from "@/components/ui/text";
import { css } from "@/styled-system/css";
import { Box, LStack, styled } from "@/styled-system/jsx";

import { useRobotChat } from "./RobotChatContext";
import { projectToolOutputs } from "./RobotChatMessageProjection";

type Props = {
  data: RobotDelegationData;
  renderParts: (
    parts: readonly StorydenUIMessage["parts"][number][],
    isUser: boolean,
  ) => ReactNode;
};

const statusLabels: Record<RobotDelegationData["status"], string> = {
  running: "Working…",
  completed: "Completed",
  failed: "Failed",
};

type PresentedStatus = RobotDelegationData["status"] | "stopped";

const presentedStatusLabels: Record<PresentedStatus, string> = {
  ...statusLabels,
  stopped: "Stopped",
};

export function RobotDelegation({ data, renderParts }: Props) {
  const { status: chatStatus } = useRobotChat();
  const messages = projectToolOutputs(data.messages).filter(
    (message) => (message.parts ?? []).length > 0,
  );
  const chatIsActive = chatStatus === "submitted" || chatStatus === "streaming";
  const presentedStatus: PresentedStatus =
    data.status === "running" && !chatIsActive ? "stopped" : data.status;

  return (
    <details
      role="group"
      aria-label={`${data.robot.name} delegation`}
      open={data.status === "failed" || undefined}
      className={css({
        w: "full",
        overflow: "hidden",
        bg: "background.inset",
        borderWidth: "thin",
        borderColor:
          data.status === "failed" ? "status.danger.border" : "border.default",
        borderRadius: "lg",
      })}
    >
      <styled.summary
        display="grid"
        gridTemplateColumns="auto minmax(0, 1fr) auto auto"
        columnGap="2"
        rowGap="0.5"
        alignItems="center"
        px="3"
        py="2.5"
        cursor="pointer"
        listStyle="none"
        _hover={{ bg: "background.surface" }}
      >
        <RobotIcon aria-hidden="true" w="4" h="4" color="text.muted" />
        <styled.span
          minW="0"
          overflow="hidden"
          textOverflow="ellipsis"
          whiteSpace="nowrap"
          color="text.default"
          fontSize="sm"
          fontWeight="medium"
        >
          {data.robot.name}
        </styled.span>
        <styled.span
          aria-live="polite"
          color={
            data.status === "failed" ? "status.danger.content" : "text.muted"
          }
          fontSize="xs"
        >
          {presentedStatusLabels[presentedStatus]}
        </styled.span>
        <ChevronDownIcon
          aria-hidden="true"
          w="3"
          h="3"
          color="text.muted"
          _detailsOpen={{ transform: "rotate(180deg)" }}
        />
        {data.request ? (
          <styled.span
            gridColumn="2 / -1"
            minW="0"
            overflow="hidden"
            textOverflow="ellipsis"
            whiteSpace="nowrap"
            color="text.muted"
            fontSize="xs"
          >
            {data.request}
          </styled.span>
        ) : null}
      </styled.summary>

      <LStack
        gap="3"
        px="3"
        py="3"
        borderTopWidth="thin"
        borderTopColor="border.muted"
      >
        {data.request ? (
          <Box>
            <Text as="span" variant="metadata">
              Delegated task
            </Text>
            <Text mt="0.5" fontSize="sm">
              {data.request}
            </Text>
          </Box>
        ) : null}

        {messages.length > 0 ? (
          <LStack gap="3" w="full" minW="0">
            {messages.map((message) => (
              <LStack key={message.id} gap="1.5" w="full" minW="0">
                {renderParts(message.parts ?? [], message.role === "user")}
              </LStack>
            ))}
          </LStack>
        ) : (
          <Text variant="supporting" fontSize="sm">
            Waiting for {data.robot.name}…
          </Text>
        )}

        {data.error ? (
          <Text variant="supporting" color="status.danger.content">
            {data.error}
          </Text>
        ) : null}
      </LStack>
    </details>
  );
}
