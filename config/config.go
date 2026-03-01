package config

var (
	Version      = "1.1.1"
	WireCmd      = "github.com/google/wire/cmd/wire@latest"
	OrbitCmd     = "github.com/go-orbit/orbit@latest"
	RepoBase     = "https://github.com/khy20040121/orbit-layout-base.git"
	RepoAdvanced = "https://github.com/khy20040121/orbit-layout-advanced.git"
	//RepoAdmin     = "https://github.com/khy20040121/orbit-layout-admin.git"
	//RepoChat      = "https://github.com/khy20040121/orbit-layout-chat.git"
	//RepoMCP       = "https://github.com/khy20040121/orbit-layout-mcp.git"
	RunExcludeDir = ".git,.idea,tmp,vendor"
	RunIncludeExt = "go,html,yaml,yml,toml,ini,json,xml,tpl,tmpl"
)
