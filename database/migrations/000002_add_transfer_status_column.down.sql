
alter table transfers add column old_status varchar default 'pending';

update transfers set old_status = status::varchar;

alter table transfers drop column status;

drop type if exists transfer_status;