"use client";

import { formatDistanceToNow } from "date-fns";

import { Unready } from "src/components/site/Unready";

import { ContentFormField } from "@/components/content/ContentComposer/ContentField";
import { MemberAvatar } from "@/components/member/MemberBadge/MemberAvatar";
import { MemberIdent } from "@/components/member/MemberBadge/MemberIdent";
import { MemberOptionsMenu } from "@/components/member/MemberOptions/MemberOptionsMenu";
import { ProfileAccountManagement } from "@/components/profile/ProfileAccountManagement/ProfileAccountManagement";
import { ProfileContent } from "@/components/profile/ProfileContent/ProfileContent";
import { ProfileSuspendedBanner } from "@/components/profile/ProfileSuspendedBanner";
import { RoleBadgeList } from "@/components/role/RoleBadge/RoleBadgeList";
import { EditAction } from "@/components/site/Action/Edit";
import { MoreAction } from "@/components/site/Action/More";
import { SaveAction } from "@/components/site/Action/Save";
import { DotSeparator } from "@/components/site/Dot";
import { LikeIcon } from "@/components/ui/icons/Like";
import { Input } from "@/components/ui/input";
import {
  Box,
  CardBox,
  Flex,
  HStack,
  LStack,
  styled,
} from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

import { Form, Props, useProfileScreen } from "./useProfileScreen";

export function ProfileScreen(props: Props) {
  const { ready, error, form, state, data, handlers } = useProfileScreen(props);

  if (!ready) {
    return <Unready error={error} />;
  }

  const { session, profile } = data;
  const { isSelf, isEditing, canViewAccount } = state;
  const isEmpty =
    !profile.bio || profile.bio === "" || profile.bio === "<body></body>";

  return (
    <LStack w="full">
      <CardBox p="0">
        <styled.form className={lstack()} p="3" onSubmit={handlers.handleSave}>
          <Flex
            direction={{ base: "column-reverse", sm: "row" }}
            w="full"
            justify="space-between"
            alignItems={{ base: "end", sm: "start" }}
          >
            {isEditing ? (
              <HStack w="full" pr={{ base: "0", sm: "24" }}>
                <MemberAvatar
                  profile={profile}
                  size="lg"
                  editable={isEditing}
                />
                <LStack gap="1">
                  <LStack gap="0">
                    <Input
                      maxW={{ base: "full", sm: "64" }}
                      size="sm"
                      height="7"
                      px="2"
                      borderBottomRadius="none"
                      fontWeight="bold"
                      {...form.register("name")}
                    />
                    <Input
                      maxW={{ base: "full", sm: "64" }}
                      size="sm"
                      height="7"
                      px="2"
                      borderTop="none"
                      borderTopRadius="none"
                      {...form.register("handle")}
                    />
                  </LStack>
                  <RoleBadgeList roles={profile.roles} />
                </LStack>
              </HStack>
            ) : (
              <MemberIdent
                profile={profile}
                size="lg"
                name="full-vertical"
                showRoles="all"
              />
            )}

            <HStack justify="end">
              {isSelf &&
                (isEditing ? (
                  <SaveAction size="sm">Save</SaveAction>
                ) : (
                  <EditAction
                    size="sm"
                    variant="ghost"
                    onClick={handlers.handleSetEditing}
                  >
                    Edit
                  </EditAction>
                ))}
              <MemberOptionsMenu profile={profile} asChild>
                <MoreAction type="button" size="sm" />
              </MemberOptionsMenu>
            </HStack>
          </Flex>

          <HStack gap="1">
            <styled.p color="fg.muted" wordBreak="keep-all">
              Joined{" "}
              <styled.time textWrap="nowrap">
                {formatDistanceToNow(new Date(profile.createdAt), {
                  addSuffix: true,
                })}
              </styled.time>
            </styled.p>
            <DotSeparator />
            <HStack
              gap="1"
              color="fg.subtle"
              wordBreak="keep-all"
              textWrap="nowrap"
            >
              <Box flexShrink="0">
                <LikeIcon w="4" />
              </Box>
              <span>{profile.like_score} likes</span>
            </HStack>
          </HStack>

          {isEmpty && !isEditing ? (
            <styled.p color="fg.subtle" fontStyle="italic">
              This profile has no bio yet...
            </styled.p>
          ) : (
            <ContentFormField<Form>
              control={form.control}
              name="bio"
              initialValue={profile.bio}
              disabled={!isEditing}
              placeholder="This profile has no bio yet..."
            />
          )}
        </styled.form>

        {profile.deletedAt && (
          <Box p="3">
            <ProfileSuspendedBanner date={new Date(profile.deletedAt)} />
          </Box>
        )}

        {canViewAccount && (
          <Box p="3">
            <ProfileAccountManagement accountId={profile.id} />
          </Box>
        )}
      </CardBox>

      <ProfileContent session={session} profile={profile} />
    </LStack>
  );
}
