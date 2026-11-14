# Catch Up

Member discussion views at `/catchup`: Unread (the default) and Unresolved. The `view`, `category` (category slug), and `page` query parameters make each view bookmarkable. Changing a view or category returns to the first page. `/unread` permanently redirects here, preserving query parameters.

Both presets use the existing thread-list API, with `ignore_pinned=true` to preserve chronological pagination. Unread uses `read_status=unread`; Unresolved uses `lifecycle_state=open` without restricting read state. State changes patch cached thread fields while preserving displayed membership and totals. Refresh or changing filters retrieves the current matching list.
