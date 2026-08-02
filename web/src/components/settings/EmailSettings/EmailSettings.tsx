import { useState } from "react";

import { Account } from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";
import { AddIcon } from "@/components/ui/icons/Add";
import { PageHeading } from "@/components/ui/page-heading";
import { SectionHeading } from "@/components/ui/section-heading";
import { Text } from "@/components/ui/text";
import { LStack, WStack } from "@/styled-system/jsx";

import { EmailCard } from "./EmailCard";
import { EmailCreateForm } from "./EmailCreateForm";

export type Props = {
  account: Account;
};

export function useEmailSettings({ account }: Props) {
  const [adding, setAdding] = useState(false);

  async function handleStartNewEmail() {
    setAdding(true);
  }

  async function handleCancelNewEmail() {
    setAdding(false);
  }

  async function handleFinishNewEmail() {
    setAdding(false);
  }

  return {
    data: {
      emails: account.email_addresses,
      adding,
    },
    handlers: {
      handleStartNewEmail,
      handleCancelNewEmail,
      handleFinishNewEmail,
    },
  };
}

export function EmailSettings(props: Props) {
  const { data, handlers } = useEmailSettings(props);

  return (
    <LStack gap="4">
      <LStack>
        <PageHeading>Email settings</PageHeading>
        <p>
          Manage your email addresses here. You can add multiple email addresses
          and use them to log in to your account. Emails are also used for
          newsletters, notifications and other communications.
        </p>
      </LStack>

      <LStack>
        <WStack>
          <SectionHeading>Email addresses</SectionHeading>
          <Button variant="subtle" onClick={handlers.handleStartNewEmail}>
            <AddIcon /> new email address
          </Button>
        </WStack>

        {data.emails.length === 0 ? (
          <Text variant="supporting">
            You do not have any email addresses associated with your account.
          </Text>
        ) : (
          data.emails.map((email) => <EmailCard key={email.id} email={email} />)
        )}

        {data.adding && (
          <EmailCreateForm
            onCancel={handlers.handleCancelNewEmail}
            onSuccess={handlers.handleFinishNewEmail}
          />
        )}
      </LStack>
    </LStack>
  );
}
