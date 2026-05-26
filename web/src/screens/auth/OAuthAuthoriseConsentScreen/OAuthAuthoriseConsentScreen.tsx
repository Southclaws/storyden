"use client";

import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useMemo, useState } from "react";

import { handle } from "@/api/client";
import { RequestError } from "@/api/common";
import {
  oAuthAuthoriseConsentSubmit,
  useOAuthAuthoriseConsent,
} from "@/api/openapi-client/auth";
import { Permission } from "@/api/openapi-schema";
import * as Alert from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { CheckCircleIcon } from "@/components/ui/icons/CheckCircle";
import { WarningIcon } from "@/components/ui/icons/Warning";
import {
  PermissionDetails,
  buildPermissionList,
} from "@/lib/permission/permission";
import { LStack, WStack, styled } from "@/styled-system/jsx";

export function OAuthAuthoriseConsentScreen() {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const requestID = searchParams.get("request_id") ?? "";
  const [submitted, setSubmitted] = useState(false);

  const returnURL = useMemo(() => {
    const params = new URLSearchParams(searchParams);
    return `${pathname}?${params.toString()}`;
  }, [pathname, searchParams]);

  const { data, error } = useOAuthAuthoriseConsent(
    { request_id: requestID },
    {
      swr: {
        enabled: requestID !== "" && !submitted,
        onError(err) {
          if (err instanceof RequestError && err.status === 401) {
            router.replace(
              `/login?return_url=${encodeURIComponent(returnURL)}`,
            );
          }
        },
      },
    },
  );

  if (!requestID) {
    return (
      <ConsentMessage
        icon={<WarningIcon />}
        title="Missing request"
        body="Open the full link from the application and try again."
      />
    );
  }

  if (error instanceof RequestError && error.status === 401) {
    return <ConsentMessage title="Redirecting" body="Sign in to continue." />;
  }

  if (submitted) {
    return (
      <ConsentMessage
        icon={<CheckCircleIcon />}
        title="Returning"
        body="You can return to the application."
      />
    );
  }

  if (error instanceof RequestError && error.message === "access_denied") {
    return (
      <Alert.Root colorPalette="red">
        <Alert.Icon asChild>
          <WarningIcon />
        </Alert.Icon>
        <Alert.Content>
          <Alert.Title>Access denied</Alert.Title>
          <Alert.Description>
            Your account does not have permission to use third-party
            applications. Contact an administrator for access.
          </Alert.Description>
        </Alert.Content>
      </Alert.Root>
    );
  }

  if (error) {
    return (
      <ConsentMessage
        icon={<WarningIcon />}
        title="Invalid request"
        body="This request is invalid, expired, or already used."
      />
    );
  }

  if (!data) {
    return <ConsentMessage title="Loading" body="Checking this request." />;
  }

  const grantedScopes = data.granted_scopes.filter(
    (scope) =>
      !["openid", "profile", "email", "offline_access"].includes(scope),
  );
  const permissions = buildPermissionList(
    ...grantedScopes.filter(isPermission),
  );

  async function submit(nextDecision: "approve" | "deny") {
    const result = await handle(
      () =>
        oAuthAuthoriseConsentSubmit({
          request_id: requestID,
          decision: nextDecision,
        }),
      { errorToast: true },
    );

    if (result) {
      setSubmitted(true);
      window.location.assign(result.location);
    }
  }

  return (
    <LStack gap="2">
      <LStack gap="1">
        <Heading size="md">{data.client_name}</Heading>
        <styled.p color="fg.muted">
          {data.inherits_user_permissions
            ? "This application will act with your current account permissions."
            : "This application is requesting access to your account."}
        </styled.p>
      </LStack>

      <Alert.Root colorPalette="orange">
        <Alert.Icon asChild>
          <WarningIcon />
        </Alert.Icon>
        <Alert.Content>
          <Alert.Title>Redirect destination</Alert.Title>
          <styled.code
            fontSize="sm"
            fontWeight="medium"
            fontFamily="mono"
            overflowWrap="anywhere"
          >
            {data.redirect_uri}
          </styled.code>
          <Alert.Description>
            Only approve if this destination matches where you started
            authentication.
          </Alert.Description>
        </Alert.Content>
      </Alert.Root>

      {permissions.length > 0 && (
        <LStack flexShrink="0" gap="2" w="full">
          <Heading size="sm" color="fg.muted">
            Permissions requested
          </Heading>

          <styled.ul
            display="flex"
            flexDir="column"
            gap="2"
            m="0"
            p="0"
            w="full"
          >
            {permissions.map((permission) => (
              <styled.li
                key={permission.value}
                borderColor="border.default"
                borderRadius="sm"
                borderWidth="thin"
                display="grid"
                gap="1"
                listStyle="none"
                p="3"
              >
                <styled.span fontSize="sm" fontWeight="medium">
                  {permission.name}
                </styled.span>
                <styled.span color="fg.muted" fontSize="sm">
                  {permission.description}
                </styled.span>
              </styled.li>
            ))}
          </styled.ul>
        </LStack>
      )}

      <WStack flexShrink="0">
        <Button onClick={() => submit("approve")}>Approve</Button>
        <Button variant="ghost" onClick={() => submit("deny")}>
          Deny
        </Button>
      </WStack>
    </LStack>
  );
}

function isPermission(scope: string): scope is Permission {
  return scope in PermissionDetails;
}

type ConsentMessageProps = {
  icon?: React.ReactNode;
  title: string;
  body: string;
};

function ConsentMessage({ icon, title, body }: ConsentMessageProps) {
  return (
    <LStack gap="3" textAlign="center" alignItems="center">
      {icon}
      <Heading size="md">{title}</Heading>
      <styled.p color="fg.muted">{body}</styled.p>
    </LStack>
  );
}
