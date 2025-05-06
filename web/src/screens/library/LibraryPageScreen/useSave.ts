import { omit } from "lodash/fp";

import { PropertyMutation, PropertyType } from "src/api/openapi-schema";

import { handle } from "@/api/client";
import { CoverImageArgs, useLibraryMutation } from "@/lib/library/library";

import { useLibraryPageContext } from "./Context";
import { Form } from "./form";
import { useEditState } from "./useEditState";

type Props = {
  handleUploadCroppedCover: () => Promise<CoverImageArgs | undefined>;
};

export function useSave({ handleUploadCroppedCover }: Props) {
  const { node } = useLibraryPageContext();
  const { revalidate, updateNode } = useLibraryMutation(node);

  const { form } = useLibraryPageContext();
  const { setEditing } = useEditState();

  const handleSubmit = form.handleSubmit(async (payload: Form) => {
    await handle(
      async () => {
        const coverConfig = await handleUploadCroppedCover();

        const isRedirecting = await updateNode(
          node.slug,
          {
            ...payload,
            properties: payload.properties.map((p) => {
              if (p.fid?.startsWith("new_field_")) {
                return omit("fid", p);
              }
              return {
                ...p,
                type: p.type as PropertyType,
              } satisfies PropertyMutation;
            }),
            url: payload.link,
          },
          coverConfig,
        );

        if (!isRedirecting) {
          // NOTE: This modifies the previous URL state, if updateNode received
          // a new slug, it will redirect to the new path automatically. This
          // causes the page to reload before the new slug is pushed to the URL.
          // So to fix this, we only call setEditing if the slug hasn't changed.
          setEditing(false);
        }
      },
      {
        promiseToast: {
          loading: "Saving...",
          success: "Page saved!",
        },
        cleanup: () => revalidate(),
      },
    );
  });

  return {
    handleSubmit,
  };
}
