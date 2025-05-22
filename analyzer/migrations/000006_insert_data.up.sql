-- +goose Up
INSERT INTO teams (team_key, name, description, country, color)
VALUES
    (1, 'Red Bull Racing', 'Команда Red Bull Racing, действующий чемпион Формулы-1', 'Austria', '#3671C6'),
    (2, 'Mercedes-AMG Petronas F1 Team', 'Официальная команда Mercedes, 8-кратный чемпион конструкторов', 'Germany', '#6CD3BF'),
    (3, 'Scuderia Ferrari', 'Легендарная итальянская команда, старейшая в Формуле-1', 'Italy', '#F91536'),
    (4, 'McLaren F1 Team', 'Британская команда с богатой историей в Формуле-1', 'United Kingdom', '#F58020'),
    (5, 'Aston Martin Aramco F1 Team', 'Фабричная команда Aston Martin, базируется в Сильверстоуне', 'United Kingdom', '#358C75'),
    (6, 'Alpine F1 Team', 'Французская команда, заводская команда Renault', 'France', '#2293D1'),
    (7, 'Williams Racing', 'Историческая британская команда, многократный чемпион', 'United Kingdom', '#64C4FF'),
    (8, 'Visa Cash App RB F1 Team', 'Младшая команда Red Bull, ранее AlphaTauri', 'Italy', '#5E8FAA'),
    (9, 'Stake F1 Team Kick Sauber', 'Швейцарская команда, ранее Alfa Romeo Racing', 'Switzerland', '#C92D4B'),
    (10, 'MoneyGram Haas F1 Team', 'Американская команда, основанная в 2016 году', 'United States', '#B6BABD');

INSERT INTO drivers (driver_number, team_key, full_name, abbreviation, country, date_of_birth, description)
VALUES
    -- Red Bull Racing (team_key = 1)
    (1, 1, 'Max Verstappen', 'VER', 'Netherlands', '1997-09-30', 'Действующий чемпион мира, доминирует в сезоне 2024'),
    (11, 1, 'Sergio Pérez', 'PER', 'Mexico', '1990-01-26', 'Опытный гонщик, второй пилот Red Bull'),

    -- Mercedes (team_key = 2)
    (44, 2, 'Lewis Hamilton', 'HAM', 'United Kingdom', '1985-01-07', '7-кратный чемпион мира, переходит в Ferrari в 2025'),
    (63, 2, 'George Russell', 'RUS', 'United Kingdom', '1998-02-15', 'Талантливый британский гонщик, лидер Mercedes'),

    -- Ferrari (team_key = 3)
    (16, 3, 'Charles Leclerc', 'LEC', 'Monaco', '1997-10-16', 'Лидер Ferrari, многократный обладатель поулов'),
    (55, 3, 'Carlos Sainz', 'SAI', 'Spain', '1994-09-01', 'Победитель Гран-при в 2024, покидает Ferrari в конце сезона'),

    -- McLaren (team_key = 4)
    (4, 4, 'Lando Norris', 'NOR', 'United Kingdom', '1999-11-13', 'Молодой талант, регулярно борется за подиумы'),
    (81, 4, 'Oscar Piastri', 'PIA', 'Australia', '2001-04-06', 'Перспективный новичок, показывает стабильные результаты'),

    -- Aston Martin (team_key = 5)
    (14, 5, 'Fernando Alonso', 'ALO', 'Spain', '1981-07-29', 'Двукратный чемпион мира, продолжает выступать на высоком уровне'),
    (18, 5, 'Lance Stroll', 'STR', 'Canada', '1998-10-29', 'Сын владельца команды, показывает переменные результаты'),

    -- Alpine (team_key = 6)
    (10, 6, 'Pierre Gasly', 'GAS', 'France', '1996-02-07', 'Опытный французский гонщик, бывший победитель Гран-при'),
    (31, 6, 'Esteban Ocon', 'OCO', 'France', '1996-09-17', 'Победитель Гран-при Венгрии-2021, покинет Alpine в конце 2024'),

    -- Williams (team_key = 7)
    (23, 7, 'Alexander Albon', 'ALB', 'Thailand', '1996-03-23', 'Лидер Williams, регулярно добывает очки'),
    (2, 7, 'Logan Sargeant', 'SAR', 'United States', '2000-12-31', 'Единственный американец в чемпионате, борется за место'),

    -- RB (team_key = 8)
    (3, 8, 'Daniel Ricciardo', 'RIC', 'Australia', '1989-07-01', 'Ветеран Формулы-1, вернулся в 2023'),
    (22, 8, 'Yuki Tsunoda', 'TSU', 'Japan', '2000-05-11', 'Японский гонщик, показывает прогресс в 2024'),

    -- Sauber (team_key = 9)
    (24, 9, 'Zhou Guanyu', 'ZHO', 'China', '1999-05-30', 'Первый китайский гонщик в Формуле-1'),
    (77, 9, 'Valtteri Bottas', 'BOT', 'Finland', '1989-08-28', 'Бывший гонщик Mercedes, опытный пилот'),

    -- Haas (team_key = 10)
    (20, 10, 'Kevin Magnussen', 'MAG', 'Denmark', '1992-10-05', 'Агрессивный гонщик, вернулся в Haas в 2022'),
    (27, 10, 'Nico Hülkenberg', 'HUL', 'Germany', '1987-08-19', 'Опытный немец, вернулся в Ф1 в 2023'),

    -- Замены и резервные гонщики 2024
    (40, 8, 'Liam Lawson', 'LAW', 'New Zealand', '2002-02-11', 'Молодой новозеландец, заменял Риккардо в 2023'),
    (45, 3, 'Oliver Bearman', 'BEA', 'United Kingdom', '2005-05-08', 'Дебютировал за Ferrari в Саудовской Аравии-2024'),
    (38, 7, 'Franco Colapinto', 'COL', 'Argentina', '2003-05-27', 'Молодой аргентинский талант, резерв Williams'),
    (50, 6, 'Jack Doohan', 'DOO', 'Australia', '2003-01-20', 'Австралийский гонщик, резерв Alpine');


-- +goose Down
delete from drivers where team_key in (1, 2,3, 4, 5, 6, 7, 8, 9, 10);
delete from teams where team_key in (1, 2,3, 4, 5, 6, 7, 8, 9, 10);
