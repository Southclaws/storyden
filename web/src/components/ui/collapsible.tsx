"use client";

import { Collapsible } from "@ark-ui/react/collapsible";

import { createStyleContext } from "@/styled-system/jsx";
import { collapsible } from "@/styled-system/recipes";
import type { ComponentProps } from "@/styled-system/types";

const { withProvider, withContext } = createStyleContext(collapsible);

export type RootProviderProps = ComponentProps<typeof RootProvider>;
export const RootProvider = withProvider(Collapsible.RootProvider, "root");

export type RootProps = ComponentProps<typeof Root>;
export const Root = withProvider(Collapsible.Root, "root");

export const Content = withContext(Collapsible.Content, "content");

export const Trigger = withContext(Collapsible.Trigger, "trigger");

export { CollapsibleContext as Context } from "@ark-ui/react/collapsible";
