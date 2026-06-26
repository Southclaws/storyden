import { getServerSession } from "@/auth/server-session";
import { UnreadyBanner } from "@/components/site/Unready";
import { EmailVerificationScreen } from "@/screens/auth/EmailVerificationScreen/EmailVerificationScreen";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

type Props = {
  searchParams: Promise<{
    returnURL?: string;
  }>;
};

export default async function Page(props: Props) {
  try {
    const searchParams = await props.searchParams;

    const account = await getServerSession({
      cache: "no-store",
    });

    return (
      <EmailVerificationScreen
        initialAccount={account}
        returnURL={searchParams.returnURL}
      />
    );
  } catch (e) {
    return <UnreadyBanner error={e} />;
  }
}
