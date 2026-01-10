<template>
  <Teleport to="body">
    <div v-if="isOpen" class="modal-overlay" @click.self="close">
      <div class="modal">
        <div class="modal-header">
          <h3 :class="type === 'error' ? 'text-danger' : 'text-primary'">{{ title }}</h3>
        </div>
        <div class="modal-body">
          <p>{{ message }}</p>
        </div>
        <div class="modal-footer">
          <button @click="close" class="btn-primary">确定</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
defineProps({
  isOpen: Boolean,
  title: { type: String, default: '提示' },
  message: { type: String, required: true },
  type: { type: String, default: 'info' } // info, error, success
})

const emit = defineEmits(['close'])

function close() {
  emit('close')
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0; 
  left: 0; 
  right: 0; 
  bottom: 0; 
  background: rgba(0,0,0,0.5); 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  z-index: 2000;
  animation: fadeIn 0.2s;
}

.modal {
  background: white;
  padding: 24px;
  border-radius: 12px;
  width: 400px;
  box-shadow: 0 20px 25px -5px rgba(0,0,0,0.1);
  animation: slideIn 0.2s;
}

.modal-header h3 {
  margin: 0 0 16px;
  font-size: 1.25rem;
}
.text-danger { color: #dc2626; }
.text-primary { color: #2563eb; }

.modal-body p {
  color: #64748b;
  font-size: 1rem;
  line-height: 1.5;
  margin-bottom: 24px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
}

.btn-primary {
  background: #2563eb;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
}
.btn-primary:hover {
  background: #1d4ed8;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideIn {
  from { transform: translateY(-20px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}
</style>
