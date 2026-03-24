package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/NeriCarcasci/spar/internal/ai/auth"
	"github.com/NeriCarcasci/spar/internal/config"
)

func RunLogin() {
	cfg, err := config.Load()
	if err != nil {
		PrintError("loading config: %v", err)
		os.Exit(1)
	}

	if auth.ProviderType(cfg.AIProvider) != auth.ProviderOpenAIOAuth {
		PrintError("OAuth login is only available when provider is \"openai-oauth\".")
		fmt.Printf("  %s\n", dimStyle.Render("Current provider: "+cfg.AIProvider))
		fmt.Printf("  %s\n", hintStyle.Render("Run \"spar settings\" to switch to OpenAI OAuth."))
		os.Exit(1)
	}

	if err := config.EnsureDirectories(); err != nil {
		PrintError("creating directories: %v", err)
		os.Exit(1)
	}

	fmt.Println(hintStyle.Render("Opening browser for OpenAI authentication..."))

	client := auth.NewOAuthClient()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := client.Login(ctx); err != nil {
		PrintError("Login failed: %v", err)
		os.Exit(1)
	}

	expires := client.TokenExpiresAt()
	fmt.Println()
	PrintSuccess("Authenticated successfully")
	fmt.Printf("  %s\n", dimStyle.Render("Token expires: "+expires.Format(time.RFC822)))
}

func RunLogout() {
	client := auth.NewOAuthClient()

	if !client.IsAuthenticated() {
		fmt.Println(hintStyle.Render("Not currently authenticated."))
		return
	}

	if err := client.Logout(); err != nil {
		PrintError("Logout failed: %v", err)
		os.Exit(1)
	}

	PrintSuccess("Logged out. AI features disabled until next login.")
}
