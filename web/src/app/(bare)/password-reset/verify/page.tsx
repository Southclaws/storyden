import { z } from "zod";

import { UnreadyBanner } from "@/components/site/Unready";
import { PasswordResetVerifyScreen } from "@/screens/auth/PasswordResetScreen/PasswordResetVerifyScreen";

type Props = {
  searchParams: Promise<{
    token: string;
  }>;
};

const QuerySchema = z.object({
  token: z.string().optional(),
});

export default async function Page(props: Props) {
  try {
    const searchParams = await props.searchParams;

    const parsed = QuerySchema.parse(searchParams);

    const { token } = parsed;

    if (!token) {
      return <p>Please check your email for a verification link.</p>;
    }

    return <PasswordResetVerifyScreen token={token} />;
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
