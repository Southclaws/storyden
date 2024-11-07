import { PublicProfileList } from "src/api/openapi-schema";
import { Empty } from "src/components/site/Empty";

import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { RoleBadgeList } from "@/components/role/RoleBadge/RoleBadgeList";
import { EmptyState } from "@/components/site/EmptyState";
import { Timestamp } from "@/components/site/Timestamp";
import * as Table from "@/components/ui/table";
import { LStack, styled } from "@/styled-system/jsx";

type Props = {
  profiles: PublicProfileList;
};

export function MemberList({ profiles }: Props) {
  if (profiles.length === 0) {
    return <EmptyState>no members were found</EmptyState>;
  }

  return (
    <Table.Root size="sm">
      <Table.Head>
        <Table.Row>
          <Table.Cell>Member</Table.Cell>
          <Table.Cell>Invited by</Table.Cell>
          <Table.Cell>Likes</Table.Cell>
          <Table.Cell>Roles</Table.Cell>
          <Table.Cell textAlign="right">Joined</Table.Cell>
        </Table.Row>
      </Table.Head>

      <Table.Body>
        {profiles.map((profile) => {
          const isBanned = Boolean(profile.deletedAt);

          return (
            <Table.Row key={profile.id}>
              <Table.Cell py="1" opacity={isBanned ? "5" : "full"}>
                <MemberBadge profile={profile} name="full-vertical" />
              </Table.Cell>

              <Table.Cell>
                {profile.invited_by ? (
                  <MemberBadge profile={profile.invited_by} name="handle" />
                ) : (
                  <styled.p color="fg.subtle" fontStyle="italic">
                    n/a
                  </styled.p>
                )}
              </Table.Cell>

              <Table.Cell>{profile.like_score}</Table.Cell>

              <Table.Cell>
                <RoleBadgeList roles={profile.roles} limit={1} />
              </Table.Cell>

              <Table.Cell>
                <LStack gap="1" alignItems="end">
                  <Timestamp created={profile.createdAt} large />
                  {isBanned && (
                    <styled.p color="tomato.10">
                      banned <Timestamp created={profile.deletedAt!} large />
                    </styled.p>
                  )}
                </LStack>
              </Table.Cell>
            </Table.Row>
          );
        })}
      </Table.Body>
    </Table.Root>
  );
}
