"use client";

import { SelectValueChangeDetails, createListCollection } from "@ark-ui/react";
import { AnimatePresence, motion } from "framer-motion";

import { Node } from "@/api/openapi-schema";
import { LibraryPageSelect } from "@/components/library/LibraryPageSelect";
import { InfoTip } from "@/components/site/InfoTip";
import { useSettingsContext } from "@/components/site/SettingsContext/SettingsContext";
import { Heading } from "@/components/ui/heading";
import { IconButton } from "@/components/ui/icon-button";
import { AdminIcon } from "@/components/ui/icons/Admin";
import { CategoryIcon } from "@/components/ui/icons/Category";
import { CheckIcon } from "@/components/ui/icons/Check";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { EditIcon } from "@/components/ui/icons/Edit";
import { LibraryIcon } from "@/components/ui/icons/Library";
import { SelectIcon } from "@/components/ui/icons/Select";
import * as Select from "@/components/ui/select";
import { Box, HStack, LStack, WStack, styled } from "@/styled-system/jsx";

import { Route, useRoute } from "../../useRoute";

const MotionBox = motion(Box);
const MotionSpan = motion.span;

const editableRoute: Record<Route["name"], boolean> = {
  index: true,
  library: false,
  admin: false,
  settings: false,
};

export function AdminZone() {
  const { isEditingEnabled, isEditing, handleToggleEditing } =
    useSettingsContext();

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
      bgColor="orange.4/50"
      borderTopRadius="md"
    >
      <HStack w="full" fontSize="xs" justify="space-between">
        <HStack gap="1">
          <AdminIcon w="4" />
          <AnimatePresence mode="wait">
            <MotionSpan
              key={isEditing ? "configure-feed" : "admin"}
              initial={{ opacity: 0, y: 2 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -2 }}
              transition={{ duration: 0.15 }}
            >
              {isEditing && route ? `Editing ${route.label}` : "Admin"}
            </MotionSpan>
          </AnimatePresence>
        </HStack>

        <HStack gap="1">
          {isRouteEditable && (
            <IconButton size="xs" variant="ghost" onClick={handleToggleEditing}>
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
              <RouteConfig route={route} />
            </LStack>
          </MotionBox>
        )}
      </AnimatePresence>
    </Box>
  );
}

function RouteConfig({ route }: { route?: Route }) {
  switch (route?.name) {
    case "index":
      return <FeedConfig />;

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

export function FeedConfig() {
  const { isEditingEnabled, isEditing, feed, updateFeed, handleToggleEditing } =
    useSettingsContext();

  if (!isEditingEnabled) {
    return null;
  }

  const collection = createListCollection({ items: sources });

  async function handleSourceTypeChange({ value }: SelectValueChangeDetails) {
    if (value.length === 0) {
      return;
    }

    const feedSourceType = value[0] as typeof feed.source.type;

    await updateFeed({
      layout: {
        type: "list",
      },
      source: {
        type: feedSourceType,
      },
    });
  }

  return (
    <LStack>
      <Select.Root
        size="xs"
        collection={collection}
        defaultValue={[feed.source.type]}
        positioning={{ sameWidth: false }}
        onValueChange={handleSourceTypeChange}
      >
        <WStack alignItems="center">
          <Select.Label>Source</Select.Label>

          <InfoTip title="Pick what is displayed on the home page">
            Change what the home page displays. For social use-cases, you can
            use Threads for a discussion feed. For a knowledge-base, curated
            directory or database, select Library and for a more traditional
            discussion board style, select Categories.
          </InfoTip>
        </WStack>
        <Select.Control>
          <Select.Trigger>
            <Select.ValueText placeholder="Select a Source" />
            <SelectIcon />
          </Select.Trigger>
        </Select.Control>
        <Select.Positioner>
          <Select.Content>
            {sources.map((item) => (
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

      <SourceConfig />
    </LStack>
  );
}

function SourceConfig() {
  const { feed } = useSettingsContext();

  switch (feed.source.type) {
    case "library":
      return <SourceLibraryConfig />;
    default:
      return null;
  }
}

function SourceLibraryConfig() {
  const { feed, updateFeed } = useSettingsContext();

  if (feed.source.type !== "library") {
    return null;
  }

  async function handleHomepageNodeChange(node: Node | undefined) {
    await updateFeed({
      layout: {
        type: "list",
      },
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
          Use a Page (optional)
        </Heading>

        <InfoTip title="Use a Page as the home screen">
          This allows you to pick a page from the library to use as the home
          page of your site. This is useful for directories if you want to
          showcase a set of child pages on the home page or for simply
          customising the home screen with layout blocks and rich content.
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
