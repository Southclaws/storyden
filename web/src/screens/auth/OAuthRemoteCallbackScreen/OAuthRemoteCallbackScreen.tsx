"use client";

import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useEffect, useMemo } from "react";

import { RequestError } from "@/api/common";
import { useOAuthRemoteCallback } from "@/api/openapi-client/auth";
import { useGetSession } from "@/api/openapi-client/misc";
import { Spinner } from "@/components/ui/Spinner";
import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { CheckCircleIcon } from "@/components/ui/icons/CheckCircle";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { LStack, VStack, styled } from "@/styled-system/jsx";

export function OAuthRemoteCallbackScreen() {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const code = searchParams.get("code") ?? "";
  const state = searchParams.get("state") ?? "";
  const providerError = searchParams.get("error") ?? "";
  const providerErrorDescription = searchParams.get("error_description") ?? "";
  const returnURL = useMemo(() => {
    const params = new URLSearchParams(searchParams);
    const query = params.toString();
    return query ? `${pathname}?${query}` : pathname;
  }, [pathname, searchParams]);

  const session = useGetSession({
    swr: {
      shouldRetryOnError: false,
      revalidateIfStale: false,
      revalidateOnFocus: false,
    },
  });
  const instanceTitle = session.data?.info.title || "this instance";
  const signedIn = !!session.data?.account;

  useEffect(() => {
    if (!session.data || signedIn) {
      return;
    }

    router.replace(`/login?return_url=${encodeURIComponent(returnURL)}`);
  }, [returnURL, router, session.data, signedIn]);

  const canComplete =
    signedIn && code !== "" && state !== "" && providerError === "";
  const { data, error, isLoading } = useOAuthRemoteCallback(
    { code, state },
    {
      swr: {
        enabled: canComplete,
        shouldRetryOnError: false,
        revalidateIfStale: false,
        revalidateOnFocus: false,
      },
    },
  );

  if (!session.data) {
    return (
      <CallbackMessage
        icon={<Spinner />}
        title="Checking session"
        body="Keep this tab open while the instance checks your sign-in state."
      />
    );
  }

  if (!signedIn) {
    return (
      <CallbackMessage
        icon={<Spinner />}
        title="Redirecting"
        body={`Sign in to ${instanceTitle} to complete authentication.`}
      />
    );
  }

  if (providerError !== "") {
    return (
      <CallbackMessage
        tone="error"
        icon={<WarningIcon />}
        title="Authentication was not completed"
        body={
          providerErrorDescription ||
          "The remote authorization server did not approve this request."
        }
      />
    );
  }

  if (code === "" || state === "") {
    return (
      <CallbackMessage
        tone="error"
        icon={<WarningIcon />}
        title="Missing callback details"
        body="Open the full authentication link and try again."
      />
    );
  }

  if (error) {
    return (
      <CallbackMessage
        tone="error"
        icon={<WarningIcon />}
        title="Authentication failed"
        body={callbackErrorMessage(error)}
      />
    );
  }

  if (isLoading || !data) {
    return (
      <CallbackMessage
        icon={<Spinner />}
        title="Finishing authentication"
        body={`Keep this tab open while ${instanceTitle} completes the flow.`}
      />
    );
  }

  return (
    <CallbackMessage
      tone="success"
      icon={<CheckCircleIcon />}
      title="Authentication complete"
      body={`You can now close this tab and return to ${instanceTitle}.`}
      action={
        <Button onClick={() => window.close()} variant="outline">
          Close tab
        </Button>
      }
    />
  );
}

type CallbackMessageProps = {
  icon?: React.ReactNode;
  title: string;
  body: string;
  tone?: "neutral" | "success" | "error";
  action?: React.ReactNode;
};

function CallbackMessage({
  icon,
  title,
  body,
  tone = "neutral",
  action,
}: CallbackMessageProps) {
  return (
    <VStack gap="4" textAlign="center">
      {icon && (
        <styled.div
          alignItems="center"
          bg={toneBackground(tone)}
          borderRadius="full"
          color={toneForeground(tone)}
          display="inline-flex"
          h="12"
          justifyContent="center"
          w="12"
          css={{
            "& svg": {
              h: "6",
              w: "6",
            },
          }}
        >
          {icon}
        </styled.div>
      )}

      <VStack gap="2" textAlign="center">
        <Heading size="md">{title}</Heading>
        <styled.p color="fg.muted">{body}</styled.p>
      </VStack>

      {action}
    </VStack>
  );
}

function callbackErrorMessage(error: unknown) {
  if (error instanceof RequestError) {
    return (
      error.problem?.detail ||
      error.problem?.title ||
      "This authentication request is invalid, expired, or already used."
    );
  }

  return "The instance could not complete this authentication request.";
}

function toneBackground(tone: CallbackMessageProps["tone"]) {
  switch (tone) {
    case "success":
      return "bg.success";
    case "error":
      return "bg.error";
    default:
      return "bg.subtle";
  }
}

function toneForeground(tone: CallbackMessageProps["tone"]) {
  switch (tone) {
    case "success":
      return "fg.success";
    case "error":
      return "fg.error";
    default:
      return "fg.muted";
  }
}
