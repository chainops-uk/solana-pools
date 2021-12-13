package services

import (
	"context"
	"github.com/everstake/solana-pools/internal/services/smodels"
	"time"
)

func (s Imp) UpdateEpoch(ctx context.Context) error {

	client := s.rpcClients["mainnet"]
	ei1, err := client.RpcClient.GetEpochInfo(ctx)
	if err != nil {
		return err
	}

	t1 := time.Now()

	<-time.After(time.Minute * 1)

	ei2, err := client.RpcClient.GetEpochInfo(ctx)
	if err != nil {
		return err
	}

	t2 := time.Now()

	if ei1.Result.Epoch != ei2.Result.Epoch {
		return err
	}

	sps := float64(ei2.Result.SlotIndex-ei1.Result.SlotIndex) / t2.Sub(t1).Seconds()

	emptyS := ei2.Result.SlotsInEpoch - ei2.Result.SlotIndex

	progress := (float64(ei2.Result.SlotIndex) / float64(ei2.Result.SlotsInEpoch)) * 100

	if sps == 0 {
		return err
	}

	s.cache.SetCurrentEpochInfo(&smodels.EpochInfo{
		Epoch:        ei2.Result.Epoch,
		SlotsInEpoch: ei2.Result.SlotsInEpoch,
		SPS:          sps,
		EndEpoch:     time.Now().Add(time.Duration((float64(emptyS) / sps) * float64(time.Second))),
		Progress:     uint8(progress),
	})

	return nil
}
