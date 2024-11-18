import { useState } from "react";

import { Account } from "@/api/openapi-schema";
import { Button } from "@/components/ui/button";
import { Heading } from "@/components/ui/heading";
import { AddIcon } from "@/components/ui/icons/Add";
import { CardBox, LStack, WStack, styled } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";

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
    <CardBox className={lstack()} gap="4">
      <LStack>
        <Heading size="md">Email settings</Heading>
        <p>
          Manage your email addresses here. You can add multiple email addresses
          and use them to log in to your account. Emails are also used for
          newsletters, notifications and other communications.
        </p>
      </LStack>

      <LStack>
        <WStack>
          <Heading size="sm">Email addresses</Heading>
          <Button
            size="xs"
            variant="subtle"
            onClick={handlers.handleStartNewEmail}
          >
            <AddIcon /> new email address
          </Button>
        </WStack>

        {data.emails.length === 0 ? (
          <styled.p color="fg.muted">
            You do not have any email addresses associated with your account.
          </styled.p>
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
    </CardBox>
  );
}
