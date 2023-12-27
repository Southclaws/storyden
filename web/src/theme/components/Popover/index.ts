import { Popover as ArkPopover } from "@ark-ui/react/popover";
import { styled } from "styled-system/jsx";
import { popover } from "styled-system/recipes";

import { createStyleContext } from "src/theme/create-style-context";

const { withProvider, withContext } = createStyleContext(popover);

const Popover = withProvider(ArkPopover.Root);
const PopoverAnchor = withContext(styled(ArkPopover.Anchor), "anchor");
const PopoverArrow = withContext(styled(ArkPopover.Arrow), "arrow");
const PopoverArrowTip = withContext(styled(ArkPopover.ArrowTip), "arrowTip");
const PopoverCloseTrigger = withContext(
  styled(ArkPopover.CloseTrigger),
  "closeTrigger",
);
const PopoverContent = withContext(styled(ArkPopover.Content), "content");
const PopoverDescription = withContext(
  styled(ArkPopover.Description),
  "description",
);
// const PopoverIndicator = withContext(styled(ArkPopover.Indicator), "indicator");
const PopoverPositioner = withContext(
  styled(ArkPopover.Positioner),
  "positioner",
);
const PopoverTitle = withContext(styled(ArkPopover.Title), "title");
const PopoverTrigger = withContext(styled(ArkPopover.Trigger), "trigger");

export {
  Popover,
  PopoverAnchor,
  PopoverArrow,
  PopoverArrowTip,
  PopoverCloseTrigger,
  PopoverContent,
  PopoverDescription,
  // PopoverIndicator,
  PopoverPositioner,
  PopoverTitle,
  PopoverTrigger,
};
