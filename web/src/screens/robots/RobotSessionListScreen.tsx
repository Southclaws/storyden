"use client";

import { formatDistanceToNow } from "date-fns";
import Link from "next/link";
import { useSearchParams } from "next/navigation";

import { useRobotSessionsList } from "@/api/openapi-client/robots";
import { RobotSessionRef, RobotSessionsListResult } from "@/api/openapi-schema";
import { MemberIdent } from "@/components/member/MemberBadge/MemberIdent";
import { BackAction } from "@/components/site/Action/Back";
import { EmptyState } from "@/components/site/EmptyState";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { UnreadyBanner } from "@/components/site/Unready";
import { CardBox } from "@/components/ui/card-box";
import { LinkButton } from "@/components/ui/link-button";
import { PageHeading } from "@/components/ui/page-heading";
import { Text } from "@/components/ui/text";
import { HStack, LStack, VStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

export function RobotSessionListScreen() {
  const searchParams = useSearchParams();
  const page = searchParams.get("page") ?? undefined;
  const { data, error } = useRobotSessionsList({
    page,
  });

  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  const currentPage = page ? parseInt(page, 10) : 1;

  return (
    <LStack className={lstack()} gap="4" w="full">
      <WStack>
        <HStack gap="2">
          <BackAction fallbackHref="/robots" />
          <PageHeading>Robot Chat Sessions</PageHeading>
        </HStack>

        <LinkButton href="/robots/chats/new" variant="subtle">
          New
        </LinkButton>
      </WStack>

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
      <PaginationControls
        path="/robots/chats"
        currentPage={currentPage}
        totalPages={data.total_pages}
        pageSize={data.page_size}
      />

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
    <CardBox>
      <Link href={`/robots/chats/${session.id}`}>
        <LStack gap="2">
          <WStack alignItems="center">
            <Text variant="supporting">{session.name}</Text>
            <styled.time fontSize="xs" color="text.subtle">
              {timeAgo}
            </styled.time>
          </WStack>

          <WStack>
            <MemberIdent profile={session.created_by} size="sm" name="handle" />
          </WStack>
        </LStack>
      </Link>
    </CardBox>
  );
}
