import { useDisclosure } from "src/utils/useDisclosure";

import { ProfileReference } from "@/api/openapi-schema";
import * as Menu from "@/components/ui/menu";

import { ReportIcon } from "../ui/icons/Report";

import { ReportMemberModal } from "./ReportMemberModal";

type Props = {
  profile: ProfileReference;
};

export function ReportMemberMenuItem({ profile }: Props) {
  const disclosure = useDisclosure();

  return (
    <>
      <Menu.Item value="report-member" onClick={disclosure.onOpen}>
        <ReportIcon />
        &nbsp; Report member
      </Menu.Item>
      <ReportMemberModal profile={profile} {...disclosure} />
    </>
  );
}
