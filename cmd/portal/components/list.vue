<template>
  <div>
    <ul>
      <li v-for="item in current" :key="item">
        {{ item }}
      </li>
    </ul>
    <div>
      <b-form-select
        v-model="appendRole"
        :options="appendCandidates"
      ></b-form-select>
    </div>
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
  },
  data() {
    return {
      appendCandidates: this.getCandidates(),
      appendRole: null
    }
  },
  methods: {
    getCandidates() {
      const res = [{ value: null, text: 'Please select an append role' }]
      for (const item of this.all) {
        if (!this.current.includes(item)) {
          res.push({ value: item, text: item })
        }
      }
      return res
    }
  }
}
</script>
