import { z } from "zod";

export const SidebarDefaultStateSchema = z.enum(["open", "closed"]);
export type SidebarDefaultState = z.infer<typeof SidebarDefaultStateSchema>;

export const SidebarSettingsSchema = z.object({
  defaultState: SidebarDefaultStateSchema.optional(),
});
export type SidebarSettings = z.infer<typeof SidebarSettingsSchema>;

export const DefaultSidebarSettings: SidebarSettings = {
  defaultState: "closed",
};
