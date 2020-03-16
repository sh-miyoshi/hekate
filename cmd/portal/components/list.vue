<template>
  <div>
    <ul>
      <li v-for="item in current" :key="item">
        {{ item }}
      </li>
    </ul>
    <ul>
      <li v-for="item in all" :key="item">
        {{ item }}
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  props: {
    current: {
      type: Array,
      required: true
    },
    all: {
      type: Array,
      required: true
    }
  }
}
</script>
