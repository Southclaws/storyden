import { SettingsIcon } from "@/components/ui/icons/Settings";
import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const SettingsID = "settings";
export const SettingsRoute = "/settings";
export const SettingsLabel = "Settings";

type Props = AnchorProps & LinkButtonStyleProps;

export function SettingsAnchor(props: Props) {
  return (
    <Anchor
      id={SettingsID}
      route={SettingsRoute}
      label={SettingsLabel}
      icon={<SettingsIcon />}
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
      icon={<SettingsIcon />}
    />
  );
}
