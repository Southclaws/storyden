import { LogoutIcon as LogoutGraphic } from "src/components/graphics/LogoutIcon";

import { LinkButtonStyleProps } from "@/components/ui/link-button";
import { Item } from "@/components/ui/menu";
import { button } from "@/styled-system/recipes";

import { AnchorProps, MenuItem } from "./Anchor";

export const LogoutID = "logout";
export const LogoutRoute = "/logout";
export const LogoutLabel = "Logout";
export const LogoutIcon = <LogoutGraphic />;

type Props = AnchorProps & LinkButtonStyleProps;

export function LogoutAnchor({ hideLabel, ...props }: Props) {
  return (
    <a
      className={button({ variant: "ghost", ...props })}
      href={LogoutRoute}
      title={LogoutLabel}
    >
      {LogoutIcon}
      {!hideLabel && (
        <>
          &nbsp;<span>{LogoutLabel}</span>
        </>
      )}
    </a>
  );
}

export function LogoutMenuItem({ hideLabel }: AnchorProps) {
  return (
    <a href={LogoutRoute}>
      <Item value={LogoutID}>
        {LogoutIcon}
        {!hideLabel && (
          <>
            &nbsp;<span>{LogoutLabel}</span>
          </>
        )}
      </Item>
    </a>
  );
}
