import Link from "next/link";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { PropertyList, PropertyListProps } from "@/components/ui/property-list";
import { css, cx } from "@/styled-system/css";

import { ColumnMenu } from "./ColumnMenu";
import { ColumnDefinition } from "./column";
import { DirectoryPropertyProps } from "./directoryCard";

type EditorProps = Pick<
  DirectoryPropertyProps,
  "node" | "onFieldValueChange"
> & {
  fid: string;
  label: string;
  value: string;
  elevated?: boolean;
};

export function DirectoryFieldEditor({
  node,
  fid,
  label,
  value,
  onFieldValueChange,
  elevated = false,
}: EditorProps) {
  const classes: string[] = [];

  if (!fid.startsWith("fixed:")) {
    classes.push(
      css({
        _focus: { outline: "none" },
        _focusVisible: { outline: "none" },
      }),
    );
  }

  if (elevated) {
    classes.push(css({ position: "relative", zIndex: "[5]" }));
  }

  return (
    <Input
      aria-label={`${label} for ${node.name}`}
      size="sm"
      variant="ghost"
      textAlign="inherit"
      className={cx(...classes)}
      value={value}
      minW="0"
      w="full"
      onChange={(event) => onFieldValueChange(node.id, fid, event.target.value)}
    />
  );
}

type PropertyListOptions = Pick<PropertyListProps, "px" | "pb"> & {
  linkValues?: boolean;
  elevated?: boolean;
};

export function DirectoryCardProperties({
  node,
  columns,
  values,
  editing,
  onFieldValueChange,
  renderPropertyLabel,
  linkValues = false,
  elevated = false,
  px,
  pb,
}: DirectoryPropertyProps & PropertyListOptions) {
  const details = columns.filter(
    (column) => (column.placement ?? "details") === "details",
  );
  function renderValue(column: ColumnDefinition) {
    const value = values[column.fid];

    if (editing) {
      return (
        <DirectoryFieldEditor
          node={node}
          fid={column.fid}
          label={column.name}
          value={value?.value ?? ""}
          onFieldValueChange={onFieldValueChange}
          elevated={elevated}
        />
      );
    }

    if (linkValues && value?.href) {
      return <Link href={value.href}>{value.value}</Link>;
    }

    return value?.value;
  }

  const items = details.map((column) => ({
    id: column.fid,
    label: renderPropertyLabel?.(column) ?? column.name,
    value: renderValue(column),
  }));
  return <PropertyList px={px} pb={pb} items={items} />;
}

export function DirectoryPropertyLabel({
  column,
  editing,
}: {
  column: ColumnDefinition;
  editing: boolean;
}) {
  if (!editing) {
    return column.name;
  }

  return (
    <ColumnMenu column={column}>
      <Button
        type="button"
        variant="plain"
        size="sm"
        px="0"
        w="full"
        h="auto"
        minH="6"
        justifyContent="start"
        textAlign="left"
        fontWeight="normal"
        color="text.muted"
        whiteSpace="normal"
        overflowWrap="anywhere"
        aria-label={`Configure ${column.name}`}
      >
        {column.name}
      </Button>
    </ColumnMenu>
  );
}
