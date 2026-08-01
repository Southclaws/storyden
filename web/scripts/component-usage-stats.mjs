#!/usr/bin/env node
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const webRoot = path.resolve(__dirname, "..");
const srcDir = path.join(webRoot, "src");
const uiDir = path.join(srcDir, "components", "ui");

const includeTests = process.argv.includes("--include-tests");
const json = process.argv.includes("--json");

const uiComponents = fs
  .readdirSync(uiDir, { withFileTypes: true })
  .filter((entry) => entry.isDirectory())
  .map((entry) => entry.name)
  .sort();

const stats = new Map(
  uiComponents.map((name) => [
    name,
    {
      component: name,
      files: new Set(),
      imports: 0,
      references: 0,
      symbols: new Set(),
    },
  ]),
);

for (const file of walkSourceFiles(srcDir)) {
  if (isInside(file, uiDir)) continue;
  if (!includeTests && isTestOrStory(file)) continue;

  const source = fs.readFileSync(file, "utf8");
  const imports = parseImports(source);

  for (const importDeclaration of imports) {
    const component = resolveUiComponent(file, importDeclaration.source);
    if (!component) continue;

    const componentStats = stats.get(component);
    if (!componentStats) continue;

    const importedSymbols = parseImportedSymbols(importDeclaration.clause);
    const body = source.replace(importDeclaration.full, "");

    componentStats.files.add(relative(file));
    componentStats.imports += 1;

    for (const symbol of importedSymbols) {
      componentStats.symbols.add(symbol.imported);
      componentStats.references += countReferences(body, symbol.local);
    }
  }
}

const rows = [...stats.values()]
  .map((row) => ({
    component: row.component,
    files: row.files.size,
    imports: row.imports,
    references: row.references,
    symbols: [...row.symbols].sort(),
    usedBy: [...row.files].sort(),
  }))
  .sort((a, b) => {
    return (
      b.references - a.references ||
      b.files - a.files ||
      a.component.localeCompare(b.component)
    );
  });

if (json) {
  console.log(JSON.stringify(rows, null, 2));
} else {
  printTable(rows);
}

function walkSourceFiles(root) {
  const files = [];
  const entries = fs.readdirSync(root, { withFileTypes: true });

  for (const entry of entries) {
    const fullPath = path.join(root, entry.name);

    if (entry.isDirectory()) {
      files.push(...walkSourceFiles(fullPath));
      continue;
    }

    if (entry.isFile() && /\.(ts|tsx)$/.test(entry.name)) {
      files.push(fullPath);
    }
  }

  return files;
}

function parseImports(source) {
  const imports = [];
  const importPattern =
    /import\s+(?:type\s+)?([\s\S]*?)\s+from\s+["']([^"']+)["']\s*;?/g;

  for (const match of source.matchAll(importPattern)) {
    imports.push({
      full: match[0],
      clause: match[1],
      source: match[2],
    });
  }

  return imports;
}

function parseImportedSymbols(clause) {
  const symbols = [];

  const namespaceMatch = clause.match(/\*\s+as\s+([A-Za-z_$][\w$]*)/);
  if (namespaceMatch) {
    symbols.push({
      imported: "*",
      local: namespaceMatch[1],
    });
  }

  const namedMatch = clause.match(/\{([\s\S]*?)\}/);
  if (namedMatch) {
    for (const rawPart of namedMatch[1].split(",")) {
      const part = rawPart.trim().replace(/^type\s+/, "");
      if (!part) continue;

      const aliasMatch = part.match(
        /^([A-Za-z_$][\w$]*)\s+as\s+([A-Za-z_$][\w$]*)$/,
      );
      if (aliasMatch) {
        symbols.push({
          imported: aliasMatch[1],
          local: aliasMatch[2],
        });
        continue;
      }

      const nameMatch = part.match(/^([A-Za-z_$][\w$]*)$/);
      if (nameMatch) {
        symbols.push({
          imported: nameMatch[1],
          local: nameMatch[1],
        });
      }
    }
  }

  const defaultClause = clause.split(",")[0]?.trim();
  if (
    defaultClause &&
    !defaultClause.startsWith("{") &&
    !defaultClause.startsWith("*") &&
    /^[A-Za-z_$][\w$]*$/.test(defaultClause)
  ) {
    symbols.push({
      imported: "default",
      local: defaultClause,
    });
  }

  return symbols;
}

function resolveUiComponent(importingFile, source) {
  const aliasPrefix = "@/components/ui/";

  if (source.startsWith(aliasPrefix)) {
    return source.slice(aliasPrefix.length).split("/")[0];
  }

  if (!source.startsWith(".")) return null;

  const resolved = path.resolve(path.dirname(importingFile), source);
  if (!isInside(resolved, uiDir)) return null;

  const relativeToUi = path.relative(uiDir, resolved);
  return relativeToUi.split(path.sep)[0] || null;
}

function countReferences(source, localName) {
  if (!localName) return 0;

  const escaped = escapeRegExp(localName);
  const references = source.match(new RegExp(`\\b${escaped}\\b`, "g"));

  return references?.length ?? 0;
}

function printTable(rows) {
  const headers = ["component", "files", "imports", "refs", "symbols"];
  const table = rows.map((row) => [
    row.component,
    String(row.files),
    String(row.imports),
    String(row.references),
    row.symbols.slice(0, 8).join(", ") +
      (row.symbols.length > 8 ? `, +${row.symbols.length - 8}` : ""),
  ]);

  const widths = headers.map((header, index) =>
    Math.max(header.length, ...table.map((row) => row[index].length)),
  );

  console.log(formatRow(headers, widths));
  console.log(widths.map((width) => "-".repeat(width)).join("  "));

  for (const row of table) {
    console.log(formatRow(row, widths));
  }

  console.log(
    `\nScanned ${relative(srcDir)} outside ${relative(uiDir)}${
      includeTests ? "" : " excluding tests/stories"
    }. Use --json for file-level detail or --include-tests to include tests/stories.`,
  );
}

function formatRow(row, widths) {
  return row
    .map((cell, index) => {
      const value =
        index === 0 || index === 4
          ? cell.padEnd(widths[index])
          : cell.padStart(widths[index]);
      return value;
    })
    .join("  ");
}

function isInside(candidate, parent) {
  const relativePath = path.relative(parent, candidate);
  return (
    relativePath === "" ||
    (!relativePath.startsWith("..") && !path.isAbsolute(relativePath))
  );
}

function isTestOrStory(file) {
  return /\.(test|stories)\.(ts|tsx)$/.test(file);
}

function relative(file) {
  return path.relative(webRoot, file);
}

function escapeRegExp(value) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}
