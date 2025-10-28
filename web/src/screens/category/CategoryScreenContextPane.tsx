"use client";

import { CategoryBadge } from "@/components/category/CategoryBadge";
import { Unready } from "@/components/site/Unready";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { SlugIcon } from "@/components/ui/icons/Slug";
import * as Table from "@/components/ui/table";
import { css, cva } from "@/styled-system/css";
import { HStack, LStack } from "@/styled-system/jsx";
import { ScrollToTop } from "@/components/ui/scroll-to-top";

import { Props, useCategoryScreen } from "./CategoryScreen";

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

export function CategoryScreenContextPane(props: Props) {
  const { ready, error, data } = useCategoryScreen(props);
  if (!ready) {
    return <Unready error={error} />;
  }

  const { category } = data;

  const tableData = [
    {
      label: "slug",
      icon: SlugIcon,
      value: category.slug,
      style: "numeric" as const,
    },
    {
      label: category.postCount === 1 ? "thread" : "threads",
      icon: DiscussionIcon,
      value: `${category.postCount}`,
    },
  ];

  return (
    <LStack gap="1">
      <CategoryBadge category={category} />
      <p className={css({ color: "fg.muted" })}>{category.description}</p>

      <Table.Root size="sm">
        <Table.Body>
          {tableData.map((item) => (
            <Table.Row key={item.label}>
              <Table.Cell fontWeight="medium" color="fg.muted">
                <HStack gap="1">
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
    </LStack>
  );
}
