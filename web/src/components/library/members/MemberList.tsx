import { PublicProfileList } from "src/api/openapi-schema";

import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { EmptyState } from "@/components/site/EmptyState";
import { Timestamp } from "@/components/site/Timestamp";
import * as Table from "@/components/ui/table";
import { useI18n } from "@/i18n/provider";
import { Box, HStack, LStack, VStack, styled } from "@/styled-system/jsx";

type Props = {
  profiles: PublicProfileList;
};

export function MemberList({ profiles }: Props) {
  const { t } = useI18n();

  if (profiles.length === 0) {
    return <EmptyState>{t("no members were found")}</EmptyState>;
  }

  return (
    <>
      {/* Desktop Table View */}
      <Box w="full" display={{ base: "none", lg: "block" }}>
        <Table.Root size="sm">
          <Table.Head>
            <Table.Row>
              <Table.Cell>{t("Member")}</Table.Cell>
              <Table.Cell>{t("Invited by")}</Table.Cell>
              <Table.Cell>{t("Likes")}</Table.Cell>
              <Table.Cell textAlign="right">{t("Joined")}</Table.Cell>
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
                    <LStack gap="1" alignItems="end">
                      <Timestamp created={profile.createdAt} large />
                      {isBanned && (
                        <styled.p color="fg.destructive">
                          {t("banned")}{" "}
                          <Timestamp created={profile.deletedAt!} large />
                        </styled.p>
                      )}
                    </LStack>
                  </Table.Cell>
                </Table.Row>
              );
            })}
          </Table.Body>
        </Table.Root>
      </Box>

      {/* Mobile Card View */}
      <VStack w="full" gap="3" display={{ base: "flex", lg: "none" }}>
        {profiles.map((profile) => {
          const isBanned = Boolean(profile.deletedAt);

          return (
            <Box
              w="full"
              key={profile.id}
              borderWidth="thin"
              borderRadius="lg"
              p="4"
              opacity={isBanned ? "5" : "full"}
              width="full"
            >
              <VStack gap="4" alignItems="stretch">
                <MemberBadge profile={profile} name="full-vertical" />

                <VStack alignItems="stretch">
                  <HStack justifyContent="space-between" gap="2">
                    <styled.span
                      color="fg.subtle"
                      fontSize="sm"
                      fontWeight="medium"
                    >
                      {t("Joined")}
                    </styled.span>
                    <Timestamp created={profile.createdAt} large />
                  </HStack>

                  <HStack justifyContent="space-between" gap="2">
                    <styled.span
                      color="fg.subtle"
                      fontSize="sm"
                      fontWeight="medium"
                    >
                      {t("Likes")}
                    </styled.span>
                    <styled.span fontSize="sm">
                      {profile.like_score}
                    </styled.span>
                  </HStack>

                  <HStack justifyContent="space-between" gap="2">
                    <styled.span
                      color="fg.subtle"
                      fontSize="sm"
                      fontWeight="medium"
                    >
                      {t("Invited by")}
                    </styled.span>
                    <Box>
                      {profile.invited_by ? (
                        <MemberBadge
                          profile={profile.invited_by}
                          name="handle"
                        />
                      ) : (
                        <styled.span
                          color="fg.subtle"
                          fontStyle="italic"
                          fontSize="sm"
                        >
                          n/a
                        </styled.span>
                      )}
                    </Box>
                  </HStack>

                  {isBanned && (
                    <Box
                      mt="1"
                      p="2"
                      bg="bg.error"
                      borderRadius="md"
                      borderWidth="thin"
                      borderColor="border.error"
                    >
                      <styled.p
                        color="fg.error"
                        fontSize="sm"
                        fontWeight="medium"
                      >
                        {t("Banned")} <Timestamp created={profile.deletedAt!} large />
                      </styled.p>
                    </Box>
                  )}
                </VStack>
              </VStack>
            </Box>
          );
        })}
      </VStack>
    </>
  );
}
