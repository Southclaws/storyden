import { produce } from "immer";
import { uniqueId } from "lodash";

import {
  NodeMutableProps,
  NodeWithChildren,
  PropertyType,
} from "@/api/openapi-schema";

import { NodeMetadata } from "./metadata";

export const applyNodeChanges = produce(
  (draft: NodeWithChildren, changes: NodeMutableProps) => {
    (Object.keys(changes) as (keyof typeof changes)[]).forEach((key) => {
      if (changes[key] === undefined || changes[key] == null) {
        return;
      }

      let _exhaustiveCheck: never;

      switch (key) {
        case "name":
          draft[key] = changes[key];
          return;

        case "slug":
          draft[key] = changes[key];
          return;

        case "description":
          draft[key] = changes[key];
          return;

        case "content":
          draft.content = changes[key];
          return;

        case "tags":
          draft.tags = changes[key].map((t) => ({
            name: t,
            colour: "white",
            item_count: 1,
          }));
          return;

        case "properties":
          draft.properties = changes[key].map((p) => ({
            fid: p.fid ?? uniqueId("new_field_"),
            sort: p.sort ?? "0",
            type: p.type ?? PropertyType.text,
            ...p,
          }));
          return;

        case "primary_image_asset_id":
          throw new Error(
            "cannot mutate `primary_image_asset_id` via applyNodeChanges",
          );

        case "hide_child_tree":
          draft.hide_child_tree = changes[key];
          draft.children = changes[key] ? [] : draft.children;
          return;

        case "parent":
          throw new Error("cannot mutate `parent` via applyNodeChanges");

        case "asset_ids":
          throw new Error("cannot mutate `asset_ids` via applyNodeChanges");

        case "asset_sources":
          throw new Error("cannot mutate `asset_sources` via applyNodeChanges");

        case "url":
          throw new Error("cannot mutate `url` via applyNodeChanges");

        case "meta":
          draft.meta = {
            ...draft.meta,
            ...changes[key],
          } satisfies NodeMetadata;
          return;

        default:
          _exhaustiveCheck = key;
          return;
      }
    });
  },
);
