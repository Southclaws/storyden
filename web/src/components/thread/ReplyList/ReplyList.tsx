import { formatDate, formatDistance, formatDistanceStrict } from "date-fns";
import { Fragment } from "react";

import { Account, Thread } from "src/api/openapi-schema";

import { VStack, styled } from "@/styled-system/jsx";

import { Reply } from "../Reply/Reply";

type Props = {
  initialSession?: Account;
  thread: Thread;
  currentPage?: number;
};

export function ReplyList({ initialSession, thread, currentPage }: Props) {
  return (
    <styled.ol
      listStyleType="none"
      m="0"
      gap="4"
      display="flex"
      flexDir="column"
      width="full"
    >
      {thread.replies.replies.map((reply, i) => {
        const previous = thread.replies[i - 1];
        const start = previous ? new Date(previous.createdAt) : undefined;
        const end = new Date(reply.createdAt);

        return (
          <Fragment key={reply.id}>
            {start && <IntervalDivider interval={{ start, end }} />}

            <styled.li listStyleType="none" m="0">
              <Reply initialSession={initialSession} thread={thread} reply={reply} currentPage={currentPage} />
            </styled.li>
          </Fragment>
        );
      })}
    </styled.ol>
  );
}

export type IntervalDividerProps = {
  interval: {
    start: Date;
    end: Date;
  };
};

export function IntervalDivider({ interval }: IntervalDividerProps) {
  const difference = interval.end.getTime() - interval.start.getTime();

  if (difference < 8640000) {
    return null;
  }

  const startLabel = formatDate(interval.start, "PP");
  const endLabel = formatDate(interval.end, "PP");

  const title = `${startLabel} - ${formatDistanceStrict(interval.start, interval.end)} - ${endLabel}`;

  return (
    <VStack w="full" color="fg.subtle" fontSize="xs">
      <time title={title}>
        {formatDistance(interval.start, interval.end)} later
      </time>
    </VStack>
  );
}
