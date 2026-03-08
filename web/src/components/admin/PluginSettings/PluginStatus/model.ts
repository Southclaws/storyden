export type PluginTab =
  | "overview"
  | "manifest"
  | "configuration"
  | "package"
  | "logs"
  | "connection";

export const DEFAULT_PLUGIN_TAB: PluginTab = "overview";

export function isPluginTab(tab: string | null, external: boolean): tab is PluginTab {
  if (!tab) {
    return false;
  }

  if (tab === "overview" || tab === "manifest" || tab === "configuration") {
    return true;
  }

  if (!external && tab === "package") {
    return true;
  }

  if (external && tab === "connection") {
    return true;
  }

  if (!external && tab === "logs") {
    return true;
  }

  return false;
}
