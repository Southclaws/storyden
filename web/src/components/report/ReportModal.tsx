import dynamic from "next/dynamic";
import { ReactNode, useState } from "react";

import { DatagraphItemKind } from "@/api/openapi-schema";
import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { UnreadyBanner } from "@/components/site/Unready";
import { UseDisclosureProps } from "@/utils/useDisclosure";

const ReportForm = dynamic(
  () => import("./ReportForm").then((module) => module.ReportForm),
  { ssr: false, loading: () => <UnreadyBanner /> },
);

export type ReportModalProps = UseDisclosureProps & {
  title: string;
  description: ReactNode;
  subject?: ReactNode;
  targetId: string;
  targetKind: DatagraphItemKind;
  submitLabel?: string;
  successMessage?: string;
  loadingMessage?: string;
};

export function ReportModal({ title, ...props }: ReportModalProps) {
  // Keep an unfinished report when dismissing and reopening the drawer.
  const [comment, setComment] = useState("");

  return (
    <ModalDrawer
      title={title}
      isOpen={props.isOpen}
      onClose={props.onClose}
      onOpenChange={props.onOpenChange}
    >
      <ReportForm
        {...props}
        initialComment={comment}
        onCommentChange={setComment}
      />
    </ModalDrawer>
  );
}
