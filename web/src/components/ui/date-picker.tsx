"use client";

import type { Assign } from "@ark-ui/react";
import {
  DatePicker as ArkDatePicker,
  DatePickerContext,
} from "@ark-ui/react/date-picker";
import { type DatePickerVariantProps, datePicker } from "styled-system/recipes";
import type { ComponentProps, HTMLStyledProps } from "styled-system/types";

import { createStyleContext } from "@/utils/create-style-context";

import { Button } from "./button";
import { IconButton } from "./icon-button";
import { CalendarIcon } from "./icons/Calendar";
import { ChevronLeftIcon, ChevronRightIcon } from "./icons/Chevron";
import { Input as UIInput } from "./input";

const { withProvider, withContext } = createStyleContext(datePicker);

export type RootProviderProps = ComponentProps<typeof RootProvider>;
export const RootProvider = withProvider<
  HTMLDivElement,
  Assign<
    Assign<HTMLStyledProps<"div">, ArkDatePicker.RootProviderBaseProps>,
    DatePickerVariantProps
  >
>(ArkDatePicker.RootProvider, "root");

export type RootProps = ComponentProps<typeof Root>;
export const Root = withProvider<
  HTMLDivElement,
  Assign<
    Assign<HTMLStyledProps<"div">, ArkDatePicker.RootBaseProps>,
    DatePickerVariantProps
  >
>(ArkDatePicker.Root, "root");

export const ClearTrigger = withContext<
  HTMLButtonElement,
  Assign<HTMLStyledProps<"button">, ArkDatePicker.ClearTriggerBaseProps>
>(ArkDatePicker.ClearTrigger, "clearTrigger");

export const Content = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, ArkDatePicker.ContentBaseProps>
>(ArkDatePicker.Content, "content");

export const Control = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, ArkDatePicker.ControlBaseProps>
>(ArkDatePicker.Control, "control");

export const Input = withContext<
  HTMLInputElement,
  Assign<HTMLStyledProps<"input">, ArkDatePicker.InputBaseProps>
>(ArkDatePicker.Input, "input");

export const Label = withContext<
  HTMLLabelElement,
  Assign<HTMLStyledProps<"label">, ArkDatePicker.LabelBaseProps>
>(ArkDatePicker.Label, "label");

export const MonthSelect = withContext<
  HTMLSelectElement,
  Assign<HTMLStyledProps<"select">, ArkDatePicker.MonthSelectBaseProps>
>(ArkDatePicker.MonthSelect, "monthSelect");

export const NextTrigger = withContext<
  HTMLButtonElement,
  Assign<HTMLStyledProps<"button">, ArkDatePicker.NextTriggerBaseProps>
>(ArkDatePicker.NextTrigger, "nextTrigger");

export const Positioner = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, ArkDatePicker.PositionerBaseProps>
>(ArkDatePicker.Positioner, "positioner");

export const PresetTrigger = withContext<
  HTMLButtonElement,
  Assign<HTMLStyledProps<"button">, ArkDatePicker.PresetTriggerBaseProps>
>(ArkDatePicker.PresetTrigger, "presetTrigger");

export const PrevTrigger = withContext<
  HTMLButtonElement,
  Assign<HTMLStyledProps<"button">, ArkDatePicker.PrevTriggerBaseProps>
>(ArkDatePicker.PrevTrigger, "prevTrigger");

export const RangeText = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, ArkDatePicker.RangeTextBaseProps>
>(ArkDatePicker.RangeText, "rangeText");

export const TableBody = withContext<
  HTMLTableSectionElement,
  Assign<HTMLStyledProps<"tbody">, ArkDatePicker.TableBodyBaseProps>
>(ArkDatePicker.TableBody, "tableBody");

export const TableCell = withContext<
  HTMLTableCellElement,
  Assign<HTMLStyledProps<"td">, ArkDatePicker.TableCellBaseProps>
>(ArkDatePicker.TableCell, "tableCell");

export const TableCellTrigger = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, ArkDatePicker.TableCellTriggerBaseProps>
>(ArkDatePicker.TableCellTrigger, "tableCellTrigger");

export const TableHead = withContext<
  HTMLTableSectionElement,
  Assign<HTMLStyledProps<"thead">, ArkDatePicker.TableHeadBaseProps>
>(ArkDatePicker.TableHead, "tableHead");

export const TableHeader = withContext<
  HTMLTableCellElement,
  Assign<HTMLStyledProps<"th">, ArkDatePicker.TableHeaderBaseProps>
>(ArkDatePicker.TableHeader, "tableHeader");

export const Table = withContext<
  HTMLTableElement,
  Assign<HTMLStyledProps<"table">, ArkDatePicker.TableBaseProps>
>(ArkDatePicker.Table, "table");

export const TableRow = withContext<
  HTMLTableRowElement,
  Assign<HTMLStyledProps<"tr">, ArkDatePicker.TableRowBaseProps>
>(ArkDatePicker.TableRow, "tableRow");

export const Trigger = withContext<
  HTMLButtonElement,
  Assign<HTMLStyledProps<"button">, ArkDatePicker.TriggerBaseProps>
>(ArkDatePicker.Trigger, "trigger");

export const ViewControl = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, ArkDatePicker.ViewControlBaseProps>
>(ArkDatePicker.ViewControl, "viewControl");

export const View = withContext<
  HTMLDivElement,
  Assign<HTMLStyledProps<"div">, ArkDatePicker.ViewBaseProps>
>(ArkDatePicker.View, "view");

export const ViewTrigger = withContext<
  HTMLButtonElement,
  Assign<HTMLStyledProps<"button">, ArkDatePicker.ViewTriggerBaseProps>
>(ArkDatePicker.ViewTrigger, "viewTrigger");

export const YearSelect = withContext<
  HTMLSelectElement,
  Assign<HTMLStyledProps<"select">, ArkDatePicker.YearSelectBaseProps>
>(ArkDatePicker.YearSelect, "yearSelect");

export const DatePicker = (props: ArkDatePicker.RootBaseProps) => {
  return (
    <Root
      positioning={{ sameWidth: true }}
      startOfWeek={1}
      selectionMode="single"
      {...props}
    >
      <Control>
        <Input index={0} asChild>
          <UIInput />
        </Input>
        <Trigger asChild>
          <IconButton
            type="button"
            variant="outline"
            aria-label="Open date picker"
          >
            <CalendarIcon />
          </IconButton>
        </Trigger>
      </Control>
      <Positioner>
        <Content>
          <View view="day">
            <DatePickerContext>
              {(api) => (
                <>
                  <ViewControl>
                    <PrevTrigger asChild>
                      <IconButton type="button" variant="ghost" size="sm">
                        <ChevronLeftIcon />
                      </IconButton>
                    </PrevTrigger>
                    <ViewTrigger asChild>
                      <Button type="button" variant="ghost" size="sm">
                        <RangeText />
                      </Button>
                    </ViewTrigger>
                    <NextTrigger asChild>
                      <IconButton type="button" variant="ghost" size="sm">
                        <ChevronRightIcon />
                      </IconButton>
                    </NextTrigger>
                  </ViewControl>
                  <Table>
                    <TableHead>
                      <TableRow>
                        {api.weekDays.map((weekDay, id) => (
                          <TableHeader key={id}>{weekDay.narrow}</TableHeader>
                        ))}
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {api.weeks.map((week, id) => (
                        <TableRow key={id}>
                          {week.map((day, id) => (
                            <TableCell key={id} value={day}>
                              <TableCellTrigger asChild>
                                <IconButton type="button" variant="ghost">
                                  {day.day}
                                </IconButton>
                              </TableCellTrigger>
                            </TableCell>
                          ))}
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </>
              )}
            </DatePickerContext>
          </View>
          <View view="month">
            <DatePickerContext>
              {(api) => (
                <>
                  <ViewControl>
                    <PrevTrigger asChild>
                      <IconButton type="button" variant="ghost" size="sm">
                        <ChevronLeftIcon />
                      </IconButton>
                    </PrevTrigger>
                    <ViewTrigger asChild>
                      <Button type="button" variant="ghost" size="sm">
                        <RangeText />
                      </Button>
                    </ViewTrigger>
                    <NextTrigger asChild>
                      <IconButton type="button" variant="ghost" size="sm">
                        <ChevronRightIcon />
                      </IconButton>
                    </NextTrigger>
                  </ViewControl>
                  <Table>
                    <TableBody>
                      {api
                        .getMonthsGrid({ columns: 4, format: "short" })
                        .map((months, id) => (
                          <TableRow key={id}>
                            {months.map((month, id) => (
                              <TableCell key={id} value={month.value}>
                                <TableCellTrigger asChild>
                                  <Button type="button" variant="ghost">
                                    {month.label}
                                  </Button>
                                </TableCellTrigger>
                              </TableCell>
                            ))}
                          </TableRow>
                        ))}
                    </TableBody>
                  </Table>
                </>
              )}
            </DatePickerContext>
          </View>
          <View view="year">
            <DatePickerContext>
              {(api) => (
                <>
                  <ViewControl>
                    <PrevTrigger asChild>
                      <IconButton type="button" variant="ghost" size="sm">
                        <ChevronLeftIcon />
                      </IconButton>
                    </PrevTrigger>
                    <ViewTrigger asChild>
                      <Button type="button" variant="ghost" size="sm">
                        <RangeText />
                      </Button>
                    </ViewTrigger>
                    <NextTrigger asChild>
                      <IconButton type="button" variant="ghost" size="sm">
                        <ChevronRightIcon />
                      </IconButton>
                    </NextTrigger>
                  </ViewControl>
                  <Table>
                    <TableBody>
                      {api.getYearsGrid({ columns: 4 }).map((years, id) => (
                        <TableRow key={id}>
                          {years.map((year, id) => (
                            <TableCell key={id} value={year.value}>
                              <TableCellTrigger asChild>
                                <Button type="button" variant="ghost">
                                  {year.label}
                                </Button>
                              </TableCellTrigger>
                            </TableCell>
                          ))}
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </>
              )}
            </DatePickerContext>
          </View>
        </Content>
      </Positioner>
    </Root>
  );
};

export type DateRangePickerProps = ArkDatePicker.RootBaseProps & {
  active?: boolean;
  triggerClassName?: string;
  hideInputs?: boolean;
};

export const DateRangePicker = ({
  active,
  triggerClassName,
  hideInputs,
  ...props
}: DateRangePickerProps) => {
  return (
    <Root
      positioning={{ sameWidth: true }}
      startOfWeek={1}
      selectionMode="range"
      {...props}
    >
      <Control>
        {!hideInputs && (
          <>
            <Input index={0} asChild>
              <UIInput size="sm" placeholder="Start date" />
            </Input>
            <Input index={1} asChild>
              <UIInput size="sm" placeholder="End date" />
            </Input>
          </>
        )}
        <Trigger asChild>
          <IconButton
            size="sm"
            type="button"
            variant={active ? "solid" : "subtle"}
            aria-label="Open date picker"
            className={triggerClassName}
          >
            <CalendarIcon />
          </IconButton>
        </Trigger>
      </Control>
      <Positioner>
        <Content>
          <View view="day">
            <DatePickerContext>
              {(api) => (
                <>
                  <ViewControl>
                    <PrevTrigger asChild>
                      <IconButton type="button" variant="ghost" size="sm">
                        <ChevronLeftIcon />
                      </IconButton>
                    </PrevTrigger>
                    <ViewTrigger asChild>
                      <Button type="button" variant="ghost" size="sm">
                        <RangeText />
                      </Button>
                    </ViewTrigger>
                    <NextTrigger asChild>
                      <IconButton type="button" variant="ghost" size="sm">
                        <ChevronRightIcon />
                      </IconButton>
                    </NextTrigger>
                  </ViewControl>
                  <Table>
                    <TableHead>
                      <TableRow>
                        {api.weekDays.map((weekDay, id) => (
                          <TableHeader key={id}>{weekDay.narrow}</TableHeader>
                        ))}
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {api.weeks.map((week, id) => (
                        <TableRow key={id}>
                          {week.map((day, id) => (
                            <TableCell key={id} value={day}>
                              <TableCellTrigger asChild>
                                <IconButton type="button" variant="ghost">
                                  {day.day}
                                </IconButton>
                              </TableCellTrigger>
                            </TableCell>
                          ))}
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </>
              )}
            </DatePickerContext>
          </View>
        </Content>
      </Positioner>
    </Root>
  );
};
