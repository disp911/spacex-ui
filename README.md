# SpaceX UI

Форк панели [3X-UI](https://github.com/MHSanaei/3x-ui) для управления Xray с
Telegram-прокси (MTProto через telemt): протокол `mtproto` создаётся на вкладке
«Подключения» так же, как любое другое подключение.

## Требования

- Linux: Ubuntu, Debian, Armbian, CentOS, Fedora, AlmaLinux, Rocky, Oracle
  Linux, Amazon Linux, Arch, Manjaro, openSUSE или Alpine.
- Архитектура: `amd64`, `arm64`, `armv7`, `armv6`, `armv5`, `386` или `s390x`.
  Telegram-прокси (telemt) входит только в сборки `amd64` и `arm64`.
- Права root.

## Установка

Последняя версия:

```bash
bash <(curl -Ls https://raw.githubusercontent.com/disp911/spacex-ui/main/install.sh)
```

В конце установщик выводит адрес панели, логин и пароль — сохраните их.

## Обновление

Обновление сохраняет все настройки, подключения и клиентов.

```bash
x-ui update
```

То же без меню `x-ui` (например, если оно повреждено):

```bash
bash <(curl -Ls https://raw.githubusercontent.com/disp911/spacex-ui/main/update.sh)
```

Когда выходит новая версия, на дашборде панели появляется оранжевая кнопка
с её номером — обновить можно и ею.

Обновить только скрипт меню `x-ui`, не трогая панель: `x-ui` → пункт
**3. Update Menu**.

## Установка определённой версии

Подходит и для отката на предыдущую версию: база с настройками остаётся на
месте. Номер версии указывается с буквой `v`; доступные версии —
на странице [Releases](https://github.com/disp911/spacex-ui/releases).

```bash
bash <(curl -Ls https://raw.githubusercontent.com/disp911/spacex-ui/main/install.sh) v2.9.12
```

То же через меню: `x-ui` → **4. Legacy Version**, затем ввести номер без `v`
(например, `2.9.12`).

## Удаление

> [!WARNING]
> Удаляются панель, Xray, telemt и база `/etc/x-ui` со всеми подключениями и
> клиентами. Если они могут понадобиться, сначала сохраните резервную копию.

Резервная копия базы:

```bash
cp /etc/x-ui/x-ui.db ~/x-ui.db.backup
```

Удаление (скрипт попросит подтверждение):

```bash
x-ui uninstall
```

## Управление панелью

`x-ui` без аргументов открывает меню со всеми действиями: сброс логина и
пароля, смена порта и пути панели, SSL-сертификаты, файрвол, BBR и другое.

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

## Где лежат файлы

| Путь | Содержимое |
| --- | --- |
| `/usr/local/x-ui/` | Панель, Xray и telemt |
| `/etc/x-ui/x-ui.db` | База: настройки, подключения, клиенты |
| `/usr/bin/x-ui` | Скрипт меню |
