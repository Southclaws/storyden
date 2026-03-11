const yaml = require("yaml");
const fs = require("fs");
const path = require("path");
const { compile } = require("json-schema-to-typescript");
const $RefParser = require("@apidevtools/json-schema-ref-parser");
const Ajv = require("ajv");

const repoRoot = path.join(__dirname, "..");
const mcpDir = path.join(repoRoot, "api");
const schemaPath = path.join(__dirname, "..", "api", "robots.yaml");
const outputDir = path.join(__dirname, "..", "web", "src", "api");
const outputPathTs = path.join(outputDir, "robots.ts");
const outputPathJson = path.join(outputDir, "robots.json");

async function generate() {
  // Load original schema with $refs intact for TypeScript generation
  const schema = yaml.parse(fs.readFileSync(schemaPath, "utf8"));

  // Dereference the schema once - we'll use it for both tool extraction and JSON output
  const dereferencedSchema = await $RefParser.dereference(schemaPath, {
    dereference: {
      circular: "ignore",
    },
  });

  // Extract tool names from definitions that have a 'title' field (these are the actual tools)
  const toolNames = Object.entries(dereferencedSchema.definitions)
    .filter(([_, def]) => def.title)
    .map(([_, def]) => def.title);

  // Create a wrapper schema that references all definitions
  const wrapperSchema = {
    $schema: "http://json-schema.org/draft-07/schema#",
    definitions: schema.definitions,
    type: "object",
    properties: Object.fromEntries(
      Object.keys(schema.definitions).map((k) => [
        k,
        { $ref: "#/definitions/" + k },
      ]),
    ),
  };

  const types = await compile(wrapperSchema, "MCPTools", {
    cwd: mcpDir,
    bannerComment: `/* eslint-disable */
/**
 * This file was automatically generated from mcp/schema.yaml
 * DO NOT MODIFY IT BY HAND. Run: node mcp/generate-ts.js
 */`,
  });

  // Generate tool name union type
  const toolNameUnion = `export type ToolName = ${toolNames
    .map((n) => `"${n}"`)
    .join(" | ")};

export const TOOL_NAMES = [${toolNames
    .map((n) => `"${n}"`)
    .join(", ")}] as const;

`;

  // Generate tool input/output type mappings
  // The dereferenced schema has the actual structures, but internal refs are still present
  const toolMappings = Object.entries(dereferencedSchema.definitions)
    .filter(
      ([_, def]) =>
        def.title && def.properties?.input && def.properties?.output,
    )
    .map(([key, def]) => {
      // The input/output might still have $ref, or might be inlined
      // Look for the ref first, or extract from the schema
      let inputRef = def.properties.input.$ref?.split("/").pop();
      let outputRef = def.properties.output.$ref?.split("/").pop();

      // If no ref (fully dereferenced), look in the definitions for matching types
      if (!inputRef) {
        inputRef = `Tool${key.replace(/^Tool/, "")}Input`;
      }
      if (!outputRef) {
        outputRef = `Tool${key.replace(/^Tool/, "")}Output`;
      }

      return { name: def.title, inputType: inputRef, outputType: outputRef };
    });

  const toolInputMap = `export type ToolInputMap = {
${toolMappings.map((t) => `  "${t.name}": ${t.inputType};`).join("\n")}
};

`;

  const toolOutputMap = `export type ToolOutputMap = {
${toolMappings.map((t) => `  "${t.name}": ${t.outputType};`).join("\n")}
};
`;

  // Generate Vercel AI SDK compatible tools type
  // Use snake_case tool names to match the actual tool names in the schema
  const vercelToolsType = `export type StorydenTools = {
${toolMappings
  .map(
    (t) =>
      `  "${t.name}": {\n    input: ${t.inputType};\n    output: ${t.outputType};\n  };`,
  )
  .join("\n")}
};
`;

  const output =
    types + "\n" + toolNameUnion + toolInputMap + toolOutputMap + vercelToolsType;

  fs.writeFileSync(outputPathTs, output);

  // Validate the dereferenced schema with AJV to ensure it's valid JSON Schema
  const ajv = new Ajv({
    strict: false,
    validateSchema: true,
  });

  try {
    // Compile the schema - this validates it's a valid JSON Schema
    ajv.compile(dereferencedSchema);
    console.log("✓ Schema is valid JSON Schema Draft 7");
  } catch (error) {
    console.error("✗ Schema validation failed:");
    console.error(error.message);
    if (error.errors) {
      console.error(JSON.stringify(error.errors, null, 2));
    }
    throw error;
  }

  fs.writeFileSync(outputPathJson, JSON.stringify(dereferencedSchema, null, 2));

  console.log(`Generated: ${outputPathTs}`);
  console.log(`Generated: ${outputPathJson}`);
  console.log(`Tool names: ${toolNames.join(", ")}`);
}

generate().catch(console.error);
