package cmd

import (
	"github.com/KimiaMontazeri/cart-tender/internal/cart-tender/cmd/api"
	"github.com/KimiaMontazeri/cart-tender/internal/config"

	"github.com/spf13/cobra"
)

// NewRootCommand creates a new cart-tender root command.
func NewRootCommand() *cobra.Command {
	var root = &cobra.Command{
		Use: "cart-tender",
	}

	cfg := config.New()

	api.Register(root, cfg)

	return root
}
