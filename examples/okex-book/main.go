package main

import (
	"context"
	"os"
	"strings"

	"github.com/c9s/bbgo/pkg/exchange/okex/okexapi"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.PersistentFlags().String("okex-api-key", "", "okex api key")
	rootCmd.PersistentFlags().String("okex-api-secret", "", "okex api secret")
	rootCmd.PersistentFlags().String("okex-api-passphrase", "", "okex api secret")
	rootCmd.PersistentFlags().String("symbol", "BNBUSDT", "symbol")
}

var rootCmd = &cobra.Command{
	Use:   "okex-book",
	Short: "okex book",

	// SilenceUsage is an option to silence usage when an error occurs.
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		symbol := viper.GetString("symbol")
		if len(symbol) == 0 {
			return errors.New("empty symbol")
		}

		key, secret, passphrase := viper.GetString("okex-api-key"),
			viper.GetString("okex-api-secret"),
			viper.GetString("okex-api-passphrase")
		if len(key) == 0 || len(secret) == 0 {
			return errors.New("empty key, secret or passphrase")
		}

		client := okexapi.NewClient()
		client.Auth(key, secret, passphrase)

		log.Infof("balances:")
		balanceSummaryList, err := client.Balances()
		if err != nil {
			return err
		}

		for _, balanceSummary := range balanceSummaryList {
			log.Infof("%+v", balanceSummary)
		}

		_ = ctx
		// cmdutil.WaitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
		return nil
	},
}

func main() {
	if _, err := os.Stat(".env.local"); err == nil {
		if err := godotenv.Load(".env.local"); err != nil {
			log.Fatal(err)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.WithError(err).Error("bind pflags error")
	}

	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		log.WithError(err).Error("cmd error")
	}
}
