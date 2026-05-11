import { useDisclosure } from "src/utils/useDisclosure";

import { Node } from "@/api/openapi-schema";
import * as Menu from "@/components/ui/menu";
import { useI18n } from "@/i18n/provider";

import { ReportIcon } from "../ui/icons/Report";

import { ReportNodeModal } from "./ReportNodeModal";

type Props = {
  node: Node;
};

export function ReportNodeMenuItem({ node }: Props) {
  const disclosure = useDisclosure();
  const { t } = useI18n();

  return (
    <>
      <Menu.Item value="report-node" onClick={disclosure.onOpen}>
        <ReportIcon />
        &nbsp; {t("Report page")}
      </Menu.Item>

      <ReportNodeModal node={node} {...disclosure} />
    </>
  );
}
