import { type MenuSelectionDetails, Portal } from "@ark-ui/react";
import {
  type AnimateLayoutChanges,
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { FormEvent, useCallback, useState } from "react";

import {
  type CategoryListOKResponse,
  type NodeListResult,
} from "@/api/openapi-schema";
import { CategoryList } from "@/components/category/CategoryList/CategoryList";
import { ModalDrawer } from "@/components/site/Modaldrawer/Modaldrawer";
import { Button } from "@/components/ui/button";
import { FormControl } from "@/components/ui/form-control";
import { FormErrorText } from "@/components/ui/form-error-text";
import { FormLabel } from "@/components/ui/form-label";
import { IconButton } from "@/components/ui/icon-button";
import { AddIcon } from "@/components/ui/icons/Add";
import { DeleteIcon } from "@/components/ui/icons/Delete";
import { EditIcon } from "@/components/ui/icons/Edit";
import { ExternalLinkIcon } from "@/components/ui/icons/ExternalLink";
import { MoreIcon } from "@/components/ui/icons/More";
import { Input } from "@/components/ui/input";
import * as Menu from "@/components/ui/menu";
import { type DragItemNavigationItem } from "@/lib/dragdrop/provider";
import { useNavigationItemEvent } from "@/lib/navigation/events";
import {
  BuiltInNavigationItemName,
  BuiltInNavigationItemTypeSchema,
  type CustomNavigationLink,
  type NavigationConfig,
  type NavigationItem,
  addBuiltInNavigationItem,
  addCustomNavigationLink,
  getNavigationItemKey,
  normaliseNavigationHref,
  removeNavigationItem,
  reorderNavigationItem,
  replaceCustomNavigationLink,
} from "@/lib/settings/navigation";
import { useNavigationMutation } from "@/lib/settings/navigation-client";
import { generateXid } from "@/utils/xid";

import { Anchor } from "../Anchors/Anchor";
import { CollectionsAnchor } from "../Anchors/Collections";
import { LinksAnchor } from "../Anchors/Link";
import { MembersAnchor } from "../Anchors/Members";
import { RolesAnchor } from "../Anchors/Roles";
import { LibraryNavigationTree } from "../LibraryNavigationTree/LibraryNavigationTree";

type Props = {
  navigation: NavigationConfig;
  currentNode?: string;
  initialNodeList?: NodeListResult;
  initialCategoryList?: CategoryListOKResponse;
  isEditing: boolean;
};

const animateNavigationLayoutChanges: AnimateLayoutChanges = ({ isSorting }) =>
  isSorting;

export function NavigationItems(props: Props) {
  const { updateNavigation } = useNavigationMutation();
  const [customLink, setCustomLink] = useState<CustomNavigationLink | null>(
    null,
  );
  const [isCreatingCustomLink, setCreatingCustomLink] = useState(false);
  const itemKeys = props.navigation.items.map(getNavigationItemKey);

  const moveItem = useCallback(
    (activeKey: string, overKey: string) => {
      const updated = reorderNavigationItem(
        props.navigation,
        activeKey,
        overKey,
      );
      if (updated !== props.navigation) {
        void updateNavigation(updated);
      }
    },
    [props.navigation, updateNavigation],
  );

  useNavigationItemEvent("navigation:reorder-item", ({ activeKey, overKey }) =>
    moveItem(activeKey, overKey),
  );

  function addBuiltIn(type: string) {
    const parsed = BuiltInNavigationItemTypeSchema.safeParse(type);
    if (!parsed.success) return;

    const updated = addBuiltInNavigationItem(props.navigation, parsed.data);
    if (updated !== props.navigation) {
      void updateNavigation(updated);
    }
  }

  function addCustomLink(link: Omit<CustomNavigationLink, "type">) {
    const updated = addCustomNavigationLink(props.navigation, link);
    if (updated !== props.navigation) {
      void updateNavigation(updated);
    }
  }

  function updateCustomLink(link: CustomNavigationLink) {
    const updated = replaceCustomNavigationLink(props.navigation, link);
    if (updated !== props.navigation) {
      void updateNavigation(updated);
    }
  }

  function removeItem(item: NavigationItem) {
    const updated = removeNavigationItem(
      props.navigation,
      getNavigationItemKey(item),
    );
    if (updated !== props.navigation) {
      void updateNavigation(updated);
    }
  }

  const renderedItems = props.navigation.items.map((item) => {
    const content = (
      <NavigationItemContent
        item={item}
        currentNode={props.currentNode}
        initialCategoryList={props.initialCategoryList}
        initialNodeList={props.initialNodeList}
      />
    );

    return props.isEditing ? (
      <SortableNavigationItem
        key={getNavigationItemKey(item)}
        item={item}
        onEditCustomLink={setCustomLink}
        onRemove={removeItem}
      >
        {content}
      </SortableNavigationItem>
    ) : (
      <div
        key={getNavigationItemKey(item)}
        className="navigation-editor__item-content"
      >
        {content}
      </div>
    );
  });

  return (
    <>
      <div className="navigation-editor__items" data-editing={props.isEditing}>
        {props.isEditing ? (
          <SortableContext
            items={itemKeys}
            strategy={verticalListSortingStrategy}
          >
            {renderedItems}
          </SortableContext>
        ) : (
          renderedItems
        )}

        {props.isEditing && (
          <CreateNavigationItemMenu
            navigation={props.navigation}
            onAddBuiltIn={addBuiltIn}
            onAddCustomLink={() => setCreatingCustomLink(true)}
          />
        )}
      </div>

      <CustomNavigationLinkModal
        link={customLink}
        isOpen={isCreatingCustomLink || customLink !== null}
        onClose={() => {
          setCreatingCustomLink(false);
          setCustomLink(null);
        }}
        onSave={(link) => {
          if (customLink) {
            updateCustomLink({
              ...link,
              type: "custom-link",
              id: customLink.id,
            });
          } else {
            addCustomLink({ ...link, id: generateXid() });
          }
          setCreatingCustomLink(false);
          setCustomLink(null);
        }}
      />
    </>
  );
}

function SortableNavigationItem({
  children,
  item,
  onEditCustomLink,
  onRemove,
}: {
  children: React.ReactNode;
  item: NavigationItem;
  onEditCustomLink: (link: CustomNavigationLink) => void;
  onRemove: (item: NavigationItem) => void;
}) {
  const key = getNavigationItemKey(item);
  const label =
    item.type === "custom-link"
      ? item.label
      : BuiltInNavigationItemName[item.type];
  const {
    attributes,
    listeners,
    setActivatorNodeRef,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: key,
    animateLayoutChanges: animateNavigationLayoutChanges,
    data: {
      type: "navigation-item",
      itemKey: key,
      label,
    } satisfies DragItemNavigationItem,
  });
  const sortableTransform = transform
    ? { ...transform, scaleX: 1, scaleY: 1 }
    : null;

  return (
    <div
      ref={setNodeRef}
      className="navigation-editor__item"
      data-dragging={isDragging}
      style={{
        transform: CSS.Transform.toString(sortableTransform),
        transition,
      }}
    >
      <button
        {...attributes}
        {...listeners}
        ref={setActivatorNodeRef}
        type="button"
        aria-label={`Move ${label}`}
        className="navigation-editor__drag-overlay"
      />

      <div aria-hidden="true" className="navigation-editor__item-content" inert>
        {children}
      </div>

      <Menu.Root
        lazyMount
        onSelect={({ value }) => {
          if (value === "edit" && item.type === "custom-link") {
            onEditCustomLink(item);
          }
          if (value === "remove") {
            onRemove(item);
          }
        }}
      >
        <Menu.Trigger asChild>
          <IconButton
            aria-label={`Configure ${label}`}
            className="navigation-editor__item-menu"
            variant="ghost"
          >
            <MoreIcon />
          </IconButton>
        </Menu.Trigger>
        <Portal>
          <Menu.Positioner>
            <Menu.Content>
              {item.type === "custom-link" && (
                <Menu.Item value="edit">
                  <EditIcon />
                  Edit link
                </Menu.Item>
              )}
              <Menu.Item value="remove">
                <DeleteIcon />
                Remove
              </Menu.Item>
            </Menu.Content>
          </Menu.Positioner>
        </Portal>
      </Menu.Root>
    </div>
  );
}

function NavigationItemContent({
  item,
  currentNode,
  initialNodeList,
  initialCategoryList,
}: {
  item: NavigationItem;
  currentNode?: string;
  initialNodeList?: NodeListResult;
  initialCategoryList?: CategoryListOKResponse;
}) {
  switch (item.type) {
    case "categories":
      return <CategoryList initialCategoryList={initialCategoryList} />;
    case "library":
      return (
        <LibraryNavigationTree
          initialNodeList={initialNodeList}
          currentNode={currentNode}
          visibility={["draft", "review", "unlisted", "published"]}
        />
      );
    case "collections":
      return <CollectionsAnchor />;
    case "links":
      return <LinksAnchor />;
    case "members":
      return <MembersAnchor />;
    case "roles":
      return <RolesAnchor />;
    case "custom-link":
      return (
        <Anchor
          id={`custom-link-${item.id}`}
          icon={<ExternalLinkIcon />}
          label={item.label}
          route={item.href}
        />
      );
  }
}

function CreateNavigationItemMenu({
  navigation,
  onAddBuiltIn,
  onAddCustomLink,
}: {
  navigation: NavigationConfig;
  onAddBuiltIn: (type: string) => void;
  onAddCustomLink: () => void;
}) {
  const existing = new Set(navigation.items.map((item) => item.type));
  const availableBuiltIns = BuiltInNavigationItemTypeSchema.options.filter(
    (type) => !existing.has(type),
  );

  function handleSelect({ value }: MenuSelectionDetails) {
    if (value === "custom-link") {
      onAddCustomLink();
      return;
    }

    onAddBuiltIn(value);
  }

  return (
    <Menu.Root lazyMount onSelect={handleSelect}>
      <Menu.Trigger asChild>
        <Button variant="outline" width="full">
          <AddIcon />
          Add navigation item
        </Button>
      </Menu.Trigger>
      <Portal>
        <Menu.Positioner>
          <Menu.Content minW="48">
            {availableBuiltIns.map((type) => (
              <Menu.Item key={type} value={type}>
                {BuiltInNavigationItemName[type]}
              </Menu.Item>
            ))}
            <Menu.Item value="custom-link">
              <ExternalLinkIcon />
              Custom link
            </Menu.Item>
          </Menu.Content>
        </Menu.Positioner>
      </Portal>
    </Menu.Root>
  );
}

function CustomNavigationLinkModal({
  link,
  isOpen,
  onClose,
  onSave,
}: {
  link: CustomNavigationLink | null;
  isOpen: boolean;
  onClose: () => void;
  onSave: (link: Pick<CustomNavigationLink, "label" | "href">) => void;
}) {
  const [error, setError] = useState<string>();

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const data = new FormData(event.currentTarget);
    const label = String(data.get("label") ?? "").trim();
    const href = normaliseNavigationHref(String(data.get("href") ?? ""));

    if (!label) {
      setError("Enter a label for the link.");
      return;
    }
    if (!href) {
      setError("Enter an internal path or a valid HTTP(S) URL.");
      return;
    }

    setError(undefined);
    onSave({ label, href });
  }

  return (
    <ModalDrawer
      title={link ? "Edit navigation link" : "Add navigation link"}
      isOpen={isOpen}
      onClose={() => {
        setError(undefined);
        onClose();
      }}
    >
      <form
        key={link?.id ?? (isOpen ? "new-open" : "new-closed")}
        className="navigation-editor__link-form"
        onSubmit={handleSubmit}
      >
        <FormControl>
          <FormLabel id="navigation-link-label-text">Label</FormLabel>
          <Input
            key={`${link?.id ?? "new"}:label`}
            id="navigation-link-label"
            name="label"
            aria-labelledby="navigation-link-label-text"
            defaultValue={link?.label}
            maxLength={80}
            required
          />
        </FormControl>

        <FormControl>
          <FormLabel id="navigation-link-href-text">Destination</FormLabel>
          <Input
            key={`${link?.id ?? "new"}:href`}
            id="navigation-link-href"
            name="href"
            aria-labelledby="navigation-link-href-text"
            defaultValue={link?.href}
            placeholder="/about or example.com"
            required
          />
        </FormControl>

        <FormErrorText>{error}</FormErrorText>

        <Button type="submit">{link ? "Save link" : "Add link"}</Button>
      </form>
    </ModalDrawer>
  );
}
