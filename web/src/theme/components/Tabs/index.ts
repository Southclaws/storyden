import { Tabs as ArkTabs } from "@ark-ui/react/tabs";
import { styled } from "styled-system/jsx";
import { tabs } from "styled-system/recipes";

import { createStyleContext } from "src/theme/create-style-context";

const { withProvider, withContext } = createStyleContext(tabs);

const Tabs = withProvider(styled(ArkTabs.Root), "root");
const TabsContent = withContext(styled(ArkTabs.Content), "content");
const TabsIndicator = withContext(styled(ArkTabs.Indicator), "indicator");
const TabsList = withContext(styled(ArkTabs.List), "list");
const TabsTrigger = withContext(styled(ArkTabs.Trigger), "trigger");

const Root = Tabs;
const Content = TabsContent;
const Indicator = TabsIndicator;
const List = TabsList;
const Trigger = TabsTrigger;

export {
  Content,
  Indicator,
  List,
  Root,
  Tabs,
  TabsContent,
  TabsIndicator,
  TabsList,
  TabsTrigger,
  Trigger,
};
