-- +goose Up

INSERT INTO meetings (
    meeting_key, name, description, circuit, location,
    start_date, end_date, year, dashboard_link
) VALUES
      -- Гран-при Великобритании 2024
      (1240, 'British Grand Prix',
       'Формула-1 Гран-при Великобритании 2024 года',
       'Сильверстоун', 'Сильверстоун, Великобритания',
       '2024-07-05 00:00:00+00', '2024-07-07 00:00:00+00',
       2024, ''),

      -- Гран-при Венгрии 2024
      (1241, 'Hungarian Grand Prix',
       'Формула-1 Гран-при Венгрии 2024 года',
       'Хунгароринг', 'Будапешт, Венгрия',
       '2024-07-19 00:00:00+00', '2024-07-21 00:00:00+00',
       2024, '');

-- +goose Down

delete from meetings where meeting_key in (1240, 1241);