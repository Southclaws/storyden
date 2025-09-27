import { Fragment } from "react";

import {
  CategoryReference,
  Thread,
  ThreadReference,
} from "src/api/openapi-schema";

import { BreadcrumbIcon } from "@/components/ui/icons/Breadcrumb";
import { LinkButton } from "@/components/ui/link-button";
import { Box, HStack, styled } from "@/styled-system/jsx";

import { DiscussionRoute } from "../site/Navigation/Anchors/Discussion";

type Props = {
  thread?: Thread;
};

type Breadcrumb =
  | {
      type: "category";
      key: string;
      category: CategoryReference;
    }
  | {
      type: "thread";
      key: string;
      thread: ThreadReference;
    };

export function Breadcrumbs({ thread }: Props) {
  const category = thread?.category
    ? [
        {
          type: "category" as const,
          key: thread.category.id,
          category: thread.category,
        },
      ]
    : [];

  // This is a list as we may do nested categories in future.
  const crumbs: Breadcrumb[] = thread
    ? [
        ...category,
        {
          type: "thread",
          key: thread.id,
          thread: thread,
        },
      ]
    : [];

  return (
    <HStack
      w="full"
      color="fg.subtle"
      overflowX="scroll"
      pt="scrollGutter"
      mt="-scrollGutter"
    >
      <LinkButton
        size="xs"
        variant="subtle"
        flexShrink="0"
        minW="min"
        href={DiscussionRoute}
      >
        Discussion
      </LinkButton>
      {crumbs.map((c) => {
        return (
          <Fragment key={c.key}>
            <Box flexShrink="0">
              <BreadcrumbIcon />
            </Box>

            <BreadcrumbButton breadcrumb={c} />
          </Fragment>
        );
      })}
    </HStack>
  );
}

function BreadcrumbButton({ breadcrumb }: { breadcrumb: Breadcrumb }) {
  switch (breadcrumb.type) {
    case "category":
      // TODO: Explore using the CategoryBadge component with subtle colour.
      return (
        <LinkButton
          size="xs"
          variant="subtle"
          flexShrink="0"
          maxW="64"
          overflow="hidden"
          href={`/d/${breadcrumb.category.slug}`}
        >
          {breadcrumb.category.name}
        </LinkButton>
      );

    case "thread":
      return (
        <LinkButton
          size="xs"
          variant="subtle"
          flexShrink="0"
          maxW="64"
          overflow="hidden"
          href={`/t/${breadcrumb.thread.slug}`}
        >
          <styled.span
            overflowX="hidden"
            textWrap="nowrap"
            textOverflow="ellipsis"
          >
            {breadcrumb.thread.title}
          </styled.span>
        </LinkButton>
      );
  }
}
