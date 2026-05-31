import { SelectValueChangeDetails, createListCollection } from "@ark-ui/react";
import { useMemo, useState } from "react";

import { handle } from "@/api/client";
import {
  accountEmailPasswordReset,
  accountPasswordResetTokenGet,
} from "@/api/openapi-client/accounts";
import { Account, ProfileReference } from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";
import * as Clipboard from "@/components/ui/clipboard";
import { IconButton } from "@/components/ui/icon-button";
import { CheckIcon } from "@/components/ui/icons/Check";
import { CopyIcon } from "@/components/ui/icons/Copy";
import { SelectIcon } from "@/components/ui/icons/Select";
import { Input } from "@/components/ui/input";
import * as Select from "@/components/ui/select";
import { WEB_ADDRESS } from "@/config";
import { HStack, VStack, styled } from "@/styled-system/jsx";

export type Props = {
  profile: ProfileReference;
  account: Account;
  hasEmail: boolean;
  onClose?: () => void;
};

export function MemberPasswordResetForm({
  profile,
  account,
  hasEmail,
  onClose,
}: Props) {
  const hasVerifiedEmailAddress = account.email_addresses?.some(
    (emailAddress) => emailAddress.verified,
  );

  if (hasEmail && hasVerifiedEmailAddress) {
    return (
      <MemberPasswordResetFormWithEmail
        profile={profile}
        account={account}
        hasEmail={hasEmail}
        onClose={onClose}
      />
    );
  } else {
    return (
      <MemberPasswordResetFormWithLink
        profile={profile}
        account={account}
        hasEmail={hasEmail}
        onClose={onClose}
      />
    );
  }
}

function MemberPasswordResetFormWithEmail({
  profile,
  account,
  onClose,
}: Props) {
  const verifiedEmailAddresses = useMemo(
    () =>
      account.email_addresses?.filter(
        (emailAddress) => emailAddress.verified,
      ) ?? [],
    [account.email_addresses],
  );
  const [selectedEmailAddressId, setSelectedEmailAddressId] = useState(
    () => verifiedEmailAddresses[0]?.id,
  );
  const selectedEmailAddress =
    verifiedEmailAddresses.find(
      (emailAddress) => emailAddress.id === selectedEmailAddressId,
    ) ?? verifiedEmailAddresses[0];
  const hasVerifiedEmailAddress = verifiedEmailAddresses.length > 0;
  const hasMultipleVerifiedEmailAddresses = verifiedEmailAddresses.length > 1;
  const emailAddressCollection = useMemo(
    () =>
      createListCollection({
        items: verifiedEmailAddresses.map((emailAddress) => ({
          label: emailAddress.email_address,
          value: emailAddress.id,
        })),
      }),
    [verifiedEmailAddresses],
  );

  async function handleSendEmail() {
    if (!selectedEmailAddress) {
      throw new Error("No verified email address found for this member");
    }

    await handle(
      async () => {
        await accountEmailPasswordReset(profile.id, {
          email_address_id: selectedEmailAddress.id,
          token_url: {
            url: `${WEB_ADDRESS}/password-reset/verify`,
            query: "token",
          },
        });

        onClose?.();
      },
      {
        promiseToast: {
          loading: "Sending password reset email...",
          success: "Password reset email sent successfully",
        },
      },
    );
  }

  function handleEmailAddressChange({ value }: SelectValueChangeDetails) {
    const [emailAddressId] = value;
    setSelectedEmailAddressId(emailAddressId);
  }

  return (
    <VStack alignItems="start" gap="4">
      <styled.p>Send a password reset email to {profile.name}?</styled.p>

      {!hasVerifiedEmailAddress && (
        <styled.p fontSize="sm" color="fg.error">
          This account has no verified email address configured.
        </styled.p>
      )}

      {selectedEmailAddress && !hasMultipleVerifiedEmailAddresses && (
        <VStack alignItems="start" gap="1" w="full">
          <styled.span fontSize="sm" fontWeight="medium">
            Recipient email
          </styled.span>
          <styled.p
            bg="bg.muted"
            borderColor="border.default"
            borderRadius="md"
            borderWidth="hairline"
            color="fg.default"
            fontSize="sm"
            px="3"
            py="2"
            w="full"
          >
            {selectedEmailAddress.email_address}
          </styled.p>
        </VStack>
      )}

      {hasMultipleVerifiedEmailAddresses && (
        <Select.Root
          collection={emailAddressCollection}
          positioning={{ sameWidth: true }}
          value={selectedEmailAddress ? [selectedEmailAddress.id] : []}
          onValueChange={handleEmailAddressChange}
          w="full"
        >
          <Select.Label>Recipient email</Select.Label>
          <Select.Control>
            <Select.Trigger w="full">
              <Select.ValueText placeholder="Select an email address" />
              <SelectIcon />
            </Select.Trigger>
          </Select.Control>
          <Select.Positioner>
            <Select.Content>
              {emailAddressCollection.items.map((item) => (
                <Select.Item key={item.value} item={item}>
                  <Select.ItemText>{item.label}</Select.ItemText>
                  <Select.ItemIndicator>
                    <CheckIcon />
                  </Select.ItemIndicator>
                </Select.Item>
              ))}
            </Select.Content>
          </Select.Positioner>
        </Select.Root>
      )}

      <styled.p fontSize="sm" color="fg.muted">
        An email will be sent to their verified email address with instructions
        to reset their password. The link will be valid for 1 hour.
      </styled.p>

      <HStack w="full">
        <Button type="button" flexGrow="1" variant="outline" onClick={onClose}>
          Cancel
        </Button>

        <Button
          type="button"
          flexGrow="1"
          onClick={handleSendEmail}
          disabled={!hasVerifiedEmailAddress}
        >
          Send Email
        </Button>
      </HStack>
    </VStack>
  );
}

function MemberPasswordResetFormWithLink({ profile, onClose }: Props) {
  const [resetUrl, setResetUrl] = useState<string | null>(null);

  async function handleGenerateToken() {
    await handle(
      async () => {
        const response = await accountPasswordResetTokenGet(profile.id);

        const url = new URL("/password-reset/verify", WEB_ADDRESS);
        url.searchParams.set("token", response.token);
        setResetUrl(url.toString());
      },
      {
        errorToast: true,
      },
    );
  }

  if (resetUrl) {
    return (
      <VStack alignItems="start" gap="4">
        <styled.p>
          <strong>Password reset link generated successfully.</strong>
        </styled.p>

        <styled.p>
          Copy this link and send it to {profile.name}. This link is valid for 1
          hour.
        </styled.p>

        <Clipboard.Root w="full" value={resetUrl}>
          <Clipboard.Control>
            <Clipboard.Input asChild>
              <Input readOnly />
            </Clipboard.Input>
            <Clipboard.Trigger asChild>
              <IconButton variant="outline">
                <Clipboard.Indicator copied={<CheckIcon />}>
                  <CopyIcon />
                </Clipboard.Indicator>
              </IconButton>
            </Clipboard.Trigger>
          </Clipboard.Control>
        </Clipboard.Root>

        <Button onClick={onClose} w="full">
          Done
        </Button>
      </VStack>
    );
  }

  return (
    <VStack alignItems="start" gap="4">
      <styled.p>Generate a password reset link for {profile.name}?</styled.p>

      <styled.p fontSize="sm" color="fg.muted">
        A password reset email cannot be sent for this member. A password reset
        link will be generated for you to copy and send through another
        communication method.
      </styled.p>

      <HStack w="full">
        <Button type="button" flexGrow="1" variant="outline" onClick={onClose}>
          Cancel
        </Button>

        <Button type="button" flexGrow="1" onClick={handleGenerateToken}>
          Generate Link
        </Button>
      </HStack>
    </VStack>
  );
}
