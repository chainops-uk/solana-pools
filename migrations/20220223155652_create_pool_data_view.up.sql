CREATE OR REPLACE VIEW pool_data_view as
SELECT pd.id,
       pd.pool_id,
       pd.epoch,
       pd.active_stake,
       pd.total_tokens_supply,
       pd.total_lamports,
       (SELECT AVG(t.apy)::numeric
        FROM pool_data t
        WHERE t.epoch between pd.epoch - 9 AND pd.epoch AND t.pool_id = pd.pool_id) as apy,
       pd.unstake_liquidity,
       pd.depossit_fee,
       pd.withdrawal_fee,
       pd.rewards_fee,
       pd.updated_at,
       pd.created_at
FROM pool_data pd;