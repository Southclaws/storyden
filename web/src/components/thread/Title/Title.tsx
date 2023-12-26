import { Thread } from "src/api/openapi/schemas";
import { Input } from "src/theme/components/Input";

import { CategoryPill } from "../CategoryPill";
import { LinkView } from "../LinkView/LinkView";

import { styled } from "@/styled-system/jsx";

import { useTitle } from "./useTitle";

export function Title(props: Thread) {
  const { editing, editingTitle, onTitleChange } = useTitle(props);

  return (
    <>
      <div>
        {editing ? (
          <Input value={editingTitle} onChange={onTitleChange} />
        ) : (
          <styled.h1 fontSize="heading.variable.1" fontWeight="bold">
            {props.title}
          </styled.h1>
        )}
      </div>
      {props.link && <LinkView link={props.link} asset={props.assets?.[0]} />}
      <CategoryPill category={props.category} />
    </>
  );
}
