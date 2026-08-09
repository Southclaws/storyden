import { PublicProfileList } from "@/api/openapi-schema";
import { MemberBadge } from "@/components/member/MemberBadge/MemberBadge";
import { EmptyState } from "@/components/site/EmptyState";
import { Timestamp } from "@/components/site/Timestamp";
import * as Table from "@/components/ui/table";
import { Text } from "@/components/ui/text";
import { Box, HStack, LStack, VStack } from "@/styled-system/jsx";

type Props = {
  profiles: PublicProfileList;
};

export function MemberList({ profiles }: Props) {
  if (profiles.length === 0) {
    return <EmptyState>no members were found</EmptyState>;
  }

  return (
    <>
      {/* Desktop Table View */}
      <Box w="full" display={{ base: "none", lg: "block" }}>
        <Table.Root size="sm">
          <Table.Head>
            <Table.Row>
              <Table.Cell>Member</Table.Cell>
              <Table.Cell>Invited by</Table.Cell>
              <Table.Cell>Likes</Table.Cell>
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
                      <Text variant="supporting" fontStyle="italic">
                        n/a
                      </Text>
                    )}
                  </Table.Cell>

                  <Table.Cell>{profile.like_score}</Table.Cell>

                  <Table.Cell>
                    <LStack gap="1" alignItems="end">
                      <Timestamp created={profile.createdAt} large />
                      {isBanned && (
                        <Text variant="metadata" color="status.danger.content">
                          banned{" "}
                          <Timestamp created={profile.deletedAt!} large />
                        </Text>
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
                    <Text as="span" variant="supporting" fontWeight="medium">
                      Joined
                    </Text>
                    <Timestamp created={profile.createdAt} large />
                  </HStack>

                  <HStack justifyContent="space-between" gap="2">
                    <Text as="span" variant="supporting" fontWeight="medium">
                      Likes
                    </Text>
                    <Text as="span" variant="supporting" color="text.default">
                      {profile.like_score}
                    </Text>
                  </HStack>

                  <HStack justifyContent="space-between" gap="2">
                    <Text as="span" variant="supporting" fontWeight="medium">
                      Invited by
                    </Text>
                    <Box>
                      {profile.invited_by ? (
                        <MemberBadge
                          profile={profile.invited_by}
                          name="handle"
                        />
                      ) : (
                        <Text as="span" variant="supporting" fontStyle="italic">
                          n/a
                        </Text>
                      )}
                    </Box>
                  </HStack>

                  {isBanned && (
                    <Box
                      mt="1"
                      p="2"
                      bg="status.danger.surface"
                      borderRadius="md"
                      borderWidth="thin"
                      borderColor="status.danger.border"
                    >
                      <Text
                        variant="metadata"
                        color="status.danger.content"
                        fontWeight="medium"
                      >
                        Banned <Timestamp created={profile.deletedAt!} large />
                      </Text>
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
