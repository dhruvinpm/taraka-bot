package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dhruvinpm/taraka-bot/core"
	localmail "github.com/dhruvinpm/taraka-bot/gmail"
	"github.com/dhruvinpm/taraka-bot/telegram"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "taraka-bot",
	Short: "TarakaBot - AI-powered export lead generation agent",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the TarakaBot",
	Run:   runStart,
}

var huntCmd = &cobra.Command{
	Use:   "hunt",
	Short: "Run a one-time hunt",
	Run:   runHunt,
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show statistics",
	Run:   runStats,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "config file")
	rootCmd.AddCommand(startCmd, huntCmd, statsCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("no config file found: %v", err)
	}
}

func buildConfig(ctx context.Context) (core.Config, error) {
	return core.Config{
		DBPath:           viper.GetString("database.path"),
		GeminiAPIKey:     viper.GetString("gemini.api_key"),
		GeminiModel:      viper.GetString("gemini.model"),
		GeminiDailyLimit: viper.GetInt("gemini.daily_limit"),
		LocalLLMURL:      viper.GetString("llm.local_url"),
		LocalLLMModel:    viper.GetString("llm.local_model"),
		SMTPConfig: localmail.SMTPConfig{
			Host:     viper.GetString("gmail.smtp_host"),
			Port:     viper.GetInt("gmail.smtp_port"),
			Username: viper.GetString("gmail.address"),
			Password: viper.GetString("gmail.app_password"),
			From:     viper.GetString("gmail.address"),
		},
		GmailAddress:     viper.GetString("gmail.address"),
		GmailCredentials: viper.GetString("gmail.oauth_credentials_file"),
		GmailToken:       viper.GetString("gmail.oauth_token_file"),
		WarmupStart:      viper.GetInt("warmup.start_per_day"),
		WarmupInc:        viper.GetInt("warmup.increment_per_week"),
		WarmupMax:        viper.GetInt("warmup.max_per_day"),
		VaultPath:        viper.GetString("identity.vault_path"),
		VaultKey:         viper.GetString("identity.vault_key"),
		DefaultCountry:   "UAE",
		DefaultNiche:     "food importer",
	}, nil
}

func runStart(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := buildConfig(ctx)
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	engine, err := core.New(ctx, cfg)
	if err != nil {
		log.Fatalf("engine init error: %v", err)
	}
	defer engine.Stop()

	engine.Start(ctx)

	token := viper.GetString("telegram.token")
	adminChatID := viper.GetInt64("telegram.admin_chat_id")

	if token != "" {
		bot, err := telegram.New(token, adminChatID, engine, engine.GetStore())
		if err != nil {
			log.Fatalf("telegram bot error: %v", err)
		}
		go bot.Start()
		defer bot.Stop()
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Println("shutting down...")
}

func runHunt(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	cfg, err := buildConfig(ctx)
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	engine, err := core.New(ctx, cfg)
	if err != nil {
		log.Fatalf("engine init error: %v", err)
	}
	defer engine.Stop()

	country := ""
	niche := ""
	if len(args) >= 1 {
		country = args[0]
	}
	if len(args) >= 2 {
		niche = args[1]
	}

	summary, err := engine.RunHunt(ctx, country, niche)
	if err != nil {
		log.Fatalf("hunt error: %v", err)
	}
	fmt.Printf("Hunt complete! Found %d leads (%d with email)\n", summary.TotalLeads, summary.EmailCount)
}

func runStats(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	cfg, err := buildConfig(ctx)
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	engine, err := core.New(ctx, cfg)
	if err != nil {
		log.Fatalf("engine init error: %v", err)
	}
	defer engine.Stop()

	stats, err := engine.GetStore().GetStats()
	if err != nil {
		log.Fatalf("stats error: %v", err)
	}
	fmt.Println("Lead Statistics:")
	total := 0
	for status, count := range stats {
		fmt.Printf("  %s: %d\n", status, count)
		total += count
	}
	fmt.Printf("Total: %d\n", total)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
