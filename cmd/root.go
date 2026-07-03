package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/gCyrille/xoctx/internal/lab"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "xoctx",
	Short: "Manage XO contexts",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(whichCmd)
	rootCmd.AddCommand(envCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.RunE = rootDefault
}

var useCmd = &cobra.Command{
	Use:   "use [lab]",
	Short: "Switch to a lab environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		if len(args) > 0 {
			name = args[0]
		} else {
			labs, err := lab.ListLabs()
			if err != nil {
				return err
			}
			if len(labs) == 0 {
				return fmt.Errorf("no labs found")
			}
			opts := make([]huh.Option[string], len(labs))
			for i, l := range labs {
				opts[i] = huh.Option[string]{Key: l, Value: l}
			}
			selectField := huh.NewSelect[string]().
				// Title("Select a lab").
				Options(opts...).
				Value(&name)
			group := huh.NewGroup(selectField)
			group.Title("Select a lab")
			if err := huh.NewForm(group).
				WithTheme(LabTheme).Run(); err != nil {
				return err
			}
		}

		profile, err := lab.LoadProfile(name)
		if err != nil {
			return err
		}

		if err := lab.SetCurrentLab(name); err != nil {
			return err
		}

		greenCheck := lipgloss.NewStyle().Foreground(lipgloss.Color("108")).Render("✓")
		output := fmt.Sprintf("%s Switched to lab '%s'\n%s", greenCheck, name, RenderSummary(profile))
		fmt.Print(padStyle.Render(output))
		return nil
	},
}

var showCmd = &cobra.Command{
	Use:   "show [lab]",
	Short: "Show a lab's environment variables",
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		if len(args) > 0 {
			name = args[0]
		} else {
			current, err := lab.CurrentLab()
			if err != nil {
				return err
			}
			if current == "" {
				return fmt.Errorf("no current lab set, pass a lab name as argument")
			}
			name = current
		}

		profile, err := lab.LoadProfile(name)
		if err != nil {
			return err
		}

		fmt.Print(padStyle.Render(RenderSummary(profile)))
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available labs",
	RunE: func(cmd *cobra.Command, args []string) error {
		labs, err := lab.ListLabs()
		if err != nil {
			return err
		}
		current, _ := lab.CurrentLab()
		fmt.Print(padStyle.Render(RenderList(labs, current)))
		return nil
	},
}

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the current lab",
	RunE: func(cmd *cobra.Command, args []string) error {
		current, err := lab.CurrentLab()
		if err != nil {
			return err
		}
		if current == "" {
			fmt.Println("none")
		} else {
			fmt.Println(current)
		}
		return nil
	},
}

var whichCmd = &cobra.Command{
	Use:   "which",
	Short: "Show current lab summary (no reload)",
	RunE: func(cmd *cobra.Command, args []string) error {
		current, err := lab.CurrentLab()
		if err != nil {
			return err
		}
		if current == "" {
			fmt.Println("No lab currently active")
			fmt.Println()
			fmt.Println("  Switch to a lab : xoctx use <name>")
			fmt.Println("  List all labs   : xoctx list")
			return nil
		}

		profile, err := lab.LoadProfile(current)
		if err != nil {
			fmt.Printf("Current lab '%s' no longer exists\n", current)
			return nil
		}

		fmt.Print(padStyle.Render(RenderSummary(profile)))
		return nil
	},
}

var envCmd = &cobra.Command{
	Use:   "env [lab]",
	Short: "Print export statements for a lab's environment variables",
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		if len(args) > 0 {
			name = args[0]
		} else {
			current, err := lab.CurrentLab()
			if err != nil {
				return err
			}
			if current == "" {
				return fmt.Errorf("no current lab set, pass a lab name as argument")
			}
			name = current
		}

		profile, err := lab.LoadProfile(name)
		if err != nil {
			return err
		}

		keys := make([]string, 0, len(profile.Vars))
		for k := range profile.Vars {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := profile.Vars[k]
			if strings.ContainsAny(v, " \t'\"") {
				fmt.Printf("export %s='%s'\n", k, strings.ReplaceAll(v, "'", "'\\''"))
			} else {
				fmt.Printf("export %s=%s\n", k, v)
			}
		}
		return nil
	},
}

func rootDefault(cmd *cobra.Command, args []string) error {
	current, err := lab.CurrentLab()
	if err != nil {
		return err
	}
	if current == "" {
		fmt.Println("No lab currently active")
		fmt.Println()
		fmt.Println("  Switch to a lab : xoctx use <name>")
		fmt.Println("  List all labs   : xoctx list")
		return nil
	}

	profile, err := lab.LoadProfile(current)
	if err != nil {
		fmt.Printf("Current lab '%s' no longer exists\n", current)
		return nil
	}

	greenCheck := lipgloss.NewStyle().Foreground(lipgloss.Color("108")).Render("✓")
	fmt.Printf("%s Reloaded lab '%s'\n", greenCheck, current)
	fmt.Print(padStyle.Render(RenderSummary(profile)))
	return nil
}
