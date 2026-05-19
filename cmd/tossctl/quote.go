package main

import (
	"fmt"
	"strings"

	"github.com/junghoonkye/tossinvest-cli/internal/domain"
	"github.com/junghoonkye/tossinvest-cli/internal/output"
	"github.com/spf13/cobra"
)

func newQuoteCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quote",
		Short: "Read quote data",
	}

	getCmd := &cobra.Command{
		Use:   "get <symbol or name>",
		Short: "Fetch quote data for a symbol or stock name",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := newAppContext(opts)
			if err != nil {
				return err
			}

			symbol := strings.Join(args, " ")
			quote, err := app.client.GetQuote(cmd.Context(), symbol)
			if err != nil {
				return err
			}

			return output.WriteQuote(cmd.OutOrStdout(), app.format, quote)
		},
	}

	var batchChart bool
	batchCmd := &cobra.Command{
		Use:   "batch <symbol> [symbol...]",
		Short: "Fetch quotes for multiple symbols at once",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := newAppContext(opts)
			if err != nil {
				return err
			}

			var quotes []domain.Quote
			for _, symbol := range args {
				quote, err := app.client.GetQuote(cmd.Context(), symbol)
				if err != nil {
					return err
				}
				quotes = append(quotes, quote)
			}

			if !batchChart {
				return output.WriteQuotes(cmd.OutOrStdout(), app.format, quotes)
			}

			var charts []domain.Chart
			for _, q := range quotes {
				chart, err := app.client.GetChart(cmd.Context(), q.ProductCode, "3m", 30)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: chart unavailable for %s: %v\n", q.Symbol, err)
					charts = append(charts, domain.Chart{})
					continue
				}
				charts = append(charts, chart)
			}
			return output.WriteQuotesWithCharts(cmd.OutOrStdout(), app.format, quotes, charts)
		},
	}
	batchCmd.Flags().BoolVar(&batchChart, "chart", false, "show sparkline chart for each symbol")

	var (
		chartInterval string
		chartCount    int
	)
	chartCmd := &cobra.Command{
		Use:   "chart <symbol or name>",
		Short: "Fetch candle chart for a symbol or stock name",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := newAppContext(opts)
			if err != nil {
				return err
			}

			symbol := strings.Join(args, " ")
			chart, err := app.client.GetChart(cmd.Context(), symbol, chartInterval, chartCount)
			if err != nil {
				return err
			}

			return output.WriteChart(cmd.OutOrStdout(), app.format, chart)
		},
	}
	chartCmd.Flags().StringVar(&chartInterval, "interval", "3m", "candle interval: 1m, 3m, 5m, 10m, 15m, 30m, 60m")
	chartCmd.Flags().IntVar(&chartCount, "count", 30, "number of candles to fetch")

	cmd.AddCommand(getCmd, batchCmd, chartCmd)

	return cmd
}
