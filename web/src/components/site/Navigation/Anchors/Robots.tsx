import { RobotIcon } from "@/components/ui/icons/Robot";
import { LinkButtonStyleProps } from "@/components/ui/link-button";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const RobotsID = "robots";
export const RobotsRoute = "/robots";
export const RobotsLabel = "Robots";

type RobotsAnchorProps = AnchorProps &
  LinkButtonStyleProps & {
    route?: string;
  };

export function RobotsAnchor({
  route = RobotsRoute,
  ...props
}: RobotsAnchorProps) {
  return (
    <Anchor
      id={RobotsID}
      route={route}
      label={RobotsLabel}
      icon={<RobotIcon />}
      {...props}
    />
  );
}

export function RobotsMenuItem() {
  return (
    <MenuItem
      id={RobotsID}
      route={RobotsRoute}
      label={RobotsLabel}
      icon={<RobotIcon />}
    />
  );
}
