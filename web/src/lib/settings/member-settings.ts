import z from "zod";

import { Account } from "@/api/openapi-schema";

import { EditorSettingsSchema } from "./editor";
import { FrontendConfiguration, Settings } from "./settings";
import { SidebarSettingsSchema } from "./sidebar";

export const MemberCustomSettingsSchema = z.object({
  editor: EditorSettingsSchema.default({}),
  sidebar: SidebarSettingsSchema.default({}),
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
  sidebar: {
    defaultState: "closed",
  },
};

export function parseMemberSettings(
  data: Account,
  global?: FrontendConfiguration,
): Member {
  const parsed = MemberCustomSettingsSchema.safeParse(data.meta);

  const meta = parsed.success ? parsed.data : DefaultMemberSettings;

  if (meta.editor.mode === undefined) {
    meta.editor.mode = global?.editor.mode ?? DefaultMemberSettings.editor.mode;
  }

  if (meta.sidebar.defaultState === undefined) {
    meta.sidebar.defaultState =
      global?.sidebar.defaultState ?? DefaultMemberSettings.sidebar.defaultState;
  }

  const settings = { ...data, meta } satisfies Member;

  return settings;
}
