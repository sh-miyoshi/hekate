<template>
  <div
    class="c-sidebar c-sidebar-dark c-sidebar-fixed c-sidebar-lg-show"
    :class="{ 'c-sidebar-minimized': !showSidebar }"
  >
    <ul v-if="showSidebar" class="c-sidebar-nav">
      <li class="c-sidebar-nav-item">
        <div class="c-sidebar-nav-link project">&nbsp;</div>
      </li>
      <li class="c-sidebar-nav-item">
        <nuxt-link
          :to="{ name: 'user-project-name-info', params: { name: project } }"
          class="c-sidebar-nav-link"
        >
          User Info
        </nuxt-link>
      </li>
      <li class="c-sidebar-nav-item">
        <nuxt-link
          :to="{
            name: 'user-project-name-password',
            params: { name: project }
          }"
          class="c-sidebar-nav-link"
        >
          Password
        </nuxt-link>
      </li>
      <li class="c-sidebar-nav-item">
        <nuxt-link
          :to="{
            name: 'user-project-name-mfa',
            params: { name: project }
          }"
          class="c-sidebar-nav-link"
        >
          MFA
        </nuxt-link>
      </li>
    </ul>
    <button
      class="c-sidebar-minimizer c-class-toggler"
      @click="showSidebar = !showSidebar"
    ></button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      showSidebar: true
    }
  }
}
</script>

<style scoped>
.project {
  font-size: 130%;
  border-bottom: 1px solid;
}
</style>
