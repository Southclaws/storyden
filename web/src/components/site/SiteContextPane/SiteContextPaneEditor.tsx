"use client";

import { Control } from "react-hook-form";

import { ContentFormField } from "@/components/content/ContentComposer/ContentField";

import { Form } from "./useSiteContextPane";

type SiteContextPaneContentFieldProps = {
  control: Control<Form>;
  initialValue?: string;
  placeholder?: string;
};

export function SiteContextPaneContentField({
  control,
  initialValue,
  placeholder,
}: SiteContextPaneContentFieldProps) {
  return (
    <ContentFormField<Form>
      control={control}
      name="content"
      initialValue={initialValue}
      placeholder={placeholder}
    />
  );
}
