<template>
  <div>
    <div v-for="item in current" :key="item">
      <div class="mb-1 input-group-text role-item">
        <i class="fa fa-remove"></i>
        {{ item }}
      </div>
    </div>
    <div>
      <b-form-select
        v-model="assignedRole"
        :options="assignedCandidates"
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
      assignedCandidates: this.getCandidates(),
      assignedRole: null
    }
  },
  methods: {
    getCandidates() {
      const res = [{ value: null, text: 'Please select an assigned role' }]
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

<style scoped>
.role-item {
  display: inline-block;
  padding: 0.175rem 0.55rem;
  width: auto;
}
</style>
