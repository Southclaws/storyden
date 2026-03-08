"use client";

import { formatDistanceToNow } from "date-fns";
import Link from "next/link";

import { useRobotSessionsList } from "@/api/openapi-client/robots";
import {
  Account,
  RobotSessionRef,
  RobotSessionsListResult,
} from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { EmptyState } from "@/components/site/EmptyState";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { UnreadyBanner } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import { IconButton } from "@/components/ui/icon-button";
import { ArrowLeftIcon } from "@/components/ui/icons/Arrow";
import { LinkButton } from "@/components/ui/link-button";
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
  initialChatSessionList: RobotSessionsListResult;
  initialChatPage?: string;
};

export function RobotSessionListScreen(props: Props) {
  const { data, error } = useRobotSessionsList(
    {
      page: props.initialChatPage,
    },
    {
      swr: {
        fallbackData: props.initialChatSessionList,
      },
    },
  );

  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  const currentPage = props.initialChatPage
    ? parseInt(props.initialChatPage, 10)
    : 1;

  return (
    <LStack className={lstack()} gap="4" w="full">
      <WStack>
        <HStack gap="2">
          <Link href="/robots">
            <IconButton variant="ghost" size="xs">
              <ArrowLeftIcon />
            </IconButton>
          </Link>
          <Heading size="md">Robot Chat Sessions</Heading>
        </HStack>

        <LinkButton href="/robots/chats/new" variant="subtle" size="xs">
          New
        </LinkButton>
      </WStack>

      <styled.p color="fg.muted">
        View all conversations with robots across your community.
      </styled.p>

      <VStack w="full">
        <RobotSessionList data={data} currentPage={currentPage} />
      </VStack>
    </LStack>
  );
}

function RobotSessionList({
  data,
  currentPage,
}: {
  data: RobotSessionsListResult;
  currentPage: number;
}) {
  if (data.sessions.length === 0) {
    return (
      <EmptyState hideContributionLabel>No robot chat sessions yet.</EmptyState>
    );
  }

  return (
    <>
      <LStack gap="3" w="full">
        {data.sessions.map((session) => (
          <RobotSessionCard key={session.id} session={session} />
        ))}
      </LStack>

      <PaginationControls
        path="/robots/chats"
        currentPage={currentPage}
        totalPages={data.total_pages}
        pageSize={data.page_size}
      />
    </>
  );
}

type RobotSessionCardProps = {
  session: RobotSessionRef;
};

function RobotSessionCard({ session }: RobotSessionCardProps) {
  const timeAgo = formatDistanceToNow(new Date(session.createdAt), {
    addSuffix: true,
  });

  return (
    <CardBox w="full" _hover={{ background: "bg.emphasized" }} cursor="pointer">
      <Link href={`/robots/chats/${session.id}`}>
        <LStack gap="2">
          <WStack alignItems="center">
            <styled.p fontSize="sm" color="fg.subtle">
              {session.name}
            </styled.p>
            <styled.time fontSize="xs" color="fg.muted">
              {timeAgo}
            </styled.time>
          </WStack>

          <WStack>
            <MemberBadge profile={session.created_by} size="sm" name="handle" />
          </WStack>
        </LStack>
      </Link>
    </CardBox>
  );
}
