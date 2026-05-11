import { DatagraphItemKind, Node } from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { useI18n } from "@/i18n/provider";
import { styled, VStack } from "@/styled-system/jsx";

import { ReportModal, ReportModalProps } from "./ReportModal";

type Props = Omit<
  ReportModalProps,
  | "title"
  | "description"
  | "subject"
  | "targetId"
  | "targetKind"
  | "submitLabel"
  | "successMessage"
  | "loadingMessage"
> & {
  node: Node;
};

export function ReportNodeModal({ node, ...disclosure }: Props) {
  const { t } = useI18n();

  return (
    <ReportModal
      title={t("Report page")}
      description={t(
        "Flag this page for moderator review. Use this if it contains incorrect, unsafe or inappropriate content.",
      )}
      subject={
        <VStack alignItems="start" gap="2">
          <styled.span
            fontWeight="medium"
            maxW="64"
            whiteSpace="pre-wrap"
            wordBreak="break-word"
          >
            {node.name}
          </styled.span>
          <MemberBadge profile={node.owner} size="sm" name="full-horizontal" />
          {node.description && (
            <styled.p
              fontSize="sm"
              color="fg.subtle"
              whiteSpace="pre-wrap"
              maxW="64"
              wordBreak="break-word"
            >
              {node.description}
            </styled.p>
          )}
        </VStack>
      }
      targetId={node.id}
      targetKind={DatagraphItemKind.node}
      submitLabel={t("Report page")}
      successMessage={t("Thanks for the report. We'll review this page.")}
      loadingMessage={t("Sending report...")}
      {...disclosure}
    />
  );
}
