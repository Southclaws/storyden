# robots

Robots, their tools and sessions, and the operational automations that invoke them.

## Navigation hierarchy

- Robots owns Trails in the product navigation. Trails uses `/robots/trails` and does not appear in the configurable site navigation.
- The Robots sidebar includes Trails with the other section destinations. Members only see it when they have `MANAGE_TRAILS`.
- When Trails exist, show the five most recent below the Robots section. The Trail list API returns definitions by most recently updated first, so preserve that order and take the first five.
- Keep recent Trails separate from recent chats. A Trail is an automation definition, while a chat is a Robot session.
- The local tab navigation provides the same Trails destination on viewports where the Robots sidebar is hidden.
- New Trail actions and Trail links stay under `/robots/trails`.

## Storybook references

- `Compositions/Navigation/Sidebar` shows the complete Robots navigation hierarchy, including recent Trails.
- `Screens/Trails/Editor` covers Trail creation and editing inside the Robots product area.
