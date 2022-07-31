# Lua CSV Reader

掃描當前目錄中所有副檔名為 `csv` 的檔案載入 `lua` 腳本中。

## 編譯

```bash
make all
```

## 工具說明

|長參數|短參數|選填|說明|預設值|範例|
|-|-|-|-|-|-|
|--file|-f|✔️|Lua 程序進入點|main.lua|-|
|--log-level|-|✔️|日誌層級debug,info,warn,error|info|-|
|--version|-v|✔️| 檢視程序建置版號|-|-|
|--help|-h|✔️| 幫助說明|-|-|

```bash
./luacsv --file main.lua
```

## Lua 檔案說明

參考 [main.lua](assets/main.lua)。

