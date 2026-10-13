import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import useSWR from "swr";

import { handle } from "@/api/client";
import { Button } from "@/components/ui/button";
import { CheckboxField } from "@/components/ui/checkbox";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormHelperText } from "@/components/ui/form-helper-text";
import { FormLabel } from "@/components/ui/form-label";
import { PageHeading } from "@/components/ui/page-heading";
import { SectionHeading } from "@/components/ui/section-heading";
import { Text } from "@/components/ui/text";
import { Textarea } from "@/components/ui/textarea";
import { resolveDiscoverySettings } from "@/lib/settings/discovery";
import { useSettingsMutation } from "@/lib/settings/mutation";
import { generateRobotsTxt } from "@/lib/settings/robots-txt";
import type { AdminSettings } from "@/lib/settings/settings";
import * as z from "@/lib/validation";
import { LStack, WStack, styled } from "@/styled-system/jsx";

const FormSchema = z.object({
  search_indexing: z.boolean(),
  training_crawling: z.boolean(),
  llms_txt: z.boolean(),
  auth_md: z.boolean(),
  ai_catalog: z.boolean(),
  api_catalog: z.boolean(),
  mcp_server_card: z.boolean(),
  robots_txt: z.string().check(z.maxLength(65536)),
});

type Form = z.infer<typeof FormSchema>;

export function DiscoverySettingsForm({
  settings,
}: {
  settings: AdminSettings;
}) {
  const { updateSettings } = useSettingsMutation();
  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: resolveDiscoverySettings(settings.services?.discovery),
  });
  const robotsSetting = form.watch("robots_txt");
  const customRobots = robotsSetting !== "";
  const [robotsDraft, setRobotsDraft] = useState<string>();
  const [resettingRobots, setResettingRobots] = useState(false);
  const busy = form.formState.isSubmitting || resettingRobots;
  const {
    data: publishedRobots,
    error: robotsError,
    isLoading: robotsLoading,
    mutate: refreshRobots,
  } = useSWR("/robots.txt", fetchRobotsTxt, { revalidateOnFocus: false });

  const onSubmit = form.handleSubmit(async (data) => {
    await handle(
      async () => {
        await updateSettings({ services: { discovery: data } });
        form.reset(data);
        await refreshRobots(undefined, {
          revalidate: true,
          throwOnError: false,
        });
        setRobotsDraft(undefined);
      },
      {
        promiseToast: {
          loading: "Saving discovery settings...",
          success: "Discovery settings saved",
        },
      },
    );
  });

  async function resetRobotsTxt() {
    setResettingRobots(true);

    try {
      const document = await generateRobotsTxt(
        form.getValues("search_indexing"),
        form.getValues("training_crawling"),
      );

      form.setValue("robots_txt", "", {
        shouldDirty: true,
        shouldValidate: true,
      });
      setRobotsDraft(document);
    } catch {
      form.setError("robots_txt", {
        message:
          "Could not load default robots.txt. Your text is unchanged; try again.",
      });
    } finally {
      setResettingRobots(false);
    }
  }

  return (
    <styled.form
      width="full"
      display="flex"
      flexDirection="column"
      gap="4"
      onSubmit={onSubmit}
    >
      <LStack gap="1">
        <WStack>
          <PageHeading>Discovery</PageHeading>
          <Button type="submit" loading={busy}>
            Save
          </Button>
        </WStack>
        <Text variant="supporting">
          Choose how search engines and other tools discover your community.
        </Text>
      </LStack>

      <SectionHeading>Search and crawlers</SectionHeading>
      <FormControl>
        <CheckboxField
          control={form.control}
          name="search_indexing"
          readOnly={busy}
        >
          Enable search indexing
        </CheckboxField>
        <FormHelperText>
          Publish sitemaps and allow search crawlers in the default robots.txt.
        </FormHelperText>
      </FormControl>
      <FormControl>
        <CheckboxField
          control={form.control}
          name="training_crawling"
          readOnly={busy}
        >
          Allow AI training crawlers
        </CheckboxField>
        <FormHelperText>
          The default robots.txt asks known AI training crawlers to stay away
          unless this is enabled.
        </FormHelperText>
      </FormControl>

      <SectionHeading>Discovery documents</SectionHeading>
      <Text variant="supporting">
        Disabled documents return “not found”. These controls do not change
        member permissions or enable or disable MCP itself.
      </Text>
      <FormControl>
        <CheckboxField control={form.control} name="llms_txt" readOnly={busy}>
          llms.txt
        </CheckboxField>
        <FormHelperText>A guide to your community for AI tools.</FormHelperText>
      </FormControl>
      <FormControl>
        <CheckboxField control={form.control} name="auth_md" readOnly={busy}>
          auth.md
        </CheckboxField>
        <FormHelperText>
          A readable guide to authenticating with the API.
        </FormHelperText>
      </FormControl>
      <FormControl>
        <CheckboxField control={form.control} name="ai_catalog" readOnly={busy}>
          AI catalog
        </CheckboxField>
        <FormHelperText>
          Lists the available AI integration services.
        </FormHelperText>
      </FormControl>
      <FormControl>
        <CheckboxField
          control={form.control}
          name="mcp_server_card"
          readOnly={busy}
        >
          MCP Server Card
        </CheckboxField>
        <FormHelperText>
          Advertises the MCP connection details when MCP is enabled.
        </FormHelperText>
      </FormControl>
      <FormControl>
        <CheckboxField
          control={form.control}
          name="api_catalog"
          readOnly={busy}
        >
          API catalog
        </CheckboxField>
        <FormHelperText>
          Helps applications discover the community API.
        </FormHelperText>
      </FormControl>

      <SectionHeading>robots.txt</SectionHeading>
      <FormControl>
        <FormLabel id="robots-txt-label">robots.txt</FormLabel>
        <Textarea
          id="robots-txt"
          aria-labelledby="robots-txt-label"
          rows={12}
          {...form.register("robots_txt")}
          value={
            robotsDraft ??
            (customRobots ? robotsSetting : (publishedRobots ?? ""))
          }
          onChange={(event) => {
            const text = event.currentTarget.value;
            setRobotsDraft(text);
            form.setValue("robots_txt", text, {
              shouldDirty: true,
              shouldValidate: true,
            });
          }}
          disabled={busy}
          placeholder={
            robotsLoading
              ? "Loading robots.txt..."
              : "Save to generate robots.txt"
          }
          spellCheck={false}
          aria-describedby="robots-txt-help robots-txt-error"
        />
        <FormHelperText id="robots-txt-help">
          {customRobots
            ? "Custom text replaces the generated crawler rules and sitemap links. Preset changes won’t overwrite your edits. Use default robots.txt and save to return to the generated policy."
            : "Generated from your search and crawler settings. Save to publish preset changes and refresh this document, or edit it to use custom text. Clearing the field restores the generated policy on save."}{" "}
          The search indexing setting still controls whether sitemaps are
          served. Crawlers choose whether to honour these requests.
        </FormHelperText>
        <FormErrorText id="robots-txt-error">
          {form.formState.errors.robots_txt?.message ??
            (robotsError
              ? "Could not load robots.txt. Save to try again; your edits are preserved."
              : undefined)}
        </FormErrorText>
      </FormControl>
      {customRobots && (
        <Button
          type="button"
          alignSelf="start"
          disabled={busy}
          loading={resettingRobots}
          onClick={resetRobotsTxt}
        >
          Use default robots.txt
        </Button>
      )}
    </styled.form>
  );
}

async function fetchRobotsTxt(url: string): Promise<string> {
  const response = await fetch(url, { cache: "no-store" });
  if (!response.ok) {
    throw new Error("Could not load robots.txt");
  }

  return response.text();
}
