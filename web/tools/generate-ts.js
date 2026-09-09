const fs = require("node:fs");
const path = require("node:path");

const $RefParser = require("@apidevtools/json-schema-ref-parser");
const Ajv = require("ajv");
const { compile } = require("json-schema-to-typescript");
const yaml = require("yaml");

const toolsDir = __dirname;
const schemaPath = path.join(toolsDir, "schema.yaml");
const outputPath = path.join(
  toolsDir,
  "..",
  "src",
  "lib",
  "webmcp",
  "tools.generated.ts",
);

function referencedTypeName(property, fallback) {
  const propertyRef = property?.$ref;
  if (!propertyRef) {
    return fallback;
  }

  const propertyType = propertyRef.split("/").pop();
  return propertyType;
}

function compactSchema(value) {
  if (Array.isArray(value)) {
    return value.map(compactSchema);
  }
  if (value === null || typeof value !== "object") {
    return value;
  }

  return Object.fromEntries(
    Object.entries(value)
      .filter(
        ([key]) => !["$schema", "$id", "title", "definitions"].includes(key),
      )
      .map(([key, child]) => [key, compactSchema(child)]),
  );
}

function getToolEntries(schema, dereferencedSchema) {
  const entries = [];
  const names = new Set();

  for (const [key, definition] of Object.entries(
    dereferencedSchema.definitions ?? {},
  )) {
    if (!definition.title) {
      continue;
    }

    const metadata = definition["x-storyden-tool"];
    if (!metadata?.title || !metadata?.path_prefix || !metadata?.annotations) {
      throw new Error(
        `Tool ${key} must declare x-storyden-tool title, path_prefix, and annotations`,
      );
    }
    for (const annotation of [
      "readOnlyHint",
      "destructiveHint",
      "idempotentHint",
      "openWorldHint",
      "untrustedContentHint",
    ]) {
      if (typeof metadata.annotations[annotation] !== "boolean") {
        throw new Error(
          `Tool ${key} must declare boolean x-storyden-tool.annotations.${annotation}`,
        );
      }
    }
    if (!definition.description) {
      throw new Error(`Tool ${key} must declare a description`);
    }
    if (!definition.properties?.input || !definition.properties?.output) {
      throw new Error(`Tool ${key} must declare input and output schemas`);
    }
    if (names.has(definition.title)) {
      throw new Error(`Duplicate WebMCP tool name: ${definition.title}`);
    }
    names.add(definition.title);

    const toolRef = schema.definitions?.[key]?.$ref;
    if (!toolRef) {
      throw new Error(
        `Tool ${key} must be declared with a $ref in schema.yaml`,
      );
    }
    const toolSchemaPath = path.resolve(toolsDir, toolRef);
    const toolSchema = yaml.parse(fs.readFileSync(toolSchemaPath, "utf8"));

    entries.push({
      key,
      name: definition.title,
      title: metadata.title,
      description: definition.description,
      pathPrefix: metadata.path_prefix,
      annotations: metadata.annotations,
      inputSchema: compactSchema(definition.properties.input),
      outputSchema: compactSchema(definition.properties.output),
      inputType: referencedTypeName(toolSchema.properties.input, `${key}Input`),
      outputType: referencedTypeName(
        toolSchema.properties.output,
        `${key}Output`,
      ),
    });
  }

  return entries;
}

function generateRuntimeTypes(entries) {
  const names = entries.map(({ name }) => JSON.stringify(name)).join(" | ");
  const nameValues = entries.map(({ name }) => JSON.stringify(name)).join(", ");
  const inputMap = entries
    .map(({ name, inputType }) => `  ${JSON.stringify(name)}: ${inputType};`)
    .join("\n");
  const outputMap = entries
    .map(({ name, outputType }) => `  ${JSON.stringify(name)}: ${outputType};`)
    .join("\n");
  const definitions = Object.fromEntries(
    entries.map(
      ({
        name,
        title,
        description,
        pathPrefix,
        annotations,
        inputSchema,
        outputSchema,
      }) => [
        name,
        {
          name,
          title,
          description,
          inputSchema,
          outputSchema,
          annotations,
          pathPrefix,
        },
      ],
    ),
  );

  return `
export type WebMCPToolName = ${names};

export const WEBMCP_TOOL_NAMES = [${nameValues}] as const;

export type WebMCPToolInputMap = {
${inputMap}
};

export type WebMCPToolOutputMap = {
${outputMap}
};

export type WebMCPToolImplementation<K extends WebMCPToolName> = (
  input: WebMCPToolInputMap[K],
) => WebMCPToolOutputMap[K] | Promise<WebMCPToolOutputMap[K]>;

export type WebMCPToolImplementations = {
  [K in WebMCPToolName]: WebMCPToolImplementation<K>;
};

export const WEBMCP_TOOL_DEFINITIONS = ${JSON.stringify(definitions, null, 2)} as const;
`;
}

async function generate() {
  const schema = yaml.parse(fs.readFileSync(schemaPath, "utf8"));
  const dereferencedSchema = await $RefParser.dereference(schemaPath, {
    dereference: { circular: "ignore" },
  });

  const ajv = new Ajv({ strict: false, validateSchema: true });
  ajv.compile(dereferencedSchema);

  const entries = getToolEntries(schema, dereferencedSchema);
  const wrapperSchema = {
    $schema: "http://json-schema.org/draft-07/schema#",
    definitions: schema.definitions,
    type: "object",
    properties: Object.fromEntries(
      Object.keys(schema.definitions).map((key) => [
        key,
        { $ref: `#/definitions/${key}` },
      ]),
    ),
  };

  const types = await compile(wrapperSchema, "WebMCPTools", {
    cwd: toolsDir,
    bannerComment: `/* eslint-disable */
/**
 * This file was automatically generated from web/tools/schema.yaml.
 * DO NOT MODIFY IT BY HAND. Run: pnpm webmcp-schema
 */`,
  });

  fs.mkdirSync(path.dirname(outputPath), { recursive: true });
  fs.writeFileSync(outputPath, types + generateRuntimeTypes(entries));

  console.log(`Generated: ${outputPath}`);
  console.log(`Tool names: ${entries.map(({ name }) => name).join(", ")}`);
}

generate().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
