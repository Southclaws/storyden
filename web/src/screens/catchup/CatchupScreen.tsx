"use client";

import { createListCollection } from "@ark-ui/react";

import { ThreadReferenceList } from "@/components/post/ThreadReferenceList";
import { EmptyState } from "@/components/site/EmptyState";
import { PaginationControls } from "@/components/site/PaginationControls/PaginationControls";
import { Unready } from "@/components/site/Unready";
import { Button } from "@/components/ui/button";
import { CheckIcon } from "@/components/ui/icons/Check";
import { SelectIcon } from "@/components/ui/icons/Select";
import { PageHeader } from "@/components/ui/page-header";
import * as Select from "@/components/ui/select";
import * as Tabs from "@/components/ui/tabs";
import { LStack } from "@/styled-system/jsx";

import { type Props, useCatchupScreen } from "./useCatchupScreen";

export function CatchupScreen(props: Props) {
  const {
    data,
    error,
    view,
    category,
    isValidating,
    refresh,
    handlePageChange,
    handleViewChange,
    handleCategoryChange,
  } = useCatchupScreen(props);

  const collection = createListCollection({
    items: [
      { label: "All categories", value: "" },
      ...props.initialCategoryList.categories.map((category) => ({
        label: category.name,
        value: category.slug,
      })),
    ],
  });

  const selectedCategory = collection.items.find(
    (item) => item.value === category,
  );
  let emptyMessage = "You're all caught up.";
  if (view === "unresolved") {
    emptyMessage = "There are no unresolved threads.";
  } else if (view === "unanswered") {
    emptyMessage = "There are no unanswered threads.";
  }

  if (category) {
    emptyMessage = `There are no ${view} threads in ${selectedCategory?.label ?? "this category"}.`;
  }

  return (
    <LStack className="catchup-page" data-sd-page="catchup" gap="4">
      <PageHeader
        title="Catch Up"
        actions={
          <Button
            type="button"
            variant="ghost"
            disabled={isValidating}
            onClick={() => void refresh()}
          >
            Refresh
          </Button>
        }
      />

      <Tabs.Root
        width="full"
        size="sm"
        variant="line"
        activationMode="manual"
        value={view}
        onValueChange={({ value }) => handleViewChange(value)}
      >
        <Tabs.Bar>
          <Tabs.List aria-label="Catch Up views">
            <Tabs.Trigger value="unread">Unread</Tabs.Trigger>
            <Tabs.Trigger value="unresolved">Unresolved</Tabs.Trigger>
            <Tabs.Trigger value="unanswered">Unanswered</Tabs.Trigger>
            <Tabs.Indicator />
          </Tabs.List>

          <Tabs.Controls>
            <Select.Root
              collection={collection}
              value={[category]}
              onValueChange={({ value }) =>
                handleCategoryChange(value[0] ?? "")
              }
            >
              <Select.Control>
                <Select.Trigger aria-label="Filter by category">
                  <Select.ValueText placeholder="All categories" />
                  <SelectIcon />
                </Select.Trigger>
              </Select.Control>
              <Select.Positioner>
                <Select.Content>
                  {collection.items.map((item) => (
                    <Select.Item key={item.value} item={item}>
                      <Select.ItemText>{item.label}</Select.ItemText>
                      <Select.ItemIndicator>
                        <CheckIcon />
                      </Select.ItemIndicator>
                    </Select.Item>
                  ))}
                </Select.Content>
              </Select.Positioner>
              <Select.HiddenSelect />
            </Select.Root>
          </Tabs.Controls>
        </Tabs.Bar>

        <Tabs.Content value={view}>
          <LStack gap="4">
            {!data ? (
              <Unready error={error} />
            ) : data.threads.length === 0 ? (
              <EmptyState icon={<CheckIcon />} hideContributionLabel>
                <p>{emptyMessage}</p>
              </EmptyState>
            ) : (
              <>
                <ThreadReferenceList
                  initialSession={props.initialSession}
                  threads={data.threads}
                />
                <PaginationControls
                  path="/catchup"
                  params={{ view, category }}
                  currentPage={data.current_page}
                  totalPages={data.total_pages}
                  pageSize={data.page_size}
                  onClick={handlePageChange}
                />
              </>
            )}
          </LStack>
        </Tabs.Content>
      </Tabs.Root>
    </LStack>
  );
}
