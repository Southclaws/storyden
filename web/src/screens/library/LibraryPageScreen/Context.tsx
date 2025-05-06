import { zodResolver } from "@hookform/resolvers/zod";
import { PropsWithChildren, createContext, useContext, useMemo } from "react";
import { FormProvider, UseFormReturn, useForm } from "react-hook-form";

import { NodeWithChildren, PropertyType } from "src/api/openapi-schema";

import { Form, FormSchema } from "./form";

type LibraryPageContext = {
  node: NodeWithChildren;
  form: UseFormReturn<Form>;
  defaultFormValues: Form;
};

const Context = createContext<LibraryPageContext | null>(null);

export function useLibraryPageContext() {
  const context = useContext(Context);
  if (!context) {
    throw new Error(
      "useLibraryPageContext must be used within a LibraryPageProvider",
    );
  }

  return context;
}

export type Props = {
  node: NodeWithChildren;
};

export function LibraryPageProvider({
  node,
  children,
}: PropsWithChildren<Props>) {
  const defaultFormValues = useMemo<Form>(
    () =>
      ({
        name: node.name,
        slug: node.slug,
        properties: node.properties.map((p, i) => ({
          fid: p.fid,
          name: p.name ?? `Field ${i}`,
          type: p.type ?? PropertyType.text,
          sort: p.sort,
          value: p.value ?? "",
        })),
        tags: node.tags.map((t) => t.name),
        link: node.link?.url,
        content: node.content,
      }) satisfies Form,
    [node],
  );

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: defaultFormValues,
  });

  return (
    <Context.Provider
      value={{
        node,
        form,
        defaultFormValues,
      }}
    >
      <FormProvider {...form}>
        {/*  */}
        {children}
      </FormProvider>
    </Context.Provider>
  );
}
