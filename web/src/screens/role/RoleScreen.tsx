"use client";

import { useRoleList } from "@/api/openapi-client/roles";
import { Account, Permission, RoleListOKResponse } from "@/api/openapi-schema";
import { RoleCard } from "@/components/role/RoleCard";
import { RoleCreateModalTrigger } from "@/components/role/RoleEdit/RoleEditModal";
import { InfoTip } from "@/components/site/InfoTip";
import { UnreadyBanner } from "@/components/site/Unready";
import { Heading } from "@/components/ui/heading";
import { CardGrid } from "@/components/ui/rich-card";
import { HStack, LStack, WStack } from "@/styled-system/jsx";
import { hasPermission } from "@/utils/permissions";

type Props = {
  session?: Account;
  initialRoles: RoleListOKResponse;
};

export function RoleScreen(props: Props) {
  const { data, error } = useRoleList({
    swr: { fallbackData: props.initialRoles },
  });
  if (!data) {
    return <UnreadyBanner error={error} />;
  }

  const canEdit = hasPermission(props.session, Permission.MANAGE_ROLES);

  return (
    <LStack>
      <WStack>
        <Heading>Roles</Heading>

        {canEdit && <RoleCreateModalTrigger />}
      </WStack>

      <HStack gap="1">
        <p>
          Roles provide granular permission control and profile customisation
          for members.
        </p>
        <InfoTip title="Aesthetic roles and badges">
          You can also use Roles as a purely aesthetic tool for providing
          members with ways to express themselves on their profile. Members can
          choose one role as a &ldquo;Badge&rdquo; which is displayed next to
          their name around the site.
        </InfoTip>
      </HStack>

      <CardGrid>
        {data.roles.map((r) => (
          <RoleCard key={r.id} role={r} editable={canEdit} />
        ))}
      </CardGrid>
    </LStack>
  );
}
