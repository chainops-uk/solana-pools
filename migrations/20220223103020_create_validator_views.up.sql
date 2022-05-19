CREATE OR REPLACE VIEW "public"."validator_view" as
SELECT v.id,
       v.image,
       v.name,
       v.delinquent,
       v.node_pk,
       vd.staking_accounts,
       vd.active_stake,
       vd.fee,
       (SELECT avg(apy) FROM validator_data where epoch >= vd.epoch-9 AND validator_id = vd.validator_id) as apy,
       (SELECT avg(score) FROM validator_data where epoch >= vd.epoch-9 AND validator_id = vd.validator_id) as score,
       (SELECT avg(skipped_slots) FROM validator_data where epoch >= vd.epoch-9 AND validator_id = vd.validator_id) as skipped_slots,
       v.data_center,
       v.created_at,
       v.updated_at
FROM validator_data vd
         JOIN validators v ON v.id::text = vd.validator_id::text
WHERE ((vd.validator_id::text, vd.updated_at) IN (SELECT validator_data.validator_id,
                                                         max(validator_data.updated_at)
                                                  FROM validator_data
                                                  GROUP BY validator_data.validator_id));

CREATE TABLE IF NOT EXISTS "public"."material_validator_data_view" (
    id                  varchar(44) Unique,
    image               text,
    name                varchar(100),
    delinquent          bool,
    node_pk             varchar(44),
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
        SET image = subquery.image,
            name = subquery.name,
            delinquent = subquery.delinquent,
            node_pk = subquery.node_pk,
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
              WHERE id = NEW.validator_id) subquery
        WHERE subquery.id = material_validator_data_view.id;
        return new;
    end if;

    INSERT INTO "public"."material_validator_data_view"(id,
                                                        image,
                                                        name,
                                                        delinquent,
                                                        node_pk,
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
           validator_view.node_pk,
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