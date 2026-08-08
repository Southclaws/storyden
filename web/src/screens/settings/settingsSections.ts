import { type Account, Permission } from "@/api/openapi-schema";
import { type Settings } from "@/lib/settings/settings";
import { hasPermission } from "@/utils/permissions";

export const DEFAULT_SETTINGS_SECTION = "interface";

export type SettingsSectionID =
  | "interface"
  | "authentication"
  | "email"
  | "access_keys"
  | "oauth";

export type SettingsSection = {
  id: SettingsSectionID;
  href: string;
  label: string;
  description: string;
};

export function getSettingsSections(
  settings: Settings,
  session?: Account,
): SettingsSection[] {
  const sections: SettingsSection[] = [
    {
      id: "interface",
      href: "/settings/interface",
      label: "Interface",
      description: "Theme and display preferences",
    },
    {
      id: "authentication",
      href: "/settings/authentication",
      label: "Authentication",
      description: "Password and sign-in methods",
    },
  ];

  if (settings.capabilities.includes("email_client")) {
    sections.push({
      id: "email",
      href: "/settings/email",
      label: "Email",
      description: "Email addresses and delivery",
    });
  }

  if (hasPermission(session, Permission.USE_PERSONAL_ACCESS_KEYS)) {
    sections.push({
      id: "access_keys",
      href: "/settings/access-keys",
      label: "Access keys",
      description: "Personal API credentials",
    });
  }

  if (
    settings.capabilities.includes("oauth") &&
    hasPermission(session, Permission.ADMINISTRATOR)
  ) {
    sections.push({
      id: "oauth",
      href: "/settings/oauth",
      label: "OAuth",
      description: "Connected application access",
    });
  }

  return sections;
}
