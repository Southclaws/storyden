import type { PublicProfile } from "@/api/openapi-schema";
import { profileGet } from "@/api/openapi-server/profiles";
import { getServerSession } from "@/auth/server-session";
import { ServerRenderedContent } from "@/components/content/ContentRenderer/ServerRenderedContent";
import { getAvatarURL } from "@/components/member/MemberBadge/MemberAvatar";
import { UnreadyBanner } from "@/components/site/Unready";
import { canonicalPage, jsonLD } from "@/lib/metadata/public";
import { getSettings } from "@/lib/settings/settings-server";
import { ProfileScreen } from "@/screens/profile/ProfileScreen";

type Props = {
  params: Promise<{ handle: string }>;
};

export default async function Page(props: Props) {
  const { handle } = await props.params;
  let session;
  let settings;
  let data;

  try {
    session = await getServerSession();
    settings = await getSettings();
    ({ data } = await profileGet(handle));
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }

  let publicProfile: PublicProfile | undefined;

  try {
    ({ data: publicProfile } = await profileGet(handle, {
      anonymous: true,
      cache: "no-store",
    }));
  } catch {
    // The signed-in profile can still render without public structured data.
  }

  const url =
    publicProfile &&
    canonicalPage(`/m/${encodeURIComponent(publicProfile.handle)}`);
  const structuredData = publicProfile && {
    "@context": "https://schema.org",
    "@type": "ProfilePage",
    url,
    dateCreated: publicProfile.joined,
    mainEntity: {
      "@type": "Person",
      name: publicProfile.name || publicProfile.handle,
      alternateName: publicProfile.handle,
      identifier: publicProfile.id,
      url,
      description: publicProfile.bio,
      image: getAvatarURL(encodeURIComponent(publicProfile.handle)),
    },
  };

  return (
    <>
      {structuredData && (
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: jsonLD(structuredData) }}
        />
      )}
      <ServerRenderedContent values={[data.bio, data.signature ?? undefined]}>
        <ProfileScreen
          initialSession={session}
          profile={data}
          initialSignatureConfig={settings.metadata.signatures}
        />
      </ServerRenderedContent>
    </>
  );
}
