<template>
  <aside>
    <div id="sidebar" class="nav-collapse">
      <ul id="nav-accordion" class="sidebar-menu">
        <li>
          <nuxt-link id="project-name" to="/project">
            <span>{{ this.$store.state.current_project }}</span>
            <i class="fa fa-edit"></i>
          </nuxt-link>
        </li>
        <li><nuxt-link to="/project/detail">Project Setting</nuxt-link></li>
        <li><nuxt-link to="/home">Role</nuxt-link></li>
        <li><nuxt-link to="/home">User</nuxt-link></li>
        <li><nuxt-link to="/home">Client</nuxt-link></li>
      </ul>
    </div>
  </aside>
</template>
