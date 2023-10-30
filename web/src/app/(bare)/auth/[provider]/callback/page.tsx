import { redirect } from "next/navigation";

import { AuthSuccessOKResponse } from "src/api/openapi/schemas";
import { server } from "src/api/server";
import ErrorBanner from "src/components/site/ErrorBanner";
import { Link } from "src/theme/components/Link";
import { deriveError } from "src/utils/error";

import { HStack, VStack } from "@/styled-system/jsx";

export type Props = {
  params: {
    provider: string;
  };
  searchParams: {
    code: string;
    state: string;
  };
};

export default async function Page(props: Props) {
  try {
    const { id } = await server<AuthSuccessOKResponse>({
      url: `/v1/auth/oauth/${props.params.provider}/callback`,
      method: "post",
      data: props.searchParams,
    });

    return redirect(`/?id=${id}`);
  } catch (e) {
    console.log(e);

    const message = deriveError(e);
    const error = (e as any).error ?? undefined;

    return (
      <VStack height="dvh" justify="center" p="10">
        <ErrorBanner message={message} error={error} />

        <HStack>
          <Link href="/register">Back to register</Link>
          <Link href="/login">Back to login</Link>
        </HStack>
      </VStack>
    );
  }
}
