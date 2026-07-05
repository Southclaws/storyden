---
name: agent-builder
description: Create or improve built-in Storyden Robots by designing focused playbooks, clear tool contracts, high-signal tool outputs, and practical eval scenarios for the current ADK-based Robot implementation.
metadata:
  owner: Storyden
---

# Storyden Agent Builder

Use this skill when creating or improving a built-in Storyden Robot, its playbook, tool descriptions, tool schemas, or tool result shapes.

Built-in Storyden Robots run inside the Storyden ADK integration. They should be self-contained: do not require runtime web access, external skills, or documents unless a configured tool explicitly provides that capability.

## Core Workflow

1. Define the Robot's job.
   - State what user outcomes it owns.
   - State what it must not attempt.
   - Keep the scope narrow enough that the available tools can finish the workflow.

2. Map the workflow to tools.
   - Prefer a few focused tools over many overlapping tools.
   - Treat tool names, descriptions, inputs, outputs, and side effects as part of the Robot prompt.
   - Move deterministic minutia into tools instead of telling the Robot to remember it.

3. Shape tool responses.
   - Return high-signal fields in domain language.
   - Include enough context for the next action.
   - Add `message` and `next_action` fields when the result should steer the Robot.
   - Hide implementation trivia unless it changes what the Robot should do.
   - Make errors actionable and specific.

4. Write the playbook.
   - Define stages: discover, plan, edit/act, validate, deliver.
   - Define when to ask versus act.
   - Define authoritative tools for facts and side effects.
   - Define failure behavior and user-facing language.
   - Keep external references distilled locally; runtime Robots may be network-restricted.

5. Add eval scenarios.
   - Use realistic tasks with multi-step workflows.
   - Include failure cases and repeated attempts.
   - Verify both final outcome and tool path when tool choice matters.
   - Review raw transcripts, not just final summaries.

## Good Robot Instructions

Good instructions:

- Define the Robot's job in one concrete paragraph.
- Define boundaries and forbidden behavior.
- Define workflow stages in the order tools should normally be used.
- Explain when to ask a user versus making a reasonable choice.
- Tie tool names to domain jobs, not internal implementation details.
- Define what success looks like.
- Define failure behavior in user-facing language.
- Keep the Robot from seeing implementation distractions that tools can handle.

Avoid:

- Long generic advice with no tool consequences.
- Instructions that duplicate deterministic validation already enforced by tools.
- Runtime links or references the Robot cannot fetch.
- Tool workflows that require the Robot to remember hidden state.

## Good Tool Descriptions

A good tool description says:

- When to use the tool.
- When not to use it.
- The authoritative source of truth it reads or writes.
- Side effects, especially destructive or external effects.
- Required preconditions.
- Common failure cases and what to do next.
- The next tool that usually follows.

Prefer domain words over implementation words. For example, say "installed supervised plugin" instead of "archive record" unless the archive itself is the user's concern.

## Tool Result Style

Return concise structured context:

- `message`: human-readable summary of what happened.
- `next_action`: what the Robot should do next when that is predictable.
- Domain identifiers only when needed for follow-up tools.
- Names and descriptions alongside IDs when the Robot must choose among objects.
- Truncation flags and counts when output may be incomplete.

For validation tools:

- Group checks by user-relevant readiness areas.
- Avoid surfacing internal build/package stages unless the Robot can act on them.
- Give the first concrete repair action.

For install/update tools:

- Own the full deterministic workflow internally.
- Persist target state early enough that retries update the same object.
- Return whether the object is active and how to inspect logs or runtime behavior.

## Current Storyden Notes

- There is no built-in delegation/handoff system yet. Do not write current Robot playbooks around delegation.
- Built-in Robot instructions and docs must be self-contained.
- Plugin Builder is a managed flow: `plugin_install` owns validation, compile, package, upload/update, and activation.
- For Plugin Builder, keep manifest runtime fields controlled by tools instead of asking the Robot to manage binaries or command names.

## Eval Scenarios

Include scenarios like:

- Create a new object from a vague user outcome and sensible defaults.
- Update an existing object without losing unrelated behavior.
- Retry after a failed side effect and confirm it updates the same target.
- Handle missing configuration without claiming success or exiting early.
- Choose the right tool from overlapping-looking options.
- Recover from a validation error using the tool's `next_action`.
