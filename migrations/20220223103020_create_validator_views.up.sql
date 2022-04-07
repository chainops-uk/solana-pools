CREATE OR REPLACE VIEW "public"."validator_view" as
SELECT v.id,
       v.image,
       v.name,
       v.delinquent,
       v.vote_pk,
       vd.staking_accounts,
       vd.active_stake,
       vd.fee,
       avg_data.apy,
       avg_data.score,
       avg_data.skipped_slots,
       v.data_center,
       v.created_at,
       v.updated_at
FROM validator_data vd
         JOIN validators v ON v.id::text = vd.validator_id::text
         JOIN (SELECT validator_data.validator_id,
                      avg(validator_data.apy)::numeric(8, 4)           AS apy,
                      avg(validator_data.score)::bigint                AS score,
                      avg(validator_data.skipped_slots)::numeric(5, 2) AS skipped_slots
               FROM validator_data
               WHERE validator_data.epoch >= (validator_data.epoch - 10)
               GROUP BY validator_data.validator_id) avg_data ON v.id::text = avg_data.validator_id::text
WHERE ((vd.validator_id::text, vd.updated_at) IN (SELECT validator_data.validator_id,
                                                         max(validator_data.updated_at)
                                                  FROM validator_data
                                                  GROUP BY validator_data.validator_id));

CREATE TABLE IF NOT EXISTS "public"."material_validator_data_view" (
    id                  varchar(44) Unique,
    image               text,
    name                varchar(100),
    delinquent          bool,
    vote_pk             varchar(44),
    staking_accounts    int8,
    active_stake        int8,
    fee                 numeric(5, 2),
    apy                 numeric(8, 4),
    score               int8,
    skipped_slots       numeric(5, 4),
    data_center         text,
    created_at          TIMESTAMP WITH TIME ZONE,
    updated_at          TIMESTAMP WITH TIME ZONE
);

CREATE OR REPLACE FUNCTION add_material_validator_data_view()
    RETURNS TRIGGER
    LANGUAGE PLPGSQL
AS
$$
BEGIN
    IF EXISTS(SELECT 1 FROM "public"."material_validator_data_view" WHERE id = NEW.validator_id) THEN
        UPDATE "public"."material_validator_data_view"
        SET id = subquery.id,
            image = subquery.image,
            name = subquery.name,
            delinquent = subquery.delinquent,
            vote_pk = subquery.vote_pk,
            staking_accounts = subquery.staking_accounts,
            active_stake = subquery.active_stake,
            fee = subquery.fee,
            apy = subquery.apy,
            score = subquery.score,
            skipped_slots = subquery.skipped_slots,
            data_center = subquery.data_center,
            updated_at = now()
        FROM (SELECT *
              FROM "public"."validator_view"
              WHERE id = NEW.validator_id
              LIMIT 1) subquery;
        return new;
    end if;

    INSERT INTO "public"."material_validator_data_view"(id,
                                                        image,
                                                        name,
                                                        delinquent,
                                                        vote_pk,
                                                        apy,
                                                        staking_accounts,
                                                        active_stake,
                                                        fee,
                                                        score,
                                                        skipped_slots,
                                                        data_center,
                                                        created_at,
                                                        updated_at)
    SELECT validator_view.id,
           validator_view.image,
           validator_view.name,
           validator_view.delinquent,
           validator_view.vote_pk,
           validator_view.apy,
           validator_view.staking_accounts,
           validator_view.active_stake,
           validator_view.fee,
           validator_view.score,
           validator_view.skipped_slots,
           validator_view.data_center,
           now(),
           now()
    FROM "public"."validator_view"
    WHERE id = NEW.validator_id
    LIMIT 1;

    RETURN NEW;
END
$$;

CREATE TRIGGER add_material_validator_data_view_trigger
    AFTER INSERT
    ON "public"."validator_data"
    FOR EACH ROW
EXECUTE PROCEDURE add_material_validator_data_view();