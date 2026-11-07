import { Account } from "@/api/openapi-schema";
import * as z from "@/lib/validation";

import { EditorSettingsSchema } from "./editor";
import { FrontendConfiguration } from "./settings";

const NotificationSettingsSchema = z.object({
  markReadOnClick: z.boolean(),
});

const MemberCustomSettingsParseSchema = z.object({
  editor: z.optional(EditorSettingsSchema),
  notifications: z.optional(z.partial(NotificationSettingsSchema)),
});

export const MemberCustomSettingsSchema = z.object({
  editor: EditorSettingsSchema,
  notifications: NotificationSettingsSchema,
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
  notifications: {
    markReadOnClick: true,
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
    ...data.meta,
    editor: {
      mode:
        rawMeta.editor?.mode ??
        global?.editor.mode ??
        DefaultMemberSettings.editor.mode,
    },
    notifications: {
      markReadOnClick:
        rawMeta.notifications?.markReadOnClick ??
        DefaultMemberSettings.notifications.markReadOnClick,
    },
  };

  const settings = { ...data, meta } satisfies Member;

  return settings;
}
