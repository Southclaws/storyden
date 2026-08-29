import { z } from "zod";

export const BuiltInNavigationItemTypeSchema = z.enum([
  "compose",
  "categories",
  "library",
  "robots",
  "collections",
  "links",
  "members",
  "roles",
]);
export type BuiltInNavigationItemType = z.infer<
  typeof BuiltInNavigationItemTypeSchema
>;

const BuiltInNavigationItemSchema = z.object({
  type: BuiltInNavigationItemTypeSchema,
});

const CustomNavigationLinkSchema = z.object({
  type: z.literal("custom-link"),
  id: z.string().trim().min(1).max(64),
  label: z.string().trim().min(1).max(80),
  href: z
    .string()
    .trim()
    .min(1)
    .max(2048)
    .refine(isSafeNavigationHref, "Enter a safe internal or HTTP(S) URL."),
});

export const NavigationItemSchema = z.discriminatedUnion("type", [
  BuiltInNavigationItemSchema,
  CustomNavigationLinkSchema,
]);
export type NavigationItem = z.infer<typeof NavigationItemSchema>;
export type CustomNavigationLink = Extract<
  NavigationItem,
  { type: "custom-link" }
>;

const NavigationConfigValueSchema = z
  .object({
    items: z.array(NavigationItemSchema),
  })
  .superRefine(({ items }, context) => {
    const seen = new Set<string>();

    items.forEach((item, index) => {
      const key = getNavigationItemKey(item);
      if (seen.has(key)) {
        context.addIssue({
          code: z.ZodIssueCode.custom,
          message: `Navigation item "${key}" may only appear once.`,
          path: ["items", index],
        });
      }
      seen.add(key);
    });
  });
export type NavigationConfig = z.infer<typeof NavigationConfigValueSchema>;

export const DefaultNavigationConfig = {
  items: [
    { type: "compose" },
    { type: "categories" },
    { type: "library" },
    { type: "collections" },
    { type: "links" },
    { type: "members" },
    { type: "roles" },
  ],
} satisfies NavigationConfig;

export const NavigationConfigSchema = NavigationConfigValueSchema.default(
  DefaultNavigationConfig,
);

export const BuiltInNavigationItemName: Record<
  BuiltInNavigationItemType,
  string
> = {
  compose: "New thread",
  categories: "Discussion",
  library: "Library",
  robots: "Robots",
  collections: "Collections",
  links: "Links",
  members: "Members",
  roles: "Roles",
};

export function getNavigationItemKey(item: NavigationItem): string {
  return item.type === "custom-link" ? `custom-link:${item.id}` : item.type;
}

export function addBuiltInNavigationItem(
  navigation: NavigationConfig,
  type: BuiltInNavigationItemType,
  afterIndex?: number,
): NavigationConfig {
  if (navigation.items.some((item) => item.type === type)) {
    return navigation;
  }

  return insertNavigationItem(navigation, { type }, afterIndex);
}

export function addCustomNavigationLink(
  navigation: NavigationConfig,
  link: Omit<CustomNavigationLink, "type">,
  afterIndex?: number,
): NavigationConfig {
  const item = NavigationItemSchema.parse({ type: "custom-link", ...link });
  if (
    navigation.items.some(
      (current) => getNavigationItemKey(current) === getNavigationItemKey(item),
    )
  ) {
    return navigation;
  }

  return insertNavigationItem(navigation, item, afterIndex);
}

export function reorderNavigationItem(
  navigation: NavigationConfig,
  activeKey: string,
  overKey: string,
): NavigationConfig {
  const activeIndex = navigation.items.findIndex(
    (item) => getNavigationItemKey(item) === activeKey,
  );
  const overIndex = navigation.items.findIndex(
    (item) => getNavigationItemKey(item) === overKey,
  );
  if (activeIndex === -1 || overIndex === -1 || activeIndex === overIndex) {
    return navigation;
  }

  const items = [...navigation.items];
  const [moved] = items.splice(activeIndex, 1);
  if (!moved) {
    return navigation;
  }

  items.splice(overIndex, 0, moved);
  return { items };
}

export function replaceCustomNavigationLink(
  navigation: NavigationConfig,
  updated: CustomNavigationLink,
): NavigationConfig {
  const parsed = CustomNavigationLinkSchema.parse(updated);
  const key = getNavigationItemKey(parsed);
  const index = navigation.items.findIndex(
    (item) => getNavigationItemKey(item) === key,
  );
  if (index === -1) {
    return navigation;
  }

  const items = [...navigation.items];
  items[index] = parsed;
  return { items };
}

export function removeNavigationItem(
  navigation: NavigationConfig,
  key: string,
): NavigationConfig {
  const items = navigation.items.filter(
    (item) => getNavigationItemKey(item) !== key,
  );
  return items.length === navigation.items.length ? navigation : { items };
}

export function isSafeNavigationHref(value: string): boolean {
  const href = value.trim();
  if (!href) return false;

  if (
    (href.startsWith("/") && !href.startsWith("//")) ||
    href.startsWith("?") ||
    href.startsWith("#")
  ) {
    return true;
  }

  try {
    const url = new URL(href);
    return url.protocol === "http:" || url.protocol === "https:";
  } catch {
    return false;
  }
}

export function normaliseNavigationHref(value: string): string | undefined {
  const href = value.trim();
  if (isSafeNavigationHref(href)) {
    return href;
  }

  if (!href || /\s/.test(href) || !href.includes(".")) {
    return undefined;
  }

  try {
    const url = new URL(`https://${href}`);
    return isSafeNavigationHref(url.toString()) ? url.toString() : undefined;
  } catch {
    return undefined;
  }
}

function insertNavigationItem(
  navigation: NavigationConfig,
  item: NavigationItem,
  afterIndex?: number,
): NavigationConfig {
  const items = [...navigation.items];
  const insertionIndex =
    afterIndex === undefined
      ? items.length
      : Math.min(Math.max(afterIndex + 1, 0), items.length);

  items.splice(insertionIndex, 0, item);
  return { items };
}
