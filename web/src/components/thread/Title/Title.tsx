import { Thread } from "src/api/openapi-schema";

import { Heading } from "@/components/ui/heading";
import { Input } from "@/components/ui/input";

import { CategoryPill } from "../CategoryPill";
import { LinkView } from "../LinkView/LinkView";

import { useTitle } from "./useTitle";

export function Title(props: Thread) {
  const { editing, editingTitle, onTitleChange } = useTitle(props);

  return (
    <>
      <div>
        {editing ? (
          <Input value={editingTitle} onChange={onTitleChange} />
        ) : (
          <Heading fontSize="heading.variable.1" fontWeight="bold">
            {props.title}
          </Heading>
        )}
      </div>
      {props.link && <LinkView link={props.link} asset={props.assets?.[0]} />}
      <CategoryPill category={props.category} />
    </>
  );
}
