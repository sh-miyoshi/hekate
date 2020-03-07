<template>
  <div
    class="c-sidebar c-sidebar-dark c-sidebar-fixed c-sidebar-lg-show"
    :class="{ 'c-sidebar-minimized': !showSidebar }"
  >
    <ul v-if="showSidebar" class="c-sidebar-nav">
      <li class="c-sidebar-nav-item">
        <nuxt-link to="/project" class="c-sidebar-nav-link project">
          <span class="mr-2">{{ this.$store.state.current_project }}</span>
          <i class="fa fa-edit"></i>
        </nuxt-link>
      </li>
      <li class="c-sidebar-nav-item">
        <nuxt-link to="/project/detail" class="c-sidebar-nav-link">
          Project Setting
        </nuxt-link>
      </li>
      <li class="c-sidebar-nav-item">
        <nuxt-link to="/user" class="c-sidebar-nav-link">User</nuxt-link>
      </li>
      <li class="c-sidebar-nav-item">
        <nuxt-link to="/home" class="c-sidebar-nav-link">Role</nuxt-link>
      </li>
      <li class="c-sidebar-nav-item">
        <nuxt-link to="/home" class="c-sidebar-nav-link">Client</nuxt-link>
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
