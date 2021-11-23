package main

import (
	"fmt"
	"github.com/everstake/solana-pools/config"
	"github.com/everstake/solana-pools/internal/dao"
	"github.com/everstake/solana-pools/internal/delivery/httpserv"
	"github.com/everstake/solana-pools/internal/services"
	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"time"
)

func main() {
	var command = &cobra.Command{
		Use:   "solana pools",
		Short: "start solana pools application",
		Long:  `start solana pools application`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log, _ := zap.NewProduction()
			defer log.Sync() // flushes buffer, if any
			cfg, err := config.NewEnv()
			if err != nil {
				log.Fatal("RUN: config.NewEnv", zap.Error(err))
			}
			d, err := dao.NewDAO(cfg)
			if err != nil {
				log.Fatal("RUN: dao.NewDAO", zap.Error(err))
			}
			s := services.NewService(cfg, d, log)
			cron := gocron.NewScheduler(time.UTC)

			up := false
			cron.Every(time.Minute * 120).Do(func() {
				if !up {
					defer func() {
						up = false
					}()
					up = true
					if err := s.UpdatePools(); err != nil {
						log.Error("UpdatePools", zap.Error(err))
					}
				}
			})
			cron.Every(time.Minute * 30).Do(func() {
				if err := s.UpdatePrice(); err != nil {
					log.Error("UpdatePrice", zap.Error(err))
				}
			})
			cron.Every(time.Minute * 30).Do(func() {
				if err := s.UpdateAPY(); err != nil {
					log.Error("UpdateAPY", zap.Error(err))
				}
			})
			cron.Every(time.Minute * 30).Do(func() {
				if err := s.UpdateValidators(); err != nil {
					log.Error("UpdateValidators", zap.Error(err))
				}
			})
			cron.SetMaxConcurrentJobs(3, gocron.RescheduleMode)
			cron.StartAsync()
			api, err := httpserv.NewAPI(cfg, s, log)
			if err != nil {
				log.Fatal("RUN: httpserv.NewAPI", zap.Error(err))
			}
			return api.Run()
		},
	}

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
