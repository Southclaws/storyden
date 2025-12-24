import { CollectionItem } from "@ark-ui/react";
import { useState } from "react";
import {
  Control,
  Controller,
  ControllerProps,
  FieldValues,
} from "react-hook-form";

import { handle } from "@/api/client";
import { categoryList } from "@/api/openapi-client/categories";
import { profileList } from "@/api/openapi-client/profiles";
import { tagList } from "@/api/openapi-client/tags";
import {
  Account,
  PublicProfile,
  PublicProfileList,
  TagNameList,
} from "@/api/openapi-schema";
import { Combotags } from "@/components/ui/combotags";
import { HStack } from "@/styled-system/jsx";

import { Form } from "./useSearch";

type Props = {
  control: Control<Form>;
};

export function AdditionalFilters({ control }: Props) {
  async function handleQueryTags(q: string): Promise<TagNameList> {
    const tags =
      (await handle(async () => {
        const { tags } = await tagList({ q });
        return tags.map((t) => t.name);
      })) ?? [];

    return tags;
  }

  async function handleQueryCategories(q: string): Promise<string[]> {
    const categories =
      (await handle(async () => {
        const { categories } = await categoryList();
        return categories
          .map((t) => t.name)
          .filter((name) => name.toLowerCase().includes(q.toLowerCase()));
      })) ?? [];

    return categories;
  }

  async function handleQueryAuthors(q: string): Promise<string[]> {
    const profiles =
      (await handle(async () => {
        const { profiles } = await profileList({ q });
        return profiles.map((v) => v.handle);
      })) ?? [];

    return profiles;
  }

  return (
    <HStack>
      <Controller<Form>
        name="authors"
        control={control}
        render={({ formState, field }) => {
          async function handleChange(v: string[]) {
            field.onChange(v);
          }
          return (
            <Combotags onQuery={handleQueryAuthors} onChange={handleChange} />
          );
        }}
      />
      <Controller<Form>
        name="categories"
        control={control}
        render={({ formState, field }) => {
          async function handleChange(v: string[]) {
            field.onChange(v);
          }
          return (
            <Combotags
              onQuery={handleQueryCategories}
              onChange={handleChange}
            />
          );
        }}
      />
      <Controller<Form>
        name="tags"
        control={control}
        render={({ fieldState, formState, field }) => {
          async function handleChange(v: string[]) {
            field.onChange(v);
          }
          return (
            <Combotags onQuery={handleQueryTags} onChange={handleChange} />
          );
        }}
      />
    </HStack>
  );
}
