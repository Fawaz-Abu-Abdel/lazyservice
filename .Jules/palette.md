## 2025-05-14 - [Tview Navigation & Scrolling]
**Learning:** In `tview.TextView`, using region tags (e.g., `["id"]`) combined with `Highlight()` and `ScrollToHighlight()` provides a smooth, native-feeling scrolling experience for keyboard-driven navigation in long lists.
**Action:** Always wrap list items in region tags when building custom text-based lists in `tview` to enable automatic scrolling and easier programmatic selection.

## 2025-05-14 - [Shortcut Consistency]
**Learning:** Adding visual hints for shortcuts (like `? Help`) without implementing the corresponding keyboard handler leads to user frustration.
**Action:** Ensure all shortcuts mentioned in footers or help menus are fully implemented in the input capture logic.
