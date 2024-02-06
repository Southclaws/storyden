import { Select as ArkSelect } from "@ark-ui/react/select";
import { styled } from "styled-system/jsx";

import { createStyleContext } from "src/theme/create-style-context";

import { select } from "@/styled-system/recipes";

const { withProvider, withContext } = createStyleContext(select);

const Select = withProvider(styled(ArkSelect.Root), "root");
const SelectClearTrigger = withContext(
  styled(ArkSelect.ClearTrigger),
  "clearTrigger",
);
const SelectContent = withContext(styled(ArkSelect.Content), "content");
const SelectControl = withContext(styled(ArkSelect.Control), "control");
const SelectIndicator = withContext(
  styled(ArkSelect.ItemIndicator),
  "indicator",
);
const SelectItem = withContext(styled(ArkSelect.Item), "item");
const SelectItemGroup = withContext(styled(ArkSelect.ItemGroup), "itemGroup");
const SelectItemGroupLabel = withContext(
  styled(ArkSelect.ItemGroupLabel),
  "itemGroupLabel",
);
const SelectItemIndicator = withContext(
  styled(ArkSelect.ItemIndicator),
  "itemIndicator",
);
const SelectItemText = withContext(styled(ArkSelect.ItemText), "itemText");
const SelectLabel = withContext(styled(ArkSelect.Label), "label");
const SelectPositioner = withContext(
  styled(ArkSelect.Positioner),
  "positioner",
);
const SelectTrigger = withContext(styled(ArkSelect.Trigger), "trigger");
const SelectValueText = withContext(styled(ArkSelect.Value), "valueText");

const Root = Select;
const ClearTrigger = SelectClearTrigger;
const Content = SelectContent;
const Control = SelectControl;
const Indicator = SelectIndicator;
const Item = SelectItem;
const ItemGroup = SelectItemGroup;
const ItemGroupLabel = SelectItemGroupLabel;
const ItemIndicator = SelectItemIndicator;
const ItemText = SelectItemText;
const Label = SelectLabel;
const Positioner = SelectPositioner;
const Trigger = SelectTrigger;
const ValueText = SelectValueText;

export {
  ClearTrigger,
  Content,
  Control,
  Indicator,
  Item,
  ItemGroup,
  ItemGroupLabel,
  ItemIndicator,
  ItemText,
  Label,
  Positioner,
  Root,
  Select,
  SelectClearTrigger,
  SelectContent,
  SelectControl,
  SelectIndicator,
  SelectItem,
  SelectItemGroup,
  SelectItemGroupLabel,
  SelectItemIndicator,
  SelectItemText,
  SelectLabel,
  SelectPositioner,
  SelectTrigger,
  SelectValueText,
  Trigger,
  ValueText,
};
