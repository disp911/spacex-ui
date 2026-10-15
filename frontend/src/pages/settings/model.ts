// The panel's settings as POST panel/setting/all returns them and
// panel/setting/update takes them back (web/entity AllSetting).

import type { MessageSchema } from '../../i18n/schema.gen'

type Text = `settings.${keyof MessageSchema['settings']}`

export interface AllSetting {
  webListen: string
  webDomain: string
  webPort: number
  webCertFile: string
  webKeyFile: string
  webBasePath: string
  sessionMaxAge: number

  pageSize: number
  expireDiff: number
  trafficDiff: number
  remarkModel: string

  timeLocation: string
  twoFactorEnable: boolean
  twoFactorToken: string

  subEnable: boolean
  subJsonEnable: boolean
  subTitle: string
  subSupportUrl: string
  subProfileUrl: string
  subAnnounce: string
  subEnableRouting: boolean
  subRoutingRules: string
  subListen: string
  subPort: number
  subPath: string
  subDomain: string
  subCertFile: string
  subKeyFile: string
  subUpdates: number
  externalTrafficInformEnable: boolean
  externalTrafficInformURI: string
  restartXrayOnClientDisable: boolean
  subEncrypt: boolean
  subShowInfo: boolean
  subURI: string
  subJsonPath: string
  subJsonURI: string
  subClashEnable: boolean
  subClashPath: string
  subClashURI: string
  subJsonFragment: string
  subJsonNoises: string
  subJsonMux: string
  subJsonRules: string

  ldapEnable: boolean
  ldapHost: string
  ldapPort: number
  ldapUseTLS: boolean
  ldapBindDN: string
  ldapPassword: string
  ldapBaseDN: string
  ldapUserFilter: string
  ldapUserAttr: string
  ldapVlessField: string
  ldapSyncCron: string
  ldapFlagField: string
  ldapTruthyValues: string
  ldapInvertFlag: boolean
  ldapInboundTags: string
  ldapAutoCreate: boolean
  ldapAutoDelete: boolean
  ldapDefaultTotalGB: number
  ldapDefaultExpiryDays: number
  ldapDefaultLimitIP: number
}

export type SettingKey = keyof AllSetting

export type TabKey = 'panel' | 'ui' | 'account' | 'ldap' | 'sub' | 'formats'

/** The tabs in order, with the settings each edits: its dot and the status line count changes and errors among them. */
export const tabs: { key: TabKey; label: Text; keys: SettingKey[] }[] = [
  {
    key: 'panel',
    label: 'settings.tabPanel',
    keys: [
      'webListen', 'webDomain', 'webPort', 'webBasePath', 'sessionMaxAge', 'webCertFile', 'webKeyFile',
      'restartXrayOnClientDisable', 'externalTrafficInformEnable', 'externalTrafficInformURI',
    ],
  },
  { key: 'ui', label: 'settings.tabUi', keys: ['pageSize', 'timeLocation', 'expireDiff', 'trafficDiff'] },
  { key: 'account', label: 'settings.tabAccount', keys: [] },
  {
    key: 'ldap',
    label: 'settings.tabLdap',
    keys: [
      'ldapEnable', 'ldapHost', 'ldapPort', 'ldapUseTLS', 'ldapBindDN', 'ldapPassword', 'ldapBaseDN', 'ldapUserFilter', 'ldapUserAttr',
      'ldapVlessField', 'ldapFlagField', 'ldapTruthyValues', 'ldapInvertFlag', 'ldapSyncCron', 'ldapInboundTags',
      'ldapAutoCreate', 'ldapAutoDelete', 'ldapDefaultTotalGB', 'ldapDefaultExpiryDays', 'ldapDefaultLimitIP',
    ],
  },
  {
    key: 'sub',
    label: 'settings.tabSub',
    keys: [
      'remarkModel', 'subEnable', 'subJsonEnable', 'subClashEnable', 'subListen', 'subDomain', 'subPort', 'subPath', 'subURI',
      'subEncrypt', 'subShowInfo', 'subTitle', 'subSupportUrl', 'subProfileUrl', 'subAnnounce', 'subEnableRouting', 'subRoutingRules',
      'subCertFile', 'subKeyFile', 'subUpdates',
    ],
  },
  {
    key: 'formats',
    label: 'settings.tabFormats',
    keys: ['subJsonPath', 'subJsonURI', 'subClashPath', 'subClashURI', 'subJsonFragment', 'subJsonNoises', 'subJsonMux', 'subJsonRules'],
  },
]

/** The URI path a subscription is served at: the reverse proxy's when one is set. */
function subPathOf(uri: string, path: string) {
  try {
    return uri ? new URL(uri).pathname : path
  } catch {
    return path
  }
}

/** What the security banner warns about in the saved settings, as translation keys. */
export function securityAlerts(s: AllSetting) {
  const alerts: Text[] = []
  if (!s.webCertFile || !s.webKeyFile) alerts.push('settings.securityTls')
  if (s.webPort === 2053) alerts.push('settings.securityPort')
  if (s.webBasePath === '/') alerts.push('settings.securityPath')
  if (s.subEnable && subPathOf(s.subURI, s.subPath) === '/sub/') alerts.push('settings.securitySubPath')
  if (s.subJsonEnable && subPathOf(s.subJsonURI, s.subJsonPath) === '/json/') alerts.push('settings.securityJsonPath')
  return alerts
}
