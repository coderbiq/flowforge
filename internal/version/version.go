package version

import "runtime/debug"

// injected 由 GoReleaser 或 Makefile 通过 -ldflags 注入。
// 必须保持 var（非 const），linker 才能改写。
var injected = "dev"

// Version 是最终解析出的版本字符串，包初始化时计算一次。
// 优先级：ldflags 注入 > go install @version 的 BuildInfo > "dev"
var Version = resolve(injected)

func resolve(ldflagsVal string) string {
	if ldflagsVal != "" && ldflagsVal != "dev" {
		return ldflagsVal
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}
