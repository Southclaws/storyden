import {
  type Account,
  InstanceCapability,
  Permission,
} from "@/api/openapi-schema";
import { hasCapability } from "@/lib/settings/capabilities";
import { type Settings } from "@/lib/settings/settings";
import { hasPermission } from "@/utils/permissions";

export const DEFAULT_ADMIN_SECTION = "brand";

export type AdminSectionID =
  | "brand"
  | "interface"
  | "moderation"
  | "system"
  | "audit_log"
  | "email"
  | "authentication"
  | "access_keys"
  | "oauth"
  | "plugins"
  | "robots";

export type AdminSection = {
  id: AdminSectionID;
  href: string;
  label: string;
  description: string;
};

export type AdminSectionGroup = {
  label: string;
  items: AdminSection[];
};

export function getAdminSectionGroups(
  settings: Settings,
  session?: Account,
): AdminSectionGroup[] {
  const canViewAdministratorSections = hasPermission(
    session,
    Permission.ADMINISTRATOR,
  );
  const oauthEnabled =
    hasCapability(InstanceCapability.oauth, settings.capabilities) &&
    canViewAdministratorSections;
  const pluginsEnabled = hasCapability(
    InstanceCapability.plugins,
    settings.capabilities,
  );

  return [
    {
      label: "Appearance",
      items: [
        {
          id: "brand",
          href: "/admin/brand",
          label: "Brand",
          description: "Site identity and presentation",
        },
        {
          id: "interface",
          href: "/admin/interface",
          label: "Interface",
          description: "Navigation and content display",
        },
      ],
    },
    {
      label: "Community",
      items: [
        {
          id: "moderation",
          href: "/admin/moderation",
          label: "Moderation",
          description: "Community rules and restrictions",
        },
      ],
    },
    {
      label: "Security",
      items: [
        {
          id: "authentication",
          href: "/admin/authentication",
          label: "Authentication",
          description: "Sign-in and registration",
        },
        {
          id: "access_keys",
          href: "/admin/access-keys",
          label: "Access keys",
          description: "Service API credentials",
        },
        ...(oauthEnabled
          ? [
              {
                id: "oauth" as const,
                href: "/admin/oauth",
                label: "OAuth",
                description: "OAuth clients and authorisations",
              },
            ]
          : []),
      ],
    },
    {
      label: "Operations",
      items: [
        {
          id: "system",
          href: "/admin/system",
          label: "System",
          description: "Infrastructure and runtime settings",
        },
        {
          id: "audit_log",
          href: "/admin/audit-log",
          label: "Audit log",
          description: "Administrative and moderation events",
        },
        ...(canViewAdministratorSections
          ? [
              {
                id: "email" as const,
                href: "/admin/email",
                label: "Email",
                description: "Delivery queue and failures",
              },
            ]
          : []),
      ],
    },
    {
      label: "Extensions",
      items: [
        ...(pluginsEnabled
          ? [
              {
                id: "plugins" as const,
                href: "/admin/plugins",
                label: "Plugins",
                description: "Installed extensions and integrations",
              },
            ]
          : []),
        {
          id: "robots",
          href: "/admin/robots",
          label: "Robots",
          description: "Models, providers, and workspaces",
        },
      ],
    },
  ];
}
