package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [shell]",
	Short: "Print shell integration script",
	Long:  "Print shell integration script to enable prompt rendering and environment exports in the current shell.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		shell := strings.ToLower(args[0])
		if shell != "zsh" {
			return fmt.Errorf("unsupported shell %q (supported: zsh)", args[0])
		}

		binPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("resolve executable path: %w", err)
		}

		fmt.Print(renderZshInit(binPath))
		return nil
	},
}

func renderZshInit(binPath string) string {
	var b strings.Builder

	fmt.Fprintf(&b, "typeset -g __XOCTX_BIN=%q\n", binPath)
	b.WriteString(`
setopt PROMPT_SUBST

xoctx_ps1() {
  [[ -n "${XOA_URL:-}" ]] || return 0

  local profiles_dir="${XOCTX_DIR:-$HOME/.xoctx}"
  local current_lab_file="${profiles_dir}/current_lab"
  [[ -f "$current_lab_file" ]] || return 0

  local lab
  lab="$(<"$current_lab_file")"
  [[ -n "$lab" ]] || return 0

  local prefix="${XO_PS1_PREFIX:-(|}"
  local suffix="${XO_PS1_SUFFIX:-)}"

  print -r -- "%F{magenta}${prefix}${lab}${suffix}%f"
}

xo_ps1() {
  xoctx_ps1 "$@"
}

xoctx() {
  local rc lab_name

  case "${1:-}" in
    list|show|current|which|env|init|help|--help|-h)
      "$__XOCTX_BIN" "$@"
      return $?
      ;;
    use|"")
      "$__XOCTX_BIN" "$@"
      rc=$?
      if [[ $rc -eq 0 ]]; then
        lab_name="$("$__XOCTX_BIN" current 2>/dev/null)"
        if [[ -n "$lab_name" && "$lab_name" != "none" ]]; then
          eval "$("$__XOCTX_BIN" env "$lab_name")"
        fi
      fi
      return $rc
      ;;
    *)
      "$__XOCTX_BIN" "$@"
      return $?
      ;;
  esac
}
`)

	return b.String()
}
