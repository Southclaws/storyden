"use client";

import { InvitationIcon } from "@/components/ui/icons/Invitation";
import { Item } from "@/components/ui/menu";

export const InvitationID = "invitation";
export const InvitationLabel = "Invite";

export function InvitationMenuItem() {
  return (
    <Item value={InvitationID}>
      <InvitationIcon />
      <span>{InvitationLabel}</span>
    </Item>
  );
}
