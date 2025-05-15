-- +goose Up
create view drivers_stats as
with grouped_laps_by_lap as (select distinct l.meeting_key, l.session_key, l.driver_number,
                                             max(l.lap_number) m_lap_number
                             from laps l
                             group by l.meeting_key, l.session_key, l.driver_number),
     grouped_laps_by_lap_and_cs as(
         select distinct l2.meeting_key, l2.session_key, l2.driver_number, l2.lap_number, max(l2.completed_sectors) m_sector
         from laps l2 inner join grouped_laps_by_lap gll
                                 on l2.meeting_key = gll.meeting_key and l2.session_key = gll.session_key
                                     and l2.driver_number = gll.driver_number
                                     and l2.lap_number = gll.m_lap_number
         group by l2.meeting_key, l2.session_key, l2.driver_number, l2.lap_number)
select distinct l3.*
from laps l3 inner join grouped_laps_by_lap_and_cs gllacs
                        on l3.driver_number = gllacs.driver_number and l3.meeting_key = gllacs.meeting_key
                            and l3.session_key = gllacs.session_key and l3.lap_number = gllacs.lap_number and l3.completed_sectors = gllacs.m_sector;

create view drivers_stats_with_positions as
    select row_number() over (order by lap_number desc, completed_sectors desc, date_end) as position, * from drivers_stats;

-- +goose Down
drop view drivers_stats_with_positions;
drop table drivers_stats;