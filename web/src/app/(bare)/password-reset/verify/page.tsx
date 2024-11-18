import { z } from "zod";

import { UnreadyBanner } from "@/components/site/Unready";
import { PasswordResetVerifyScreen } from "@/screens/auth/PasswordResetScreen/PasswordResetVerifyScreen";

type Props = {
  searchParams: Promise<{
    token: string;
  }>;
};

export const QuerySchema = z.object({
  token: z.string(),
});

export default async function Page(props: Props) {
  try {
    const searchParams = await props.searchParams;

    const parsed = QuerySchema.safeParse(searchParams);

    if (!parsed.success) {
      throw new Error(
        "Missing password reset token, please try resetting your password again.",
      );
    }

    const { token } = parsed.data;

    return <PasswordResetVerifyScreen token={token} />;
  } catch (error) {
    return <UnreadyBanner error={error} />;
  }
}
