"use client";

import { formatDistanceToNow } from "date-fns";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect } from "react";

import { useRobotSessionGet } from "@/api/openapi-client/robots";
import { Account, RobotSession } from "@/api/openapi-schema";
import { toStorydenUIMessages } from "@/api/robots-types";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { FullPageChatInput } from "@/components/robots/RobotChat/FullPageChatInput";
import { FullPageMessageList } from "@/components/robots/RobotChat/FullPageMessageList";
import { RobotListMenu } from "@/components/robots/RobotListMenu";
import { RobotWorkspaceSelect } from "@/components/robots/RobotWorkspaceSelect";
import {
  RobotChatContext,
  useRobotChat,
} from "@/components/site/CommandPalette/RobotChat/RobotChatContext";
import { UnreadyBanner } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import { IconButton } from "@/components/ui/icon-button";
import { ArrowLeftIcon } from "@/components/ui/icons/Arrow";
import { css } from "@/styled-system/css";
import { HStack, LStack, WStack, styled } from "@/styled-system/jsx";

type Props = {
  initialSession: Account;
  initialChatSession: RobotSession | null;
  initialChatBefore?: string;
  initialChatLimit?: string;
  initialSelectedRobotID?: string;
};

const containerStyles = css({
  height: "viewportHeight",
  maxHeight: "viewportHeight",
  display: "flex",
  flexDirection: "column",
  justifyContent: "space-between",
});

export function RobotSessionScreen(props: Props) {
  const isNewSession = props.initialChatSession === null;

  // Only fetch if we have an existing session
  const { data, error } = useRobotSessionGet(
    props.initialChatSession?.id ?? "",
    {
      before: props.initialChatBefore,
      limit: props.initialChatLimit,
    },
    {
      swr: {
        fallbackData: props.initialChatSession ?? undefined,
        isPaused: () => isNewSession,
      },
    },
  );

  if (!isNewSession && !data) {
    return <UnreadyBanner error={error} />;
  }

  const session = data ?? undefined;
  const messages = toStorydenUIMessages(session?.message_list.messages ?? []);
  const initialSelectedRobotID =
    props.initialSelectedRobotID ?? session?.active_robot_id;

  return (
    <div className={containerStyles}>
      <RobotChatContext
        initialSessionID={session?.id}
        initialMessages={messages}
        initialNextBefore={session?.message_list.next_before}
        initialSelectedRobotID={initialSelectedRobotID}
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
  const router = useRouter();
  const searchParams = useSearchParams();
  const { sessionId, isSessionConfirmed, selectedRobot } = useRobotChat();
  const selectedRobotID = selectedRobot?.id;
  const currentRobotID = searchParams.get("robot") ?? undefined;

  useEffect(() => {
    if (isNewSession && isSessionConfirmed && sessionId) {
      const query = selectedRobotID ? `?robot=${selectedRobotID}` : "";

      console.debug(
        `[RobotSessionScreen] Session confirmed, updating URL to: /robots/chats/${sessionId}`,
      );
      router.replace(`/robots/chats/${sessionId}${query}`);
    }
  }, [isNewSession, isSessionConfirmed, selectedRobotID, sessionId, router]);

  useEffect(() => {
    if (isNewSession || !session?.id) {
      return;
    }

    if (selectedRobotID === currentRobotID) {
      return;
    }

    const query = selectedRobotID ? `?robot=${selectedRobotID}` : "";

    router.replace(`/robots/chats/${session.id}${query}`);
  }, [currentRobotID, isNewSession, router, selectedRobotID, session?.id]);

  return (
    <>
      <ChatPageHeader session={session} isNewSession={isNewSession} />
      <FullPageMessageList />
      <FullPageChatInput />
      <WStack mt="1">
        <HStack>
          <RobotListMenu />
          <RobotWorkspaceSelect size="xs" variant="outline" />
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
          <Link href="/robots/chats">
            <IconButton variant="ghost" size="xs">
              <ArrowLeftIcon />
            </IconButton>
          </Link>
          <Heading size="md">{title}</Heading>
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
    <HStack color="fg.muted" fontSize="xs" gap="1">
      <styled.span>chat started by</styled.span>
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
