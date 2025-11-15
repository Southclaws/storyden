import z from "zod";

import { Account } from "@/api/openapi-schema";

import { EditorSettingsSchema } from "./editor";
import { FrontendConfiguration, Settings } from "./settings";

export const MemberCustomSettingsSchema = z.object({
  editor: EditorSettingsSchema.default({}),
});
export type MemberCustomSettings = z.infer<typeof MemberCustomSettingsSchema>;

// Member extends Account with custom member settings by typing `meta`.
export type Member = Account & {
  meta: MemberCustomSettings;
};

export const DefaultMemberSettings: MemberCustomSettings = {
  editor: {
    mode: "richtext",
  },
};

export function parseMemberSettings(
  global: FrontendConfiguration,
  data: Account,
): Member {
  const parsed = MemberCustomSettingsSchema.safeParse(data.meta);

  const meta = parsed.success ? parsed.data : DefaultMemberSettings;

  if (meta.editor.mode === undefined) {
    meta.editor.mode = global.editor.mode;
  }

  const settings = { ...data, meta } satisfies Member;

  return settings;
}
