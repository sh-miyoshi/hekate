<template>
  <div class="container index-window">
    <div class="c-wrapper">
      <div class="c-body">
        <h1 class="h1">HEKATE</h1>
        <p><i>simple identity management server</i></p>
        <div class="row">
          <div class="col-sm-6">
            <div class="card floa">
              <div class="card-header text-left">
                <i class="fa fa-user mr-1"></i>
                Admin Console
              </div>
              <div class="card-body">
                <div class="mb-3 text-left">
                  Centrally manage all aspects of the Hekate
                </div>
                <button class="btn btn-primary" @click="login">
                  Login
                </button>
              </div>
            </div>
          </div>
          <div class="col-sm-6">
            <div class="card">
              <div class="card-header text-left">
                <i class="fa fa-file-text-o mr-1"></i>
                Documentation
              </div>
              <div class="card-body">
                <div class="mb-3 text-left">
                  User guide, API specification, ...
                </div>
                <a href="https://github.com/sh-miyoshi/hekate"
                  >Go to Hekate Project</a
                >
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <footer class="c-footer"></footer>
  </div>
</template>

<script>
export default {
  layout: 'index',
  methods: {
    login() {
      this.$auth.Login('master')
    }
  }
}
</script>

<style scoped></style>
