import { formatDate, formatDistance, formatDistanceStrict } from "date-fns";
import { Fragment } from "react";

import { Post } from "src/api/openapi-schema";

import { VStack, styled } from "@/styled-system/jsx";

import { PostView } from "./PostView/PostView";

type Props = {
  slug?: string;
  posts: Post[];
};

export function PostListView(props: Props) {
  return (
    <styled.ol
      listStyleType="none"
      m="0"
      gap="4"
      display="flex"
      flexDir="column"
      width="full"
    >
      {props.posts.map((p, i) => {
        const previous = props.posts[i - 1];
        const start = previous ? new Date(previous.createdAt) : undefined;
        const end = new Date(p.createdAt);

        return (
          <Fragment key={p.id}>
            {start && <IntervalDivider interval={{ start, end }} />}

            <styled.li listStyleType="none" m="0">
              <PostView {...p} />
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
