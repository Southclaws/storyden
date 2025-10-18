import { ReportIcon } from "@/components/ui/icons/Report";

import { Anchor, AnchorProps, MenuItem } from "./Anchor";

export const ReportsID = "reports";
export const ReportsRoute = "/reports";
export const ReportsLabel = "Reports";

export function ReportsAnchor(props: AnchorProps) {
  return (
    <Anchor
      id={ReportsID}
      route={ReportsRoute}
      label={ReportsLabel}
      icon={<ReportIcon />}
      {...props}
    />
  );
}

export function ReportsMenuItem() {
  return (
    <MenuItem
      id={ReportsID}
      route={ReportsRoute}
      label={ReportsLabel}
      icon={<ReportIcon />}
    />
  );
}
