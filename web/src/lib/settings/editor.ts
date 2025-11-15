import z from "zod";

export const EditorModeSchema = z.enum(["richtext", "markdown"]);
export type EditorMode = z.infer<typeof EditorModeSchema>;

export const EditorSettingsSchema = z.object({
  mode: EditorModeSchema.default("richtext"),
});
export type EditorSettings = z.infer<typeof EditorSettingsSchema>;
