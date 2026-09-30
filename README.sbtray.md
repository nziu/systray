# sbtray patches

This is the [sbtray](https://github.com/nziu/sbtray) fork of `fyne.io/systray`.
It carries four Linux StatusNotifierItem (SNI) changes on top of a pinned
upstream base (see Upstream). Only `systray_unix.go` is modified; the Windows
and macOS backends are untouched.

## 1. Icon by theme name (`SetIconName`, `SetIconThemePath`)

Publishes the SNI `IconName`/`IconThemePath` properties (both writable) and
emits `NewIcon`. Hosts that support themed icons (GNOME, KDE) render the
resolved file — typically an SVG — at the panel's real size, so it stays sharp
on any scale factor. When `IconPixmap` is left empty (`iconPixmaps` returns an
empty array) the host relies solely on the name.

## 2. Multi-size `IconPixmap` (`SetIconPixmaps`)

`SetIconPixmaps([][]byte)` publishes several pre-rendered images as the SNI
`IconPixmap` array at once. Hosts pick the smallest entry that is at least the
requested physical size (GNOME's AppIndicator extension does this in
`pixmapsUtils.getBestPixmap`), so supplying the common panel sizes lets them
render at an exact size instead of rescaling a single large pixmap — which is
what looks blurry. It takes precedence over `SetIcon`.

## 3. No `Activate` in introspection without a handler

GNOME's AppIndicator extension computes
`supportsActivation = !!interfaceInfo.lookup_method('Activate')`. When true, a
single primary click is deferred by the double-click timeout (`double-click`,
default 400 ms) before the menu opens — the "left click is laggy, right click is
not" symptom. `sniIntrospection()` drops `Activate`/`SecondaryActivate` from the
introspection data unless `SetOnTapped`/`SetOnSecondaryTapped` registered a
handler. The methods stay callable (they already return `UnknownMethod` when
unset), so a menu-only tray gets an instant left-click menu.

## 4. Republish the icon once the properties exist

`SetIcon`/`SetIconName` record their value before the D-Bus properties are
exported, so an icon set from `onReady` can fall into the window between
`createPropSpec()` capturing the old value and `instance.props` being assigned.
`republishIcon()`, called at the end of `nativeStart`, pushes the recorded state
so the icon is never lost.

## Upstream

- Base: `fyne.io/systray` @ `f60f01b` (2026-08-14, after the v1.12.2 release)
- Patch commits:
  - `528cad2` `SetIconName` / `SetIconThemePath` (section 1)
  - `5738b4c` hide `Activate` without a handler; `republishIcon` (sections 3-4)
  - `673211f` `SetIconPixmaps` (section 2)
  - `44d8710` fix ARGB channel order in `argbForImage`
- Versioning: this fork uses its own `v0.x.y` tags, decoupled from the upstream
  version. `v0.1.0` is the first tag.
- License: see `LICENSE`
