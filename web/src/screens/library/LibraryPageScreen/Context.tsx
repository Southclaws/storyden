import { zodResolver } from "@hookform/resolvers/zod";
import { PropsWithChildren, createContext, useContext, useMemo } from "react";
import { FormProvider, UseFormReturn, useForm } from "react-hook-form";

import { NodeWithChildren, PropertyType } from "src/api/openapi-schema";

import { WithMetadata, hydrateNode } from "@/lib/library/metadata";

import { Form, FormSchema } from "./form";

type LibraryPageContext = {
  node: WithMetadata<NodeWithChildren>;
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
  const nodeWithMeta = hydrateNode(node);

  const defaultFormValues = useMemo<Form>(
    () =>
      ({
        name: nodeWithMeta.name,
        slug: nodeWithMeta.slug,
        properties: nodeWithMeta.properties.map((p, i) => ({
          fid: p.fid,
          name: p.name ?? `Field ${i}`,
          type: p.type ?? PropertyType.text,
          sort: p.sort,
          value: p.value ?? "",
        })),
        childPropertySchema: nodeWithMeta.child_property_schema,
        tags: nodeWithMeta.tags.map((t) => t.name),
        link: nodeWithMeta.link?.url,
        content: nodeWithMeta.content,
        meta: nodeWithMeta.meta,
      }) satisfies Form,
    [nodeWithMeta],
  );

  const form = useForm<Form>({
    resolver: zodResolver(FormSchema),
    defaultValues: defaultFormValues,
  });

  return (
    <Context.Provider
      value={{
        node: nodeWithMeta,
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
