import { EditableBadge } from "@/components/ui/surface";
import { Grid, HStack, styled } from "@/styled-system/jsx";

import { ColumnDefinition } from "./column";
import {
  DirectoryPropertyProps,
  getDirectoryBadgeColumns,
} from "./directoryCard";

export function DirectoryPropertyBadges({
  node,
  columns,
  values,
  editing,
  onFieldValueChange,
  placement,
}: DirectoryPropertyProps & { placement: ColumnDefinition["placement"] }) {
  const badges = getDirectoryBadgeColumns(columns, values, editing, placement);
  return badges.map((column) => {
    const value = values[column.fid]?.value ?? "";
    return (
      <EditableBadge
        key={column.fid}
        editing={editing}
        value={value}
        inputLabel={`${column.name} for ${node.name}`}
        placeholder={column.name}
        onValueChange={(nextValue) =>
          onFieldValueChange(node.id, column.fid, nextValue)
        }
        size="sm"
        colorPalette="gray"
        title={`${column.name}: ${value}`}
        maxW="full"
        h="auto"
        py="0.5"
      >
        <styled.span srOnly>{column.name}: </styled.span>
        {value}
      </EditableBadge>
    );
  });
}

export function DirectoryBadgeFooter(props: DirectoryPropertyProps) {
  return (
    <Grid
      w="full"
      minW="0"
      gridTemplateColumns="minmax(0, max-content) minmax(0, max-content)"
      alignItems="end"
      justifyContent="space-between"
      gap="2"
    >
      <HStack minW="0" flexWrap="wrap" gap="1">
        <DirectoryPropertyBadges {...props} placement="footer-start" />
      </HStack>
      <HStack minW="0" flexWrap="wrap" justify="end" gap="1">
        <DirectoryPropertyBadges {...props} placement="footer-end" />
      </HStack>
    </Grid>
  );
}
