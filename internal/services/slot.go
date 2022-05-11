package services

import (
	"context"
	"errors"
	"github.com/everstake/solana-pools/internal/dao/cache"
	"time"
)

func (s Imp) GetAvgSlotTimeMS() (float64, error) {
	arr, err := s.Cache.GetSlotArr()
	if err != nil && !errors.Is(err, cache.KeyWasNotFound) {
		return 0, err
	}
	if errors.Is(err, cache.KeyWasNotFound) {
		return 550, nil
	}

	var sum float64
	var count int
	for _, f := range arr {
		if f != 0 {
			count++
		}
		sum += f
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

	arr, err := s.Cache.GetSlotArr()
	if err != nil && !errors.Is(err, cache.KeyWasNotFound) {
		return err
	}

	for i, _ := range arr {
		if arr[i] == 0 {

			arr[i] = (1 / sps) * 1000

			if i < 23 {
				arr[i+1] = 0
			} else {
				arr[0] = 0
			}
			break
		}

	}

	s.Cache.SetSlotArr(arr)

	return nil
}
