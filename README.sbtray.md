# sbtray patches

This is the [sbtray](https://github.com/nziu/sbtray) fork of `fyne.io/systray`.
It tracks upstream `master` and carries three Linux StatusNotifierItem (SNI)
changes on top. Only `systray_unix.go` is modified; the Windows and macOS
backends are untouched.

## 1. Icon by theme name (`SetIconName`, `SetIconThemePath`)

Publishes the SNI `IconName`/`IconThemePath` properties (both writable) and
emits `NewIcon`. Hosts that support themed icons (GNOME, KDE) render the
resolved file — typically an SVG — at the panel's real size, so it stays sharp
on any scale factor. When `IconPixmap` is left empty (`iconPixmapValue` returns
an empty array) the host relies solely on the name.

## 2. No `Activate` in introspection without a handler

GNOME's AppIndicator extension computes
`supportsActivation = !!interfaceInfo.lookup_method('Activate')`. When true, a
single primary click is deferred by the double-click timeout (`double-click`,
default 400 ms) before the menu opens — the "left click is laggy, right click is
not" symptom. `sniIntrospection()` drops `Activate`/`SecondaryActivate` from the
introspection data unless `SetOnTapped`/`SetOnSecondaryTapped` registered a
handler. The methods stay callable (they already return `UnknownMethod` when
unset), so a menu-only tray gets an instant left-click menu.

## 3. Republish the icon once the properties exist

`SetIcon`/`SetIconName` record their value before the D-Bus properties are
exported, so an icon set from `onReady` can fall into the window between
`createPropSpec()` capturing the old value and `instance.props` being assigned.
`republishIcon()`, called at the end of `nativeStart`, pushes the recorded state
so the icon is never lost.

## Upstream

- Base: `fyne.io/systray` upstream `master`
- License: see `LICENSE`
