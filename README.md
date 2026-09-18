# SpaceX UI

## Установка

```bash
bash <(curl -Ls https://raw.githubusercontent.com/disp911/spacex-ui/main/install.sh)
```

## Установка архивной версии

```bash
bash <(curl -Ls https://raw.githubusercontent.com/disp911/spacex-ui/main/install.sh) v2.9.12
```

## Обновление

Обновление сохраняет все настройки, подключения и клиентов.

```bash
bash <(curl -Ls https://raw.githubusercontent.com/disp911/spacex-ui/main/update.sh)
```
или в меню

```bash
x-ui update
```

## Удаление

> [!WARNING]
> Удаляются панель, Xray, telemt и база `/etc/x-ui` со всеми подключениями и
> клиентами. Сохраните резервную копию при необходимости.

Резервная копия базы:

```bash
cp /etc/x-ui/x-ui.db ~/x-ui.db.backup
```

Удаление:

```bash
x-ui uninstall
```

## Управление панелью

| Команда | Действие |
| --- | --- |
| `x-ui` | Меню управления |
| `x-ui start` | Запустить панель |
| `x-ui stop` | Остановить панель |
| `x-ui restart` | Перезапустить панель |
| `x-ui restart-xray` | Перезапустить только Xray |
| `x-ui status` | Состояние панели |
| `x-ui settings` | Текущие настройки: порт, путь, адрес входа |
| `x-ui enable` | Включить автозапуск при загрузке системы |
| `x-ui disable` | Выключить автозапуск |
| `x-ui log` | Журнал панели |
| `x-ui banlog` | Журнал блокировок Fail2ban |
| `x-ui update` | Обновить панель до последней версии |
| `x-ui update-all-geofiles` | Обновить geo-файлы |
| `x-ui legacy` | Установить определённую версию |
| `x-ui install` | Установить панель |
| `x-ui uninstall` | Удалить панель |

## Расположение файлов

| Путь | Содержимое |
| --- | --- |
| `/usr/local/x-ui/` | Панель, Xray и telemt |
| `/etc/x-ui/x-ui.db` | База: настройки, подключения, клиенты |
| `/usr/bin/x-ui` | Скрипт меню |
