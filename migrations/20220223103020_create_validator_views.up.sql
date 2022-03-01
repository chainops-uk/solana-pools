CREATE OR REPLACE VIEW validator_view as
SELECT v.id,
       v.image,
       v.name,
       v.delinquent,
       v.vote_pk,
       AVG(vd.apy)::numeric(8,4) as apy,
       AVG(vd.score)::int8 as score,
       AVG(vd.skipped_slots)::numeric(5,2) as skipped_slots,
       v.data_center,
       v.created_at,
       v.updated_at
FROM public.validator_data vd
         INNER JOIN validators v on v.id = vd.validator_id
WHERE vd.epoch BETWEEN vd.epoch - 10 AND vd.epoch
GROUP BY v.id, v.image, v.name, v.delinquent, v.vote_pk, v.data_center, v.created_at, v.updated_at;