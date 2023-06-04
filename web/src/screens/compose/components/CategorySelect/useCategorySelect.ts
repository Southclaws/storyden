import { useEffect } from "react";
import { useFormContext } from "react-hook-form";

import { useCategoryList } from "src/api/openapi/categories";

import { ThreadCreate } from "../ComposeForm/useComposeForm";

export function useCategorySelect() {
  const ctx = useFormContext<ThreadCreate>();
  const { data, error } = useCategoryList();

  useEffect(() => {
    if (data?.categories[0]?.id) {
      const defaultCategory = data.categories[0].id;

      ctx.setValue("category", defaultCategory);
    }
  }, [ctx, data]);

  return {
    control: ctx.control,
    fieldError: ctx.formState.errors.category,
    categories: data?.categories,
    error,
  };
}
