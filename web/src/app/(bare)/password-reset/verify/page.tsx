import { z } from "zod";

import { UnreadyBanner } from "@/components/site/Unready";
import { PasswordResetVerifyScreen } from "@/screens/auth/PasswordResetScreen/PasswordResetVerifyScreen";

// TODO: Cache Components adoption. Refactor this route so this opt-out can be removed.
// See: https://nextjs.org/docs/app/guides/migrating-to-cache-components
export const instant = false;

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
