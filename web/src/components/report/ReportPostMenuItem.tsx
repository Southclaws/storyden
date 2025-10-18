import { useDisclosure } from "src/utils/useDisclosure";

import { ProfileReference } from "@/api/openapi-schema";
import * as Menu from "@/components/ui/menu";

import { ReportIcon } from "../ui/icons/Report";

import { ReportPostModal } from "./ReportPostModal";

type Props = {
  targetKind: "thread" | "reply" | "post";
  targetId: string;
  author: ProfileReference;
  menuLabel: string;
  headline?: string;
  body?: string;
};

export function ReportPostMenuItem({
  targetKind,
  targetId,
  author,
  headline,
  body,
  menuLabel,
}: Props) {
  const disclosure = useDisclosure();

  return (
    <>
      <Menu.Item value={`report-${targetKind}`} onClick={disclosure.onOpen}>
        <ReportIcon />
        &nbsp;
        {menuLabel}
      </Menu.Item>

      <ReportPostModal
        targetKind={targetKind}
        targetId={targetId}
        author={author}
        headline={headline}
        body={body}
        {...disclosure}
      />
    </>
  );
}

export function truncateBody(value: string | undefined, maxLength = 280) {
  if (!value) return undefined;
  const trimmed = value.trim();
  if (trimmed.length <= maxLength) {
    return trimmed;
  }
  return `${trimmed.slice(0, maxLength).trim()}…`;
}
