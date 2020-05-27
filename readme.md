# Feed Parser Bot

## Run
Execute this to launch the bot\
`TOKEN=%token% DEBUG_MODE=%debug_mode% HOST=%host% PORT=%port% USER=%user% PASS=%pass% DBNAME=%db_name% SCRAP=%scrap% YTToken=%YTToken% go run main.go`

%token% - Telegram Bot API Token\
%debug_mode% - Boolean that indicates if program should be run in debug mode\
%host%, %port%, %user%, %pass%, %db_name% - PostgreSQL credentials\
%scrap% - Boolean that indicates if program should scrap News Sources on start\
%YTToken% - Yandex.Translate API Token