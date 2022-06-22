package services

import (
	"context"
	"github.com/everstake/solana-pools/internal/dao/dmodels"
	"github.com/everstake/solana-pools/internal/dao/postgres"
	"time"
)

func (s Imp) GetAvgSlotTimeMS() (float64, error) {
	arr, err := s.DAO.GetSlotTime(&postgres.SlotTimeCondition{
		Limit:  24,
		Offset: 0,
	})
	if err != nil {
		return 0, err
	}

	var sum float64
	var count int
	for _, f := range arr {
		if f.SlotTime != 0 {
			count++
		}
		sum += f.SlotTime
	}

	if count == 0 || sum == 0 {
		return 550, nil
	}

	if count < 3 {
		count++
		sum += 550
	}

	return sum / float64(count), nil
}

func (s Imp) UpdateSlotTimeMS() error {

	client := s.rpcClients["mainnet"]
rep:
	ei1, err := client.RpcClient.GetEpochInfo(context.Background())
	if err != nil {
		return err
	}

	t1 := time.Now()

	<-time.After(time.Hour * 1)

	ei2, err := client.RpcClient.GetEpochInfo(context.Background())
	if err != nil {
		return err
	}

	t2 := time.Now()

	if ei1.Result.Epoch != ei2.Result.Epoch {
		goto rep
	}

	sps := float64(ei2.Result.SlotIndex-ei1.Result.SlotIndex) / t2.Sub(t1).Seconds()
	if sps == 0 {
		return err
	}

	if err := s.DAO.CreateSlotTime(&dmodels.SlotTime{
		SlotTime:  (1 / sps) * 1000,
		Epoch:     ei2.Result.Epoch,
		CreatedAt: time.Now(),
	}); err != nil {
		return err
	}

	return nil
}
