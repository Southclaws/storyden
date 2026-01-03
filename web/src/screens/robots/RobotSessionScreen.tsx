"use client";

import { formatDistanceToNow } from "date-fns";
import Link from "next/link";

import { useRobotSessionGet } from "@/api/openapi-client/robots";
import {
  Account,
  RobotSession,
  RobotSessionMessage,
} from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { EmptyState } from "@/components/site/EmptyState";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { UnreadyBanner } from "@/components/site/Unready";
import { Badge } from "@/components/ui/badge";
import { Heading } from "@/components/ui/heading";
import { ArrowLeftIcon } from "@/components/ui/icons/Arrow";
import { IconButton } from "@/components/ui/icon-button";
import {
  CardBox,
  HStack,
  LStack,
  VStack,
  WStack,
  styled,
} from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

type Props = {
  initialSession: Account;
  initialChatSession: RobotSession;
  initialChatPage?: string;
};

export function RobotSessionScreen(props: Props) {
  const { data, error } = useRobotSessionGet(
    props.initialChatSession.id,
    {
      page: props.initialChatPage,
    },
    {
      swr: {
        fallbackData: props.initialChatSession,
      },
    },
  );

  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  const currentPage = props.initialChatPage
    ? parseInt(props.initialChatPage, 10)
    : 1;

  const timeAgo = formatDistanceToNow(new Date(data.createdAt), {
    addSuffix: true,
  });

  return (
    <CardBox className={lstack()} gap="4">
      <WStack alignItems="center">
        <HStack gap="2">
          <Link href="/robots/chats">
            <IconButton variant="ghost" size="sm">
              <ArrowLeftIcon />
            </IconButton>
          </Link>
          <Heading size="md">Chat Session</Heading>
        </HStack>
        <styled.time fontSize="sm" color="fg.muted">
          {timeAgo}
        </styled.time>
      </WStack>

      <HStack>
        <styled.p fontSize="sm" color="fg.subtle">
          Started by:
        </styled.p>
        <MemberBadge profile={data.created_by} size="sm" name="handle" />
      </HStack>

      <MessageList
        sessionId={data.id}
        messages={data.message_list.messages ?? []}
        currentPage={currentPage}
        totalPages={data.message_list.total_pages}
        pageSize={data.message_list.page_size}
      />
    </CardBox>
  );
}

type MessageListProps = {
  sessionId: string;
  messages: RobotSessionMessage[];
  currentPage: number;
  totalPages: number;
  pageSize: number;
};

function MessageList({
  sessionId,
  messages,
  currentPage,
  totalPages,
  pageSize,
}: MessageListProps) {
  if (messages.length === 0) {
    return (
      <EmptyState hideContributionLabel>
        No messages in this session.
      </EmptyState>
    );
  }

  return (
    <>
      <VStack gap="3" w="full" alignItems="stretch">
        {messages.map((message) => (
          <MessageCard key={message.id} message={message} />
        ))}
      </VStack>

      <PaginationControls
        path={`/robots/chats/${sessionId}`}
        currentPage={currentPage}
        totalPages={totalPages}
        pageSize={pageSize}
      />
    </>
  );
}

type MessageCardProps = {
  message: RobotSessionMessage;
};

function MessageCard({ message }: MessageCardProps) {
  const timeAgo = formatDistanceToNow(new Date(message.created_at), {
    addSuffix: true,
  });

  const isRobotMessage = !!message.robot;

  return (
    <CardBox
      borderLeftWidth="thick"
      borderLeftStyle="solid"
      borderLeftColor={isRobotMessage ? "accent.default" : "border.default"}
    >
      <LStack gap="2">
        <WStack alignItems="center">
          {isRobotMessage ? (
            <HStack gap="2">
              <Badge size="sm" variant="solid">
                {message.robot?.name ?? "Robot"}
              </Badge>
            </HStack>
          ) : (
            <HStack gap="2">
              {message.author && (
                <MemberBadge
                  profile={message.author}
                  size="sm"
                  name="handle"
                />
              )}
            </HStack>
          )}
          <styled.time fontSize="xs" color="fg.muted">
            {timeAgo}
          </styled.time>
        </WStack>

        <styled.p fontSize="sm" fontFamily="mono" color="fg.subtle">
          {message.id}
        </styled.p>
      </LStack>
    </CardBox>
  );
}
