package main

import "github.com/spf13/cobra"

var deckCmd = &cobra.Command{
	Use:   "deck",
	Short: "Manage flashcard decks",
}

func init() {
	rootCmd.AddCommand(deckCmd)
}
