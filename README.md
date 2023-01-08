# **clnr**

## Утилита для очистки Linux <img src="https://media.tenor.com/fP_RQeMnWecAAAAj/penguin-wiping-floor.gif" width="20">

Что на данный момент умеет программа:

- <img src="https://media.tenor.com/PBuEkZA9cVwAAAAi/sceptical-trashcan.gif" width="15"> &nbsp; Очищать кэш оперативной памяти
- <img src="https://media.tenor.com/gGY6gCZu42kAAAAi/doggy-dog.gif" width="17"> &nbsp; Очищать временные файлы Linux
- <img src="https://media.tenor.com/VRQnbam6nfwAAAAi/wiping-squidward.gif" width="20"> &nbsp; Очищать раздел подкачки

***Примечание:** данная утилита работает только на Linux (Работа на macOS не гарантируется)*

## Установка

```bash
git clone https://github.com/arcxevodov/clnr.git
cd clnr
make
sudo make install
```

## Флаги

```bash
-r:  Очистка кэша оперативной памяти
-s:  Перезагрузка Swap
-t:  Очистка временных файлов    
```
