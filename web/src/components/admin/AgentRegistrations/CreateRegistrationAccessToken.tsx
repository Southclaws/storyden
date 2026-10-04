import { createListCollection } from "@ark-ui/react";

import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { Button } from "@/components/ui/button";
import * as Clipboard from "@/components/ui/clipboard";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormLabel } from "@/components/ui/form-label";
import { Input } from "@/components/ui/input";
import { SelectField } from "@/components/ui/select";
import { Text } from "@/components/ui/text";
import { Textarea } from "@/components/ui/textarea";
import { getAPIAddress } from "@/config";
import { css } from "@/styled-system/css";
import { LStack, WStack } from "@/styled-system/jsx";

import { useCreateRegistrationAccessToken } from "./useCreateRegistrationAccessToken";

const lifetimes = createListCollection({
  items: [
    { label: "1 hour", value: "1" },
    { label: "1 day", value: "24" },
    { label: "7 days", value: "168" },
    { label: "30 days", value: "720" },
  ],
});

export function CreateRegistrationAccessToken({
  onClose,
}: {
  onClose: () => void;
}) {
  const { form, handleSubmit, createdSecret } =
    useCreateRegistrationAccessToken();
  const submitting = form.formState.isSubmitting;
  const agentPrompt = `Register a bot account with \`sd\` on \`${getAPIAddress()}\`. Choose a handle and a local context name.

Run \`sd auth register ${getAPIAddress()} --handle <handle> --name <context> --registration-token-stdin --format json\`, supplying this token on stdin: \`${createdSecret ?? ""}\`

Then use \`sd --context <context> --help\` to learn how to use the Storyden command line tool.
`;

  return (
    <ModalDrawer
      isOpen
      onClose={submitting ? undefined : onClose}
      dismissable={!submitting}
      className={css({ width: { base: "full", md: "breakpoint-sm" } })}
      title={
        createdSecret
          ? "Registration token created"
          : "Create registration token"
      }
    >
      {createdSecret ? (
        <LStack gap="4">
          <Text>
            This token is shown only once. Copy it and share it with the person
            or agent allowed to register.
          </Text>
          <Clipboard.Root value={createdSecret} w="full">
            <Clipboard.Label>Registration token</Clipboard.Label>
            <Clipboard.Control>
              <Clipboard.Input asChild>
                <Input readOnly aria-label="Registration token" />
              </Clipboard.Input>
              <Clipboard.Trigger asChild>
                <Button variant="outline">
                  <Clipboard.Indicator copied="Copied">
                    Copy
                  </Clipboard.Indicator>
                </Button>
              </Clipboard.Trigger>
            </Clipboard.Control>
          </Clipboard.Root>
          <Text variant="supporting">
            This token authorises account registration only; the new account’s
            roles determine its permissions.
          </Text>
          <Clipboard.Root value={agentPrompt} w="full">
            <Clipboard.Label>Prompt for your agent</Clipboard.Label>
            <Textarea
              readOnly
              aria-label="Prompt for your agent"
              value={agentPrompt}
              rows={8}
            />
            <Clipboard.Trigger asChild>
              <Button variant="outline" aria-label="Copy agent prompt">
                <Clipboard.Indicator copied="Copied">
                  Copy agent prompt
                </Clipboard.Indicator>
              </Button>
            </Clipboard.Trigger>
          </Clipboard.Root>
          <Button onClick={onClose}>Done</Button>
        </LStack>
      ) : (
        <form onSubmit={handleSubmit}>
          <LStack gap="4">
            <Text variant="supporting">
              Pre-approve a limited number of new agent accounts. Each
              successful registration uses one allowance.
            </Text>
            <FormControl>
              <FormLabel>Label</FormLabel>
              <Input
                id="registration-token-label"
                aria-label="Label"
                placeholder="e.g. Engineering team"
                maxLength={200}
                {...form.register("label")}
              />
              <FormErrorText>
                {form.formState.errors.label?.message}
              </FormErrorText>
            </FormControl>
            <FormControl>
              <FormLabel>Valid for</FormLabel>
              <SelectField
                control={form.control}
                name="lifetime"
                collection={lifetimes}
                placeholder="Choose a lifetime"
                ariaLabel="Valid for"
              />
              <FormErrorText>
                {form.formState.errors.lifetime?.message}
              </FormErrorText>
            </FormControl>
            <FormControl>
              <FormLabel>Maximum registrations</FormLabel>
              <Input
                id="registration-token-limit"
                aria-label="Maximum registrations"
                type="number"
                min={1}
                max={2147483647}
                step={1}
                {...form.register("max_registrations", { valueAsNumber: true })}
              />
              <FormErrorText>
                {form.formState.errors.max_registrations?.message}
              </FormErrorText>
            </FormControl>
            <FormErrorText>{form.formState.errors.root?.message}</FormErrorText>
            <WStack>
              <Button
                type="button"
                variant="outline"
                onClick={onClose}
                disabled={submitting}
              >
                Cancel
              </Button>
              <Button type="submit" variant="solid" loading={submitting}>
                Create token
              </Button>
            </WStack>
          </LStack>
        </form>
      )}
    </ModalDrawer>
  );
}
