import { useDisclosure } from "src/utils/useDisclosure";

import { Node } from "@/api/openapi-schema";
import * as Menu from "@/components/ui/menu";

import { ReportIcon } from "../ui/icons/Report";

import { ReportNodeModal } from "./ReportNodeModal";

type Props = {
  node: Node;
};

export function ReportNodeMenuItem({ node }: Props) {
  const disclosure = useDisclosure();

  return (
    <>
      <Menu.Item value="report-node" onClick={disclosure.onOpen}>
        <ReportIcon />
        &nbsp; Report page
      </Menu.Item>

      <ReportNodeModal node={node} {...disclosure} />
    </>
  );
}
