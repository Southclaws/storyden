import slugify from "@sindresorhus/slugify";
import { useEffect } from "react";

import { useLibraryPageContext } from "./Context";

export function useAutoSlug() {
  const { form } = useLibraryPageContext();

  const { name } = form.watch();

  useEffect(() => {
    if (!form.getFieldState("slug").isDirty) {
      const autoSlug = slugify(name);
      form.setValue("slug", autoSlug);
    }
  }, [form, name]);
}
