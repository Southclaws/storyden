import { invitationGet } from "@/api/openapi-server/invitations";
import { getServerSession } from "@/auth/server-session";
import { MemberIdent } from "@/components/member/MemberBadge/MemberIdent";
import { UnreadyBanner } from "@/components/site/Unready";
import { LinkButton } from "@/components/ui/link-button";
import { getSettings } from "@/lib/settings/settings-server";
import { Box, Divider, LStack, VStack, styled } from "@/styled-system/jsx";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

type Props = {
  params: Promise<{
    invitation_id: string;
  }>;
};

export default async function Page({ params }: Props) {
  try {
    const { invitation_id: invitationID } = await params;
    const [{ data: invitation }, settings, session] = await Promise.all([
      invitationGet(invitationID, { cache: "no-store" }),
      getSettings(),
      getServerSession({ cache: "no-store" }),
    ]);

    return (
      <LStack gap="4" textAlign="center" alignItems="center">
        <Divider />

        <VStack gap="1">
          <Box>
            <MemberIdent
              profile={invitation.creator}
              size="md"
              name="full-vertical"
            />
          </Box>
          <styled.p>
            has invited you to <strong>{settings.title}</strong>
          </styled.p>
        </VStack>

        {session ? (
          <VStack
            w="full"
            gap="4"
            borderWidth="thin"
            borderStyle="solid"
            borderColor="border.warning"
            bgColor="bg.warning"
            color="fg.warning"
            borderRadius="md"
            p="4"
          >
            <VStack gap="1" textWrap="balance">
              <styled.p fontWeight="semibold">
                You&apos;re already signed in as{" "}
                <strong>{session.handle}</strong> on {settings.title}.
              </styled.p>
              <styled.p fontSize="sm">
                You cannot accept an invitation while already signed in.
              </styled.p>
            </VStack>

            <LinkButton w="full" href="/">
              Home
            </LinkButton>
          </VStack>
        ) : (
          <LinkButton
            w="full"
            href={`/register?invitation_id=${invitation.id}`}
          >
            Accept
          </LinkButton>
        )}
      </LStack>
    );
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}

export async function generateMetadata({ params }: Props) {
  try {
    const { invitation_id: invitationID } = await params;
    const [{ data: invitation }, settings] = await Promise.all([
      invitationGet(invitationID, { cache: "no-store" }),
      getSettings(),
    ]);

    return {
      title: `${invitation.creator.name} invited you to ${settings.title}`,
      description: `Accept your invitation to join ${settings.title}.`,
    };
  } catch {
    return {
      title: "Invitation",
      description: "Accept your invitation to join the community.",
    };
  }
}
