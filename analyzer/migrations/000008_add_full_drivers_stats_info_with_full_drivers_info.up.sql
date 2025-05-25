-- +goose Up
alter view full_driver_data rename to all_drivers_stats;

create view full_driver_data as (
select ads.*, d.full_name, d.abbreviation, t.name, t.color from all_drivers_stats ads
    inner join drivers d
        inner join teams t on d.team_key = t.team_key
    on ads.driver_number = d.driver_number
);

-- +goose Down

drop view full_driver_data;
alter view all_drivers_stats rename to full_driver_data;

