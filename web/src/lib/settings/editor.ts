import * as z from "@/lib/validation";

export const EditorModeSchema = z.enum(["richtext", "markdown"]);
export type EditorMode = z.infer<typeof EditorModeSchema>;

export const EditorSettingsSchema = z.object({
  mode: z.optional(EditorModeSchema),
});
export type EditorSettings = z.infer<typeof EditorSettingsSchema>;
