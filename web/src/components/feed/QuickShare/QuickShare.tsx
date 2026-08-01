import { match } from "ts-pattern";

import { LinkReference } from "@/api/openapi-schema";
import { CategorySelect } from "@/components/category/CategorySelect/CategorySelect";
import { FormComposeField } from "@/components/content/ContentComposer";
import { Button } from "@/components/ui/button";
import { FormErrorText } from "@/components/ui/form-error-text";
import { CreateIcon } from "@/components/ui/icons/Create";
import { Card } from "@/components/ui/rich-card";
import { Spinner } from "@/components/ui/spinner";
import { CardBox, Flex, HStack, WStack } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { getAssetURL } from "@/utils/asset";

import { Props, useQuickShare } from "./useQuickShare";

export function QuickShare(props: Props) {
  const {
    form,
    state: { formRef, hydratedLink, resetKey },
    handlers,
  } = useQuickShare(props);

  // TODO: Render a prompt to sign up to contribute if not logged in.
  if (!props.initialSession) {
    return null;
  }

  return (
    <CardBox bgColor="bg.default">
      <form
        className={lstack({
          gap: "2",
        })}
        ref={formRef}
        onFocus={handlers.handleFocus}
        onSubmit={handlers.handlePost}
      >
        <FormComposeField
          control={form.control}
          name="body"
          placeholder="Share a thought, a link, something cool..."
          resetKey={resetKey}
        />

        <WStack
          w="full"
          justifyContent={
            props.showCategorySelect ? "space-between" : "flex-end"
          }
        >
          {props.showCategorySelect && (
            <HStack alignItems="center">
              <CategorySelect control={form.control} name="category" />

              <FormErrorText>
                {form.formState.errors["category"]?.message}
              </FormErrorText>
            </HStack>
          )}

          <Button
            type="submit"
            size="sm"
            variant="subtle"
            loading={form.formState.isSubmitting}
          >
            <CreateIcon />
            Share
          </Button>
        </WStack>
      </form>

      {match(hydratedLink)
        .with(null, () => null)
        .with("loading", () => <Spinner />)
        .otherwise((link: LinkReference) => (
          <Card
            id={link.slug}
            shape="row"
            title={link.title || "(No site title found)"}
            text={link.description || "(No site description found)"}
            image={getAssetURL(link.primary_image?.path)}
            url={link.url}
          />
        ))}
    </CardBox>
  );
}
