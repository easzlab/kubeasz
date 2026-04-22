<template>
  <el-config-provider :locale="elLocale">
    <router-view />
  </el-config-provider>
</template>

<script setup>
import { computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import en from 'element-plus/dist/locale/en.mjs'
import zhTw from 'element-plus/dist/locale/zh-tw.mjs'
import fr from 'element-plus/dist/locale/fr.mjs'

const { locale } = useI18n()

const elLocales = {
  'zh-CN': zhCn,
  en,
  'zh-TW': zhTw,
  fr
}

const elLocale = computed(() => elLocales[locale.value] || en)

watch(locale, (val) => {
  localStorage.setItem('ksk8s_locale', val)
})
</script>
