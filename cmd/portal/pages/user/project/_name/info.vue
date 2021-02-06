<template>
  <div>
    <h4>User Info</h4>
    <p>user id</p>
    <p>user name</p>
  </div>
</template>

<script>
export default {
  layout: 'user',
  middleware: 'userAuth'
}
</script>
