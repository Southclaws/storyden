"use client";

import { debounce } from "lodash";
import { parseAsArrayOf, parseAsString, useQueryState } from "nuqs";
import { useEffect, useRef, useState } from "react";

import { profileList } from "@/api/openapi-client/profiles";
import {
  MultiSelectPicker,
  MultiSelectPickerItem,
} from "@/components/ui/MultiSelectPicker";
import { deriveError } from "@/utils/error";

export function InvitedByFilter() {
  const [invitedBy, setInvitedBy] = useQueryState(
    "invited_by",
    parseAsArrayOf(parseAsString).withDefault([]),
  );
  const [searchResults, setSearchResults] = useState<MultiSelectPickerItem[]>(
    [],
  );
  const [queryError, setQueryError] = useState<string | null>(null);

  const value = [...invitedBy].map((handle) => ({
    label: handle,
    value: handle,
  }));

  const debouncedSearch = useRef(
    debounce(async (query: string) => {
      try {
        setQueryError(null);

        if (query.length === 0) {
          setSearchResults([]);
          return;
        }

        const result = await profileList({ q: query });
        const items = result.profiles.map((profile) => ({
          label: profile.name,
          value: profile.handle,
        }));
        setSearchResults(items);
      } catch (error) {
        setQueryError(deriveError(error));
        setSearchResults([]);
      }
    }, 300),
  ).current;

  useEffect(() => {
    return () => {
      debouncedSearch.cancel();
    };
  }, []);

  function handleQuery(query: string) {
    if (query.length === 0) {
      setSearchResults([]);
    } else {
      debouncedSearch(query);
    }
  }

  const handleChange = async (items: MultiSelectPickerItem[]) => {
    const values = items.map((item) => item.value);
    await setInvitedBy(values.length > 0 ? values : null);
  };

  return (
    <MultiSelectPicker
      inputPlaceholder="Invited by"
      value={value}
      size="sm"
      triggerProps={{
        width: "full",
        minW: "32",
        flexShrink: "1",
      }}
      onQuery={handleQuery}
      queryResults={searchResults}
      onChange={handleChange}
      queryError={queryError}
    />
  );
}
