<template>
  <div class="wrapper">
    <h2>New Project Info</h2>
    <div class="form-panel">
      <label for="name" class="col-sm-2 control-label elem">
        Name
        <span class="required">*</span>
      </label>
      <div class="col-sm-5 elem">
        <input v-model="name" type="text" class="form-control" />
        <span class="help-block">
          TODO(help message)
        </span>
      </div>

      <div class="divider"></div>

      <button class="btn btn-theme">Create</button>
      <nuxt-link to="/project">Cancel</nuxt-link>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      name: ''
    }
  }
}
</script>

<style scoped>
.required {
  color: #ee2222;
  font-size: 18px;
}

.elem {
  float: left;
}

.divider {
  clear: both;
  border-bottom: 1px solid #eff2f7;
  padding-bottom: 15px;
  margin-bottom: 15px;
}
</style>
