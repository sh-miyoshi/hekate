<template>
  <div class="container">
    <div class="text-center py-4">
      <h1>
        <template v-if="error.statusCode === 404">
          404 Page Not Found
        </template>
        <template v-else-if="error.statusCode === 500">
          500 Internal Server Error
        </template>
        <template v-else>
          Unknown Error
        </template>
      </h1>
      <hr class="col-xs-12" />
      <div v-if="error.message !== ''" class="m-2">
        {{ error.message }}
      </div>
      <div>
        <a href="/" class="btn btn-primary">TOP</a>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    error: {
      type: Object,
      required: true
    }
  },
  layout: 'index'
}
</script>
