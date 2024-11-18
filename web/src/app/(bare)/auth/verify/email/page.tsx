import { getServerSession } from "@/auth/server-session";
import { UnreadyBanner } from "@/components/site/Unready";
import { EmailVerificationScreen } from "@/screens/auth/EmailVerificationScreen/EmailVerificationScreen";

type Props = {
  searchParams: Promise<{
    returnURL?: string;
  }>;
};

export default async function Page(props: Props) {
  try {
    const searchParams = await props.searchParams;

    const account = await getServerSession();

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
