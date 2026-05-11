import { DatagraphItemKind, ProfileReference } from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { useI18n } from "@/i18n/provider";

import { ReportModal, ReportModalProps } from "./ReportModal";

export type ReportMemberModalProps = Omit<ReportModalProps, "title" | "description" | "subject" | "targetId" | "targetKind" | "submitLabel" | "successMessage" | "loadingMessage"> & {
  profile: ProfileReference;
};

export function ReportMemberModal({ profile, ...disclosure }: ReportMemberModalProps) {
  const { t } = useI18n();

  return (
    <ReportModal
      title={t("Report {{name}}", { name: profile.name })}
      description={
        t(
          "Tell us what's happening with this member. Reports help moderators keep the community safe.",
        )
      }
      subject={<MemberBadge profile={profile} name="full-vertical" />}
      targetId={profile.id}
      targetKind={DatagraphItemKind.profile}
      submitLabel={t("Report member")}
      successMessage={t(
        "Thanks for the report. Our moderators will review this member.",
      )}
      loadingMessage={t("Sending member report...")}
      {...disclosure}
    />
  );
}
