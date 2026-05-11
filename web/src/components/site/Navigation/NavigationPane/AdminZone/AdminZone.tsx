"use client";

import {
  CheckboxCheckedChangeDetails,
  SelectValueChangeDetails,
  createListCollection,
} from "@ark-ui/react";
import { AnimatePresence, motion } from "framer-motion";
import { useSWRConfig } from "swr";

import { type Account, type Node } from "@/api/openapi-schema";
import { LibraryPageSelect } from "@/components/library/LibraryPageSelect";
import { InfoTip } from "@/components/site/InfoTip";
import * as Checkbox from "@/components/ui/checkbox";
import { Heading } from "@/components/ui/heading";
import { IconButton } from "@/components/ui/icon-button";
import { AdminIcon } from "@/components/ui/icons/Admin";
import { CategoryIcon } from "@/components/ui/icons/Category";
import { CheckIcon } from "@/components/ui/icons/Check";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { EditIcon } from "@/components/ui/icons/Edit";
import { LayoutGridIcon } from "@/components/ui/icons/LayoutGrid";
import { LayoutListIcon } from "@/components/ui/icons/LayoutList";
import { LibraryIcon } from "@/components/ui/icons/Library";
import { SelectIcon } from "@/components/ui/icons/Select";
import * as Select from "@/components/ui/select";
import { useI18n } from "@/i18n/provider";
import { type FeedConfig } from "@/lib/settings/feed";
import {
  useFeedConfig,
  useFeedEditorState,
  useFeedMutation,
} from "@/lib/settings/feed-client";
import { type Settings } from "@/lib/settings/settings";
import { Box, HStack, LStack, WStack, styled } from "@/styled-system/jsx";

import { Route, useRoute } from "../../useRoute";

const MotionBox = motion.create(Box);
const MotionSpan = motion.span;

const editableRoute: Record<Route["name"], boolean> = {
  index: true,
  library: false,
  admin: false,
  settings: false,
};

type Props = {
  initialSession?: Account;
  initialSettings?: Settings;
};

type UpdateFeed = (feed: FeedConfig) => Promise<void>;

export function AdminZone({ initialSession, initialSettings }: Props) {
  const { t } = useI18n();
  const feed = useFeedConfig(initialSettings, false);
  const { updateFeed } = useFeedMutation();
  const { isEditingEnabled, isEditing, handleToggleEditing } =
    useFeedEditorState({ initialSession, initialSettings });

  const route = useRoute();
  const isRouteEditable = route && editableRoute[route.name];

  if (!isEditingEnabled) {
    return null;
  }

  return (
    <Box
      w="full"
      pl="3"
      pr="2"
      py="1"
      bgColor="bg.warning"
      borderTopRadius="md"
    >
      <HStack w="full" fontSize="xs" justify="space-between">
        <HStack gap="1" color="fg.warning">
          <AdminIcon w="4" />
          <AnimatePresence mode="wait">
            <MotionSpan
              key={isEditing ? "configure-feed" : "admin"}
              initial={{ opacity: 0, y: 2 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -2 }}
              transition={{ duration: 0.15 }}
            >
              {isEditing && route
                ? t("Editing {{label}}", { label: t(route.label) })
                : t("Admin")}
            </MotionSpan>
          </AnimatePresence>
        </HStack>

        <HStack gap="1">
          {isRouteEditable && (
            <IconButton
              size="xs"
              variant="ghost"
              onClick={handleToggleEditing}
              type="button"
              aria-label={
                isEditing ? t("Close feed editor") : t("Open feed editor")
              }
              aria-pressed={isEditing}
              title={
                isEditing ? t("Close feed editor") : t("Open feed editor")
              }
            >
              <EditIcon w="4" />
            </IconButton>
          )}
        </HStack>
      </HStack>

      <AnimatePresence initial={false}>
        {isEditing && (
          <MotionBox
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            overflow="hidden"
            transition={{ duration: 0.25, ease: "easeInOut" }}
          >
            <LStack py="2">
              <RouteConfig feed={feed} updateFeed={updateFeed} route={route} />
            </LStack>
          </MotionBox>
        )}
      </AnimatePresence>
    </Box>
  );
}

function RouteConfig({
  route,
  feed,
  updateFeed,
}: {
  route?: Route;
  feed: FeedConfig;
  updateFeed: UpdateFeed;
}) {
  switch (route?.name) {
    case "index":
      return <FeedConfig feed={feed} updateFeed={updateFeed} />;

    default:
      return null;
  }
}

const sources = [
  {
    label: "Threads",
    value: "threads" as const,
    icon: <DiscussionIcon width="4" />,
  },
  {
    label: "Library",
    value: "library" as const,
    icon: <LibraryIcon width="4" />,
  },
  {
    label: "Categories",
    value: "categories" as const,
    icon: <CategoryIcon width="4" />,
  },
];

const layouts = [
  {
    label: "List",
    value: "list" as const,
    icon: <LayoutListIcon width="4" />,
  },
  {
    label: "Grid",
    value: "grid" as const,
    icon: <LayoutGridIcon width="4" />,
  },
];

export function FeedConfig({
  feed,
  updateFeed,
}: {
  feed: FeedConfig;
  updateFeed: UpdateFeed;
}) {
  const { t } = useI18n();
  const translatedSources = sources.map((item) => ({
    ...item,
    label: t(item.label),
  }));
  const translatedLayouts = layouts.map((item) => ({
    ...item,
    label: t(item.label),
  }));
  const sourceCollection = createListCollection({ items: translatedSources });
  const layoutCollection = createListCollection({ items: translatedLayouts });

  const canUpdateLayout =
    feed.source.type === "categories" ||
    (feed.source.type === "library" && feed.source.node === undefined);

  async function handleSourceTypeChange({ value }: SelectValueChangeDetails) {
    if (value.length === 0) {
      return;
    }

    const feedSourceType = value[0] as typeof feed.source.type;

    // Create proper source config based on type
    let sourceConfig: typeof feed.source;
    switch (feedSourceType) {
      case "threads":
        sourceConfig = { type: "threads", quickShare: "enabled" };
        break;
      case "library":
        sourceConfig = { type: "library" };
        break;
      case "categories":
        sourceConfig = {
          type: "categories",
          threadListMode: "uncategorised",
          quickShare: "enabled",
        };
        break;
      default:
        return;
    }

    await updateFeed({
      layout: feed.layout,
      source: sourceConfig,
    });
  }

  async function handleLayoutTypeChange({ value }: SelectValueChangeDetails) {
    if (value.length === 0) {
      return;
    }

    const feedLayoutType = value[0] as typeof feed.layout.type;

    await updateFeed({
      layout: {
        type: feedLayoutType,
      },
      source: feed.source,
    });
  }

  return (
    <LStack>
      <Select.Root
        size="xs"
        collection={sourceCollection}
        defaultValue={[feed.source.type]}
        positioning={{ sameWidth: false }}
        onValueChange={handleSourceTypeChange}
      >
        <WStack alignItems="center">
          <Select.Label>{t("Source")}</Select.Label>

          <InfoTip title={t("Pick what is displayed on the home page")}>
            {t(
              "Change what the home page displays. For social use-cases, you can use Threads for a discussion feed. For a knowledge-base, curated directory or database, select Library and for a more traditional discussion board style, select Categories.",
            )}
          </InfoTip>
        </WStack>
        <Select.Control>
          <Select.Trigger>
            <Select.ValueText placeholder={t("Select a source")} />
            <SelectIcon />
          </Select.Trigger>
        </Select.Control>
        <Select.Positioner>
          <Select.Content>
            {translatedSources.map((item) => (
              <Select.Item key={item.value} item={item}>
                <Select.ItemText mr="2">
                  <HStack gap="1">
                    <styled.span w="4">{item.icon}</styled.span>
                    <styled.span>{item.label}</styled.span>
                  </HStack>
                </Select.ItemText>
                <Select.ItemIndicator>
                  <CheckIcon />
                </Select.ItemIndicator>
              </Select.Item>
            ))}
          </Select.Content>
        </Select.Positioner>
      </Select.Root>

      {canUpdateLayout && (
        <Select.Root
          size="xs"
          collection={layoutCollection}
          defaultValue={[feed.layout.type]}
          positioning={{ sameWidth: false }}
          onValueChange={handleLayoutTypeChange}
        >
          <WStack alignItems="center">
            <Select.Label>{t("Layout")}</Select.Label>

          </WStack>
          <Select.Control>
            <Select.Trigger>
              <Select.ValueText placeholder={t("Select a layout")} />
              <SelectIcon />
            </Select.Trigger>
          </Select.Control>
          <Select.Positioner>
            <Select.Content>
              {translatedLayouts.map((item) => (
                <Select.Item key={item.value} item={item}>
                  <Select.ItemText mr="2">
                    <HStack gap="1">
                      <styled.span w="4">{item.icon}</styled.span>
                      <styled.span>{item.label}</styled.span>
                    </HStack>
                  </Select.ItemText>
                  <Select.ItemIndicator>
                    <CheckIcon />
                  </Select.ItemIndicator>
                </Select.Item>
              ))}
            </Select.Content>
          </Select.Positioner>
        </Select.Root>
      )}

      <SourceConfig feed={feed} updateFeed={updateFeed} />
    </LStack>
  );
}

function SourceConfig({
  feed,
  updateFeed,
}: {
  feed: FeedConfig;
  updateFeed: UpdateFeed;
}) {
  switch (feed.source.type) {
    case "threads":
      return <SourceThreadsConfig feed={feed} updateFeed={updateFeed} />;
    case "library":
      return <SourceLibraryConfig feed={feed} updateFeed={updateFeed} />;
    case "categories":
      return <SourceCategoriesConfig feed={feed} updateFeed={updateFeed} />;
    default:
      return null;
  }
}

function SourceThreadsConfig({
  feed,
  updateFeed,
}: {
  feed: FeedConfig;
  updateFeed: UpdateFeed;
}) {
  const { t } = useI18n();

  if (feed.source.type !== "threads") {
    return null;
  }

  async function handleQuickShareChange({
    checked,
  }: CheckboxCheckedChangeDetails) {
    await updateFeed({
      layout: feed.layout,
      source: {
        type: "threads",
        quickShare: checked ? "enabled" : "disabled",
      },
    });
  }

  return (
    <LStack gap="1">
      <WStack alignItems="center">
        <Heading fontWeight="medium" size="xs">
          {t("Quick Share")}
        </Heading>

        <InfoTip title={t("Show quick share box")}>
          {t(
            "Display a quick share box at the top of the thread list to allow users to quickly create new threads.",
          )}
        </InfoTip>
      </WStack>

      <Checkbox.Root
        size="sm"
        checked={feed.source.quickShare === "enabled"}
        onCheckedChange={handleQuickShareChange}
      >
        <Checkbox.Control>
          <Checkbox.Indicator>
            <CheckIcon />
          </Checkbox.Indicator>
        </Checkbox.Control>
        <Checkbox.Label>{t("Show Quick Share")}</Checkbox.Label>
        <Checkbox.HiddenInput />
      </Checkbox.Root>
    </LStack>
  );
}

function SourceLibraryConfig({
  feed,
  updateFeed,
}: {
  feed: FeedConfig;
  updateFeed: UpdateFeed;
}) {
  const { t } = useI18n();

  if (feed.source.type !== "library") {
    return null;
  }

  async function handleHomepageNodeChange(node: Node | undefined) {
    await updateFeed({
      layout: feed.layout,
      source: {
        type: "library",
        node: node?.id,
      },
    });
  }

  return (
    <LStack gap="1">
      <WStack alignItems="center">
        <Heading fontWeight="medium" size="xs">
          {t("Use a Page (optional)")}
        </Heading>

        <InfoTip title={t("Use a Page as the home screen")}>
          {t(
            "This allows you to pick a page from the library to use as the home page of your site. This is useful for directories if you want to showcase a set of child pages on the home page or for simply customising the home screen with layout blocks and rich content.",
          )}
        </InfoTip>
      </WStack>
      <LibraryPageSelect
        onChange={handleHomepageNodeChange}
        value={feed.source.node}
        defaultValue={feed.source.node}
      />
    </LStack>
  );
}

const threadListModes = [
  {
    label: "None",
    value: "none" as const,
  },
  {
    label: "All threads",
    value: "all" as const,
  },
  {
    label: "Uncategorised only",
    value: "uncategorised" as const,
  },
];

function SourceCategoriesConfig({
  feed,
  updateFeed,
}: {
  feed: FeedConfig;
  updateFeed: UpdateFeed;
}) {
  const { t } = useI18n();
  const { mutate } = useSWRConfig();

  if (feed.source.type !== "categories") {
    return null;
  }

  const threadListModeCollection = createListCollection({
    items: threadListModes.map((item) => ({ ...item, label: t(item.label) })),
  });

  async function handleThreadListModeChange({
    value,
  }: SelectValueChangeDetails) {
    if (value.length === 0 || feed.source.type !== "categories") {
      return;
    }

    const mode = value[0] as "none" | "all" | "uncategorised";

    await updateFeed({
      layout: feed.layout,
      source: {
        type: "categories",
        threadListMode: mode,
        quickShare: feed.source.quickShare,
      },
    });

    // Invalidate all /threads calls to trigger re-fetch with new params
    await mutate(
      (key) => typeof key === "string" && key.startsWith("/threads"),
    );
  }

  async function handleQuickShareChange({
    checked,
  }: CheckboxCheckedChangeDetails) {
    if (feed.source.type !== "categories") {
      return;
    }

    await updateFeed({
      layout: feed.layout,
      source: {
        type: "categories",
        threadListMode: feed.source.threadListMode,
        quickShare: checked ? "enabled" : "disabled",
      },
    });
  }

  return (
    <LStack gap="1">
      <WStack alignItems="center">
        <Heading fontWeight="medium" size="xs">
          {t("Thread list display")}
        </Heading>

        <InfoTip title={t("Choose what threads to show")}>
          {t(
            'Control which threads are displayed below the categories. Select "None" to only show categories, "All threads" to show threads from all categories, or "Uncategorised only" to show only threads without a category.',
          )}
        </InfoTip>
      </WStack>

      <Select.Root
        size="xs"
        collection={threadListModeCollection}
        defaultValue={[feed.source.threadListMode]}
        positioning={{ sameWidth: false }}
        onValueChange={handleThreadListModeChange}
      >
        <Select.Control>
          <Select.Trigger>
            <Select.ValueText placeholder={t("Select thread list mode")} />
            <SelectIcon />
          </Select.Trigger>
        </Select.Control>
        <Select.Positioner>
          <Select.Content>
            {threadListModeCollection.items.map((item) => (
              <Select.Item key={item.value} item={item}>
                <Select.ItemText>{item.label}</Select.ItemText>
                <Select.ItemIndicator>
                  <CheckIcon />
                </Select.ItemIndicator>
              </Select.Item>
            ))}
          </Select.Content>
        </Select.Positioner>
      </Select.Root>

      <WStack alignItems="center">
        <Heading fontWeight="medium" size="xs">
          {t("Quick Share")}
        </Heading>

        <InfoTip title={t("Show quick share box")}>
          {t(
            "Display a quick share box at the top of the thread list to allow users to quickly create new threads.",
          )}
        </InfoTip>
      </WStack>

      <Checkbox.Root
        size="sm"
        checked={feed.source.quickShare === "enabled"}
        onCheckedChange={handleQuickShareChange}
      >
        <Checkbox.Control>
          <Checkbox.Indicator>
            <CheckIcon />
          </Checkbox.Indicator>
        </Checkbox.Control>
        <Checkbox.Label>{t("Show Quick Share")}</Checkbox.Label>
        <Checkbox.HiddenInput />
      </Checkbox.Root>
    </LStack>
  );
}
