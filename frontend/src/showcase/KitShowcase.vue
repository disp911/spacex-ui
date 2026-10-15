<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useTheme } from '../stores/theme'
import {
  XActionBar, XActionList, XActionRow, XBadge, XBusyOverlay, XButton, XCard, XCheckbox, XChip, XCodeArea, XCodeInput, XConfirm, XCounter, XDataTable, XEmpty, XErrorBox, XExample, XField, XFieldNote, XInfoRow, XInput, XMenu, XModal, XMultiSelect, XNavItem, XNest, XNestCard, XNotice, XNumber, XOrderToggles,
  XPager, XPickList, XQrCode, XRecordCard, XRowLink, XSection, XSegmented, XSelect, XSettingBlock, XSettingDivider, XSettingRow, XSettingsSection, XSettingsSkeleton, XStatusDot, XStatusText, XSwitch, XTabs, XTagInput, XTextarea, XTooltip, XUsagePill,
  type MenuItem,
} from '../kit'

/*
 * Every kit component in every state, in both themes. The state table forces
 * hover, focus and pressed with the is-hover / is-focus / is-active classes
 * that the components accept next to the real pseudo-classes.
 */

const store = useTheme()
const theme = computed({ get: () => store.theme, set: (v) => store.set(v) })

const states = [
  { key: 'default', label: 'Обычное', cls: '' },
  { key: 'hover', label: 'Наведение', cls: 'is-hover' },
  { key: 'focus', label: 'Фокус', cls: 'is-focus' },
  { key: 'active', label: 'Нажатие', cls: 'is-active' },
  { key: 'disabled', label: 'Недоступно', cls: '' },
  { key: 'invalid', label: 'Ошибка', cls: '' },
  { key: 'modified', label: 'Изменено', cls: '' },
] as const
type State = (typeof states)[number]['key']
const has = (s: State, ...list: State[]) => list.includes(s)

const languages = [
  { value: 'ru', label: '🇷🇺 Русский' },
  { value: 'en', label: '🇬🇧 English' },
  { value: 'fa', label: '🇮🇷 فارسی' },
  { value: 'zh', label: '🇨🇳 中文', disabled: true },
]
const strategies = ['AsIs', 'IPIfNonMatch', 'IPOnDemand'].map((v) => ({ value: v, label: v }))
const badRouting = 'happ://routing/…'
const geo = ['geoip:private', 'geoip:ru', 'geosite:category-ads-all', 'geosite:google', 'geosite:telegram']

const form = ref({
  enabled: true,
  remark: 'Germany-VLESS',
  port: 443 as number | null,
  lang: 'ru' as string | null,
  strategy: 'IPIfNonMatch' as string | null,
  password: 'correct-horse',
  tags: ['geoip:private', 'geoip:ru'],
  note: 'Текст объявления для подписки',
  sniffing: true,
  routeOnly: false,
  agree: true,
})
const view = ref<'all' | 'online' | 'depleted'>('all')
const mode = ref<'search' | 'filter'>('search')
const page = ref(5)
const modalOpen = ref(false)
const sheetSize = ref<'sm' | 'md' | 'lg' | 'xl'>('md')
const boxedCheck = ref(true)
const successDemo = ref(false)
const logOpen = ref(false)
const confirmOpen = ref(false)
const confirmBusy = ref(false)
function runConfirm() {
  confirmBusy.value = true
  setTimeout(() => {
    confirmBusy.value = false
    confirmOpen.value = false
  }, 1500)
}
const level = ref('info')
const levelItems = computed<MenuItem[]>(() =>
  [
    { key: 'debug', label: 'Debug', dot: 'var(--upload)' },
    { key: 'info', label: 'Info', dot: 'var(--acc)' },
    { key: 'warning', label: 'Warning', dot: 'var(--warn)' },
    { key: 'err', label: 'Error', dot: 'var(--err)' },
  ].map((i) => ({ ...i, checked: i.key === level.value })),
)
const lastAction = ref('—')
const codeText = ref(JSON.stringify({ log: { loglevel: 'warning' }, dns: { servers: ['https://1.1.1.1/dns-query'] }, routing: { domainStrategy: 'AsIs' } }, null, 2))
const setTab = ref('panel')
const setTabs = [
  { key: 'panel', label: 'Панель', dot: 'err' as const },
  { key: 'ui', label: 'Интерфейс' },
  { key: 'acct', label: 'Учетная запись' },
  { key: 'ldap', label: 'LDAP', dot: 'mod' as const },
  { key: 'sub', label: 'Подписка' },
  { key: 'fmt', label: 'Formats' },
]
const setOpen = ref(true)
const setVals = ref({ port: 2054, session: 30, json: true, path: '/panel-x7/', rules: 'geosite:category-ads-all,geosite:ru,geoip:private,geoip:ru', tags: ['inbound-443'], frag: true })
const setInbounds = ['inbound-443', 'inbound-8443', 'vless-reality', 'trojan-ws'].map((v) => ({ value: v, label: v }))
const busyOpen = ref(false)
const restartAsk = ref(false)
const setTheme = ref('dark')
const demoCode = ref('')
const demoTfa = ref(false)
const demoParts = ref(['Inbound', 'Email'])
const demoSep = ref('-')
const demoTags = ref(['geoip:private'])
watch(busyOpen, (on) => on && setTimeout(() => (busyOpen.value = false), 3000))
const typeChecks = ref({ direct: true, blocked: true, proxy: false })
const tableColumns = [
  { key: 'time', label: 'Дата', width: 110, dim: true },
  { key: 'to', label: 'Куда' },
  { key: 'out', label: 'Исходящий', width: 120 },
  { key: 'email', label: 'Email', width: 100, strong: true },
]
const tableRows = [
  { time: '20.09 23:58:58', to: 'wa.me:443', out: 'direct', email: 'alice', tone: undefined },
  { time: '20.09 23:58:56', to: 'cdn.jsdelivr.net:443', out: 'proxy-de', email: 'bob', tone: 'upload' as const },
  { time: '20.09 23:58:50', to: 'ads.doubleclick.net:443', out: 'blocked', email: 'carol', tone: 'err' as const },
]
const pickItems = ['v26.4.25', 'v26.3.27', 'v25.12.8', 'v25.11.11'].map((v, i) => ({ key: v, label: v, current: i === 0 }))
const sectionOpen = ref(true)
const sw = ref({ a: false, b: true, loading: false })

const menuItems: MenuItem[] = [
  { key: 'edit', label: 'Изменить', icon: 'edit' },
  { key: 'refresh', label: 'Сбросить трафик', icon: 'refresh' },
  { key: 'export', label: 'Экспорт ссылок', icon: 'arrowRight', disabled: true },
  { key: 'delete', label: 'Удалить', icon: 'delete', danger: true, divided: true },
]

function applySwitch() {
  sw.value.loading = true
  setTimeout(() => (sw.value.loading = false), 1200)
}

const swatches = {
  'Поверхности': ['bg', 'side', 'card', 'nest', 'dlg', 'field', 'ctrl', 'chip', 'seg', 'track'],
  'Текст': ['fg', 'fg2', 'mut', 'dim', 'faint'],
  'Акцент': ['acc', 'acc-hover', 'acc-press', 'acc-soft', 'acc-ink'],
  'Статусы': ['ok', 'warn', 'err', 'unl', 'on', 'off', 'mod'],
}
const radii = [4, 6, 8, 10, 12, 16, 22]
const fontSizes = [10, 11, 12, 13, 14, 16, 18, 20, 24]
const spacing = [2, 4, 6, 8, 10, 12, 14, 16, 20, 24, 32, 40, 48, 56]
</script>

<template>
  <div class="kit">
    <header class="kit__top">
      <div>
        <h1 class="kit__title">spacex-ui · kit</h1>
        <p class="kit__sub">Компоненты v3 во всех состояниях. Ширина окна меньше 768 px показывает мобильные размеры.</p>
      </div>
      <XSegmented
        v-model="theme"
        label="Тема"
        :options="[{ value: 'dark', label: 'Тёмная', icon: 'moon' }, { value: 'light', label: 'Светлая', icon: 'sun' }]"
      />
    </header>

    <!-- Tokens -->
    <XCard title="Токены">
      <div class="kit__tokens">
        <div v-for="(names, group) in swatches" :key="group" class="kit__swatch-group">
          <h3 class="kit__h3">{{ group }}</h3>
          <div class="kit__swatches">
            <div v-for="n in names" :key="n" class="kit__swatch">
              <span class="kit__swatch-color" :style="{ background: `var(--${n})` }" />
              <code>--{{ n }}</code>
            </div>
          </div>
        </div>
        <div class="kit__swatch-group">
          <h3 class="kit__h3">Радиусы</h3>
          <div class="kit__swatches">
            <div v-for="r in radii" :key="r" class="kit__radius" :style="{ borderRadius: `var(--r-${r})` }">{{ r }}</div>
          </div>
        </div>
        <div class="kit__swatch-group">
          <h3 class="kit__h3">Размеры шрифта</h3>
          <div class="kit__fonts">
            <span v-for="f in fontSizes" :key="f" :style="{ fontSize: `var(--fs-${f})` }">{{ f }} Подключения</span>
            <code class="kit__mono">IBM Plex Mono · 203.0.113.7:443 · 48,2 МБ</code>
          </div>
        </div>
        <div class="kit__swatch-group">
          <h3 class="kit__h3">Отступы</h3>
          <div class="kit__spacing">
            <div v-for="s in spacing" :key="s" class="kit__space">
              <span :style="{ width: `var(--sp-${s})` }" />
              <code>{{ s }}</code>
            </div>
          </div>
        </div>
      </div>
    </XCard>

    <!-- States -->
    <XCard title="Состояния" flush>
      <div class="kit__table-wrap">
        <table class="kit__table">
          <thead>
            <tr>
              <th>Компонент</th>
              <th v-for="s in states" :key="s.key">{{ s.label }}</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <th><b>Button</b><code>primary</code></th>
              <td v-for="s in states" :key="s.key">
                <XButton v-if="!has(s.key, 'invalid', 'modified')" variant="primary" icon="plus" :class="s.cls" :disabled="s.key === 'disabled'">Создать</XButton>
                <XButton v-else-if="s.key === 'invalid'" variant="primary" loading>Сохранение…</XButton>
              </td>
            </tr>
            <tr>
              <th><b>Button</b><code>secondary</code></th>
              <td v-for="s in states" :key="s.key">
                <XButton v-if="!has(s.key, 'invalid', 'modified')" icon="globe" :class="s.cls" :disabled="s.key === 'disabled'">Geo-файлы</XButton>
              </td>
            </tr>
            <tr>
              <th><b>Button</b><code>soft</code></th>
              <td v-for="s in states" :key="s.key">
                <XButton v-if="!has(s.key, 'invalid', 'modified')" variant="soft" :class="s.cls" :disabled="s.key === 'disabled'">Подтвердить</XButton>
              </td>
            </tr>
            <tr>
              <th><b>Button</b><code>danger</code></th>
              <td v-for="s in states" :key="s.key">
                <XButton v-if="!has(s.key, 'invalid', 'modified')" variant="danger" :class="s.cls" :disabled="s.key === 'disabled'">Сбросить</XButton>
              </td>
            </tr>
            <tr>
              <th><b>Button</b><code>icon · ghost · sm</code></th>
              <td v-for="s in states" :key="s.key">
                <div v-if="!has(s.key, 'invalid', 'modified')" class="kit__row">
                  <XButton icon="refresh" label="Обновить" :class="s.cls" :disabled="s.key === 'disabled'" />
                  <XButton variant="ghost" size="sm" icon="dots" label="Ещё" :class="s.cls" :disabled="s.key === 'disabled'" />
                  <XButton variant="ghost-danger" size="sm" icon="delete" label="Удалить" :class="s.cls" :disabled="s.key === 'disabled'" />
                </div>
              </td>
            </tr>
            <tr>
              <th><b>Input</b><code>text</code></th>
              <td v-for="s in states" :key="s.key">
                <XInput
                  :model-value="s.key === 'invalid' ? 'panel' : 'Germany-VLESS'"
                  :class="s.cls"
                  :disabled="s.key === 'disabled'"
                  :invalid="s.key === 'invalid'"
                  :modified="s.key === 'modified'"
                />
              </td>
            </tr>
            <tr>
              <th><b>Input</b><code>mono · search</code></th>
              <td v-for="s in states" :key="s.key">
                <XInput
                  mono
                  icon="search"
                  placeholder="Поиск"
                  :model-value="s.key === 'focus' ? 'cdn' : ''"
                  :class="s.cls"
                  :disabled="s.key === 'disabled'"
                  :invalid="s.key === 'invalid'"
                  :modified="s.key === 'modified'"
                />
              </td>
            </tr>
            <tr>
              <th><b>Password</b></th>
              <td v-for="s in states" :key="s.key">
                <XInput
                  type="password"
                  model-value="secret-value"
                  :class="s.cls"
                  :disabled="s.key === 'disabled'"
                  :invalid="s.key === 'invalid'"
                  :modified="s.key === 'modified'"
                />
              </td>
            </tr>
            <tr>
              <th><b>Number</b><code>unit</code></th>
              <td v-for="s in states" :key="s.key">
                <XNumber
                  :model-value="s.key === 'invalid' ? 70000 : 2053"
                  unit="порт"
                  :class="s.cls"
                  :disabled="s.key === 'disabled'"
                  :invalid="s.key === 'invalid'"
                  :modified="s.key === 'modified'"
                />
              </td>
            </tr>
            <tr>
              <th><b>Select</b></th>
              <td v-for="s in states" :key="s.key">
                <XSelect
                  :options="languages"
                  model-value="ru"
                  :class="s.cls"
                  :disabled="s.key === 'disabled'"
                  :invalid="s.key === 'invalid'"
                  :modified="s.key === 'modified'"
                />
              </td>
            </tr>
            <tr>
              <th><b>Tag input</b></th>
              <td v-for="s in states" :key="s.key">
                <XTagInput
                  :model-value="['geoip:private', 'geoip:ru']"
                  mono
                  :class="s.cls"
                  :disabled="s.key === 'disabled'"
                  :invalid="s.key === 'invalid'"
                  :modified="s.key === 'modified'"
                />
              </td>
            </tr>
            <tr>
              <th><b>Textarea</b></th>
              <td v-for="s in states" :key="s.key">
                <XTextarea
                  :model-value="s.key === 'invalid' ? badRouting : 'Текст объявления'"
                  :rows="2"
                  :class="s.cls"
                  :disabled="s.key === 'disabled'"
                  :invalid="s.key === 'invalid'"
                  :modified="s.key === 'modified'"
                />
              </td>
            </tr>
            <tr>
              <th><b>Switch</b><code>38×22 · моб. 44×26</code></th>
              <td v-for="s in states" :key="s.key">
                <div class="kit__row">
                  <XSwitch
                    :model-value="s.key !== 'invalid'"
                    :class="s.cls"
                    :disabled="s.key === 'disabled'"
                    :invalid="s.key === 'invalid'"
                    :modified="s.key === 'modified'"
                    label="Включить"
                  />
                  <XSwitch v-if="!has(s.key, 'invalid', 'modified')" :class="s.cls" :disabled="s.key === 'disabled'" label="Выключено" />
                </div>
              </td>
            </tr>
            <tr>
              <th><b>Checkbox</b></th>
              <td v-for="s in states" :key="s.key">
                <div v-if="s.key !== 'modified'" class="kit__row">
                  <XCheckbox :model-value="true" :class="s.cls" :disabled="s.key === 'disabled'" :invalid="s.key === 'invalid'" />
                  <XCheckbox :class="s.cls" :disabled="s.key === 'disabled'" :invalid="s.key === 'invalid'" />
                  <XCheckbox indeterminate :class="s.cls" :disabled="s.key === 'disabled'" />
                </div>
              </td>
            </tr>
            <tr>
              <th><b>Section</b><code>заголовок</code></th>
              <td v-for="s in states" :key="s.key">
                <XSection v-if="!has(s.key, 'disabled', 'invalid', 'modified')" title="Сертификаты" :open="false" :class="s.cls" />
              </td>
            </tr>
            <tr>
              <th><b>Chip</b><code>interactive</code></th>
              <td v-for="s in states" :key="s.key">
                <XChip v-if="!has(s.key, 'disabled', 'invalid', 'modified')" tone="ok" mono interactive :class="s.cls">45,2 ГБ / 100 ГБ</XChip>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </XCard>

    <!-- Variants -->
    <XCard title="Варианты">
      <div class="kit__variants">
        <div class="kit__variant">
          <h3 class="kit__h3">Chip · статусы</h3>
          <div class="kit__row kit__row--wrap">
            <XChip tone="ok" icon="user" mono>10</XChip>
            <XChip tone="off" mono>1</XChip>
            <XChip tone="err" mono>2</XChip>
            <XChip tone="warn" mono>2</XChip>
            <XChip tone="on" dot mono>4</XChip>
            <XChip tone="unl" icon="infinity" mono>3</XChip>
            <XChip tone="neutral" mono>12 ГБ / ∞</XChip>
          </div>
        </div>
        <div class="kit__variant">
          <h3 class="kit__h3">Chip · срок</h3>
          <div class="kit__row kit__row--wrap">
            <XChip tone="ok">через 41 д</XChip>
            <XChip tone="warn">через 3 д</XChip>
            <XChip tone="err">5 д назад</XChip>
            <XChip tone="off">через 60 д</XChip>
            <XChip tone="ok" dot>Онлайн</XChip>
            <XChip>Офлайн</XChip>
          </div>
        </div>
        <div class="kit__variant">
          <h3 class="kit__h3">Usage pill</h3>
          <div class="kit__stack">
            <XUsagePill used="18,4 ГБ" limit="50 ГБ" :ratio="0.37" />
            <XUsagePill used="46,8 ГБ" limit="50 ГБ" :ratio="0.94" tone="warn" />
            <XUsagePill used="50 ГБ" limit="50 ГБ" :ratio="1" tone="err" exhausted />
            <XUsagePill used="0 Б" limit="30 ГБ" :ratio="0" disabled />
            <XUsagePill used="100,5 ГБ" limit="∞" :ratio="null" />
            <XUsagePill used="через 9 д" limit="30 д" :ratio="0.7" textual />
          </div>
        </div>
        <div class="kit__variant">
          <h3 class="kit__h3">Status dot · badge · counter</h3>
          <div class="kit__row kit__row--wrap">
            <XStatusDot tone="ok" label="ok" />
            <XStatusDot tone="warn" label="warn" />
            <XStatusDot tone="err" label="err" />
            <XStatusDot tone="unl" label="unlimited" />
            <XStatusDot tone="off" label="off" />
            <XStatusDot tone="ok" pulse label="online" />
            <XStatusDot tone="err" blink label="blink" />
            <XStatusDot size="sm" blink label="sm" />
          </div>
          <div class="kit__row kit__row--wrap">
            <XBadge mono>vless</XBadge>
            <XBadge mono>shadowsocks</XBadge>
            <XBadge kind="outline" mono>xhttp</XBadge>
            <XBadge kind="outline" mono>grpc</XBadge>
            <XBadge kind="accent" mono>reality</XBadge>
            <XBadge kind="accent">NEW</XBadge>
            <XCounter>12</XCounter>
            <XCounter accent>3</XCounter>
          </div>
        </div>
        <div class="kit__variant">
          <h3 class="kit__h3">Segmented · pager</h3>
          <div class="kit__row kit__row--wrap">
            <XSegmented
              v-model="view"
              label="Фильтр"
              :options="[{ value: 'all', label: 'Все' }, { value: 'online', label: 'Онлайн' }, { value: 'depleted', label: 'Исчерпано' }]"
            />
            <XSegmented
              v-model="mode"
              icon-only
              label="Режим"
              :options="[{ value: 'search', label: 'Поиск', icon: 'search' }, { value: 'filter', label: 'Фильтр', icon: 'filter' }]"
            />
          </div>
          <XPager v-model="page" :pages="12" />
        </div>
        <div class="kit__variant">
          <h3 class="kit__h3">Notice</h3>
          <div class="kit__stack">
            <XNotice tone="info">Изменения применятся после перезапуска Xray.</XNotice>
            <XNotice tone="warn" title="Сертификат истекает">Осталось 6 дней. Продлите его до 18 октября.</XNotice>
            <XNotice tone="err" title="Ошибка при запуске Xray">failed to listen on 0.0.0.0:443: address already in use</XNotice>
            <XNotice tone="ok">Настройки сохранены.</XNotice>
            <XNotice tone="warn" title="Соединение не защищено (HTTP)" compact>Не вводите конфиденциальные данные</XNotice>
          </div>
        </div>
        <div class="kit__variant">
          <h3 class="kit__h3">Кнопки с меню · кнопка-чекбокс · подпись</h3>
          <div class="kit__row kit__row--wrap">
            <XButton caret><span class="x-caps">Строк</span> 100</XButton>
            <XButton size="sm" mono caret>за 2 мин</XButton>
            <XCheckbox v-model="boxedCheck" boxed>SysLog</XCheckbox>
            <XCheckbox v-model="typeChecks.direct" boxed tone="ok">Direct</XCheckbox>
            <XCheckbox v-model="typeChecks.blocked" boxed tone="err">Blocked</XCheckbox>
            <XCheckbox v-model="typeChecks.proxy" boxed tone="upload">Proxy</XCheckbox>
            <XButton icon="download" phone-icon-only>Скачать</XButton>
            <XButton variant="soft" size="sm" phone-block>Обновить панель</XButton>
            <XButton variant="primary" :success="successDemo" @click="successDemo = !successDemo">Успех</XButton>
          </div>
          <div class="kit__row kit__row--wrap">
            <span class="x-caps">Подпись 10</span>
            <span class="x-caps x-caps--lg">Подпись 11</span>
          </div>
        </div>
        <div class="kit__variant">
          <h3 class="kit__h3">Строка-ссылка · пункт меню</h3>
          <div class="kit__stack" style="max-width: 280px; padding: 0 12px">
            <XRowLink icon="doc" label="Журнал панели" />
            <XRowLink icon="database" label="Резервная копия" />
          </div>
          <div class="kit__stack" style="max-width: 236px">
            <XNavItem icon="dashboard" label="Дашборд" to="/kit" />
            <XNavItem icon="sliders" label="Настройки" />
            <XNavItem icon="logout" label="Выход" danger />
          </div>
        </div>
        <div class="kit__variant">
          <h3 class="kit__h3">Список действий</h3>
          <XActionList style="max-width: 400px">
            <XActionRow title="Экспорт базы данных" icon="download">Скачать файл <code>.db</code> с резервной копией.</XActionRow>
            <XActionRow title="Импорт базы данных" icon="upload" loading>Загрузить файл <code>.db</code> с устройства.</XActionRow>
          </XActionList>
        </div>
        <div class="kit__variant">
          <h3 class="kit__h3">Пусто · ошибка</h3>
          <XEmpty title="Нет записей" hint="Для выбранного уровня и глубины журнал пуст" />
          <XErrorBox text="failed to listen on 0.0.0.0:443&#10;bind: address already in use" />
        </div>
        <div class="kit__variant">
          <h3 class="kit__h3">Таблица журнала · карточки на телефоне</h3>
          <XDataTable :columns="tableColumns" :rows="tableRows" :row-key="(r) => r.to" :row-tone="(r) => r.tone" :height="140" />
          <div class="kit__stack" style="max-width: 360px">
            <XRecordCard v-for="r in tableRows" :key="r.to" :title="r.to" time="23:58:58" name="alice" detail="inbound-443" :badge="r.out" :tone="r.tone" />
            <XRecordCard loading />
          </div>
        </div>
        <div class="kit__variant">
          <h3 class="kit__h3">Текст с номерами строк</h3>
          <XCodeArea v-model="codeText" :height="160" label="config.json" />
        </div>
        <div class="kit__variant">
          <h3 class="kit__h3">Строка значения</h3>
          <div class="kit__stack" style="max-width: 400px">
            <XInfoRow label="Текущая версия">3.0.0</XInfoRow>
            <XInfoRow label="Последняя версия" accent><XStatusDot size="sm" blink />3.0.1</XInfoRow>
            <XInfoRow label="Панель обновлена" accent><XChip tone="ok">Обновлено</XChip></XInfoRow>
          </div>
        </div>
        <div class="kit__variant">
          <h3 class="kit__h3">Список выбора · загрузка</h3>
          <XPickList :items="pickItems" current-label="Установлена" style="max-width: 400px" />
          <XPickList :items="[]" loading style="max-width: 400px" />
        </div>
        <div class="kit__variant">
          <h3 class="kit__h3">Полоса прокрутки</h3>
          <div class="kit__scroll" tabindex="0">
            <div v-for="n in 12" :key="n" class="mono">{{ String(n).padStart(2, '0') }} · 2026/10/14 12:00:{{ String(n).padStart(2, '0') }} INFO - XRAY: accepted tcp:www.example.com:443 [inbound-443 -> direct]</div>
          </div>
        </div>
      </div>
    </XCard>

    <!-- Live -->
    <XCard title="Страница настроек">
      <div class="kit__stack">
        <XTabs v-model="setTab" :items="setTabs" label="Разделы" />
        <XNotice tone="warn" title="Предупреждение системы безопасности" :items="['Подключение к панели не защищено. Установите SSL-сертификат.', 'Порт панели по умолчанию небезопасен.']" closable>
          Некоторые настройки уязвимы.
        </XNotice>
        <div class="kit__row kit__row--wrap">
          <XStatusText tone="mod">Не сохранено: 2</XStatusText>
          <XStatusText tone="err">Исправьте ошибки: 1 · не сохранено: 2</XStatusText>
          <XButton @click="busyOpen = true">Перезапуск: оверлей на 3 с</XButton>
          <XButton icon="refresh" @click="restartAsk = true">Подтверждение перезапуска</XButton>
        </div>
        <XSettingsSection v-model:open="setOpen" title="Основное" :count="6" :modified="2" :errors="1" note="IP по умолчанию">
          <XSettingRow title="Порт панели" description="Порт, на котором работает панель" narrow modified>
            <XNumber v-model="setVals.port" :min="1" :max="65535" modified />
          </XSettingRow>
          <XSettingRow title="Продолжительность сессии" description="Значение в минутах" narrow error="Значение должно быть не меньше 60">
            <XNumber v-model="setVals.session" unit="мин" invalid />
          </XSettingRow>
          <XSettingRow title="Включить JSON-подписку" description="Отдельный адрес с подпиской в формате JSON" inline>
            <XSwitch v-model="setVals.json" label="JSON" />
          </XSettingRow>
          <XSettingDivider label="Расширенный шаблон" />
          <XSettingRow title="Корневой путь URL-адреса панели" description="Должен начинаться и заканчиваться символом /" modified>
            <XInput v-model="setVals.path" mono modified />
            <template #note><XFieldNote>После перезапуска: <code>http://127.0.0.1:2053/panel-x7Q2/</code></XFieldNote></template>
          </XSettingRow>
          <XSettingRow title="Тема" compact>
            <XSegmented v-model="setTheme" :options="[{ value: 'light', label: 'Светлая', icon: 'sun' }, { value: 'dark', label: 'Тёмная', icon: 'moon' }]" block label="Тема" />
          </XSettingRow>
          <XSettingRow title="Правила маршрутизации" description="Длинное значение: растёт до трёх строк, Enter не переносит" tall>
            <XTextarea v-model="setVals.rules" mono :grow="3" />
          </XSettingRow>
          <XSettingRow title="Подключения для LDAP" description="Новые клиенты попадут в выбранные подключения">
            <XMultiSelect v-model="setVals.tags" :options="setInbounds" mono placeholder="Все подключения" empty="Подключений пока нет" />
          </XSettingRow>
          <XSettingRow title="Фрагментация" inline>
            <XSwitch v-model="setVals.frag" label="Фрагментация" />
            <template v-if="setVals.frag" #nested>
              <XNest title="Параметры фрагментации" :columns="140">
                <XField v-slot="f" label="Packets"><XInput v-bind="f" model-value="tlshello" mono /></XField>
                <XField v-slot="f" label="Length"><XInput v-bind="f" model-value="100-200" mono /></XField>
                <XField v-slot="f" label="Interval"><XInput v-bind="f" model-value="10-20" mono /></XField>
              </XNest>
            </template>
          </XSettingRow>
        </XSettingsSection>
        <XActionBar style="position: static; border-radius: 12px">
          <template #state><XStatusText tone="mod">Не сохранено: 2</XStatusText></template>
          <XButton variant="primary">Сохранить</XButton>
          <XButton icon="refresh" disabled>Перезапуск</XButton>
        </XActionBar>
        <XSettingsSection title="Отдельная форма и 2FA" :count="2" note="Отдельная форма">
          <XSettingBlock>
            <div class="kit__row kit__row--wrap">
              <XCodeInput v-model="demoCode" label="Код" style="width: 260px" />
              <XQrCode value="otpauth://totp/SpaceX%20UI:Administrator?secret=JBSWY3DPEHPK3PXP" copy="JBSWY3DPEHPK3PXP" copy-label="Скопировать токен" :size="120" />
            </div>
            <template #note>
              <XFieldNote tone="err">Неверное имя пользователя или пароль</XFieldNote>
              <span>Применяется кнопкой «Подтвердить», отдельно от «Сохранить».</span>
            </template>
            <template #actions><XButton variant="soft" icon="lock">Подтвердить</XButton></template>
          </XSettingBlock>
          <XSettingRow title="Модель примечания" description="Части по порядку, разделитель — сеткой символов" tall>
            <div style="display: flex; flex-direction: column; gap: 8px">
              <XExample label="Пример примечания:" :value="'#' + demoParts.join(demoSep)" />
              <div style="display: flex; gap: 8px">
                <XOrderToggles v-model="demoParts" :options="[{ value: 'Inbound', label: 'Inbound' }, { value: 'Email', label: 'Email' }, { value: 'Other', label: 'Other' }]" mono />
                <div style="flex: none; width: 98px">
                  <XSelect v-model="demoSep" :options="[' ', '-', '_', '@', ':', '~', '|', ',', '.', '/'].map((v) => ({ value: v, label: v === ' ' ? 'пробел' : v, glyph: v === ' ' ? '␣' : undefined }))" :columns="5" mono />
                </div>
              </div>
            </div>
          </XSettingRow>
          <XSettingRow title="Включить 2FA" description="Подпись состояния слева от переключателя" inline>
            <XSwitch v-model="demoTfa" :text="demoTfa ? 'Включена' : 'Выключена'" label="2FA" />
          </XSettingRow>
        </XSettingsSection>
        <XNest title="Список внутри блока">
          <div style="display: flex; flex-direction: column; gap: 10px">
            <XNestCard title="Noise №1" modified :columns="130">
              <template #actions><XButton variant="ghost-danger" size="sm" icon="delete">Удалить</XButton></template>
              <XField v-slot="a" label="Packet"><XInput v-bind="a" model-value="10-20" mono /></XField>
              <XField v-slot="a" label="Delay (ms)"><XInput v-bind="a" model-value="10-16" mono /></XField>
            </XNestCard>
            <div><XButton variant="dashed" icon="plus">Добавить Noise</XButton></div>
            <XTagInput v-model="demoTags" :presets="[{ value: 'geoip:private', label: 'Private IP' }, { value: 'geoip:ru', label: '🇷🇺 Russia' }]" presets-title="Пресеты" placeholder="Выберите пресет или введите значение" more-placeholder="Добавить…" />
          </div>
        </XNest>
        <XSettingsSkeleton :rows="[2]" />
        <XCard>
          <XEmpty tone="err" icon="alert" title="Произошла ошибка при получении параметров." />
        </XCard>
      </div>
    </XCard>
    <XConfirm v-model:open="restartAsk" icon="refresh" title="Перезапуск панели" cta="Да" @confirm="restartAsk = false">
      Вы уверены, что хотите перезапустить панель? Если панель будет недоступна, проверьте лог сервера
    </XConfirm>
    <XBusyOverlay :open="busyOpen" title="Перезапуск панели" text="Ожидание ответа панели…" />

    <div class="kit__live">
      <XCard title="Форма">
        <template #actions>
          <XTooltip text="Подсказка появляется при наведении и при фокусе с клавиатуры">
            <XButton variant="ghost" size="sm" icon="info" label="Справка" />
          </XTooltip>
          <XMenu v-slot="{ props }" :items="menuItems" @select="lastAction = $event">
            <XButton v-bind="props" variant="ghost" size="sm" icon="dots" label="Действия" />
          </XMenu>
        </template>
        <div class="kit__form">
          <XField layout="row" label="Включить" hint="Подключение принимает клиентов">
            <XSwitch v-model="form.enabled" label="Включить" />
          </XField>
          <div class="kit__grid">
            <XField v-slot="f" label="Примечание" required>
              <XInput v-bind="f" v-model="form.remark" />
            </XField>
            <XField v-slot="f" label="Порт" :error="form.port != null && form.port > 65535 ? 'Порт от 1 до 65535' : undefined">
              <XNumber v-bind="f" v-model="form.port" :min="1" :max="65535" :invalid="form.port != null && form.port > 65535" style="width: 100%" />
            </XField>
            <XField v-slot="f" label="Язык">
              <XSelect v-bind="f" v-model="form.lang" :options="languages" />
            </XField>
            <XField v-slot="f" label="Стратегия домена" hint="Как Xray определяет IP-адрес">
              <XSelect v-bind="f" v-model="form.strategy" :options="strategies" mono />
            </XField>
          </div>
          <XField v-slot="f" label="Пароль">
            <XInput v-bind="f" v-model="form.password" type="password" mono autocomplete="new-password" />
          </XField>
          <XField v-slot="f" label="Правила" hint="Enter — из списка, запятая — свой вариант">
            <XTagInput v-bind="f" v-model="form.tags" :presets="geo" mono placeholder="geoip:… или geosite:…" />
          </XField>
          <XField v-slot="f" label="Объявление">
            <XTextarea v-bind="f" v-model="form.note" />
          </XField>
          <XSection v-model:open="sectionOpen" title="Sniffing" description="Определение протокола по первым байтам соединения.">
            <template #meta><XSwitch v-model="form.sniffing" label="Sniffing" /></template>
            <XField layout="row" label="Только маршрутизация">
              <XSwitch v-model="form.routeOnly" label="Только маршрутизация" />
            </XField>
            <XCheckbox v-model="form.agree">Переопределять адрес назначения</XCheckbox>
          </XSection>
          <div class="kit__row kit__row--end">
            <XButton>Отмена</XButton>
            <XButton variant="primary" icon="check">Сохранить</XButton>
          </div>
        </div>
      </XCard>

      <XCard title="Окна и поведение">
        <div class="kit__stack">
          <p class="kit__p">Последнее действие меню: <code>{{ lastAction }}</code></p>
          <div class="kit__row kit__row--wrap">
            <XSegmented
              v-model="sheetSize"
              label="Размер окна"
              :options="[{ value: 'sm', label: 'sm 440' }, { value: 'md', label: 'md 560' }, { value: 'lg', label: 'lg 720' }, { value: 'xl', label: 'xl 800' }]"
            />
            <XButton variant="primary" @click="modalOpen = true">Открыть окно</XButton>
            <XButton @click="logOpen = true">Окно на весь экран телефона</XButton>
            <XButton @click="confirmOpen = true">Подтверждение</XButton>
          </div>
          <div class="kit__row kit__row--wrap">
            <XMenu v-slot="{ props }" :items="levelItems" title="Уровень" placement="bottom-start" sheet @select="level = $event">
              <XButton v-bind="props" caret>Меню-лист на телефоне</XButton>
            </XMenu>
            <XMenu v-slot="{ props }" :items="levelItems" title="Уровень" placement="bottom-start" scrim @select="level = $event">
              <XButton v-bind="props" caret>Меню с затемнением</XButton>
            </XMenu>
          </div>
          <XField layout="row" label="Оптимистичный переключатель" hint="Спиннер на время запроса">
            <XSwitch v-model="sw.b" :loading="sw.loading" label="Применить" @update:model-value="applySwitch" />
          </XField>
        </div>
      </XCard>
    </div>

    <XModal v-model:open="modalOpen" title="Создать подключение" :size="sheetSize" cancel>
      <XField v-slot="f" label="Примечание">
        <XInput v-bind="f" v-model="form.remark" autofocus />
      </XField>
      <XField v-slot="f" label="Протокол">
        <XSelect v-bind="f" v-model="form.strategy" :options="strategies" />
      </XField>
      <XField layout="row" label="Включить">
        <XSwitch v-model="form.enabled" label="Включить" />
      </XField>
      <template #footer>
        <XButton variant="primary" @click="modalOpen = false">Создать</XButton>
      </template>
    </XModal>

    <XConfirm v-model:open="confirmOpen" title="Восстановить базу данных?" cta="Восстановить" :busy="confirmBusy" busy-label="Восстановление…" @confirm="runConfirm">
      Текущая база данных будет заменена файлом backup.db, после этого панель перезапустится.
    </XConfirm>

    <XModal v-model:open="logOpen" title="Журнал панели" size="xl" flush full>
      <template #subtitle>Подзаголовок · <span class="kit__mono">flush</span>: тело без отступов, <span class="kit__mono">full</span>: весь экран телефона</template>
      <XEmpty title="Нет записей" hint="Тело окна рисует сам экран" />
    </XModal>
  </div>
</template>

<style scoped>
.kit {
  display: flex;
  flex-direction: column;
  gap: 24px;
  max-width: 1280px;
  margin: 0 auto;
  padding: 32px 24px 56px;
}

.kit__top {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
}

.kit__title {
  margin: 0;
  font: 600 var(--fs-24) / var(--lh-tight) var(--font-sans);
  letter-spacing: -.4px;
}

.kit__sub,
.kit__p {
  margin: 4px 0 0;
  color: var(--mut);
}

.kit__h3 {
  margin: 0 0 10px;
  color: var(--dim);
  font: 500 var(--fs-11) / 1 var(--font-mono);
  letter-spacing: .07em;
  text-transform: uppercase;
}

.kit__tokens {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 24px;
}

.kit__swatches {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.kit__swatch {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
  font-size: var(--fs-10);
  color: var(--dim);
}

.kit__swatch-color {
  width: 44px;
  height: 32px;
  border: 1px solid var(--line);
  border-radius: var(--r-8);
}

.kit__radius {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border: 1px solid var(--acc-line);
  background: var(--acc-soft);
  color: var(--acc);
  font: 500 var(--fs-11) var(--font-mono);
}

.kit__fonts {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.kit__mono {
  margin-top: 6px;
  color: var(--fg2);
}

.kit__spacing {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.kit__space {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: var(--fs-10);
  color: var(--dim);
}

.kit__space span {
  height: 8px;
  border-radius: var(--r-pill);
  background: var(--acc);
}

.kit__table-wrap {
  overflow-x: auto;
}

.kit__table {
  width: 100%;
  min-width: 1180px;
  border-collapse: collapse;
}

.kit__table th,
.kit__table td {
  padding: 12px 14px;
  border-bottom: 1px solid var(--div);
  text-align: left;
  vertical-align: middle;
}

.kit__table td {
  width: 13%;
  border-left: 1px solid var(--div);
}

.kit__table thead th {
  color: var(--dim);
  font: 500 var(--fs-10) / 1 var(--font-mono);
  letter-spacing: .07em;
  text-transform: uppercase;
}

.kit__table tbody th {
  min-width: 150px;
}

.kit__table tbody th b {
  display: block;
  font-weight: 600;
}

.kit__table tbody th code {
  color: var(--dim);
  font-size: var(--fs-11);
}

.kit__table tbody tr:last-child > * {
  border-bottom: 0;
}

.kit__row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.kit__row--wrap {
  flex-wrap: wrap;
}

.kit__row--end {
  justify-content: flex-end;
}

.kit__stack {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 10px;
}

.kit__stack > .x-notice,
.kit__stack > .x-field {
  align-self: stretch;
}

.kit__variants {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(340px, 1fr));
  gap: 32px;
}

.kit__variant {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.kit__scroll {
  max-width: 400px;
  height: 140px;
  padding: 8px 12px;
  overflow: auto;
  border: 1px solid var(--line);
  border-radius: var(--r-12);
  color: var(--mut);
  font-size: var(--fs-12);
  line-height: var(--lh-code);
  white-space: nowrap;
}

.kit__live {
  display: grid;
  grid-template-columns: minmax(0, 3fr) minmax(0, 2fr);
  gap: 24px;
  align-items: start;
}

.kit__form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.kit__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

@media (max-width: 960px) {
  .kit__live {
    grid-template-columns: minmax(0, 1fr);
  }
}

@media (max-width: 767px) {
  .kit {
    padding: 20px 16px 40px;
  }

  .kit__grid {
    grid-template-columns: minmax(0, 1fr);
  }

  .kit__variants {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
