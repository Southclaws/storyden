import { Cog6ToothIcon } from "@heroicons/react/24/outline";

import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const SettingsID = "settings";
export const SettingsRoute = "/settings";
export const SettingsLabel = "Settings";
export const SettingsIcon = <Cog6ToothIcon />;

type Props = AnchorProps & LinkButtonStyleProps;

export function SettingsAnchor(props: Props) {
  return (
    <Anchor
      id={SettingsID}
      route={SettingsRoute}
      label={SettingsLabel}
      icon={SettingsIcon}
      {...props}
    />
  );
}

export function SettingsMenuItem() {
  return (
    <MenuItem
      id={SettingsID}
      route={SettingsRoute}
      label={SettingsLabel}
      icon={SettingsIcon}
    />
  );
}
