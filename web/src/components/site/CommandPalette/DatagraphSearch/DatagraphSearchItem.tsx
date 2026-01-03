import { Command } from "cmdk";

import { DatagraphItemKind, DatagraphMatch } from "@/api/openapi-schema";
import { CalendarIcon } from "@/components/ui/icons/Calendar";
import { CollectionIcon } from "@/components/ui/icons/Collection";
import { DiscussionIcon } from "@/components/ui/icons/Discussion";
import { LibraryIcon } from "@/components/ui/icons/Library";
import { ProfileIcon } from "@/components/ui/icons/Profile";
import { ReplyIcon } from "@/components/ui/icons/Reply";
import { Box, HStack, LStack, styled } from "@/styled-system/jsx";

import {
  badgeColourCSS,
  getDatagraphKindColour,
  getDatagraphKindLabel,
} from "../../../datagraph/DatagraphItemCard";
import { Badge } from "../../../ui/badge";

type Props = {
  result: DatagraphMatch;
  handleNavigate: (path: string) => void;
};

export function DatagraphSearchItem({ result, handleNavigate }: Props) {
  const path = buildPermalink(result);

  function handleSelect() {
    handleNavigate(path);
  }

  const label = getDatagraphKindLabel(result.kind);
  const colour = getDatagraphKindColour(result.kind);
  const cssVars = badgeColourCSS(colour);

  return (
    <Command.Item value={result.name} onSelect={handleSelect}>
      <HStack gap="2" justify="space-between" w="full">
        <HStack gap="2" minW="0" flex="1">
          <Box w="4" flexShrink="0">
            {getIconForKind(result.kind)}
          </Box>
          <LStack gap="0" minW="0" flex="1">
            <styled.h1 lineClamp={1} fontSize="sm" fontWeight="semibold">
              {result.name}
            </styled.h1>
            {result.description && (
              <styled.span lineClamp={1} fontSize="xs" fontWeight="normal">
                {result.description}
              </styled.span>
            )}
          </LStack>
        </HStack>
        <Badge
          size="sm"
          style={cssVars}
          backgroundColor="var(--colors-color-palette-bg)"
          borderColor="var(--colors-color-palette-bo)"
          color="var(--colors-color-palette-fg)"
          flexShrink="0"
        >
          {label}
        </Badge>
      </HStack>
    </Command.Item>
  );
}

function buildPermalink(result: DatagraphMatch): string {
  switch (result.kind) {
    case DatagraphItemKind.thread:
      return `/t/${result.slug}`;

    case DatagraphItemKind.post:
    case DatagraphItemKind.reply:
      return `/t/locate/${result.id}`;

    case DatagraphItemKind.node:
      return `/l/${result.slug}`;

    case DatagraphItemKind.collection:
      return `/c/${result.slug}`;

    case DatagraphItemKind.profile:
      return `/m/${result.slug}`;

    case DatagraphItemKind.event:
      return `/e/${result.slug}`;

    default:
      return `/`;
  }
}

const getIconForKind = (kind: DatagraphItemKind) => {
  switch (kind) {
    case DatagraphItemKind.thread:
      return <DiscussionIcon />;
    case DatagraphItemKind.post:
      return <DiscussionIcon />;
    case DatagraphItemKind.reply:
      return <ReplyIcon />;
    case DatagraphItemKind.node:
      return <LibraryIcon />;
    case DatagraphItemKind.collection:
      return <CollectionIcon />;
    case DatagraphItemKind.profile:
      return <ProfileIcon />;
    case DatagraphItemKind.event:
      return <CalendarIcon />;
    default:
      return null;
  }
};
