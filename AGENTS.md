# Repository Instructions

Read info.md files and source code only within the scope explicitly specified by the user. 

Do not proactively explore the repository outside that scope.

Expand the scope only when there is a clear technical dependency or when the requested change cannot be implemented safely with the information available in the specified area.

When expanding the scope, inspect only the minimum additional files or directories required.

Before every non-trivial action, first analyze the task and discuss likely difficulties, nuances, and specification gaps with the user. Do not start that action until the user initiates it with an explicit formal command.



Interfaces and implementations in `entities/` must follow the rules described in `_docs/interface_rules.md`. 

Describe package-specific or implementation-specific behavior in `info.md` inside that package's directory (not in `_docs`).

When running Go tests in this workspace, set writable caches explicitly: `GOCACHE=/tmp/base_go_gocache CCACHE_DIR=/tmp/base_go_ccache go test ...`.

Prefer existing repository patterns over new conventions. Keep changes scoped to the requested package or behavior.



Do not show diffs.
