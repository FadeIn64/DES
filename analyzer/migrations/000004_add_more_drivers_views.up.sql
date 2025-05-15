-- +goose Up
create view drivers_positions_with_intervals as(
select dswp.position, dswp.meeting_key, dswp.session_key, dswp.driver_number, dswp.lap_number, di.interval, di.prediction_laps_to_overtake
from drivers_stats_with_positions dswp
         inner join drivers_intervals di
                    on dswp.meeting_key = di.meeting_key and dswp.session_key = di.session_key and dswp.driver_number = di.driver_number);

create view driver_positions_with_intervals_and_last_lap as(
with complete_laps_group_by_driver as(
   select cl.meeting_key, cl.session_key, cl.driver_number, max(cl.lap_number) m_lap_number from complete_laps cl
   group by cl.meeting_key, cl.session_key, cl.driver_number
), last_complete_laps as(
   select cl2.meeting_key, cl2.session_key, cl2.driver_number, cl2.lap_duration
   from complete_laps cl2 inner join complete_laps_group_by_driver clgbd
     on cl2.meeting_key = clgbd.meeting_key and cl2.session_key = clgbd.session_key
         and cl2.driver_number = clgbd.driver_number and cl2.lap_number = clgbd.m_lap_number
)
select dpwi.*, lcl.lap_duration last_lap_duration
from drivers_positions_with_intervals dpwi inner join last_complete_laps lcl
     on dpwi.meeting_key = lcl.meeting_key and dpwi.session_key = lcl.session_key
         and dpwi.driver_number = lcl.driver_number
);

create view full_driver_data as(
   with pit_laps as(
       select cl.meeting_key, cl.session_key, cl.driver_number, count(is_pit_out_lap) pitsops, max(lap_number) last_pit_lap
       from complete_laps cl
       where is_pit_out_lap = true
       group by cl.meeting_key, cl.session_key, cl.driver_number
   )
   select dvpwiall.*, pl.pitsops, pl.last_pit_lap
   from driver_positions_with_intervals_and_last_lap dvpwiall left join pit_laps pl
    on dvpwiall.meeting_key = pl.meeting_key and dvpwiall.session_key = pl.session_key
        and dvpwiall.driver_number = pl.driver_number
);



-- +goose Down
drop view full_driver_data;
drop view driver_positions_with_intervals_and_last_lap;
drop view drivers_positions_with_intervals;
