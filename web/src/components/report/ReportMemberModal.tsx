import { DatagraphItemKind, ProfileReference } from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";

import { ReportModal, ReportModalProps } from "./ReportModal";

export type ReportMemberModalProps = Omit<ReportModalProps, "title" | "description" | "subject" | "targetId" | "targetKind" | "submitLabel" | "successMessage" | "loadingMessage"> & {
  profile: ProfileReference;
};

export function ReportMemberModal({ profile, ...disclosure }: ReportMemberModalProps) {
  return (
    <ReportModal
      title={`Report ${profile.name}`}
      description={
        "Tell us what's happening with this member. Reports help moderators keep the community safe."
      }
      subject={<MemberBadge profile={profile} name="full-vertical" />}
      targetId={profile.id}
      targetKind={DatagraphItemKind.profile}
      submitLabel="Report member"
      successMessage="Thanks for the report. Our moderators will review this member."
      loadingMessage="Sending member report..."
      {...disclosure}
    />
  );
}
