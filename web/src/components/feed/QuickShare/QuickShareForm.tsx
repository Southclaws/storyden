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
import { Box, HStack, WStack } from "@/styled-system/jsx";
import { lstack } from "@/styled-system/patterns";
import { getAssetURL } from "@/utils/asset";

import { Props, useQuickShare } from "./useQuickShare";

export function QuickShareForm(
  props: Props & {
    isActive: boolean;
    focusOnReady: boolean;
    onActivate: () => void;
  },
) {
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
        onSubmit={(event) => {
          if (!props.isActive) {
            event.preventDefault();
            props.onActivate();
            return;
          }
          return handlers.handlePost(event);
        }}
      >
        <Box className="quick-share__composer" w="full" minH="8">
          {props.isActive ? (
            <ContentComposerField
              control={form.control}
              focusOnReady={props.focusOnReady}
              name="body"
              placeholder="Share a thought, a link, something cool..."
              resetKey={resetKey}
            />
          ) : (
            <Button
              className="content-composer content-composer--idle quick-share__trigger"
              type="button"
              variant="ghost"
              aria-label="Quickly create a new thread"
              w="full"
              minH="8"
              px="0"
              alignItems="start"
              justifyContent="start"
              fontWeight="normal"
              fontSize="md"
              color="text.muted"
              textAlign="start"
              whiteSpace="normal"
              onFocus={props.onActivate}
              onClick={props.onActivate}
            >
              Share a thought, a link, something cool...
            </Button>
          )}
        </Box>

        <WStack
          className="quick-share__actions"
          w="full"
          justifyContent={
            props.showCategorySelect ? "space-between" : "flex-end"
          }
        >
          {props.showCategorySelect && (
            <HStack alignItems="center">
              <CategorySelect
                control={form.control}
                name="category"
                initialCategoryList={props.initialCategoryList}
              />

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
            image={getAssetURL(link.primary_image?.path, "m")}
            url={link.url}
          />
        ))}
    </>
  );
}
