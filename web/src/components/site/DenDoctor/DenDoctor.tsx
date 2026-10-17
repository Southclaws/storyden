"use client";

import { useSyncExternalStore } from "react";

import * as Alert from "@/components/ui/alert";
import { WarningIcon } from "@/components/ui/icons/Warning";
import { LinkButton } from "@/components/ui/link-button";
import { Text } from "@/components/ui/text";
import { API_ADDRESS, WEB_ADDRESS } from "@/config";
import { diagnoseInstallation } from "@/lib/den-doctor/diagnostics";
import { LStack, WStack, styled } from "@/styled-system/jsx";

type Props = {
  apiAddress?: string;
  webAddress?: string;
  browserOrigin?: string;
  showHealthy?: boolean;
};

export function DenDoctor({
  apiAddress = API_ADDRESS,
  webAddress = WEB_ADDRESS,
  browserOrigin,
  showHealthy = false,
}: Props) {
  const currentOrigin = useSyncExternalStore(
    subscribeToLocation,
    readBrowserOrigin,
    () => undefined,
  );
  const origin = browserOrigin ?? currentOrigin;
  const issues = diagnoseInstallation(apiAddress, webAddress, origin);
  const api = new URL(apiAddress);
  const web = new URL(webAddress);
  const tone = issues.some((issue) => issue.severity === "error")
    ? "danger"
    : issues.length
      ? "warning"
      : "info";

  if (!issues.length && !showHealthy) {
    return null;
  }

  return (
    <Alert.Root tone={tone} mb="4">
      {issues.length > 0 && (
        <Alert.Icon asChild>
          <WarningIcon aria-hidden="true" />
        </Alert.Icon>
      )}
      <Alert.Content
        minW="0"
        width="full"
        css={{ "& code": { fontFamily: "mono", overflowWrap: "anywhere" } }}
      >
        <Alert.Title as="h2">Den Doctor</Alert.Title>
        <Alert.Description>
          {issues.length
            ? "Check these addresses before inviting people to your community."
            : "No address problems detected. These checks compare your configured URLs; they do not test server connectivity."}
        </Alert.Description>
        {issues.length > 0 && (
          <LStack as="ul" gap="3" py="2" alignItems="stretch">
            {issues.map((issue) => (
              <li key={issue.code}>
                <LStack gap="1" alignItems="stretch">
                  <Text
                    as="h3"
                    variant="supporting"
                    fontWeight="semibold"
                    color="inherit"
                  >
                    {issue.title}
                  </Text>
                  <Text variant="supporting" color="inherit">
                    {issue.explanation}
                  </Text>
                  <Text variant="supporting" color="inherit">
                    {issue.remedy}
                  </Text>
                </LStack>
              </li>
            ))}
          </LStack>
        )}
        <details>
          <styled.summary cursor="pointer" fontSize="sm">
            Show configuration details
          </styled.summary>
          <LStack gap="2" py="3" alignItems="stretch" overflowWrap="anywhere">
            <Text variant="supporting" color="inherit">
              Website (<code>PUBLIC_WEB_ADDRESS</code>):{" "}
              <code>{`${web.origin}${web.pathname}`}</code>
            </Text>
            <Text variant="supporting" color="inherit">
              API (<code>PUBLIC_API_ADDRESS</code>):{" "}
              <code>{`${api.origin}${api.pathname}`}</code>
            </Text>
            {origin && (
              <Text variant="supporting" color="inherit">
                This browser: <code>{origin}</code>
              </Text>
            )}
            <Text variant="supporting" color="inherit">
              Change these addresses in your deployment configuration, then
              restart Storyden. They cannot be changed from this page.
            </Text>
          </LStack>
        </details>
        {issues.length > 0 && (
          <WStack gap="2" pt="2" flexWrap="wrap">
            {issues.some((issue) => issue.code === "website_origin") && (
              <LinkButton
                href={web.origin}
                variant="outline"
                maxW="full"
                whiteSpace="normal"
                h="auto"
                minH="6"
                py="0.5"
              >
                Open configured website
              </LinkButton>
            )}
            <LinkButton
              href="https://www.storyden.org/docs/introduction/troubleshooting"
              variant="outline"
              alignSelf="start"
              maxW="full"
              whiteSpace="normal"
              h="auto"
              minH="6"
              py="0.5"
            >
              Read troubleshooting guide
            </LinkButton>
          </WStack>
        )}
      </Alert.Content>
    </Alert.Root>
  );
}

function readBrowserOrigin() {
  return window.location.origin;
}

function subscribeToLocation(callback: () => void) {
  window.addEventListener("popstate", callback);
  return () => window.removeEventListener("popstate", callback);
}
