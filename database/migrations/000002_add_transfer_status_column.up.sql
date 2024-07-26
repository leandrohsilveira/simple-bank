create type transfer_status as enum ('pending', 'rejected', 'transfering', 'completed', 'failed');

alter table transfers add column status transfer_status default 'pending'::transfer_status not null;

DO $$
BEGIN
  IF EXISTS(SELECT *
    FROM information_schema.columns
    WHERE table_name='transfers' and column_name='old_status')
  THEN
    update transfers set status = old_status::transfer_status;
    alter table transfers drop column old_status;  
  END IF;
END $$;
