import { createListCollection } from "@ark-ui/react";
import { uniqueId } from "lodash";
import { CheckIcon, ChevronsUpDownIcon, XIcon } from "lucide-react";

import { Node } from "@/api/openapi-schema";
import * as Combobox from "@/components/ui/combobox";
import { Combotags } from "@/components/ui/combotags";
import { IconButton } from "@/components/ui/icon-button";
import { Input } from "@/components/ui/input";
import * as TagsInput from "@/components/ui/tags-input";

export type Props = {
  editing: boolean;
  node: Node;
};

export function LibraryPageTagsList(props: Props) {
  function handleChange() {
    //
  }

  const example = [
    //
    "tag1",
    "tag2",
    "tag3",
    "tag" + uniqueId(),
    "really really really long tag name that should break the UI and if it does then fuck michael",
    "tag" + uniqueId(),
    "tag" + uniqueId(),
    "tag" + uniqueId(),
    "tag" + uniqueId(),
    "tag" + uniqueId(),
    "tag" + uniqueId(),
    "tag" + uniqueId(),
    "tag" + uniqueId(),
    "tag" + uniqueId(),
    "tag" + uniqueId(),
    "tag" + uniqueId(),
    "tag" + uniqueId(),
    "tag" + uniqueId(),
    "tag" + uniqueId(),
    "tag" + uniqueId(),
    "tag" + uniqueId(),
    "tag" + uniqueId(),
    "tag" + uniqueId(),
  ];

  async function handleQuery(query: string) {
    return example.filter((item) => item.includes(query));
  }

  return (
    <>
      <Combotags onQuery={handleQuery} />
    </>
  );
}
