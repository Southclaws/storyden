import z from "zod";

import { Account } from "@/api/openapi-schema";

import { EditorSettingsSchema } from "./editor";
import { FrontendConfiguration, Settings } from "./settings";
import { SidebarSettingsSchema } from "./sidebar";

const MemberCustomSettingsParseSchema = z.object({
  editor: EditorSettingsSchema.optional(),
  sidebar: SidebarSettingsSchema.optional(),
});

export const MemberCustomSettingsSchema = z.object({
  editor: EditorSettingsSchema,
  sidebar: SidebarSettingsSchema,
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
  const parsed = MemberCustomSettingsParseSchema.safeParse(data.meta ?? {});

  const rawMeta = parsed.success ? parsed.data : {};
  if (!parsed.success && data.meta) {
    console.warn("Failed to parse member settings meta:", parsed.error);
  }

  const meta: MemberCustomSettings = {
    editor: {
      mode:
        rawMeta.editor?.mode ??
        global?.editor.mode ??
        DefaultMemberSettings.editor.mode,
    },
    sidebar: {
      defaultState:
        rawMeta.sidebar?.defaultState ??
        global?.sidebar.defaultState ??
        DefaultMemberSettings.sidebar.defaultState,
    },
  };

  const settings = { ...data, meta } satisfies Member;

  return settings;
}
