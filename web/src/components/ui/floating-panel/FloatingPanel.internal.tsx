"use client";

import type { Assign } from "@ark-ui/react";
import { FloatingPanel } from "@ark-ui/react/floating-panel";

import { floatingPanel } from "@/styled-system/recipes";
import type { JsxStyleProps } from "@/styled-system/types";
import { createStyleContext } from "@/utils/create-style-context";

const { withRootProvider, withContext } = createStyleContext(floatingPanel);

export interface RootProps extends FloatingPanel.RootProps {}
export interface RootProviderProps extends FloatingPanel.RootProviderProps {}

export const Root = withRootProvider<RootProps>(FloatingPanel.Root);
export const RootProvider = withRootProvider<RootProviderProps>(
  FloatingPanel.RootProvider,
);

export const Trigger = withContext<
  HTMLButtonElement,
  Assign<JsxStyleProps, FloatingPanel.TriggerProps>
>(FloatingPanel.Trigger, "trigger");

export const Positioner = withContext<
  HTMLDivElement,
  Assign<JsxStyleProps, FloatingPanel.PositionerProps>
>(FloatingPanel.Positioner, "positioner");

export const Content = withContext<
  HTMLDivElement,
  Assign<JsxStyleProps, FloatingPanel.ContentProps>
>(FloatingPanel.Content, "content");

export const DragTrigger = withContext<
  HTMLDivElement,
  Assign<JsxStyleProps, FloatingPanel.DragTriggerProps>
>(FloatingPanel.DragTrigger, "dragTrigger");

export const Header = withContext<
  HTMLDivElement,
  Assign<JsxStyleProps, FloatingPanel.HeaderProps>
>(FloatingPanel.Header, "header");

export const Title = withContext<
  HTMLHeadingElement,
  Assign<JsxStyleProps, FloatingPanel.TitleProps>
>(FloatingPanel.Title, "title");

export const Control = withContext<
  HTMLDivElement,
  Assign<JsxStyleProps, FloatingPanel.ControlProps>
>(FloatingPanel.Control, "control");

export const CloseTrigger = withContext<
  HTMLButtonElement,
  Assign<JsxStyleProps, FloatingPanel.CloseTriggerProps>
>(FloatingPanel.CloseTrigger, "closeTrigger");

export const Body = withContext<
  HTMLDivElement,
  Assign<JsxStyleProps, FloatingPanel.BodyProps>
>(FloatingPanel.Body, "body");

export const ResizeTrigger = withContext<
  HTMLDivElement,
  Assign<JsxStyleProps, FloatingPanel.ResizeTriggerProps>
>(FloatingPanel.ResizeTrigger, "resizeTrigger");

export const StageTrigger = withContext<
  HTMLButtonElement,
  Assign<JsxStyleProps, FloatingPanel.StageTriggerProps>
>(FloatingPanel.StageTrigger, "stageTrigger");

export { FloatingPanelContext as Context } from "@ark-ui/react/floating-panel";
