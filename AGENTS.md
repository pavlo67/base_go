# Repository Instructions

Before making code or documentation changes, read the relevant files in `_docs`.

Before every non-trivial action, first analyze the task and discuss likely difficulties, nuances, and specification gaps with the user. Do not start that action until the user initiates it with an explicit formal command.

For interface and implementation packages, follow `_docs/interface_rules.md`.

When a change touches an area covered by `_docs`, check the final diff against those rules before finishing. Add or change `_docs` rules only after explicit agreement with the user.

Describe package-specific or implementation-specific behavior in `info.md` inside that package's directory.

When running Go tests in this workspace, set writable caches explicitly: `GOCACHE=/tmp/base_go_gocache CCACHE_DIR=/tmp/base_go_ccache go test ...`.

Prefer existing repository patterns over new conventions. Keep changes scoped to the requested package or behavior.
