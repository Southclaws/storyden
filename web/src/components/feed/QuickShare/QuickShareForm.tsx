import { match } from "ts-pattern";

import { LinkReference } from "@/api/openapi-schema";
import { CategorySelect } from "@/components/category/CategorySelect/CategorySelect";
import { ContentComposerField } from "@/components/content/ContentComposer";
import { Button } from "@/components/ui/button";
import { FormErrorText } from "@/components/ui/form-error-text";
import { CreateIcon } from "@/components/ui/icons/Create";
import { Spinner } from "@/components/ui/spinner";
import { Card } from "@/components/ui/surface";
import { cx } from "@/styled-system/css";
import { HStack, WStack } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { getAssetURL } from "@/utils/asset";

import { Props, useQuickShare } from "./useQuickShare";

export function QuickShareForm(props: Props & { focusOnReady: boolean }) {
  const {
    form,
    state: { hydratedLink, resetKey },
    handlers,
  } = useQuickShare(props);

  return (
    <>
      <form
        className={cx(
          "quick-share__form",
          lstack({
            gap: "2",
          }),
        )}
        onSubmit={handlers.handlePost}
      >
        <ContentComposerField
          control={form.control}
          focusOnReady={props.focusOnReady}
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
    </>
  );
}
