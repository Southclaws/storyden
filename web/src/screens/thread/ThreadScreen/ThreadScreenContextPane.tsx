"use client";

import { formatDate } from "date-fns";
import {
  CalendarIcon,
  HandHeartIcon,
  PencilLineIcon,
  TagIcon,
  UsersIcon,
} from "lucide-react";

import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { Unready } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import * as Table from "@/components/ui/table";
import { cva } from "@/styled-system/css";
import { HStack, LStack, styled } from "@/styled-system/jsx";
import { button } from "@/styled-system/recipes";

import { Props, useThreadScreen } from "./useThreadScreen";

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
      },
    },
  },
});

export function ThreadScreenContextPane(props: Props) {
  const { ready, error, form, isEditing, isEmpty, resetKey, data, handlers } =
    useThreadScreen(props);
  if (!ready) {
    return <Unready error={error} />;
  }

  const { thread } = data;

  const tableData = [
    {
      label: "ID",
      icon: TagIcon,
      value: thread.id,
      style: "numeric" as const,
    },
    {
      label: "author",
      icon: PencilLineIcon,
      value: (
        <MemberBadge
          profile={thread.author}
          size="xs"
          avatar="hidden"
          name="full-horizontal"
        />
      ),
    },
    {
      label: "started",
      icon: CalendarIcon,
      value: formatDate(thread.createdAt, "MMM d, yyyy"),
    },
    {
      label: "replies",
      icon: UsersIcon,
      value: `${thread.reply_status.replies}`,
    },
    {
      label: "participating",
      icon: HandHeartIcon,
      value: thread.reply_status.replied ? "Yes" : "No",
    },
  ];

  return (
    <LStack gap="1">
      <Heading>{thread.title}</Heading>
      <styled.p color="fg.muted">{thread.description}</styled.p>

      <Table.Root size="sm">
        <Table.Body>
          {tableData.map((item) => (
            <Table.Row key={item.label}>
              <Table.Cell fontWeight="medium" color="fg.muted">
                <HStack gap="1">
                  <item.icon width="14" />
                  <span>{item.label}</span>
                </HStack>
              </Table.Cell>
              <Table.Cell
                className={valueStyles({ style: item.style })}
                display="flex"
                justifyContent="flex-end"
                alignItems="center"
                textAlign="right"
              >
                {item.value}
              </Table.Cell>
            </Table.Row>
          ))}
        </Table.Body>
      </Table.Root>

      <styled.p>
        <styled.a
          color="fg.muted"
          className={button({
            variant: "subtle",
            size: "xs",
          })}
          href="#"
        >
          scroll to top
        </styled.a>
      </styled.p>
    </LStack>
  );
}
