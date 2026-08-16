"use client";

import { formatDistanceToNow } from "date-fns";
import { type CSSProperties, useEffect } from "react";

import { useRobotSessionGet } from "@/api/openapi-client/robots";
import { RobotSession } from "@/api/openapi-schema";
import { toStorydenUIMessages } from "@/api/robots-types";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { FullPageChatInput } from "@/components/robots/RobotChat/FullPageChatInput";
import { FullPageMessageList } from "@/components/robots/RobotChat/FullPageMessageList";
import { RobotWorkspaceSelect } from "@/components/robots/RobotWorkspaceSelect";
import { BackAction } from "@/components/site/Action/Back";
import {
  RobotChatContext,
  useRobotChat,
} from "@/components/site/CommandPalette/RobotChat/RobotChatContext";
import { UnreadyBanner } from "@/components/site/Unready";
import { PageHeading } from "@/components/ui/page-heading";
import { css } from "@/styled-system/css";
import { HStack, LStack, WStack, styled } from "@/styled-system/jsx";
import { token } from "@/styled-system/tokens";

type Props = {
  sessionId?: string;
  initialChatBefore?: string;
  initialChatLimit?: string;
  initialRobotID?: string;
};

const mobileViewportStyle = {
  "--robot-session-screen-mobile-height": `calc(${token("sizes.viewportHeight")} - var(--navigation-topbar-height))`,
} as CSSProperties;

const containerStyles = css({
  height: {
    base: "var(--robot-session-screen-mobile-height)",
    md: "viewportHeight",
  },
  maxHeight: {
    base: "var(--robot-session-screen-mobile-height)",
    md: "viewportHeight",
  },
  minHeight: "0",
  display: "flex",
  flexDirection: "column",
  justifyContent: "space-between",
  overflow: "hidden",
});

export function RobotSessionScreen(props: Props) {
  const isNewSession = props.sessionId === undefined;

  const { data, error } = useRobotSessionGet(
    props.sessionId ?? "",
    {
      before: props.initialChatBefore,
      limit: props.initialChatLimit,
    },
    { swr: { enabled: !isNewSession } },
  );

  if (!isNewSession && !data) {
    return <UnreadyBanner error={error} />;
  }

  const session = data ?? undefined;
  const messages = toStorydenUIMessages(session?.message_list.messages ?? []);
  const rootRobotID = isNewSession
    ? props.initialRobotID
    : session?.root_robot_id;
  return (
    <div
      className={containerStyles}
      data-testid="robot-session-screen"
      style={mobileViewportStyle}
    >
      <RobotChatContext
        initialSessionID={session?.id}
        initialMessages={messages}
        initialNextBefore={session?.message_list.next_before}
        initialStreamOffset={session?.stream_offset}
        initialActiveTurnID={session?.active_turn_id}
        initialRootRobotID={rootRobotID}
        initialSelectedWorkspaceID={session?.active_workspace?.workspace_id}
      >
        <ChatPageContent session={session} isNewSession={isNewSession} />
      </RobotChatContext>
    </div>
  );
}

function ChatPageContent({
  session,
  isNewSession,
}: {
  session?: RobotSession;
  isNewSession: boolean;
}) {
  const { sessionId, isSessionConfirmed } = useRobotChat();

  useEffect(() => {
    if (isNewSession && isSessionConfirmed && sessionId) {
      console.debug(
        `[RobotSessionScreen] Session confirmed, updating URL to: /robots/chats/${sessionId}`,
      );
      window.history.replaceState(
        window.history.state,
        "",
        `/robots/chats/${encodeURIComponent(sessionId)}`,
      );
    }
  }, [isNewSession, isSessionConfirmed, sessionId]);

  return (
    <>
      <ChatPageHeader session={session} isNewSession={isNewSession} />
      <FullPageMessageList />
      <FullPageChatInput />
      <WStack mt="1" flexShrink="0">
        <HStack>
          <RobotWorkspaceSelect variant="outline" />
        </HStack>
        {session && <StatusText session={session} />}
      </WStack>
    </>
  );
}

function ChatPageHeader({
  session,
  isNewSession,
}: {
  session?: RobotSession;
  isNewSession: boolean;
}) {
  const title = isNewSession ? "New Chat" : (session?.name ?? "Chat");

  return (
    <LStack flexShrink="0">
      <WStack alignItems="center" flexShrink="0">
        <HStack gap="2">
          <BackAction fallbackHref="/robots/chats" />
          <PageHeading>{title}</PageHeading>
        </HStack>
      </WStack>
    </LStack>
  );
}

function StatusText({ session }: { session: RobotSession }) {
  const timeAgo = formatDistanceToNow(new Date(session.createdAt), {
    addSuffix: true,
  });

  return (
    <HStack color="text.subtle" fontSize="xs" gap="1">
      <span>chat started by</span>
      <MemberBadge
        profile={session.created_by}
        avatar="hidden"
        size="xs"
        name="handle"
      />
      <styled.time>{timeAgo}</styled.time>
    </HStack>
  );
}
