"use client";

import { useId } from "react";

import { cx } from "@/styled-system/css";
import { toggleGroup, toggleMatrix } from "@/styled-system/recipes";

import { CheckIcon } from "../icons/Check";
import { RemoveIcon } from "../icons/Remove";
import * as Table from "../table";

export type ToggleMatrixRow = {
  value: string;
  label: string;
  description?: string;
  disabled?: boolean;
};

export type ToggleMatrixColumn = {
  value: string;
  label: string;
  disabled?: boolean;
};

export type ToggleMatrixValue = Record<
  string,
  Record<string, boolean | "mixed">
>;

export type ToggleMatrixValueChangeDetails = {
  value: ToggleMatrixValue;
  row: string;
  column: string;
  pressed: boolean;
};

export type ToggleMatrixProps = {
  label: string;
  rowHeading: string;
  rows: ToggleMatrixRow[];
  columns: ToggleMatrixColumn[];
  value: ToggleMatrixValue;
  onValueChange: (details: ToggleMatrixValueChangeDetails) => void;
  disabled?: boolean;
  className?: string;
  "aria-describedby"?: string;
};

export function ToggleMatrix({
  label,
  rowHeading,
  rows,
  columns,
  value,
  onValueChange,
  disabled,
  className,
  "aria-describedby": describedBy,
}: ToggleMatrixProps) {
  const id = useId();
  const styles = toggleMatrix();
  const toggleStyles = toggleGroup({ size: "lg", variant: "ghost" });

  function toggle(row: string, column: string) {
    const pressed = value[row]?.[column] !== true;

    onValueChange({
      value: { ...value, [row]: { ...value[row], [column]: pressed } },
      row,
      column,
      pressed,
    });
  }

  return (
    <Table.Root
      size="sm"
      aria-label={label}
      aria-describedby={describedBy}
      className={cx(styles.root, className)}
    >
      <Table.Head>
        <Table.Row>
          <Table.Header scope="col" color="text.muted">
            {rowHeading}
          </Table.Header>
          {columns.map((column) => (
            <Table.Header
              key={column.value}
              scope="col"
              className={styles.columnHeader}
            >
              {column.label}
            </Table.Header>
          ))}
        </Table.Row>
      </Table.Head>
      <Table.Body>
        {rows.map((row, index) => (
          <Table.Row key={row.value}>
            <Table.Header scope="row" className={styles.rowHeader}>
              {row.label}
              {row.description && (
                <p
                  id={`${id}-${index}-description`}
                  className={styles.description}
                >
                  {row.description}
                </p>
              )}
            </Table.Header>
            {columns.map((column) => {
              const state = value[row.value]?.[column.value] ?? false;

              return (
                <Table.Cell key={column.value} className={styles.cell}>
                  <button
                    type="button"
                    className={cx(toggleStyles.item, styles.item)}
                    aria-label={`${row.label}: ${column.label}`}
                    aria-describedby={
                      row.description ? `${id}-${index}-description` : undefined
                    }
                    aria-pressed={state}
                    data-state={state === false ? "off" : "on"}
                    disabled={disabled || row.disabled || column.disabled}
                    onClick={() => toggle(row.value, column.value)}
                  >
                    {state === true ? (
                      <CheckIcon aria-hidden="true" />
                    ) : (
                      <RemoveIcon aria-hidden="true" />
                    )}
                  </button>
                </Table.Cell>
              );
            })}
          </Table.Row>
        ))}
      </Table.Body>
    </Table.Root>
  );
}
