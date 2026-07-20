package main

import (
	"context"
	"errors"
	"strings"

	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/spf13/cobra"
)

var (
	searchTerms []string
	searchDeck  string
)

var errSearchTermsRequired = errors.New("search requires at least one term (positional query or --term)")

func collectSearchTerms(cmd *cobra.Command, args []string) ([]string, error) {
	terms, _ := cmd.Flags().GetStringArray("term")
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		terms = append(terms, args[0])
	}

	filtered := make([]string, 0, len(terms))
	for _, term := range terms {
		term = strings.TrimSpace(term)
		if term != "" {
			filtered = append(filtered, term)
		}
	}
	if len(filtered) == 0 {
		return nil, errSearchTermsRequired
	}
	return filtered, nil
}

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search cards across decks by text",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		terms, err := collectSearchTerms(cmd, args)
		if err != nil {
			return err
		}

		deck, _ := cmd.Flags().GetString("deck")

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			results, err := repo.SearchCards(ctx, terms, deck)
			if err != nil {
				return handleRepoError(err)
			}
			return formatter().PrintSearchResults(cmd.OutOrStdout(), results)
		})
	},
}

func init() {
	searchCmd.Flags().StringArrayVar(&searchTerms, "term", nil, "Search term (repeatable; terms are OR-matched)")
	searchCmd.Flags().StringVar(&searchDeck, "deck", "", "Limit search to one deck")
	rootCmd.AddCommand(searchCmd)
}
