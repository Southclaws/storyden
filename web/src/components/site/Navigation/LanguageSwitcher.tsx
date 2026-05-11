"use client";

import { Button } from "@/components/ui/button";
import {
  Content,
  Item,
  Positioner,
  Root as MenuRoot,
  Trigger,
} from "@/components/ui/menu";
import { useI18n } from "@/i18n/provider";

export function LanguageSwitcher() {
  const { locale, setLocale, t } = useI18n();

  const label =
    locale === "zh" ? t("lang.switcher.short.zh") : t("lang.switcher.short.en");

  return (
    <MenuRoot positioning={{ placement: "bottom-end" }}>
      <Trigger asChild>
        <Button variant="ghost" size="sm" aria-label={t("lang.switcher.label")}>
          {label}
        </Button>
      </Trigger>

      <Positioner zIndex="tooltip">
        <Content zIndex="tooltip">
          <Item value="en" onClick={() => setLocale("en")}>
            {t("lang.english")}
          </Item>
          <Item value="zh" onClick={() => setLocale("zh")}>
            {t("lang.chinese")}
          </Item>
        </Content>
      </Positioner>
    </MenuRoot>
  );
}
