/*
Copyright © 2024 Jyotirmoy Bandyopadhayaya <hi@b68.dev>
*/
package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cli/go-gh/v2"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-new",
	Short: "A gh-new is a CLI tool to quickly create a new GitHub repository from the command line.",
	Long:  `A gh-new is a CLI tool to quickly create a new GitHub repository from the command line.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			println("🚨 Please provide a repository name")
			println("👉 gh new -h")
			os.Exit(1)
		}

		repo_name := args[0]

		if repo_name == "." {
			// Get the current directory name
			dir, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			repo_name = dir
			// Get the last part of the directory
			repo_name = filepath.Base(dir)
		}

		owner, _ := cmd.Flags().GetString("owner")
		private, _ := cmd.Flags().GetBool("private")
		msg, _ := cmd.Flags().GetString("msg")
		push, _ := cmd.Flags().GetBool("push")

		if owner == "" {
			args := []string{"api", "user", "--jq", ".login"}

			username, _, _ := gh.Exec(args...)

			usernameString := username.String()

			owner = strings.Split(usernameString, "\n")[0]
		}

		if msg == "" {
			msg = "init: Project initialized"
		}

		createRepo(repo_name, owner, private, msg, push)
	},
}

func createRepo(repo_name string, owner string, private bool, msg string, push bool) {
	repoInitCmd := exec.Command("git", "init")
	repoInitCmd.Run()
	println("🔧 Repository initialized")

	repoAddCmd := exec.Command("git", "add", ".")
	repoAddCmd.Run()
	println("📦 Files added to repository")

	repoCommitCmd := exec.Command("git", "commit", "-m", msg)
	repoCommitCmd.Run()
	println("📝 Repository committed")

	if private {
		repoCreateCmd := exec.Command("gh", "repo", "create", owner+"/"+repo_name, "--private", "--source", ".")
		repoCreateCmd.Run()
		println("🚀 Private Repository created on GitHub")
	} else {
		repoCreateCmd := exec.Command("gh", "repo", "create", owner+"/"+repo_name, "--public", "--source", ".")
		repoCreateCmd.Run()
		println("🚀 Public Repository created on GitHub")
	}

	fetchLocalBranchCmd := exec.Command("git", "branch", "--show-current")
	fetchLocalBranchOut, _ := fetchLocalBranchCmd.Output()
	fetchLocalBranchOutStr := strings.TrimSuffix(string(fetchLocalBranchOut), "\n")
	println("🔍 Current branch: " + fetchLocalBranchOutStr)

	if push {
		repoPushCmd := exec.Command("git", "push", "origin", fetchLocalBranchOutStr)
		repoPushCmd.Run()
		println("🚀 Repository pushed to GitHub")
	}

	println("🎉 Repository created successfully")
	println("👉 https://github.com/" + owner + "/" + repo_name)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("private", "p", false, "Set the repository to private")
	rootCmd.Flags().StringP("owner", "o", "", "Set the owner of the repository")
	rootCmd.Flags().StringP("msg", "m", "", "Set the message of initial commit")
	rootCmd.Flags().BoolP("push", "u", false, "Push the repository to GitHub")
}
