"use client";

import { useRouter } from "next/navigation";
import { use, useEffect, useState } from "react";

import { oAuthProviderCallback } from "src/api/openapi-client/auth";
import { Unready } from "src/components/site/Unready";

import { handle } from "@/api/client";
import { OAuthCallback } from "@/api/openapi-schema";

export const dynamic = "force-dynamic";

export type Props = {
  params: Promise<{
    provider: string;
  }>;
  searchParams: Promise<OAuthCallback>;
};

export default function Page(props: Props) {
  const router = useRouter();
  const [error, setError] = useState<unknown | null>(null);

  const params = use(props.params);
  const searchParams = use(props.searchParams);

  useEffect(() => {
    handle(
      async () => {
        if (error != null) {
          return;
        }

        if (params.provider === undefined || searchParams === undefined) {
          return;
        }

        const { id } = await oAuthProviderCallback(
          params.provider,
          searchParams,
        );

        router.push(`/?id=${id}`);
      },
      {
        errorToast: false,
        onError: async (e) => {
          setError(e);
        },
      },
    );
  }, [router, error, params, searchParams]);

  return <Unready error={error} />;
}
