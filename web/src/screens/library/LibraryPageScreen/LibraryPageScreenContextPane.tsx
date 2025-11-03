"use client";

import { formatDate } from "date-fns";

import { useNodeGet } from "@/api/openapi-client/nodes";
import { NodeWithChildren } from "@/api/openapi-schema";
import { LibraryPageCommentsList } from "@/components/library/comments/LibraryPageCommentsList";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { Unready } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import { AuthorIcon } from "@/components/ui/icons/Author";
import { CalendarIcon } from "@/components/ui/icons/Calendar";
import { MembersIcon } from "@/components/ui/icons/Members";
import { ParticipatingIcon } from "@/components/ui/icons/Participating";
import { SlugIcon } from "@/components/ui/icons/Slug";
import { ScrollToTop } from "@/components/ui/scroll-to-top";
import * as Table from "@/components/ui/table";
import { css, cva } from "@/styled-system/css";
import { HStack, LStack } from "@/styled-system/jsx";

type Props = {
  slug: string;
  initialNode: NodeWithChildren;
};

const valueStyles = cva({
  base: {},
  defaultVariants: {
    style: "base",
  },
  variants: {
    style: {
      base: {},
      numeric: {
        fontVariant: "tabular-nums",
        fontFamily: "mono",
        textOverflow: "ellipsis",
        overflowX: "hidden",
        display: "flex",
        alignItems: "center",
        height: "full",
      },
    },
  },
});

export function LibraryPageScreenContextPane(props: Props) {
  const { data, error } = useNodeGet(props.slug);
  if (!data) {
    return <Unready error={error} />;
  }

  const tableData = [
    {
      label: "slug",
      icon: SlugIcon,
      style: "numeric" as const,
      value: (
        <span className={valueStyles({ style: "numeric" })}>{data.id}</span>
      ),
    },
    {
      label: "author",
      icon: AuthorIcon,
      value: (
        <HStack>
          <MemberBadge
            profile={data.owner}
            size="sm"
            // avatar="hidden"
            name="handle"
          />
        </HStack>
      ),
    },
    {
      label: "created",
      icon: CalendarIcon,
      value: formatDate(data.createdAt, "MMM d, yyyy"),
    },
  ];

  return (
    <LStack gap="1">
      <Heading>{data.name}</Heading>
      <p className={css({ color: "fg.muted" })}>{data.description}</p>

      <Table.Root size="sm" tableLayout="fixed" w="full" overflow="hidden">
        <Table.Body>
          {tableData.map((item) => (
            <Table.Row key={item.label}>
              <Table.Cell fontWeight="medium" color="fg.muted">
                <HStack gap="1" flexShrink="0">
                  <item.icon width="4" />
                  <span>{item.label}</span>
                </HStack>
              </Table.Cell>
              <Table.Cell
                className={valueStyles({ style: item.style })}
                display="flex"
                justifyContent="flex-end"
                alignItems="center"
                textAlign="right"
                overflow="hidden"
                textOverflow="ellipsis"
                width="full"
                maxWidth="full"
                minW="0"
              >
                {item.value}
              </Table.Cell>
            </Table.Row>
          ))}
        </Table.Body>
      </Table.Root>

      <p>
        <ScrollToTop />
      </p>

      {/* <LibraryPageCommentsList node={data} /> */}
    </LStack>
  );
}
