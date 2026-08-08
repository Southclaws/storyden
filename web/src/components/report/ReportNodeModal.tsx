import { DatagraphItemKind, Node } from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { Text } from "@/components/ui/text";
import { VStack, styled } from "@/styled-system/jsx";

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
  return (
    <ReportModal
      title="Report page"
      description="Flag this page for moderator review. Use this if it contains incorrect, unsafe or inappropriate content."
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
            <Text
              variant="supporting"
              whiteSpace="pre-wrap"
              maxW="64"
              wordBreak="break-word"
            >
              {node.description}
            </Text>
          )}
        </VStack>
      }
      targetId={node.id}
      targetKind={DatagraphItemKind.node}
      submitLabel="Report page"
      successMessage="Thanks for the report. We'll review this page."
      loadingMessage="Sending report..."
      {...disclosure}
    />
  );
}
